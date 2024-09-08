package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	target   = flag.String("target", "", "Target server to proxy")
	listen   = flag.String("listen", ":8080", "Address to listen on")
	username = flag.String("username", "", "Username for basic access authentication")
	password = flag.String("password", "", "Password for basic access authentication")
)

func main() {
	flag.Parse()

	if *target == "" {
		fmt.Println("Error: target is required")
		os.Exit(1)
	}

	proxy := NewProxy(*target)

	if *username != "" || *password != "" {
		auth := fmt.Sprintf("%s:%s", *username, *password)
		authStr := "Basic " + strings.TrimSpace(string(os.ExpandEnv(auth)))
		fmt.Println("HTTP Proxy Listening with Basic Auth:", authStr)
		http.ListenAndServe(*listen, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !checkAuth(r.Header.Get("Authorization"), authStr) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			proxy.ServeHTTP(w, r)
		}))
	} else {
		fmt.Println("HTTP Proxy Listening...")
		http.ListenAndServe(*listen, proxy)
	}
}

type Proxy struct {
	target *string
}

func NewProxy(target string) *Proxy {
	return &Proxy{&target}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxy := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(*p.target)
		},
	}

	client := &http.Client{Transport: proxy}
	client.Do(r)
}

func checkAuth(header, expected string) bool {
	if header == "" {
		return false
	}
	return header == expected
}
