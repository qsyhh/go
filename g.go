package main

import (
    "flag"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    target := flag.String("target", "https://enka.network", "Target server to proxy")
    listen := flag.String("listen", "0.0.0.0:7860", "Address to listen on")
    scheme := flag.String("scheme", "https", "Protocol scheme for the target URL")
    flag.Parse()

    if *target == "" || *listen == "" {
        log.Fatal("请提供 -target 和 -listen 参数")
    }

    u, err := url.Parse(*target)
    if err != nil {
        log.Fatalf("解析目标URL失败： %v", err)
    }
    u.Host = "enka.network"
    u.Scheme = *scheme
    if u.Scheme == "" {
        log.Fatalf("目标URL的协议方案为空")
    }

    proxy := httputil.NewSingleHostReverseProxy(u)
    proxy.Transport = &http.Transport{
        DisableKeepAlives: true,
    }

    // 创建一个自定义的处理器，用于修改请求的URL
    handler := func(w http.ResponseWriter, r *http.Request) {
        r.URL.Path = "/" // 设置路径为根路径
        r.URL.RawQuery = "" // 移除查询参数
        proxy.ServeHTTP(w, r)
    }

    log.Printf("启动反向代理，监听 %s，转发到 %s", *listen, *target)
    log.Fatal(http.ListenAndServe(*listen, handler))
}

