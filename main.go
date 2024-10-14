package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Magic kkkk
	body, _ := io.ReadAll(r.Body)
	values, _ := url.ParseQuery(string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Recovery user type/prefix
	user_type := strings.Split(values["username"][0], ":")

	var targetURL *url.URL

	// Check user type/prefix and send to correct rhsso
	if user_type[0] == "XXXF" {
		targetURL, _ = url.Parse("https://xxpf.xxx.com.br")
	} else if user_type[0] == "XXJ" {
		targetURL, _ = url.Parse("https://xxpj.xxx.com.br")
	} else {
		targetURL, _ = url.Parse("https://xxpf.xxx.com.br")
		return
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		IdleConnTimeout: 90 * time.Second,
	}

	// Configure reverse proxy with correct informations
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = targetURL.Host

	// Send request to backend
	proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", proxyHandler)

	fmt.Println("Proxy up and running... ;) ")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error to start server: %v\n", err)
	}
}
