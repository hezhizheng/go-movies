package utils

import (
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"log"
)

func TimingSpider() {

	log.Println("cron TimingSpider start:")

	// v3 用法 干
	c := cron.New(cron.WithSeconds())

	// 每天定时执行的条件
	spec := viper.GetString(`cron.timing_spider`)

	c.AddFunc(spec, func() {
		go StartSpider()
		log.Println("corn Spider ing:")
	})

	go c.Start()

	// 关闭计划任务, 但是不能关闭已经在执行中的任务.
	// defer c.Stop()

	// 阻塞
	//select {}
}
