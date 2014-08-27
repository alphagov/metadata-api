package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"
	"gopkg.in/unrolled/render.v1"
)

var renderer = render.New(render.Options{})

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	renderer.JSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

func main() {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/healthcheck", HealthCheckHandler)

	middleware := negroni.New()
	middleware.Use(negronilogrus.NewCustomMiddleware(
		logrus.InfoLevel, &logrus.JSONFormatter{}, "metadata-api"))
	middleware.UseHandler(httpMux)

	middleware.Run(":3000")
}
