package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
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
