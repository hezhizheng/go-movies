package controller

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go_movies/utils"
	"go_movies/utils/spider"
	"log"
	"net/http"
)

func Debug(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//film := "6,7,8,9,10,11,12,20,21,31"
	//t := strings.Contains(film, "21")
	//fmt.Fprint(w, "DEBUG", t)
	//return
	//go utils.DelAllListCacheKey()

	//c := utils.RedisDB.Exists("detail_links:id:10").Val()
	//fmt.Fprint(w, "DEBUG", c)
	//return
	go spider.StartApi()

	return

	//x1 := `HD高清$https://youku.cdn10-okzy.com/share/aec3d6abde16f9aad1970fad0fed2cb3$$$HD高清$https://youku.cdn10-okzy.com/20200406/13525_b65f0b6c/index.m3u8`
	x1 := `第01集$https://iqiyi.cdn9-okzy.com/share/1d01bd2e16f57892f0954902899f0692#第02集$https://iqiyi.cdn9-okzy.com/share/664c7298d2b73b3c7fe2d1e8d1781c06#第03集$https://iqiyi.cdn9-okzy.com/share/6e2290dbf1e11f39d246e7ce5ac50a1e$$$第01集$https://iqiyi.cdn9-okzy.com/20200406/8392_bc901e4e/index.m3u8#第02集$https://iqiyi.cdn9-okzy.com/20200406/8394_d3645c2c/index.m3u8#第03集$https://iqiyi.cdn9-okzy.com/20200406/8393_f1911261/index.m3u8`
	y, y1 := spider.FormatVodPlayUrl(x1)
	fmt.Fprint(w, "DEBUG", y, y1)
	return
	var nav spider.ResData
	//utils.Json.Unmarshal(b, &nav)

	x := `{
    "code": 1,
    "msg": "数据列表",
    "page": 1,
    "pagecount": 2064,
    "limit": "25",
    "total": 51582,
    "list": [
        {
            "vod_id": 51592,
            "vod_name": "扎克风暴",
            "type_id": 31,
            "type_name": "欧美动漫",
            "vod_en": "zhakefengbao",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "完结",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51591,
            "vod_name": "危险性游戏",
            "type_id": 11,
            "type_name": "剧情片",
            "vod_en": "weixianxingyouxi",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "HD",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51589,
            "vod_name": "挑逗性游戏",
            "type_id": 6,
            "type_name": "动作片",
            "vod_en": "tiaodouxingyouxi",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "HD",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51588,
            "vod_name": "继怪怪守护神",
            "type_id": 30,
            "type_name": "日韩动漫",
            "vod_en": "jiguaiguaishouhushen",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "更新至1集",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51587,
            "vod_name": "大象的眼泪",
            "type_id": 11,
            "type_name": "剧情片",
            "vod_en": "daxiangdeyanlei",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "HD",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51590,
            "vod_name": "外星小子哆布哆",
            "type_id": 29,
            "type_name": "国产动漫",
            "vod_en": "waixingxiaoziduobuduo",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "完结",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51586,
            "vod_name": "在爱的名义下",
            "type_id": 23,
            "type_name": "日本剧",
            "vod_en": "zaiaidemingyixia",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "完结",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51585,
            "vod_name": "黑名单上的人",
            "type_id": 16,
            "type_name": "欧美剧",
            "vod_en": "heimingdanshangderen",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "完结",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51584,
            "vod_name": "格莱普尼尔",
            "type_id": 30,
            "type_name": "日韩动漫",
            "vod_en": "gelaipunier",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "更新至1集",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51583,
            "vod_name": "恐怖之夜：噩梦电台",
            "type_id": 10,
            "type_name": "恐怖片",
            "vod_en": "kongbuzhiyeemengdiantai",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "HD",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51582,
            "vod_name": "偶像星愿第二季",
            "type_id": 30,
            "type_name": "日韩动漫",
            "vod_en": "ouxiangxingyuandierji",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "更新至02集",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51581,
            "vod_name": "戈梅拉岛",
            "type_id": 37,
            "type_name": "伦理片",
            "vod_en": "gemeiladao",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "HD",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51580,
            "vod_name": "金鱼公主/金鱼姬",
            "type_id": 23,
            "type_name": "日本剧",
            "vod_en": "jinyugongzhujinyuji",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "完结",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51578,
            "vod_name": "搜查会议在客厅 再来一碗",
            "type_id": 23,
            "type_name": "日本剧",
            "vod_en": "souchahuiyizaiketingzailaiyiwan",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "更新至01集",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51579,
            "vod_name": "恶魔之酸",
            "type_id": 10,
            "type_name": "恐怖片",
            "vod_en": "emozhisuan",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "HD",
            "vod_play_from": "ckm3u8,kuyun"
        },
        {
            "vod_id": 51577,
            "vod_name": "应援",
            "type_id": 23,
            "type_name": "日本剧",
            "vod_en": "yingyuan",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "更新至05集",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51576,
            "vod_name": "印度支那",
            "type_id": 11,
            "type_name": "剧情片",
            "vod_en": "yinduzhina",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "HD",
            "vod_play_from": "ckm3u8,kuyun"
        },
        {
            "vod_id": 51574,
            "vod_name": "华沙间谍",
            "type_id": 16,
            "type_name": "欧美剧",
            "vod_en": "huashajiandie",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "完结",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51573,
            "vod_name": "还会与你相见3次",
            "type_id": 23,
            "type_name": "日本剧",
            "vod_en": "huanhuiyunixiangjian3ci",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "完结",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51572,
            "vod_name": "从零到我爱你",
            "type_id": 11,
            "type_name": "剧情片",
            "vod_en": "conglingdaowoaini",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "HD",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51571,
            "vod_name": "沙特奇趣录大电影",
            "type_id": 31,
            "type_name": "欧美动漫",
            "vod_en": "shateqiquludadianying",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "HD",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51570,
            "vod_name": "疯狂老爹",
            "type_id": 7,
            "type_name": "喜剧片",
            "vod_en": "fengkuanglaodie",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "HD",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51569,
            "vod_name": "楼下女友请签收",
            "type_id": 13,
            "type_name": "国产剧",
            "vod_en": "louxianvyouqingqianshou",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "更新至08集",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51568,
            "vod_name": "毛泽东的亲家张文秋",
            "type_id": 11,
            "type_name": "剧情片",
            "vod_en": "maozedongdeqinjiazhangwenqiu",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "HD",
            "vod_play_from": "kuyun,ckm3u8"
        },
        {
            "vod_id": 51495,
            "vod_name": "这段恋爱是罪过吗？/这份爱是罪恶吗",
            "type_id": 23,
            "type_name": "日本剧",
            "vod_en": "zheduanlianaishizuiguomazhefenaishizuiema",
            "vod_time": "2020-04-05 23:21:15",
            "vod_remarks": "更新至01集",
            "vod_play_from": "kuyun,ckm3u8"
        }
    ],
    "class": [
        {
            "type_id": 1,
            "type_name": "电影"
        },
        {
            "type_id": 2,
            "type_name": "连续剧"
        },
        {
            "type_id": 3,
            "type_name": "综艺"
        },
        {
            "type_id": 4,
            "type_name": "动漫"
        },
        {
            "type_id": 5,
            "type_name": "资讯"
        },
        {
            "type_id": 6,
            "type_name": "动作片"
        },
        {
            "type_id": 7,
            "type_name": "喜剧片"
        },
        {
            "type_id": 8,
            "type_name": "爱情片"
        },
        {
            "type_id": 9,
            "type_name": "科幻片"
        },
        {
            "type_id": 10,
            "type_name": "恐怖片"
        },
        {
            "type_id": 11,
            "type_name": "剧情片"
        },
        {
            "type_id": 12,
            "type_name": "战争片"
        },
        {
            "type_id": 13,
            "type_name": "国产剧"
        },
        {
            "type_id": 14,
            "type_name": "香港剧"
        },
        {
            "type_id": 15,
            "type_name": "韩国剧"
        },
        {
            "type_id": 16,
            "type_name": "欧美剧"
        },
        {
            "type_id": 17,
            "type_name": "公告"
        },
        {
            "type_id": 18,
            "type_name": "头条"
        },
        {
            "type_id": 20,
            "type_name": "纪录片"
        },
        {
            "type_id": 21,
            "type_name": "微电影"
        },
        {
            "type_id": 22,
            "type_name": "台湾剧"
        },
        {
            "type_id": 23,
            "type_name": "日本剧"
        },
        {
            "type_id": 24,
            "type_name": "海外剧"
        },
        {
            "type_id": 25,
            "type_name": "内地综艺"
        },
        {
            "type_id": 26,
            "type_name": "港台综艺"
        },
        {
            "type_id": 27,
            "type_name": "日韩综艺"
        },
        {
            "type_id": 28,
            "type_name": "欧美综艺"
        },
        {
            "type_id": 29,
            "type_name": "国产动漫"
        },
        {
            "type_id": 30,
            "type_name": "日韩动漫"
        },
        {
            "type_id": 31,
            "type_name": "欧美动漫"
        },
        {
            "type_id": 32,
            "type_name": "港台动漫"
        },
        {
            "type_id": 33,
            "type_name": "海外动漫"
        },
        {
            "type_id": 34,
            "type_name": "福利片"
        },
        {
            "type_id": 35,
            "type_name": "解说"
        },
        {
            "type_id": 36,
            "type_name": "电影解说"
        },
        {
            "type_id": 37,
            "type_name": "伦理片"
        }
    ]
}`

	err := utils.Json.Unmarshal([]byte(x), &nav)

	log.Println(err)
	fmt.Fprint(w, "err", err)
	fmt.Fprint(w, "DEBUG", nav)
}
