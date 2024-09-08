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
	listen := flag.String("listen", "0.0.0.0:7860", "监听地址和端口")
	flag.Parse()

	if *target == "" {
		fmt.Println("请指定目标服务器的地址")
		return
	}

	// 确保目标 URL 包含协议
	if !strings.HasPrefix(*target, "http://") && !strings.HasPrefix(*target, "https://") {
		*target = "http://" + *target
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
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Connection", "keep-alive")

		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.URL.Path = targetURL.Path
		if targetURL.RawQuery != "" {
			req.URL.RawQuery = targetURL.RawQuery
		}
	}

	// 启动服务器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 将请求转发到目标服务器
		proxy.ServeHTTP(w, r)
	})

	fmt.Printf("反向代理服务器启动在 %s\n", *listen)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		fmt.Println("启动服务器失败:", err)
	}
}
