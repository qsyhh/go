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
	r.URL.Path = targetURL.Path
	if targetURL.RawQuery != "" {
	    r.URL.RawQuery = targetURL.RawQuery
	}
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	r.Header.Set("Accept", "*/*")
	r.Header.Set("Accept-Language", "en-US,en;q=0.5")
	r.Header.Set("Accept-Encoding", "gzip, deflate")
	r.Header.Set("Connection", "keep-alive")
        proxy.ServeHTTP(w, r)
	
		log.Printf("已转发请求: %s %s", r.Method, r.URL.String())
    })


	// 启动服务器
	log.Printf("反向代理服务器启动在 %s", *listen)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
