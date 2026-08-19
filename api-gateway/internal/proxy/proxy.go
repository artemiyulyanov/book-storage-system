package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Route struct {
	PathPrefix string
	Target     *url.URL
}

type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{}
}

func (router *Router) RegisterService(pathPrefix, targetURL string) error {
	target, err := url.Parse(targetURL)

	if err != nil {
		return err
	}

	if !strings.HasSuffix(pathPrefix, "/") {
		pathPrefix += "/"
	}

	router.routes = append(router.routes, Route{
		PathPrefix: pathPrefix,
		Target:     target,
	})

	return nil
}

func (router *Router) Handler() http.Handler {
	mux := http.NewServeMux()

	for _, route := range router.routes {
		proxy := newReverseProxy(route.Target)

		trimmedPrefix := strings.TrimSuffix(route.PathPrefix, "/")
		mux.Handle(route.PathPrefix, http.StripPrefix(trimmedPrefix, proxy))
	}

	return mux
}

func newReverseProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error: %s -> %s: %v", r.URL.Path, target.String(), err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error":"service unavailable"}`))
	}

	return proxy
}
