# Go Movies

> 一个基于golang的爬虫电影站

## 使用安装 

### 开启爬虫
- /go-spider
- 消耗：Windows 下 cup 10% 左右，内存 30mb 左右 
- 网络正常的情况下，爬虫完毕耗时大概21分钟左右（存在部分爬取失败的记录）

## Tools
- [https://github.com/gocolly/colly](https://github.com/gocolly/colly) 爬虫框架
- html/template 模板引擎
- 数据库 redis 缓存/持久 [https://github.com/Go-redis/redis](https://github.com/Go-redis/redis)
- 路由 [https://github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
- json解析 jsoniter [github.com/json-iterator/go](github.com/json-iterator/go)

## 目录结构参考beego设置

## TODO
-[ ] 跨平台打包,模板路径不正确
-[ ] goroutine 并发数控制
-[ ] 爬取数据的完整性
-[ ] redis查询问题？
