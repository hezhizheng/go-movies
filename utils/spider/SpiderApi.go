package spider

import (
	"fmt"
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

const ApiHost = "https://api.okzy.tv/api.php/provide/vod"
const AcList = "list"
const AcDetail = "detail"

type Lists struct {
	VodId         int    `json:"vod_id"` // 如果json中vod_id不是“1”，而是 1 ，这里一定要声明为 int ！！！fuck 不愧是静态强类型
	VodName       string `json:"vod_name"`
	TypeId        int    `json:"type_id"`
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

var (
	Smutex sync.Mutex
	wg     sync.WaitGroup
)

func StartApi() {
	list(1)
}

func list(pg int) {
	// 执行时间标记
	startTime := time.Now()
	defer ants.Release()
	antPool, _ := ants.NewPool(300)
	for _, subCategoryId := range subCategoryIds() {
		wg.Add(1)
		pageCount, t := pageCount(subCategoryId, pg)

		antPool.Submit(func() {
			actionList(t, pg, pageCount)
			wg.Done()
		})

	}
	wg.Wait()

	// 结束时间标记
	endTime := time.Since(startTime)

	ExecSecondsS := strconv.FormatFloat(endTime.Seconds(), 'f', -1, 64)
	ExecMinutesS := strconv.FormatFloat(endTime.Minutes(), 'f', -1, 64)
	ExecHoursS := strconv.FormatFloat(endTime.Hours(), 'f', -1, 64)

	fmt.Println("执行完成......")

	// 删除已缓存的页面
	go DelAllListCacheKey()

	// 钉钉通知
	SendDingMsg("本次爬虫执行时间为：" + ExecSecondsS + "秒 \n 即" + ExecMinutesS + "分钟 \n 即" + ExecHoursS + "小时 \n " + runtime.GOOS)

}

func actionList(subCategoryId string, pg int, pageCount int) {

	for j := pg; j <= pageCount; j++ {

		url := ApiHost + "?ac=" + AcList + "&t=" + subCategoryId + "&pg=" + strconv.Itoa(j)
		log.Println("当前page"+strconv.Itoa(j), url)

		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer func() {
			// 用完需要释放资源
			fasthttp.ReleaseResponse(resp)
			fasthttp.ReleaseRequest(req)
		}()

		// 默认是application/x-www-form-urlencoded
		//req.Header.SetContentType("application/json")
		req.Header.SetMethod("GET")

		req.SetRequestURI(url)

		//requestBody := []byte(`{"code":1}`)
		//req.SetBody(requestBody)

		if err := fasthttp.Do(req, resp); err != nil {
			fmt.Println("请求失败:", err.Error())
			return
		}

		body := resp.Body()

		var nav ResData
		err := utils.Json.Unmarshal(body, &nav)
		if err != nil {
			log.Println(err)
		}

		for _, value := range nav.List {
			// 模板时间
			timeTemplate := "2006-01-02 15:04:05"
			stamp1, _ := time.ParseInLocation(timeTemplate, value.VodTime, time.Local)

			utils.RedisDB.ZAdd("detail_links:id:"+strconv.Itoa(value.TypeId), &redis.Z{
				Score:  float64(stamp1.Unix()),
				Member: `/?m=vod-detail-id-` + strconv.Itoa(value.VodId) + `.html`,
			})

			film := []int{6, 7, 8, 9, 10, 11, 12, 20, 21, 37}
			tv := []int{13, 14, 15, 16, 22, 23, 24}
			cartoon := []int{29, 30, 31, 32, 33}

			if inType(value.TypeId, film) {
				utils.RedisDB.ZAdd("detail_links:id:1", &redis.Z{
					Score:  float64(stamp1.Unix()),
					Member: `/?m=vod-detail-id-` + strconv.Itoa(value.VodId) + `.html`,
				})
			}

			if inType(value.TypeId, tv) {
				utils.RedisDB.ZAdd("detail_links:id:2", &redis.Z{
					Score:  float64(stamp1.Unix()),
					Member: `/?m=vod-detail-id-` + strconv.Itoa(value.VodId) + `.html`,
				})
			}

			if inType(value.TypeId, cartoon) {
				utils.RedisDB.ZAdd("detail_links:id:4", &redis.Z{
					Score:  float64(stamp1.Unix()),
					Member: `/?m=vod-detail-id-` + strconv.Itoa(value.VodId) + `.html`,
				})
			}

			// 获取详情
			Detail(strconv.Itoa(value.VodId))

		}
	}

}

func pageCount(subCategoryId string, pg int) (int, string) {
	url := ApiHost + "?ac=" + AcList + "&t=" + subCategoryId + "&pg=" + strconv.Itoa(pg)

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.Header.SetMethod("GET")

	req.SetRequestURI(url)

	//requestBody := []byte(`{"code":1}`)
	//req.SetBody(requestBody)

	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Println("请求失败:", err.Error())
		return 0, subCategoryId
	}

	body := resp.Body()

	var nav ResData
	err := utils.Json.Unmarshal(body, &nav)
	if err != nil {
		log.Println(err)
	}

	log.Println("获取总页数", url, "total", nav.PageCount)
	return nav.PageCount, subCategoryId
}

// id与旧的网页爬虫对应不上
func Detail(id string) {
	// movies_detail:/?m=vod-detail-id-10051.html:movie_name:第102次相亲

	url := ApiHost + "?ac=" + AcDetail + "&ids=" + id + "&pg=1"

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	// 默认是application/x-www-form-urlencoded
	//req.Header.SetContentType("application/json")
	req.Header.SetMethod("GET")

	req.SetRequestURI(url)

	//requestBody := []byte(`{"code":1}`)
	//req.SetBody(requestBody)

	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Println("请求失败:", err.Error())
		return
	}

	body := resp.Body()

	var nav ResData
	err := utils.Json.Unmarshal(body, &nav)
	if err != nil {
		log.Println(err)
	}

	if len(nav.List) <= 0 {
		// todo 验证
		log.Println("没有list", url)
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
		k := map[string]string{
			"episode":   strconv.Itoa(ik + 1),
			"play_link": kuyunValue}
		Smutex.Lock()
		kuyunAry = append(kuyunAry, k)
		Smutex.Unlock()
	}

	for ic, ckm3u8Value := range ckm3u8 {
		c := map[string]string{
			"episode":   strconv.Itoa(ic + 1),
			"play_link": ckm3u8Value}
		Smutex.Lock()
		ckm3u8Ary = append(ckm3u8Ary, c)
		Smutex.Unlock()
	}

	for im, mp4Value := range mp4 {
		m := map[string]string{
			"episode":   strconv.Itoa(im + 1),
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

// 定义主类与子类的关系
func CategoriesStr() string {
	categories := `[
    {
        "link": "/?m=vod-type-id-1.html",
        "name": "电影片",
        "type_id": "1",
        "sub": [
            {
                "link": "/?m=vod-type-id-6.html",
                "name": "动作片",
                "type_id": "6",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-7.html",
                "name": "喜剧片",
                "type_id": "7",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-8.html",
                "name": "爱情片",
                "type_id": "8",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-9.html",
                "name": "科幻片",
                "type_id": "9",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-10.html",
                "name": "恐怖片",
                "type_id": "10",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-11.html",
                "name": "剧情片",
                "type_id": "11",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-12.html",
                "name": "战争片",
                "type_id": "12",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-20.html",
                "name": "纪录片",
                "type_id": "20",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-21.html",
                "name": "微电影",
                "type_id": "21",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-37.html",
                "name": "伦理片",
                "type_id": "37",
                "sub": null
            }
        ]
    },
    {
        "link": "/?m=vod-type-id-2.html",
        "name": "连续剧",
        "type_id": "2",
        "sub": [
            {
                "link": "/?m=vod-type-id-13.html",
                "name": "国产剧",
                "type_id": "13",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-14.html",
                "name": "香港剧",
                "type_id": "14",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-15.html",
                "name": "韩国剧",
                "type_id": "15",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-16.html",
                "name": "欧美剧",
                "type_id": "16",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-22.html",
                "name": "台湾剧",
                "type_id": "22",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-23.html",
                "name": "日本剧",
                "type_id": "23",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-24.html",
                "name": "海外剧",
                "type_id": "24",
                "sub": null
            }
        ]
    },
    {
        "link": "/?m=vod-type-id-4.html",
        "name": "动漫片",
        "type_id": "4",
        "sub": [
            {
                "link": "/?m=vod-type-id-29.html",
                "name": "国产动漫",
                "type_id": "29",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-30.html",
                "name": "日韩动漫",
                "type_id": "30",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-31.html",
                "name": "欧美动漫",
                "type_id": "31",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-32.html",
                "name": "港台动漫",
                "type_id": "32",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-33.html",
                "name": "海外动漫",
                "type_id": "33",
                "sub": null
            }
        ]
    }
]`

	return categories
}

// 获取
func subCategoryIds() []string {
	var nav []Categories
	categoriesStr := CategoriesStr()

	if utils.RedisDB.Exists("categories").Val() == 0 {
		utils.RedisDB.Set("categories", categoriesStr, 0).Err()
	}

	err := utils.Json.Unmarshal([]byte(categoriesStr), &nav)
	if err != nil {
		log.Println(err)
	}

	CategoryIds := make([]string, 0)
	for _, value := range nav {
		for _, subValue := range value.Sub {
			Smutex.Lock()
			CategoryIds = append(CategoryIds, subValue.TypeId)
			Smutex.Unlock()
		}
	}

	return CategoryIds
}

func FormatVodPlayUrl(VodPlayUrl string) ([]string, []string) {

	SplitVodPlayUrl := strings.Split(VodPlayUrl, "$$$")

	r, _ := regexp.Compile("https?://([\\w-]+\\.)+[\\w-]+(/[\\w-./?%&=]*)?")

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

	AllListCacheKey := utils.RedisDB.Keys("movie_lists_key:detail_links:*").Val()

	// 删除已经缓存的数据
	for _, val := range AllListCacheKey {
		utils.RedisDB.Del(val)
	}
}
