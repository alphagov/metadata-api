package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"
	"gopkg.in/unrolled/render.v1"

	"github.com/alphagov/metadata-api/content_api"
)

var (
	bearerToken = getEnvDefault("BEARER_TOKEN", "foo")
	contentAPI  = getEnvDefault("CONTENT_API", "content-api")
	port        = getEnvDefault("PORT", "3000")

	renderer = render.New(render.Options{})
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	renderer.JSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

func InfoHandler(contentAPI, bearerToken string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Path[len("/info"):]

		if len(slug) <= 1 || slug == "/" {
			renderError(w, http.StatusNotFound, "not found")
			return
		}

		artefact, err := content_api.FetchArtefact(contentAPI, bearerToken, slug)
		if err != nil {
			renderError(w, http.StatusInternalServerError, err.Error())
			return
		}

		metadata := &Metadata{
			Artefact:     artefact,
			ResponseInfo: &ResponseInfo{Status: "ok"},
		}

		renderer.JSON(w, http.StatusOK, metadata)
	}
}

func main() {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/healthcheck", HealthCheckHandler)
	httpMux.HandleFunc("/info", InfoHandler(contentAPI, bearerToken))

	middleware := negroni.New()
	middleware.Use(negronilogrus.NewCustomMiddleware(
		logrus.InfoLevel, &logrus.JSONFormatter{}, "metadata-api"))
	middleware.UseHandler(httpMux)

	middleware.Run(":" + port)
}

func renderError(w http.ResponseWriter, status int, errorString string) {
	renderer.JSON(w, status, &Metadata{ResponseInfo: &ResponseInfo{Status: errorString}})
}

func getEnvDefault(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	return val
}
