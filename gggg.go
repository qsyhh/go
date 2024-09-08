package main

import (
    "flag"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

type proxyHandler struct {
    target string
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    u, err := url.Parse(h.target)
    if err != nil {
        log.Fatalf("解析目标URL失败： %v", err)
    }
    if u.Scheme == "" {
        log.Printf("目标URL的协议方案为空")
        u.Scheme == "https:\/\/"
    }
    if u.Host == "" {
        log.Fatalf("目标URL的主机名为空")
    }
    log.Printf("解析目标URL：", u)
    proxy := httputil.NewSingleHostReverseProxy(u)
    proxy.Transport = &http.Transport{
        DisableKeepAlives: true,
    }

    r.URL.Path = "/" // 设置路径为根路径
    r.URL.RawQuery = "" // 移除查询参数
    proxy.ServeHTTP(w, r)
}

func main() {
    target := flag.String("target", "", "Target server to proxy")
    listen := flag.String("listen", "", "Address to listen on")
    flag.Parse()

    if *target == "" || *listen == "" {
        log.Fatal("请提供 -target 和 -listen 参数")
    }

    log.Printf("启动反向代理，监听 %s，转发到 %s", *listen, *target)
    log.Fatal(http.ListenAndServe(*listen, &proxyHandler{target: *target}))
}
