<p align="center">
<img src="https://cdn.jsdelivr.net/gh/hezhizheng/static-image-hosting@main/image-hosting/20211124102659_GRITHQADFCQQUIUN.png">
</p>

> golang + redis 实现的影站(低级爬虫)。无管理后台，效果站：[http://film.hzz.cool](http://film.hzz.cool) 支持手机端访问播放

## Github地址
[https://github.com/hezhizheng/go-movies](https://github.com/hezhizheng/go-movies)

## features
- 使用 Golang `embed` 内嵌静态资源(html、js、css等)，运行只依赖编译后可执行的二进制文件与redis
- 支持docker启动方式
- 简单资源分类、搜索的支持
- 内置自动爬虫、自动更新最新资源的定时任务，基本满足日常看片需求。
- 钉钉机器人通知

## Tip
- 暂时只维护API请求版本 [API接口说明](https://help.tiankongapi.com/)，后续可能追加其他资源支持
- API版本首次启动会全量请求并存储到redis，之后每小时定时爬取最近更新的影视资源

## 目录结构

```
|-- Dockerfile
|-- LICENSE.txt
|-- config
|   |-- app.go
|   `-- app.go.backup                        #程序配置文件
|-- controller                               #controller层，基本的页面渲染
|   |-- DebugController.go
|   |-- IndexController.go
|   `-- SpiderController.go
|-- docker-compose.yml
|-- go.mod
|-- go.sum
|-- main.go
|-- models                                   # 定义一些redis查询的方法
|   |-- Category.go
|   |-- Movies.go
|   `-- readme.md
|-- readme.md
|-- routes
|   `-- route.go
|-- runner.conf                              # fresh 配置文件
|-- services                                 # 普通业务处理类
|   |-- CategoryService.go
|   |-- MoviesService.go
|   `-- readme.md
|-- static2                                  # js、css、image等静态资源文件夹
|-- utils                                    # 一些工具类
|   |-- Cron.go
|   |-- Dingrobot.go
|   |-- Helper.go
|   |-- JsonUtil.go
|   |-- Pagination.go
|   |-- RedisUtil.go
|   |-- Spider.go
|   |-- SpiderTask.go
|   `-- spider                               # 爬虫api版主要功能代码
|       |-- SpiderTaskPolicy.go
|       `-- tian_kong
|           |-- CategoriesStr.go
|           `-- SpiderApi.go
`-- views                                    # html模板目录
    `-- tmpl
        `-- temp_global_var.go               # 定义 embed 的全局变量跟一些模板调用的函数
```

## 首页效果
<center>
<figure class="half">
    <img src=https://cdn.jsdelivr.net/gh/hezhizheng/static-image-hosting@main/image-hosting/20211124102515_SFWJDVHPPVCYNKDO.png>
    <img src=https://cdn.jsdelivr.net/gh/hezhizheng/static-image-hosting@main/image-hosting/20211124102515_EWYNHTWMKQRHUCDX.png>
</figure>
</center>

## 使用安装 (go version >= 1.18)
```
# 下载
git clone https://github.com/hezhizheng/go-movies

# 进入目录
cd go-movies

# 配置文件(默认使用redis db10的库，可自行修改app.go中的配置)
cp ./config/app.go.backup ./config/app.go

# 配置说明
app.spider_path: 爬虫路由
app.spider_path_name: 爬虫路由名称
app.debug_path: debug的路由
app.debug_path_name: debug的路由名称
cron.timing_spider: 定时爬虫的CRON表达式
ding.access_token: 钉钉机器人token
app.spider_mod: 固定参数为 TianKongApi
app.spider_mod: 开发模式建议设置为`true` 避免修改静态资源需要重启服务

# 启动 (首次启动会自动开启爬虫任务)
go run main.go 
or
# 安装 fresh 工具
fresh

# 如安装依赖包失败，请使用代理
export GOPROXY=https://goproxy.io,direct
or
export GOPROXY=https://goproxy.cn,direct

访问
http://127.0.0.1:8899
```

### 开启爬虫
- 已内置定时爬虫，首次全量请求之后，每小时定时爬取最近更新的影视资源(可自行修改配置文件cron.timing_spider表达式来控制间隔)
- 主动运行：直接访问链接 http://127.0.0.1:8899/movies-spider
- 耗时：具体时间受目标网站/接口的响应速度影响

## Tools
- 数据库 redis 缓存/持久 [https://github.com/Go-redis/redis](https://github.com/Go-redis/redis)
  - Zset：每个分类为一个有序集合
    - score：电影更新的时间戳
    - member：电影对应的实际URL
  - Hash：电影具体信息(名称、封面图等)、每页数据的缓存
- 路由 [https://github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
- json解析 jsoniter [github.com/json-iterator/go](github.com/json-iterator/go)
- 跨平台打包：https://github.com/mitchellh/gox
- web server 框架：https://github.com/valyala/fasthttp
- 配置文件读取：https://github.com/spf13/viper
- 热重启：https://github.com/gravityblast/fresh


## 编译可执行文件(跨平台)

```
# 用法参考 https://github.com/mitchellh/gox
# 生成文件可直接执行 
gox -osarch="windows/amd64" -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"
gox -osarch="darwin/amd64" -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"
gox -osarch="linux/amd64" -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"
```

- 提供win64的已编译的文件下载 [release](https://github.com/hezhizheng/go-movies/releases)

`使用请确保redis为开启状态，默认使用 DB10，启动成功之后会自动执行爬虫，可自行访问 http://127.0.0.1:8899/movies-spider 进行爬虫`

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
gox -osarch="linux/amd64" -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"

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
gox -osarch="linux/amd64" -ldflags "-s -w"

# 运行
sudo docker-compose up -d

打开游览器访问 http://127.0.0.1:8899 即可看见网站效果
```




## TODO
- [x] 跨平台编译,模板路径不正确
  - ~~使用 https://github.com/rakyll/statik 处理 js、css、image等静态资源~~
  - ~~使用 https://github.com/shiyanhui/hero 替换 html/template 模板引擎~~
  - 使用 golang 1.16 embed 特性
- [x] redis查询慢问题
  - 使用 hash 缓存页面数据
  - 使用scan 代替 keys*
- [x] 增加配置文件读取
  - 使用 https://github.com/spf13/viper
- [x] Docker 部署
- [x] goroutine 并发数控制
  - 使用 https://github.com/panjf2000/ants


## Other
许多Go的原理还没弄懂，有精力会慢慢深究下。写得很潦草，多多包涵。

## 🔋 JetBrains 开源证书支持

`go-movies` 项目一直以来都是在 JetBrains 公司旗下的 GoLand 集成开发环境中进行开发，基于 **free JetBrains Open Source license(s)** 正版免费授权，在此表达我的谢意。

<a href="https://www.jetbrains.com/?from=go-movies" target="_blank"><img src="https://i.loli.net/2020/08/16/Brce57tU4SQWspm.png" width="250" align="middle"/></a>

## License
[MIT](./LICENSE.txt)
