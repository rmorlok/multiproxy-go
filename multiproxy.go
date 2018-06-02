// Based on: https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c
package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"strings"
	"os"
	"fmt"
	"net/url"
	"io"
)

var hosts []url.URL

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	orig_url := *req.URL

	for i := 0; i < len(hosts); i++ {
		// The host we are trying
		h := hosts[i]

		var path string

		if h.Path == "" || h.Path == "/" {
			path = orig_url.Path
		} else {
			if strings.HasSuffix(h.Path, "/") {
				path = fmt.Sprintf("%s%s", h.Path, orig_url.Path[1:])
			} else {
				path = fmt.Sprintf("%s%s", h.Path, orig_url.Path)
			}
		}

		// Combine the data to create a new URL
		u := url.URL{
			Scheme: h.Scheme,
			Opaque: orig_url.Opaque,
			User: h.User,
			Host: h.Host,
			Path: path,
			RawQuery: orig_url.RawQuery,
		}

		log.Printf("Checking path %s", u.String())

		// Update the request with the new URL
		req.URL = &u

		resp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode > 400 {
			// Didn't find what we are looking for
			continue
		} else {
			// Found it!
			log.Printf("Found at path %s...", u.String())
			copyHeader(w.Header(), resp.Header)
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
			return
		}
	}

	// Didn't find the URL at any of our hosts
	log.Printf("Did not find host for %s", orig_url.Path)
	w.WriteHeader(404)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func main() {

	var pemPath string
	flag.StringVar(&pemPath, "pem", "server.pem", "path to pem file")

	var keyPath string
	flag.StringVar(&keyPath, "key", "server.key", "path to key file")

	var proto string
	flag.StringVar(&proto, "proto", "http", "Proxy protocol (http or https)")

	var ip string
	flag.StringVar(&ip, "ip", "0.0.0.0", "IP to bind to")

	var port int
	flag.IntVar(&port, "port", 8888, "Port to serve traffic from")

	flag.Parse()

	if proto != "http" && proto != "https" {
		log.Fatal("Protocol must be either http or https")
		os.Exit(1)
	}

	if len(flag.Args()) == 0 {
		log.Fatal("Must specify at least one proxy host")
		os.Exit(1)
	}

	for i := 0; i < len(flag.Args()); i++ {
		u, err := url.Parse(flag.Args()[i])
		if err != nil {
			log.Fatal(fmt.Sprintf("Host '%s' is invalid", flag.Args()[i]))
			os.Exit(1)
		}
		hosts = append(hosts, *u)
	}

	log.Printf("Setting up multi-proxy to the following hosts: %s", strings.Join(flag.Args(), ", "))
	log.Printf("Serving proxy on %s:%d", ip, port)

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", ip, port),
		Handler: http.HandlerFunc(handleHTTP),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	if proto == "http" {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
	}
}
