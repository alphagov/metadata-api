package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"
	"gopkg.in/unrolled/render.v1"
)

type ResponseInfo struct {
	Status string `json:"status"`
}

type Metadata struct {
	ResponseInfo *ResponseInfo `json:"_response_info"`
}

var renderer = render.New(render.Options{})

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	renderer.JSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[len("/info"):]

	if len(slug) <= 1 || slug == "/" {
		renderer.JSON(w, http.StatusNotFound, &Metadata{ResponseInfo: &ResponseInfo{Status: "not found"}})
		return
	}

	renderer.JSON(w, http.StatusOK, &Metadata{ResponseInfo: &ResponseInfo{Status: "ok"}})
}

func main() {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/healthcheck", HealthCheckHandler)
	httpMux.HandleFunc("/info", InfoHandler)

	middleware := negroni.New()
	middleware.Use(negronilogrus.NewCustomMiddleware(
		logrus.InfoLevel, &logrus.JSONFormatter{}, "metadata-api"))
	middleware.UseHandler(httpMux)

	middleware.Run(":3000")
}
