# Go Movies

> golang + redis 实现的影站(低级爬虫)。无管理后台，效果站： [https://go-movies.hzz.cool/](https://go-movies.hzz.cool/) 支持手机端访问播放

> 静态文件与go文件统一编译，运行只依赖编译后可执行的二进制文件与redis

> 内置自动爬虫的定时任务，基本满足日常看片需求。

## Tip
- 支持网页爬虫与API请求的形式，可通过 config/app.go 配置定义选择使用的版本
- 由于目标网站会封锁直接通过网页爬虫的IP,在没有找到稳定IP池的情况下，推荐优先使用API版本（PS：网页爬虫版可用，但可能会被封IP）
- master同时维护网页爬虫与API的两个版本 [API接口说明.txt](http://www.jisudhw.com/help/API%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E.txt)
- 由于爬虫版与API版的电影ID不匹配，建议使用不同的redis DB库分别存储

## Github地址
[https://github.com/hezhizheng/go-movies](https://github.com/hezhizheng/go-movies)

## 首页效果
![img](https://i.loli.net/2019/12/05/Qzqv4HWoMp2DByi.png)

## 使用安装 
```
# 下载
git clone https://github.com/hezhizheng/go-movies

# 进入目录
cd go-movies

# 生成配置文件(默认使用redis db10的库，可自行修改app.go中的配置)
cp ./config/app.go.backup ./config/app.go

# 配置说明
app.spider_path: 爬虫路由
app.spider_path_name: 爬虫路由名称
app.debug_path: debug的路由
app.debug_path_name: debug的路由名称
cron.timing_spider: 定时爬虫的CRON表达式
ding.access_token: 钉钉机器人token
app.spider_mod: 爬虫模式 参数：api/WebPage (api为直接请求API模式，WebPage为网页爬虫模式)

# 启动 (首次启动会自动开启爬虫任务)
go run main.go 
or
# 安装 bee 工具
bee run

# 如安装依赖包失败，请使用代理
export GOPROXY=https://goproxy.io,direct
or
export GOPROXY=https://goproxy.cn,direct

访问
http://127.0.0.1:8899
```

### 开启爬虫
- ~~直接访问链接 http://127.0.0.1:8899/movies-spider (开启定时任务，定时爬取就好)~~
  - 已内置定时爬虫，默认凌晨一点开启爬虫(可修改配置文件cron.timing_spider表达式)
- 耗时：具体时间受目标网站/接口的响应速度影响

## Tools
- [https://github.com/gocolly/colly](https://github.com/gocolly/colly) 爬虫框架
- 模板引擎：https://github.com/shiyanhui/hero
- 数据库 redis 缓存/持久 [https://github.com/Go-redis/redis](https://github.com/Go-redis/redis)
  - Zset：每个分类为一个有序集合
    - score：电影更新的时间戳
    - member：电影对应的实际URL
  - Hash：电影具体信息(名称、封面图等)、每页数据的缓存
- 路由 [https://github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
- json解析 jsoniter [github.com/json-iterator/go](github.com/json-iterator/go)
- 跨平台打包：https://github.com/mitchellh/gox
- 静态资源处理：https://github.com/rakyll/statik
- web server 框架：https://github.com/valyala/fasthttp

## 注意
```
# 修改静态文件/static  、 views/hero 需要先安装包的依赖，执行以下编译命令，更多用法可参考官方redame.md

# https://github.com/rakyll/statik
statik -src=xxxPath/go_movies/static -f 

# https://github.com/shiyanhui/hero
hero -source="./views/hero"
```

## 编译可执行文件(跨平台)

```
# 用法参考 https://github.com/mitchellh/gox
# 生成文件可直接执行 Linux
gox -osarch="linux/amd64" 
......
```
- ~~提供win64、Linux64的已编译的文件下载~~ （请自行编译）

`使用请确保redis为开启状态，默认使用 DB10，启动成功之后会自动执行爬虫，可自行访问 http://127.0.0.1:8899/movies-spider 进行爬虫`

[微云](https://share.weiyun.com/5iLGksd)  （推荐微云+[proxyee-down](https://github.com/proxyee-down-org/proxyee-down),  原来go版本已经在开发中了...）

![img](https://i.loli.net/2020/01/04/OxsqRunwliy31zN.png)

## Docker 部署（使用docker-compose可直接忽略该步骤）

```
# 安装 redis 镜像(已有可以忽略) 
sudo docker pull redis:latest

# 启动redis容器
# 根据实际情况分配端口 -p 宿主机端口:容器端口
sudo docker run -itd --name redis-test -p 6379:6379 redis

# 修改 app.go 的redis 连接地址为容器名称
"addr":"redis-test"

# 编译go-movies
gox -osarch="linux/amd64"

# 构造镜像
sudo docker build -t go-movies-docker-scratch .

# 启动容器
sudo docker run --link redis-test:redis -p 8899:8899 -d go-movies-docker-scratch

```

## docker-compose 一键启动
```
# 修改 app.go 的redis 连接地址为容器名称，这里需要跟docker-compose.yml中保持一致
"addr":"redis-test"

# 编译go-movies
gox -osarch="linux/amd64"

# 运行
sudo docker-compose up -d

打开游览器访问 http://127.0.0.1:8899 即可看见网站效果
```

## 目录结构参考beego设置

## TODO
- [x] 跨平台编译,模板路径不正确
  - 使用 https://github.com/rakyll/statik 处理 js、css、image等静态资源
  - 使用 https://github.com/shiyanhui/hero 替换 html/template 模板引擎
- [x] redis查询问题
  - 缓存页面数据
- [x] 增加配置文件读取
  - 使用 https://github.com/spf13/viper
- [x] Docker 部署
- [x] goroutine 并发数控制
  - 使用 https://github.com/panjf2000/ants
- [ ] 爬取数据的完整性


## Other
许多Go的原理还没弄懂，有精力会慢慢深究下。写得很潦草，多多包涵。

## 🔋 JetBrains 开源证书支持

`go-movies` 项目一直以来都是在 JetBrains 公司旗下的 GoLand 集成开发环境中进行开发，基于 **free JetBrains Open Source license(s)** 正版免费授权，在此表达我的谢意。

<a href="https://www.jetbrains.com/?from=go-movies" target="_blank"><img src="https://i.loli.net/2020/08/16/Brce57tU4SQWspm.png" width="250" align="middle"/></a>

## License
[MIT](./LICENSE.txt)
