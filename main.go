package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"
	"gopkg.in/unrolled/render.v1"

	"github.com/alphagov/metadata-api/content_api"
	"github.com/alphagov/metadata-api/need_api"
)

var (
	config            *Config
	logging           *logrus.Logger
	loggingMiddleware *negronilogrus.Middleware

	appDomain = getEnvDefault("GOVUK_APP_DOMAIN", "alphagov.co.uk")
	port      = getEnvDefault("HTTP_PORT", "3000")

	contentAPI = "contentapi." + appDomain
	needAPI    = "need-api." + appDomain

	renderer = render.New(render.Options{})
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	renderer.JSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

func InfoHandler(contentAPI, needAPI string, config *Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var needs []*need_api.Need

		slug := r.URL.Path[len("/info"):]

		if len(slug) <= 1 || slug == "/" {
			renderError(w, http.StatusNotFound, "not found")
			return
		}

		artefact, err := content_api.FetchArtefact(contentAPI, config.BearerTokenContentAPI, slug)
		if err != nil {
			renderError(w, http.StatusInternalServerError, err.Error())
			return
		}

		for _, needID := range artefact.Details.NeedIDs {
			need, err := need_api.FetchNeed(needAPI, config.BearerTokenNeedAPI, needID)
			if err != nil {
				renderError(w, http.StatusInternalServerError, err.Error())
				return
			}
			needs = append(needs, need)
		}

		metadata := &Metadata{
			Artefact:     artefact,
			Needs:        needs,
			ResponseInfo: &ResponseInfo{Status: "ok"},
		}

		renderer.JSON(w, http.StatusOK, metadata)
	}
}

func init() {
	var err error

	loggingMiddleware = negronilogrus.NewCustomMiddleware(
		logrus.InfoLevel, &logrus.JSONFormatter{}, "metadata-api")
	logging = loggingMiddleware.Logger

	config, err = ReadConfig("config.json")
	if err != nil {
		logging.Fatalln("Couldn't load configuration", err)
	}
}

func main() {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/healthcheck", HealthCheckHandler)
	httpMux.HandleFunc("/info", InfoHandler(contentAPI, needAPI, config))

	middleware := negroni.New()
	middleware.Use(loggingMiddleware)
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
