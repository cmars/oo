package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	target = flag.String("target", "", "target URL to proxy")
)

func main() {
	flag.Parse()
	if *target == "" {
		log.Fatalf("-target is required")
	}

	proxy := glazed.NewProxy(*target)
	handler := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			targetURL := proxy.Target()

			if !proxy.isAuthenticated(req) {
				proxy.authenticate(req)
			}

			req.URL.Host = targetURL.Host

			req.Host = targetURL.Host
			req.URL.Scheme = targetURL.Scheme
			req.URL.Path = singleJoiningSlash(targetURL.Path, req.URL.Path)

			targetQuery := targetURL.RawQuery
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
		},
	}
	http.ListenAndServe(":8080", proxy)
}
