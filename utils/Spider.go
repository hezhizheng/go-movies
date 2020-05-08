package utils

import (
	"github.com/go-redis/redis/v7"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 爬取网站的域名，如访问不了，以下几个域名建议重复替换使用
//const host = "http://www.jisudhw.com"
const host = "http://www.okzy.co"

//const host = "http://www.okzyw.com"

// redis key

// 分类key
const CategoriesKey = "categories"

// 电影详情
const moviesDetail = "movies_detail:"

type Categories struct {
	Link string       `json:"link"`
	Name string       `json:"name"`
	Sub  []Categories `json:"sub"`
}

type Movies struct {
	Link      string `json:"link"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Cover     string `json:"cover"`
	UpdatedAt string `json:"updated_at"`
}

type MoviesDetail struct {
	Link     string                 `json:"link"`
	Name     string                 `json:"name"`
	Cover    string                 `json:"cover"`
	Quality  string                 `json:"quality"`
	Score    string                 `json:"score"`
	KuYun    string                 `json:"ku_yun"`
	CK       string                 `json:"ckm3u8"`
	Download string                 `json:"download"`
	Detail   map[string]interface{} `json:"detail"`
}

type Spider struct {
	SpiderTask
}

var (
	Smutex sync.Mutex
	wg     sync.WaitGroup
)

func (spider *Spider) Start() {
	go StartSpider()
}

func (spider *Spider) PageDetail(id string) {
	// /?m=vod-detail-id-56275.html
	url := "/?m=vod-detail-id-" + id + ".html"
	go MoviesInfo(url)
}

func StartSpider() {
	star := time.Now()
	defer ants.Release()

	// 获取所有分类
	Categories := SpiderOKCategories()

	antPoolStartSpiderSubCate, _ := ants.NewPool(25)

	antPoolStartSpider := antPoolStartSpiderSubCate

	// 爬取所有主类
	SpiderCategories(Categories, antPoolStartSpider)
	// 爬取主类对应的子类
	SpiderSubCategories(Categories, antPoolStartSpiderSubCate)

	wg.Wait()

	end := time.Since(star)

	ExecSecondsS := strconv.FormatFloat(end.Seconds(), 'f', -1, 64)
	ExecMinutesS := strconv.FormatFloat(end.Minutes(), 'f', -1, 64)
	ExecHoursS := strconv.FormatFloat(end.Hours(), 'f', -1, 64)

	// 清楚缓存的页面数据
	go DelAllListCacheKey()

	log.Println("本次爬虫执行时间为：" + ExecSecondsS + "秒 \n 即" + ExecMinutesS + "分钟 \n 即" + ExecHoursS + "小时 \n ")

	// 钉钉通知
	sendDingMsg("本次爬虫执行时间为：" + ExecSecondsS + "秒 \n 即" + ExecMinutesS + "分钟 \n 即" + ExecHoursS + "小时 \n " + runtime.GOOS)

}

func SpiderCategories(Categories []Categories, antPoolStartSpider *ants.Pool) {
	for _, v := range Categories {
		cateUrl := v.Link
		wg.Add(1)

		antPoolStartSpider.Submit(func() {
			SpiderOKMovies(cateUrl)
			wg.Done()
		})
	}
}

func SpiderSubCategories(Categories []Categories, antPoolStartSpiderSubCate *ants.Pool) {
	for _, v := range Categories {
		childrenCates := v.Sub

		for _, childrenCate := range childrenCates {
			wg.Add(1)
			childrenCateUrl := childrenCate.Link
			antPoolStartSpiderSubCate.Submit(func() {
				SpiderOKMovies(childrenCateUrl)
				wg.Done()
			})
		}
	}
}

// 爬取所有类别
func SpiderOKCategories() []Categories {

	c := colly.NewCollector(
		colly.Async(true),
	)

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	retryCount := 0

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(res *colly.Response, err error) {
		log.Println("Something went wrong:", err)
		if retryCount < 3 {
			retryCount += 1
			_retryErr := res.Request.Retry()
			log.Println("retry wrong:", _retryErr)
		}
	})

	// 主类
	Cate := make([]Categories, 0)

	// 导航栏、分类
	c.OnHTML("ul#sddm li", func(e *colly.HTMLElement) {

		categoryLink := e.ChildAttr("a", "href")

		categoryName := e.ChildText("a[onmouseout]")

		// 子类
		SubCate := make([]Categories, 0)

		e.ForEach("div a", func(i int, element *colly.HTMLElement) {

			subCategoryLink := element.Attr("href")
			subCategoryName := element.Text

			_subCate := Categories{
				Link: subCategoryLink,
				Name: subCategoryName,
			}

			if subCategoryName != categoryName {
				// 追加
				Smutex.Lock()
				SubCate = append(SubCate, _subCate)
				Smutex.Unlock()
			}

		})

		// 主类别
		_cate := Categories{
			Link: categoryLink,
			Name: categoryName,
			Sub:  SubCate,
		}

		// 去掉首页、福利、综艺片、解说 链接
		if categoryName != "" && categoryName != "福利片" && categoryName != "综艺片" && categoryName != "解说" {
			// 追加
			Smutex.Lock()
			Cate = append(Cate, _cate)
			Smutex.Unlock()
		}

	})

	// 在OnHTML之后被调用
	c.OnScraped(func(_ *colly.Response) {

		categories, _ := Json.MarshalIndent(Cate, "", " ")

		Smutex.Lock()
		err := RedisDB.Set(CategoriesKey, string(categories), 0).Err()
		Smutex.Unlock()
		log.Println(err)

	})

	visitError := c.Visit(host)

	log.Println(visitError)

	c.Wait()

	return Cate
}

// 爬取所有类别的电影
func SpiderOKMovies(cateUrl string) {

	c := colly.NewCollector(
		colly.Async(true),
	)

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	retryCount := 0
	c.OnError(func(res *colly.Response, err error) {
		log.Println("Something went wrong:", err)
		if retryCount < 3 {
			retryCount += 1
			_retryErr := res.Request.Retry()
			log.Println("retry wrong:", _retryErr)
		}
	})

	var lastPageInt int

	c.OnHTML(".pages input[type=button]", func(e *colly.HTMLElement) {

		lastPageStr := e.Attr("onclick")

		lastPageStrSplit := strings.Split(lastPageStr, ",")[1]

		// 最后一页
		lastPage, _ := strconv.Atoi(strings.Split(lastPageStrSplit, ")")[0])

		lastPageInt = lastPage // todo lastPage

		for j := 1; j <= lastPageInt; j++ {
			pageUrl := CategoryToPageUrl(cateUrl, strconv.Itoa(j))
			// 爬取所有主类下面的电影
			ForeachPage(cateUrl, pageUrl)
		}
	})

	visitError := c.Visit(host + cateUrl)

	log.Println(visitError)

	c.Wait()
}

// 获取电影详情信息
func ForeachPage(cateUrl string, url string) {

	c := colly.NewCollector(
		colly.Async(true),
	)

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	retryCount := 0
	c.OnError(func(res *colly.Response, err error) {
		log.Println("Something went wrong:", err)
		if retryCount < 3 {
			retryCount += 1
			_retryErr := res.Request.Retry()
			log.Println("retry wrong:", _retryErr)
		}
	})

	// 导航栏、分类
	c.OnHTML(".xing_vb li", func(e *colly.HTMLElement) {

		spanClass := e.ChildAttr("span", "class")

		// 列表数据
		if spanClass == "tt" {
			link := e.ChildAttr("a", "href")
			updateAt := e.ChildText(".xing_vb6")

			// 模板时间
			timeTemplate := "2006-01-02"
			stamp1, _ := time.ParseInLocation(timeTemplate, updateAt, time.Local)

			Smutex.Lock()
			RedisDB.ZAdd("detail_links:id:"+TransformId(cateUrl), &redis.Z{
				Score:  float64(stamp1.Unix()),
				Member: link,
			})
			Smutex.Unlock()

			// 获取详情
			MoviesInfo(link)
		}
	})

	visitError := c.Visit(host + url)

	log.Println(visitError)
	log.Println("当前页面")
	log.Println(url)
	c.Wait()
}

func MoviesInfo(url string) MoviesDetail {

	c := colly.NewCollector(
		colly.Async(true),
	)

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	retryCount := 0
	c.OnError(func(res *colly.Response, err error) {
		log.Println("Something went wrong:", err)
		if retryCount < 3 {
			retryCount += 1
			_retryErr := res.Request.Retry()
			log.Println("retry wrong:", _retryErr)
		}
	})

	// 所有电影
	md := MoviesDetail{}

	detail := make(map[string]interface{})

	var kuyunAry []map[string]string

	var ckm3u8Ary []map[string]string

	var downloadAry []map[string]string

	c.OnHTML(".warp", func(e *colly.HTMLElement) {

		cover := e.ChildAttr("div .vodImg>img", "src")
		name := e.ChildText("div .vodh>h2")
		quality := e.ChildText("div .vodh span")
		score := e.ChildText("div .vodh label")

		_type := ""
		e.ForEach("div .vodinfobox ul li", func(i int, element *colly.HTMLElement) {
			if i == 3 {
				_type = element.ChildText("span")
			}
		})

		// 有些页面 1 是 ckm3u8  2 是 kuyun  wtf!

		e.ForEach("div #1 ul li", func(i int, element *colly.HTMLElement) {

			playLink := element.ChildAttr("input", "value")

			Episode := strconv.Itoa(i + 1)
			Episode = transformEpisode(_type, Episode, element.Text)

			if strings.Index(playLink, "m3u8") == -1 {
				kuyun := map[string]string{
					"episode":   Episode,
					"play_link": playLink}

				Smutex.Lock()
				kuyunAry = append(kuyunAry, kuyun)
				Smutex.Unlock()
			} else {
				ckm3u8 := map[string]string{
					"episode":   Episode,
					"play_link": playLink}
				Smutex.Lock()
				ckm3u8Ary = append(ckm3u8Ary, ckm3u8)
				Smutex.Unlock()
			}

		})

		e.ForEach("div #2 ul li", func(i int, element *colly.HTMLElement) {

			playLink := element.ChildAttr("input", "value")

			Episode := strconv.Itoa(i + 1)
			Episode = transformEpisode(_type, Episode, element.Text)

			if strings.Index(playLink, "m3u8") == -1 {
				kuyun := map[string]string{
					"episode":   Episode,
					"play_link": playLink}

				Smutex.Lock()
				kuyunAry = append(kuyunAry, kuyun)
				Smutex.Unlock()
			} else {
				ckm3u8 := map[string]string{
					"episode":   Episode,
					"play_link": playLink}
				Smutex.Lock()
				ckm3u8Ary = append(ckm3u8Ary, ckm3u8)
				Smutex.Unlock()
			}
		})

		e.ForEach("div #down_1 ul li", func(i int, element *colly.HTMLElement) {

			playLink := element.ChildAttr("input", "value")

			Episode := strconv.Itoa(i + 1)
			Episode = transformEpisode(_type, Episode, element.Text)

			download := map[string]string{
				"episode":   Episode,
				"play_link": playLink}

			Smutex.Lock()
			downloadAry = append(downloadAry, download)
			Smutex.Unlock()
		})

		kuyunAryJson, _ := Json.MarshalIndent(kuyunAry, "", " ")
		ckm3u8AryJson, _ := Json.MarshalIndent(ckm3u8Ary, "", " ")
		downloadAryJson, _ := Json.MarshalIndent(downloadAry, "", " ")

		// detail["alias"] = e.ChildText("div .vodinfobox>ul>li:eq(0)") // WTF 不支持这样的选择器
		// xpath 还是靠谱
		// 别名
		c.OnXML("/html/body/div[5]/div[1]/div/div/div[2]/div[2]/ul/li[1]/span", func(e *colly.XMLElement) {
			detail["alias"] = e.Text
		})

		// 导演
		c.OnXML("/html/body/div[5]/div[1]/div/div/div[2]/div[2]/ul/li[2]/span", func(e *colly.XMLElement) {
			detail["director"] = e.Text
		})

		// 主演
		c.OnXML("/html/body/div[5]/div[1]/div/div/div[2]/div[2]/ul/li[3]/span", func(e *colly.XMLElement) {
			detail["starring"] = e.Text
		})

		// 类型
		detail["type"] = _type

		// 地区
		c.OnXML("/html/body/div[5]/div[1]/div/div/div[2]/div[2]/ul/li[5]/span", func(e *colly.XMLElement) {
			detail["area"] = e.Text
		})

		c.OnXML("/html/body/div[5]/div[1]/div/div/div[2]/div[2]/ul/li[6]/span", func(e *colly.XMLElement) {
			detail["language"] = e.Text
		})

		// 上映时间
		c.OnXML("/html/body/div[5]/div[1]/div/div/div[2]/div[2]/ul/li[7]/span", func(e *colly.XMLElement) {
			detail["released"] = e.Text
		})

		// 片长
		c.OnXML("/html/body/div[5]/div[1]/div/div/div[2]/div[2]/ul/li[8]/span", func(e *colly.XMLElement) {
			detail["length"] = e.Text
		})

		// 更新时间
		c.OnXML("/html/body/div[5]/div[1]/div/div/div[2]/div[2]/ul/li[9]/span", func(e *colly.XMLElement) {
			detail["update"] = e.Text
		})

		// 总播放量
		c.OnXML("/html/body/div[5]/div[1]/div/div/div[2]/div[2]/ul/li[10]/span", func(e *colly.XMLElement) {
			detail["total_playback"] = e.Text
		})

		// 剧情简介
		c.OnXML("/html/body/div[5]/div[3]/div[2]", func(e *colly.XMLElement) {
			detail["vod_play_info"] = e.Text
		})

		if detail["vod_play_info"] == "" || detail["vod_play_info"] == nil {
			c.OnXML("/html/body/div[5]/div[2]/div[2]/text()", func(e *colly.XMLElement) {
				detail["vod_play_info"] = e.Text
			})
		}

		md = MoviesDetail{
			Link:     url,
			Name:     name,
			Cover:    cover,
			Quality:  quality,
			Score:    score,
			Detail:   detail,
			KuYun:    string(kuyunAryJson),
			CK:       string(ckm3u8AryJson),
			Download: string(downloadAryJson),
		}

	})

	// 在OnHTML之后被调用
	c.OnScraped(func(_ *colly.Response) {

		_moviesInfo := make(map[string]interface{})

		_moviesInfo["link"] = md.Link
		_moviesInfo["cover"] = md.Cover
		_moviesInfo["name"] = md.Name
		_moviesInfo["quality"] = md.Quality
		_moviesInfo["score"] = md.Score
		_moviesInfo["kuyun"] = md.KuYun
		_moviesInfo["ckm3u8"] = md.CK
		_moviesInfo["download"] = md.Download

		_detail, _ := Json.MarshalIndent(md.Detail, "", " ")

		_moviesInfo["detail"] = string(_detail)

		if md.Name != "" {
			Smutex.Lock()
			t := RedisDB.HMSet(moviesDetail+url+":movie_name:"+md.Name, _moviesInfo).Err()
			log.Println(t)
			Smutex.Unlock()
		}

	})

	visitError := c.Visit(host + url)

	log.Println(visitError)

	c.Wait()

	return md
}

// /?m=vod-type-id-1.html  => /?m=vod-type-id-1-pg-1
func CategoryToPageUrl(categoryUrl string, page string) string {
	// 主类链接： /?m=vod-type-id-1.html
	// 主类的页面链接 /?m=vod-type-id-1-pg-
	categoryUrlStrSplit := strings.Split(categoryUrl, ".html")[0]

	pageUrl := categoryUrlStrSplit + "-pg-" + page + ".html"

	return pageUrl
}

// 获取url中的链接
func TransformId(Url string) string {
	UrlStrSplit := strings.Split(Url, "-id-")[1]

	return strings.TrimRight(UrlStrSplit, ".html")
}

func DelAllListCacheKey() {

	AllListCacheKey := RedisDB.Keys("movie_lists_key:detail_links:*").Val()

	// 删除已经缓存的数据
	for _, val := range AllListCacheKey {
		RedisDB.Del(val)
	}
}

func isFilm(_type string) bool {
	return strings.Contains(_type, "片")
}

// 电影只处理国语跟广东话、其他语言暂不处理
func transformEpisode(_type, episode, linkName string) string {

	if isFilm(_type) == true {
		if strings.Contains(linkName, "粤语") == true {
			episode = "粤语"
		}
		if strings.Contains(linkName, "国语") == true {
			episode = "国语"
		}
	}

	return episode
}

func sendDingMsg(msg string) {
	accessToken := viper.GetString(`ding.access_token`)
	if accessToken == "" {
		return
	}
	webhook := "https://oapi.dingtalk.com/robot/send?access_token=" + accessToken
	robot := NewRobot(webhook)

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
