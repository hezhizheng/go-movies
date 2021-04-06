package spider

import (
	"github.com/spf13/viper"
	"go_movies/utils"
	"go_movies/utils/spider/tian_kong"
)

// 定义 mod 的映射关系
var spiderModMap = map[string]utils.SpiderTask{
	"api":     &SpiderApi{},
	"TianKongApi" : &tian_kong.SpiderApi{},
	"WebPage": &utils.Spider{}}

func Create() utils.SpiderTask {

	mod := viper.GetString(`app.spider_mod`)

	return spiderModMap[mod]
}
