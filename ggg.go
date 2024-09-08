package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	// 解析命令行参数
	target := flag.String("target", "", "目标服务器的地址")
	listen := flag.String("listen", "0.0.0.0:7860", "监听地址和端口")
	flag.Parse()

	if *target == "" {
		log.Fatal("请指定目标服务器的地址")
	}

	// 确保目标 URL 包含协议
	if !strings.HasPrefix(*target, "http://") && !strings.HasPrefix(*target, "https://") {
		*target = "https://" + *target
	}

	// 解析目标 URL
	targetURL, err := url.Parse(*target)
	if err != nil {
		log.Fatalf("目标 URL 解析失败: %v", err)
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 设置请求头
	proxy.Director = func(req *http.Request) {
		req.Header.Set("User-Agent", "Miao-Plugin/3.1")

		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.URL.Path = targetURL.Path
		if targetURL.RawQuery != "" {
			req.URL.RawQuery = targetURL.RawQuery
		}
	}

	// 处理请求并记录日志
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("接收到请求: %s %s", r.Method, r.URL.String())
		proxy.ServeHTTP(w, r)
		log.Printf("已转发请求: %s %s", r.Method, r.URL.String())
	})

	// 启动服务器
	log.Printf("反向代理服务器启动在 %s", *listen)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
