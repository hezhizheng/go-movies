package request

import (
	"crypto/tls"
	"errors"
	"github.com/valyala/fasthttp"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/**
fasthttp client 示例
@link http://liuqh.icu/2022/04/13/go/package/34-fasthttp/
@link https://zsmhub.github.io/post/golang/http%E5%AE%A2%E6%88%B7%E7%AB%AF/
*/

type RequestForm struct {
	Uri      string      // 请求uri
	Params   url.Values  // 请求参数
	Header   url.Values  // 请求头
	RespBody string      // 响应结果
	Resp     interface{} // 响应结果结构化
}

type client struct{}

type IClient interface {
	Get(dto *RequestForm) ([]byte, error)
	Post(dto *RequestForm) ([]byte, error)
}

const ttl = 5 * time.Second

var (
	defaultDialer = &fasthttp.TCPDialer{Concurrency: 200} // tcp 并发200

	FastClient = CreateFastHttpClient()

	Client IClient = new(client)
)

func CreateFastHttpClient() fasthttp.Client {
	return fasthttp.Client{
		MaxConnsPerHost:     300,
		MaxIdleConnDuration: 10 * time.Second, // 空闲链接时间应短，避免请求服务的 keep-alive 过短主动关闭，默认10秒
		MaxConnDuration:     10 * time.Minute,
		ReadTimeout:         30 * time.Second,
		WriteTimeout:        30 * time.Second,
		MaxResponseBodySize: 1024 * 1024 * 10,
		MaxConnWaitTimeout:  time.Minute,
		Dial: func(addr string) (net.Conn, error) {
			idx := 3 // 重试三次
			for {
				idx--
				conn, err := defaultDialer.DialTimeout(addr, 10*time.Second) // tcp连接超时时间10s
				if err != fasthttp.ErrDialTimeout || idx == 0 {
					return conn, err
				}
			}
		},

		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
			CurvePreferences:   []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}
}

func (c *client) Get(dto *RequestForm) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	paramsStr := HttpBuildQuery(dto.Params)

	req.SetRequestURI(dto.Uri + "?" + paramsStr)
	req.Header.SetMethod(http.MethodGet)

	if dto.Header != nil {
		for k, v := range dto.Header {
			req.Header.Set(k, strings.Join(v, ","))
		}
	}

	if err := FastClient.DoTimeout(req, resp, ttl); err != nil {
		return nil, err
	}

	dto.RespBody = string(resp.Body())

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, errors.New(dto.RespBody)
	}

	return resp.Body(), nil
}

func (c *client) Post(dto *RequestForm) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	// application/json 编码方式
	//req.SetBody(helper.JsonMarshalByte(dto.Params))
	req.Header.SetContentType("application/json")

	// application/x-www-form-urlencoded 编码方式
	// paramsStr, err := helper.HttpBuildQuery(dto.Params)
	// if err != nil {
	//     return err
	// }
	// req.SetBody([]byte(paramsStr))
	// req.Header.SetContentType("application/x-www-form-urlencoded")

	req.SetRequestURI(dto.Uri)
	req.Header.SetMethod(http.MethodPost)

	if dto.Header != nil {
		for k, v := range dto.Header {
			req.Header.Set(k, strings.Join(v, ","))
		}
	}

	if err := FastClient.DoTimeout(req, resp, ttl); err != nil {
		return nil, err
	}

	dto.RespBody = string(resp.Body())

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, errors.New(dto.RespBody)
	}

	return resp.Body(), nil
}

func HttpBuildQuery(queryData url.Values) string {
	// var uri url.URL
	//    q := uri.Query()
	//    q.Add("name", "张三")
	//    q.Add("age", "20")
	//    q.Add("sex", "1")
	//    queryStr := q.Encode()
	//    fmt.Println(queryStr)
	return queryData.Encode()
}
