package tian_kong

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
                "name": "犯罪片",
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
                "name": "动漫片",
                "type_id": "20",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-21.html",
                "name": "剧情片",
                "type_id": "21",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-32.html",
                "name": "记录片",
                "type_id": "32",
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
                "name": "台湾剧",
                "type_id": "15",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-16.html",
                "name": "美国剧",
                "type_id": "16",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-24.html",
                "name": "日本剧",
                "type_id": "24",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-25.html",
                "name": "海外剧",
                "type_id": "25",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-26.html",
                "name": "韩国剧",
                "type_id": "26",
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
                "link": "/?m=vod-type-id-28.html",
                "name": "国产动漫",
                "type_id": "28",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-29.html",
                "name": "日本动漫",
                "type_id": "29",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-30.html",
                "name": "欧美动漫",
                "type_id": "30",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-31.html",
                "name": "海外动漫",
                "type_id": "31",
                "sub": null
            }
        ]
    }
]`

	return categories
}
