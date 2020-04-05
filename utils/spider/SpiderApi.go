package spider

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"go_movies/utils"
	"log"
)

const API_HOST = "https://api.okzy.tv/api.php/provide/vod"
const AC_LIST = "list"
const AC_DETAIL = "detail"

type Lists struct {
	VodId       int    `json:"vod_id"` // 如果json中vod_id不是“1”，而是 1 ，这里一定要声明为 int ！！！fuck 不愧是静态强类型
	VodName     string `json:"vod_name"`
	TypeId      int    `json:"type_id"`
	TypeName    string `json:"type_name"`
	VodEn       string `json:"vod_en"`
	VodTime     string `json:"vod_time"`
	VodRemarks  string `json:"vod_remarks"`
	VodPlayFrom string `json:"vod_play_from"`
}

type ResData struct {
	Msg       string  `json:"msg"`
	Code      int     `json:"code"`
	Page      int     `json:"page"`
	PageCount int     `json:"pagecount"`
	Limit     string  `json:"limit"`
	Total     int     `json:"total"`
	List      []Lists `json:"list"`
}

type Categories struct {
	Link string       `json:"link"`
	Name string       `json:"name"`
	Sub  []Categories `json:"sub"`
}

func StartApi() {
	list()
}

func list() {
	url := API_HOST + "?ac=" + AC_LIST

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

	fmt.Println("result:\r\n", nav)

	for _, val := range nav.List {
		log.Println("TypeName", val.TypeName)
	}

}

// id与旧的网页爬虫对应不上
func detail(id string) {

}

// 定义主类与子类的关系
func CategoriesStr() string {
	categories := `[
    {
        "link": "/?m=vod-type-id-1.html",
        "name": "电影片",
        "sub": [
            {
                "link": "/?m=vod-type-id-6.html",
                "name": "动作片",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-7.html",
                "name": "喜剧片",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-8.html",
                "name": "爱情片",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-9.html",
                "name": "科幻片",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-10.html",
                "name": "恐怖片",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-11.html",
                "name": "剧情片",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-12.html",
                "name": "战争片",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-20.html",
                "name": "纪录片",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-21.html",
                "name": "微电影",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-37.html",
                "name": "伦理片",
                "sub": null
            }
        ]
    },
    {
        "link": "/?m=vod-type-id-2.html",
        "name": "连续剧",
        "sub": [
            {
                "link": "/?m=vod-type-id-13.html",
                "name": "国产剧",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-14.html",
                "name": "香港剧",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-15.html",
                "name": "韩国剧",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-16.html",
                "name": "欧美剧",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-22.html",
                "name": "台湾剧",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-23.html",
                "name": "日本剧",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-24.html",
                "name": "海外剧",
                "sub": null
            }
        ]
    },
    {
        "link": "/?m=vod-type-id-4.html",
        "name": "动漫片",
        "sub": [
            {
                "link": "/?m=vod-type-id-29.html",
                "name": "国产动漫",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-30.html",
                "name": "日韩动漫",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-31.html",
                "name": "欧美动漫",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-32.html",
                "name": "港台动漫",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-33.html",
                "name": "海外动漫",
                "sub": null
            }
        ]
    }
]`

	return categories
}
