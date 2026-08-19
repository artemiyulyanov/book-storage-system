package proxy

import (
	"api-gateway/internal/middleware"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Route struct {
	PathPrefix       string
	Target           *url.URL
	ProtectedMethods map[string]bool
}

type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{}
}

func (router *Router) RegisterService(pathPrefix, targetURL string, protectedMethods []string) error {
	target, err := url.Parse(targetURL)

	if err != nil {
		return err
	}

	if !strings.HasSuffix(pathPrefix, "/") {
		pathPrefix += "/"
	}

	methodsSet := make(map[string]bool, len(protectedMethods))
	if protectedMethods != nil {
		for _, m := range protectedMethods {
			methodsSet[m] = true
		}
	}

	router.routes = append(router.routes, Route{
		PathPrefix:       pathPrefix,
		Target:           target,
		ProtectedMethods: methodsSet,
	})

	return nil
}

func (router *Router) Handler(jwtSecret string) http.Handler {
	mux := http.NewServeMux()

	for _, route := range router.routes {
		proxy := newReverseProxy(route.Target)
		trimmedPrefix := strings.TrimSuffix(route.PathPrefix, "/")

		var handler http.Handler = http.StripPrefix(trimmedPrefix, proxy)
		handler = middleware.JWTAuthForMethods(jwtSecret, route.ProtectedMethods)(handler)

		mux.Handle(route.PathPrefix, handler)
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
