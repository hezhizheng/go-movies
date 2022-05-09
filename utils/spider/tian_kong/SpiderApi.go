package tian_kong

import (
	"github.com/go-redis/redis/v7"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"go_movies/utils"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const ApiHost = "https://api.tiankongapi.com/api.php/provide/vod"
const AcList = "list"
const AcDetail = "detail"

type SpiderApi struct {
	utils.SpiderTask
}

type Lists struct {
	VodId         int    `json:"vod_id"` // 如果json中vod_id不是“1”，而是 1 ，这里一定要声明为 int ！！！fuck 不愧是静态强类型
	VodName       string `json:"vod_name"`
	TypeId        int    `json:"type_id"`
	TypeId1       int    `json:"type_id_1"`
	TypeName      string `json:"type_name"`
	VodEn         string `json:"vod_en"`
	VodTime       string `json:"vod_time"`
	VodRemarks    string `json:"vod_remarks"`
	VodPlayFrom   string `json:"vod_play_from"`
	VodPlayUrl    string `json:"vod_play_url"`
	VodPDownUrl   string `json:"vod_down_url"`
	VodPic        string `json:"vod_pic"`
	VodArea       string `json:"vod_area"`
	VodDirector   string `json:"vod_director"`
	VodLang       string `json:"vod_lang"`
	VodYear       string `json:"vod_year"`
	VodSub        string `json:"vod_sub"`
	VodDuration   string `json:"vod_duration"`
	VodActor      string `json:"vod_actor"`
	VodContent    string `json:"vod_content"`
	VodPointsPlay int    `json:"vod_points_play"`
	VodScore      string `json:"vod_score"`
}

type ResData struct {
	Msg       string  `json:"msg"`
	Code      int     `json:"code"`
	Page      string  `json:"page"`
	PageCount int     `json:"pagecount"`
	Limit     string  `json:"limit"`
	Total     int     `json:"total"`
	List      []Lists `json:"list"`
}

type Categories struct {
	Link   string       `json:"link"`
	Name   string       `json:"name"`
	TypeId string       `json:"type_id"`
	Sub    []Categories `json:"sub"`
}

type FastHttp struct {
	f    fasthttp.Client
	req  *fasthttp.Request
	resp *fasthttp.Response
}

type CatePageCount struct {
	categoryId string
	PageCount  int
}

var (
	Smutex  sync.Mutex
	wg      sync.WaitGroup
	timeOut = 10 * time.Second // 请求超时时间
)

func (spiderApi *SpiderApi) Start() {
	go StartApi()
}

func (spiderApi *SpiderApi) PageDetail(id string) {
	go Detail(id, 0)
}

func (spiderApi *SpiderApi) DoRecentUpdate() {
	DoRecentUpdate()
}

func StartApi() {
	allMoviesDoneKeyExists := utils.RedisDB.Exists("all_movies_done").Val()
	if allMoviesDoneKeyExists > 0 {
		return
	}
	list(1)
}

func list(pg int) {
	// 执行时间标记
	startTime := time.Now()
	defer ants.Release()
	antPool, _ := ants.NewPool(200)

	//_f := initFastHttp()

	catePageCounts := getCategoryPageCount()

	log.Println(catePageCounts)

	for _, catePageCount := range catePageCounts {
		wg.Add(1)
		categoryId := catePageCount.categoryId
		PageCount := catePageCount.PageCount

		antPool.Submit(func() {
			// 这里不能直接使用 catePageCount.categoryId 、catePageCount.PageCount
			// 在 submit 之前赋值变量传进来
			actionList(categoryId, pg, PageCount)
			wg.Done()
		})

	}

	wg.Wait()

	// 结束时间标记
	endTime := time.Since(startTime)

	ExecSecondsS := strconv.FormatFloat(endTime.Seconds(), 'f', 2, 64)
	ExecMinutesS := strconv.FormatFloat(endTime.Minutes(), 'f', 2, 64)
	ExecHoursS := strconv.FormatFloat(endTime.Hours(), 'f', 2, 64)

	log.Println("执行完成......")

	// 删除已缓存的页面
	go DelAllListCacheKey()

	// 全量 done -> set done 永久Redis 标识 -> new corntab every min ( done key exist && recent_update_key expire ) -> set recent_update_key 1h expire -> do recent 3h update
	// 一周进行一次全量爬取，资源网站的电影ID是会变的，fuck!!!
	utils.RedisDB.SetNX("all_movies_done", "done", time.Second*604800)

	// 钉钉通知
	SendDingMsg("本次爬虫执行时间为：" + ExecSecondsS + "秒 \n 即" + ExecMinutesS + "分钟 \n 即" + ExecHoursS + "小时 \n " + runtime.GOOS)
}

func DoRecentUpdate() {
	allMoviesDoneKeyExists := utils.RedisDB.Exists("all_movies_done").Val()
	recentUpdateKeyExists := utils.RedisDB.Exists("recent_update_key").Val()

	if allMoviesDoneKeyExists > 0 && recentUpdateKeyExists == 0 {
		startTime := time.Now()

		utils.RedisDB.SetNX("recent_update_key", "done", time.Second*3600).Err()

		actionRecentUpdateList()

		// 结束时间标记
		endTime := time.Since(startTime)
		ExecSecondsS := strconv.FormatFloat(endTime.Seconds(), 'f', 2, 64)
		ExecMinutesS := strconv.FormatFloat(endTime.Minutes(), 'f', 2, 64)
		ExecHoursS := strconv.FormatFloat(endTime.Hours(), 'f', 2, 64)

		SendDingMsg("最近更新执行完成，耗时：" + ExecSecondsS + "秒 \n 即" + ExecMinutesS + "分钟 \n 即" + ExecHoursS + "小时 \n " + runtime.GOOS)
	}
}

func actionRecentUpdateList() {
	defer ants.Release()
	antPool, _ := ants.NewPool(200)

	film := GetAssignCategoryIds("film")
	tv := GetAssignCategoryIds("tv")
	cartoon := GetAssignCategoryIds("cartoon")
	saveTypes := GetIntSubCategoryIds()

	pageCount := RecentUpdatePageCount(0)
	//pageCount := 5
	wg2 := sync.WaitGroup{}
	wg2.Add(pageCount)
	for _j := 1; _j <= pageCount; _j++ {
		j := _j
		// 使用 goroutine 执行的任务
		task := func() {
			//log.Println("jj",j)
			//time.Sleep(time.Second * 2)
			//return
			url := ApiHost + "?h=6" + "&pg=" + strconv.Itoa(j)
			_, resp, gErr := fasthttp.Get(nil, url)
			if gErr != nil {
				log.Println("actionRecentUpdateList 请求失败:", gErr.Error())
				return
			}

			var nav ResData
			err := utils.Json.Unmarshal(resp, &nav)
			if err != nil {
				log.Println("actionRecentUpdateList json 序列化错误", err)
				return
			}

			for _, value := range nav.List {
				if inType(value.TypeId, saveTypes) {
					//log.Println("value.TypeId",value.TypeId)
					//continue
					// 模板时间
					timeTemplate := "2006-01-02 15:04:05"
					stamp1, _ := time.ParseInLocation(timeTemplate, value.VodTime, time.Local)

					utils.RedisDB.ZAdd("detail_links:id:"+strconv.Itoa(value.TypeId), &redis.Z{
						Score:  float64(stamp1.Unix()),
						Member: `/?m=vod-detail-id-` + strconv.Itoa(value.VodId) + `.html`,
					})

					if inType(value.TypeId, film) {
						utils.RedisDB.ZAdd("detail_links:id:1", &redis.Z{
							Score:  float64(stamp1.Unix()),
							Member: `/?m=vod-detail-id-` + strconv.Itoa(value.VodId) + `.html`,
						})
						// 获取详情
						Detail(strconv.Itoa(value.VodId), 0)
					}

					if inType(value.TypeId, tv) {
						utils.RedisDB.ZAdd("detail_links:id:3", &redis.Z{
							Score:  float64(stamp1.Unix()),
							Member: `/?m=vod-detail-id-` + strconv.Itoa(value.VodId) + `.html`,
						})
						// 获取详情
						Detail(strconv.Itoa(value.VodId), 0)
					}

					if inType(value.TypeId, cartoon) {
						utils.RedisDB.ZAdd("detail_links:id:24", &redis.Z{
							Score:  float64(stamp1.Unix()),
							Member: `/?m=vod-detail-id-` + strconv.Itoa(value.VodId) + `.html`,
						})
						// 获取详情
						Detail(strconv.Itoa(value.VodId), 0)
					}
				}
			}
		}

		// 提交任务
		submitErr := antPool.Submit(func() {
			defer wg2.Done()
			task()
		})

		if submitErr != nil {
			log.Println("antPool submitErr：", submitErr)
		}
	}
	wg2.Wait()

	if pageCount > 0 {
		go DelAllListCacheKey()
	}

	log.Println("all actionRecentUpdateList done")
}

func RecentUpdatePageCount(retry int) int {
	if retry >= 4 {
		return 0
	}
	url := ApiHost + "?h=6&pg=1"

	_, resp, gErr := fasthttp.Get(nil, url)
	if gErr != nil {
		retry++
		log.Println("RecentUpdatePageCount 请求失败:", retry, url, gErr.Error())
		return RecentUpdatePageCount(retry)
	}

	var nav ResData
	err := utils.Json.Unmarshal(resp, &nav)
	if err != nil {
		log.Println(err, "RecentUpdatePageCount json解析失败")
		return 0
	}

	PageCount := nav.PageCount
	log.Println("获取最近更新总页数", url, "PageCount", PageCount)
	return PageCount
}

func actionList(subCategoryId string, pg int, pageCount int) {
	//return

	film := GetAssignCategoryIds("film")
	tv := GetAssignCategoryIds("tv")
	cartoon := GetAssignCategoryIds("cartoon")
	saveTypes := GetIntSubCategoryIds()

	for j := pg; j <= pageCount; j++ {
		url := ApiHost + "?ac=" + AcList + "&t=" + subCategoryId + "&pg=" + strconv.Itoa(j)
		log.Println("当前page"+strconv.Itoa(j), url, pageCount)

		_, resp, gErr := fasthttp.Get(nil, url)
		if gErr != nil {
			log.Println("actionList 请求失败:", url, gErr.Error())
			return
		}

		var nav ResData
		err := utils.Json.Unmarshal(resp, &nav)
		if err != nil {
			log.Println("actionList json 解析失败", url, err)
			return
		}

		for _, value := range nav.List {
			// 模板时间
			timeTemplate := "2006-01-02 15:04:05"
			stamp1, _ := time.ParseInLocation(timeTemplate, value.VodTime, time.Local)

			if inType(value.TypeId, saveTypes) {
				utils.RedisDB.ZAdd("detail_links:id:"+strconv.Itoa(value.TypeId), &redis.Z{
					Score:  float64(stamp1.Unix()),
					Member: `/?m=vod-detail-id-` + strconv.Itoa(value.VodId) + `.html`,
				})

				if inType(value.TypeId, film) {
					utils.RedisDB.ZAdd("detail_links:id:1", &redis.Z{
						Score:  float64(stamp1.Unix()),
						Member: `/?m=vod-detail-id-` + strconv.Itoa(value.VodId) + `.html`,
					})
					// 获取详情
					Detail(strconv.Itoa(value.VodId), 0)
				}

				if inType(value.TypeId, tv) {
					utils.RedisDB.ZAdd("detail_links:id:3", &redis.Z{
						Score:  float64(stamp1.Unix()),
						Member: `/?m=vod-detail-id-` + strconv.Itoa(value.VodId) + `.html`,
					})
					// 获取详情
					Detail(strconv.Itoa(value.VodId), 0)
				}

				if inType(value.TypeId, cartoon) {
					utils.RedisDB.ZAdd("detail_links:id:24", &redis.Z{
						Score:  float64(stamp1.Unix()),
						Member: `/?m=vod-detail-id-` + strconv.Itoa(value.VodId) + `.html`,
					})
					// 获取详情
					Detail(strconv.Itoa(value.VodId), 0)
				}
			}
		}
	}
	return
}

func pageCount(subCategoryId string, retry int) (int, string) {
	if retry >= 4 {
		return 0, subCategoryId
	}
	url := ApiHost + "?ac=" + AcList + "&t=" + subCategoryId + "&pg=1"

	_, resp, err := fasthttp.Get(nil, url)
	if err != nil {
		retry++
		log.Println("pageCount 请求失败:", retry, url, err.Error())
		return pageCount(subCategoryId, retry)
	}

	var nav ResData
	jErr := utils.Json.Unmarshal(resp, &nav)
	if jErr != nil {
		log.Println(jErr, "pageCount json 解析失败")
		return 0, subCategoryId
	}

	PageCount := nav.PageCount
	log.Println("获取总页数", url, "PageCount", PageCount, "subCategoryId", subCategoryId)
	return PageCount, subCategoryId
}

// id与旧的网页爬虫对应不上
func Detail(id string, retry int) {
	// movies_detail:/?m=vod-detail-id-10051.html:movie_name:第102次相亲
	url := ApiHost + "?ac=" + AcDetail + "&ids=" + id + "&pg=1"

	//retryMax := 3
	//if retry >= retryMax {
	//	log.Println("重试已结束", url, retry)
	//	return
	//}

	_, resp, gErr := fasthttp.Get(nil, url)
	if gErr != nil {
		log.Println("Detail 请求失败:", gErr.Error(), url)
		return
	}

	var nav ResData
	err := utils.Json.Unmarshal(resp, &nav)
	if err != nil {
		log.Println(err, "detail json 解析失败")
		return
	}

	if len(nav.List) <= 0 {
		log.Println("没有list", url)

		// 重试
		//for {
		//	if retry > retryMax {
		//		log.Println("超过最大重试次数，重试机制已跳出", url, retry)
		//		break
		//	}
		//	retry++
		//	Detail(id, retry)
		//	log.Println("正在重试...", url, retry)
		//}

		return
	}
	listDetail := nav.List[0]

	_moviesInfo := make(map[string]interface{})

	var kuyunAry []map[string]string

	var ckm3u8Ary []map[string]string
	//
	var downloadAry []map[string]string

	kuyun, ckm3u8 := FormatVodPlayUrl(listDetail.VodPlayUrl)

	mp4 := FormatVodPDownUrl(listDetail.VodPDownUrl)

	for ik, kuyunValue := range kuyun {
		episode := strconv.Itoa(ik + 1)
		episode = Lang(listDetail.TypeId1, kuyunValue, listDetail.VodPlayUrl, episode)
		k := map[string]string{
			"episode":   episode,
			"play_link": kuyunValue}
		Smutex.Lock()
		kuyunAry = append(kuyunAry, k)
		Smutex.Unlock()
	}

	for ic, ckm3u8Value := range ckm3u8 {
		episode := strconv.Itoa(ic + 1)
		episode = Lang(listDetail.TypeId1, ckm3u8Value, listDetail.VodPlayUrl, episode)
		c := map[string]string{
			"episode":   episode,
			"play_link": ckm3u8Value}
		Smutex.Lock()
		ckm3u8Ary = append(ckm3u8Ary, c)
		Smutex.Unlock()
	}

	for im, mp4Value := range mp4 {
		episode := strconv.Itoa(im + 1)
		episode = Lang(listDetail.TypeId1, mp4Value, listDetail.VodPDownUrl, episode)
		m := map[string]string{
			"episode":   episode,
			"play_link": mp4Value}
		Smutex.Lock()
		downloadAry = append(downloadAry, m)
		Smutex.Unlock()
	}

	kuyunAryJson, _ := utils.Json.MarshalIndent(kuyunAry, "", " ")
	ckm3u8AryJson, _ := utils.Json.MarshalIndent(ckm3u8Ary, "", " ")
	downloadAryJson, _ := utils.Json.MarshalIndent(downloadAry, "", " ")

	link := `/?m=vod-detail-id-` + strconv.Itoa(listDetail.VodId) + `.html`
	_moviesInfo["link"] = link
	_moviesInfo["cover"] = listDetail.VodPic
	_moviesInfo["name"] = listDetail.VodName
	_moviesInfo["quality"] = listDetail.VodRemarks
	_moviesInfo["score"] = listDetail.VodScore
	_moviesInfo["kuyun"] = string(kuyunAryJson)
	_moviesInfo["ckm3u8"] = string(ckm3u8AryJson)
	_moviesInfo["download"] = string(downloadAryJson)

	mDetail := make(map[string]interface{})
	mDetail["alias"] = listDetail.VodSub
	mDetail["director"] = listDetail.VodDirector
	mDetail["starring"] = listDetail.VodActor
	mDetail["type"] = listDetail.TypeName
	mDetail["area"] = listDetail.VodArea
	mDetail["language"] = listDetail.VodLang
	mDetail["released"] = listDetail.VodYear
	mDetail["length"] = listDetail.VodDuration
	mDetail["update"] = listDetail.VodTime
	mDetail["total_playback"] = listDetail.VodPointsPlay
	mDetail["vod_play_info"] = listDetail.VodContent

	_detail, _ := utils.Json.MarshalIndent(mDetail, "", " ")

	_moviesInfo["detail"] = string(_detail)

	t := utils.RedisDB.HMSet("movies_detail:"+link+":movie_name:"+listDetail.VodName, _moviesInfo).Err()
	log.Println("当前详情", url, t)
}

// 获取所有类别ID
func subCategoryIds() []string {
	var nav []Categories
	categoriesStr := CategoriesStr()

	if utils.RedisDB.Exists("categories").Val() == 0 {
		utils.RedisDB.Set("categories", categoriesStr, 0).Err()
	}

	CategoryIds := make([]string, 0)

	err := utils.Json.Unmarshal([]byte(categoriesStr), &nav)
	if err != nil {
		log.Println("subCategoryIds json 解析失败", err)
		return CategoryIds
	}

	for _, value := range nav {
		for _, subValue := range value.Sub {
			Smutex.Lock()
			CategoryIds = append(CategoryIds, subValue.TypeId)
			Smutex.Unlock()
		}
	}

	return CategoryIds
}

// 获取每个类别对应的总数
func getCategoryPageCount() []CatePageCount {
	subCategoryIds := subCategoryIds()

	var CatePageCounts []CatePageCount

	for _, subCategoryId := range subCategoryIds {

		pageCount, t := pageCount(subCategoryId, 0)

		CatePageCount := CatePageCount{
			categoryId: t,
			PageCount:  pageCount,
		}

		Smutex.Lock()
		CatePageCounts = append(CatePageCounts, CatePageCount)
		Smutex.Unlock()
	}

	return CatePageCounts
}

// 区分电影的广东话跟国语
func Lang(vodType int, resourcesUrl, allResource, episode string) string {

	if vodType == 1 {
		cantonese := "粤语$" + resourcesUrl
		mandarin := "国语$" + resourcesUrl
		hdCantonese := "HD粤语高清$" + resourcesUrl
		hdMandarin := "HD国语高清$" + resourcesUrl
		if strings.Contains(allResource, cantonese) ||
			strings.Contains(allResource, hdCantonese) {
			episode = "粤语"
		} else if strings.Contains(allResource, mandarin) ||
			strings.Contains(allResource, hdMandarin) {
			episode = "国语"
		}
	}

	return episode
}

func FormatVodPlayUrl(VodPlayUrl string) ([]string, []string) {

	SplitVodPlayUrl := strings.Split(VodPlayUrl, "$$$")

	r, _ := regexp.Compile("https?://([\\w-]+\\.)+[\\w-]+(/[\\w-./?%&=]*)?")

	// 这里剧集好像是 kuyun 在前面 [0] m3u8 在后面 [1]  ,电影则是相反的。。。
	// 暂时先不处理，直接在播放列表通过播放地址的后缀区分
	kuyun := r.FindAllString(SplitVodPlayUrl[0], -1)

	ckm3u8 := []string{""}
	if len(SplitVodPlayUrl) >= 2 {
		ckm3u8 = r.FindAllString(SplitVodPlayUrl[1], -1)
	}

	return kuyun, ckm3u8
}

func FormatVodPDownUrl(VodPDownUrl string) []string {

	// todo: 对中文之后的直接过滤掉了，干！
	// (https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]
	//r, _ := regexp.Compile("https?://([\\w-]+\\.)+[\\w-]+(/[\\w-./?%&=]*)?")
	//
	//mp4 := r.FindAllString(VodPDownUrl, -1)
	//
	//return mp4

	c := strings.Split(VodPDownUrl, "$")

	shift := c[1:] // 去掉第一个元素，一般是切割出来没用的

	urls := make([]string, 0)

	// http://xz3-7.okzyxz.com/20190524/23916_07fb2078/死亡地带S01E01.mp4#第02集
	// 处理链接后面的#号符
	for _, v := range shift {
		split := strings.Split(v, "#")
		Smutex.Lock()
		urls = append(urls, split[0])
		Smutex.Unlock()
	}

	return urls
}

func inType(s int, d []int) bool {
	for _, k := range d {
		if s == k {
			return true
		}
	}
	return false
}

func SendDingMsg(msg string) {
	accessToken := viper.GetString(`ding.access_token`)
	if accessToken == "" {
		return
	}
	webhook := "https://oapi.dingtalk.com/robot/send?access_token=" + accessToken
	robot := utils.NewRobot(webhook)

	title := "goMovies 爬虫通知API"
	text := "#### goMovies 爬虫通知API  \n " + msg
	atMobiles := []string{""}
	isAtAll := true

	err := robot.SendMarkdown(title, text, atMobiles, isAtAll)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("已发送钉钉通知")
}

func DelAllListCacheKey() {
	utils.RedisDB.Del("paginate")
}

func GetAssignCategoryIds(_type string) []int {

	ids := make([]int, 0)

	var nav []Categories
	categoriesStr := CategoriesStr()

	err := utils.Json.Unmarshal([]byte(categoriesStr), &nav)
	if err != nil {
		log.Println("GetAssignCategoryIds json 解析失败", err)
		return ids
	}

	film := nav[0].Sub
	tv := nav[1].Sub
	cartoon := nav[2].Sub

	switch _type {
	case "film":
		for _, v := range film {
			intId, _ := strconv.Atoi(v.TypeId)
			ids = append(ids, intId)
		}
	case "tv":
		for _, v := range tv {
			intId, _ := strconv.Atoi(v.TypeId)
			ids = append(ids, intId)
		}
	case "cartoon":
		for _, v := range cartoon {
			intId, _ := strconv.Atoi(v.TypeId)
			ids = append(ids, intId)
		}
	}
	return ids
}

func GetIntSubCategoryIds() []int {
	ids := make([]int, 0)

	var nav []Categories
	categoriesStr := CategoriesStr()

	err := utils.Json.Unmarshal([]byte(categoriesStr), &nav)
	if err != nil {
		log.Println("getIntSubCategoryIds json 解析失败", err)
		return ids
	}

	for _, value := range nav {
		for _, subValue := range value.Sub {
			Smutex.Lock()
			intId, _ := strconv.Atoi(subValue.TypeId)
			ids = append(ids, intId)
			Smutex.Unlock()
		}
	}

	return ids
}
