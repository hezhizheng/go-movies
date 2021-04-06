package main

import (
	"bytes"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/viper"
	"go_movies/config"
	"go_movies/routes"
	_ "go_movies/statik"
	"go_movies/utils"
	"go_movies/utils/spider"
	"log"
	"net/http"
	"os"
)

// Reads from the routes slice to translate the values to httprouter.Handle
// 遍历路由
func TraversingRouter() *httprouter.Router {

	AllRoutes := routes.AllRoutes()

	router := httprouter.New()
	for _, route := range AllRoutes {
		var handle httprouter.Handle

		handle = route.HandlerFunc

		log.Println(route.Path)
		router.Handle(route.Method, route.Path, handle)
	}

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	// 配置静态文件访问
	router.ServeFiles("/static/*filepath", statikFS)
	return router
}

// 初始化配置文件
func init() {
	viper.SetConfigType("json") // 设置配置文件的类型

	readConfig := errors.New("未定义配置文件")

	if _, err := os.Stat("./app.json"); os.IsNotExist(err) {
		readConfig = viper.ReadConfig(bytes.NewBuffer(config.AppJsonConfig))
	} else {
		viper.SetConfigName("app")
		viper.AddConfigPath(".")
		readConfig = viper.ReadInConfig()
	}

	if err := readConfig; err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("no such config file")
		} else {
			// Config file was found but another error was produced
			log.Println("read config error")
		}
		log.Fatal(err) // 读取配置文件失败致命错误
	}
}

// 首次启动自动开启爬虫
func firstSpider() {

	hasHK := utils.RedisDB.Exists("detail_links:id:14").Val()
	log.Println("hasHK", hasHK)
	// 不存在首页的key 则认为是第一次启动
	if hasHK == 0 {
		spider.Create().Start()
	}
}

func main() {
	// 注册所有路由
	router := TraversingRouter()

	// 初始化 redis 连接
	utils.InitRedisDB()
	defer utils.CloseRedisDB()

	port := viper.GetString(`app.port`)
	log.Println("监听端口", "http://127.0.0.1"+port)

	firstSpider()

	// 启动定时爬虫任务
	utils.TimingSpider(func() {
		spider.Create().Start()
		return
	})

	utils.RecentUpdate(func() {
		spider.Create().DoRecentUpdate()
		return
	})

	log.Println(http.ListenAndServe(port, router))

}
