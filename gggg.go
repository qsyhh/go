package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	// 解析命令行参数
	target := flag.String("target", "", "目标服务器的地址")
	listen := flag.String("listen", "", "监听地址和端口")
	flag.Parse()

	if *target == "" {
		fmt.Println("请指定目标服务器的地址")
		return
	}

	// 确保目标 URL 包含协议
	if !strings.HasPrefix(*target, "http://") && !strings.HasPrefix(*target, "https://") {
		*target = "https://" + *target
	}

	// 解析目标 URL
	targetURL, err := url.Parse(*target)
	if err != nil {
		fmt.Println("目标 URL 解析失败:", err)
		return
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

	// 启动服务器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	fmt.Printf("反向代理服务器启动在 %s\n", *listen)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		fmt.Println("启动服务器失败:", err)
	}
}
