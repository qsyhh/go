package main

import (
	"flag"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    // 目标服务器的 URL
	listen := flag.String("listen", "0.0.0.0:7860", "监听地址和端口")
	flag.Parse()

    targetURL, err := url.Parse("https://enka.network/")
    if err != nil {
        log.Fatalf("Failed to parse target URL: %v", err)
    }

    // 创建反向代理
    proxy := httputil.NewSingleHostReverseProxy(targetURL)

    // 自定义处理逻辑（如果需要）
    proxy.ModifyResponse = func(resp *http.Response) error {
        // 可以在这里对响应进行处理
        return nil
    }

    // 处理传入的请求并将其转发到目标服务器
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("接收到请求: %s %s", r.Method, r.URL.String())

        // 设置请求的 URL 为目标服务器的 URL
        r.URL.Host = targetURL.Host
        r.URL.Scheme = targetURL.Scheme
		r.Header.Set("User-Agent", "Miao-Plugin/3.1")
        r.Header.Set("X-Forwarded-Host", r.Host)
        proxy.ServeHTTP(w, r)
	
		log.Printf("已转发请求: %s %s", r.Method, r.URL.String())
    })


	// 启动服务器
	log.Printf("反向代理服务器启动在 %s", *listen)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
