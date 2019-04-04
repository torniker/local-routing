package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var routes map[string]string

func main() {
	portPtr := flag.Int("port", 80, "port to run")
	flag.Parse()
	port := 80
	if portPtr != nil {
		port = *portPtr
	}
	args := flag.Args()
	routes = make(map[string]string)
	for _, a := range args {
		r := strings.Split(a, "->")
		if len(r) == 2 {
			routes[r[0]] = r[1]
			fmt.Printf("routing %v to %v\n", r[0], r[1])
		}
	}
	if len(routes) == 0 {
		fmt.Println("no routes provided")
		return
	}
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		panic(err)
	}
}

type responseErr struct {
	Error string `json:"error"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// routes := getEvnRoutes()
	if _, ok := routes[r.Host]; !ok {
		w.WriteHeader(http.StatusNotFound)
		body := responseErr{
			Error: fmt.Sprintf("could not find route for host: %v", r.Host),
		}
		json.NewEncoder(w).Encode(body)
		return
	}
	remote, err := url.Parse(routes[r.Host])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		body := responseErr{
			Error: fmt.Sprintf("could not parse url: %v", r.Host),
		}
		json.NewEncoder(w).Encode(body)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func getEvnRoutes() map[string]string {
	routes := make(map[string]string)
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "ROUTE") {
			continue
		}
		r := strings.Split(strings.Split(e, "=")[1], "->")
		routes[r[0]] = r[1]
	}
	return routes
}
func proxyHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		p.ServeHTTP(w, r)
	}
}
