package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	targetURL, err := url.Parse("http://enka.network")
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	log.Println("Starting proxy server on 0.0.0.0:7860")
	err = http.ListenAndServe("0.0.0.0:7860", nil)
	if err != nil {
		log.Fatal(err)
	}
}


