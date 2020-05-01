package spider

import (
	"github.com/spf13/viper"
	"go_movies/utils"
)

func Create() utils.SpiderTask {

	mod := viper.GetString(`app.spider_mod`)

	switch mod {
	case "api":
		//return &SpiderApi{}
		return new(SpiderApi)
	case "WebPage":
		return new(utils.Spider)
	default:
		return new(SpiderApi)
	}
}
