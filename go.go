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
    scheme := flag.String("scheme", "", "Protocol scheme for the target URL")
    flag.Parse()

    if *target == "" || *listen == "" {
        log.Fatal("请提供 -target 和 -listen 参数")
    }

    u, err := url.Parse(*target)
    if err != nil {
        log.Fatalf("解析目标URL失败： %v", err)
    }
    u.Scheme = *scheme
    if u.Scheme == "" {
        log.Fatalf("目标URL的协议方案为空")
    }

    proxy := httputil.NewSingleHostReverseProxy(u)

    // 添加自定义的http.Handler
    handler := struct {
        http.Handler
    }{
        Handler: func(w http.ResponseWriter, r *http.Request) {
            r.Host = u.Host // 设置正确的Host字段
            proxy.ServeHTTP(w, r)
        },
    }

    log.Printf("启动反向代理，监听 %s，转发到 %s", *listen, *target)
    log.Fatal(http.ListenAndServe(*listen, handler))
}
