# Go Movies

> 一个基于golang的爬虫电影站，效果站： [https://go-movies.hezhizheng.com/](https://go-movies.hezhizheng.com/)

![img](https://i.loli.net/2019/12/05/Qzqv4HWoMp2DByi.png)

## 使用安装 
```
# 下载
git clone https://github.com/hezhizheng/go-movies

# 进入目录
cd go-movies

# 启动
go run main.go 
or
# 安装 bee 工具
bee run

# 如安装依赖包失败，请使用代理
export GOPROXY=https://goproxy.io,direct
or
export GOPROXY=https://goproxy.io,direct

访问
http://127.0.0.1:8899
```

### 开启爬虫
- 直接访问链接http://127.0.0.1:8899/movies-spider(开启定时任务，定时爬取就好)
- 消耗：Windows 下 cup 10% 左右，内存 30mb 左右(爬虫完毕都会降下来) 
- 网络正常的情况下，爬虫完毕耗时大概21分钟左右（存在部分资源爬取失败的情况）

## Tools
- [https://github.com/gocolly/colly](https://github.com/gocolly/colly) 爬虫框架
- html/template 模板引擎
- 数据库 redis 缓存/持久 [https://github.com/Go-redis/redis](https://github.com/Go-redis/redis)
- 路由 [https://github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
- json解析 jsoniter [github.com/json-iterator/go](github.com/json-iterator/go)

## 目录结构参考beego设置

## TODO
- [ ] 跨平台打包,模板路径不正确
- [ ] goroutine 并发数控制
- [ ] 爬取数据的完整性
- [ ] redis查询问题？

## Other
许多Go的原理还没弄懂，有精力会慢慢深究下。写得很潦草，多多包涵。
