package main

import (
    "flag"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    target := flag.String("target", "", "Target server to proxy")
    listen := flag.String("listen", "", "Address to listen on")
    flag.Parse()

    if *target == "" || *listen == "" {
        log.Fatal("请提供 -target 和 -listen 参数")
    }

    u, err := url.Parse(*target)
    if err != nil {
        log.Fatalf("解析目标URL失败： %v", err)
    }

    proxy := httputil.NewSingleHostReverseProxy(u)

    log.Printf("启动反向代理，监听 %s，转发到 %s", *listen, *target)
    log.Fatal(http.ListenAndServe(*listen, proxy))
}

