package spider

import (
	"go_movies/utils"
	"go_movies/utils/spider/haohua"

	"github.com/spf13/viper"
)

// 定义 mod 的映射关系
var spiderModMap = map[string]utils.SpiderTask{
	"haohua": &haohua.SpiderApi{}, // 豪华资源
}

func Create() utils.SpiderTask {
	mod := viper.GetString(`app.spider_mod`)
	return spiderModMap[mod]
}
