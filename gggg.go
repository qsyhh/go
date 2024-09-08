package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	target := os.Getenv("TARGET")
	listen := os.Getenv("LISTEN")

	if target == "" || listen == "" {
		log.Fatal("请设置环境变量 TARGET 和 LISTEN")
	}

	u, err := url.Parse(target)
	if err != nil {
		log.Fatalf("解析目标URL失败： %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	log.Printf("启动反向代理，监听 %s，转发到 %s", listen, target)
	log.Fatal(http.ListenAndServe(listen, proxy))
}
