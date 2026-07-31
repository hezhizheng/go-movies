package haohua

// 定义主类与子类的关系
func CategoriesStr() string {
	categories := `[
    {
        "link": "/?m=vod-type-id-2.html",
        "name": "电影",
        "type_id": "2",
        "sub": [
            {
                "link": "/?m=vod-type-id-9.html",
                "name": "动作片",
                "type_id": "9",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-10.html",
                "name": "爱情片",
                "type_id": "10",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-11.html",
                "name": "喜剧片",
                "type_id": "11",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-12.html",
                "name": "科幻片",
                "type_id": "12",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-13.html",
                "name": "恐怖片",
                "type_id": "13",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-14.html",
                "name": "剧情片",
                "type_id": "14",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-15.html",
                "name": "战争片",
                "type_id": "15",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-16.html",
                "name": "记录片",
                "type_id": "16",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-23.html",
                "name": "动画片",
                "type_id": "23",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-34.html",
                "name": "灾难片",
                "type_id": "34",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-35.html",
                "name": "悬疑片",
                "type_id": "35",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-36.html",
                "name": "犯罪片",
                "type_id": "36",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-37.html",
                "name": "奇幻片",
                "type_id": "37",
                "sub": null
            }
        ]
    },
    {
        "link": "/?m=vod-type-id-1.html",
        "name": "电视剧",
        "type_id": "1",
        "sub": [
            {
                "link": "/?m=vod-type-id-3.html",
                "name": "欧美剧",
                "type_id": "3",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-4.html",
                "name": "香港剧",
                "type_id": "4",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-5.html",
                "name": "韩剧",
                "type_id": "5",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-6.html",
                "name": "日剧",
                "type_id": "6",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-7.html",
                "name": "马泰剧",
                "type_id": "7",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-20.html",
                "name": "内地剧",
                "type_id": "20",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-28.html",
                "name": "台湾剧",
                "type_id": "28",
                "sub": null
            }
        ]
    },
    {
        "link": "/?m=vod-type-id-17.html",
        "name": "动漫",
        "type_id": "17",
        "sub": [
            {
                "link": "/?m=vod-type-id-24.html",
                "name": "中国动漫",
                "type_id": "24",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-25.html",
                "name": "日本动漫",
                "type_id": "25",
                "sub": null
            },
            {
                "link": "/?m=vod-type-id-26.html",
                "name": "欧美动漫",
                "type_id": "26",
                "sub": null
            }
        ]
    }
]`

	return categories
}
