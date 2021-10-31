package main

import (
	"bytes"
	"embed"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"go_movies/config"
	"go_movies/routes"
	"go_movies/utils"
	"go_movies/utils/spider"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed static2/*
var embedStatic2 embed.FS



// Reads from the routes slice to translate the values to httprouter.Handle
// 遍历路由
func traversingRouter() *httprouter.Router {

	AllRoutes := routes.AllRoutes()

	router := httprouter.New()
	for _, route := range AllRoutes {
		var handle httprouter.Handle

		handle = route.HandlerFunc

		log.Println(route.Path)
		router.Handle(route.Method, route.Path, handle)
	}

	if viper.GetString(`app.debug`) == "false" {
		// live 模式 打包用
		fsys, _ := fs.Sub(embedStatic2, "static2")
		router.ServeFiles("/static2/*filepath", http.FS(fsys))
	}else{
		// dev 开发用 避免修改静态资源需要重启服务
		router.ServeFiles("/static2/*filepath", http.Dir("static2"))
	}
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

	hasHome := utils.RedisDB.Exists("detail_links:id:1").Val()
	log.Println("hasHome", hasHome)
	// 不存在首页的key 则认为是第一次启动
	if hasHome == 0 {
		spider.Create().Start()
	}
}

func main() {
	// 注册所有路由
	router := traversingRouter()

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
