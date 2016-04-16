package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/alphagov/performanceplatform-client.go"
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"
	"github.com/quipo/statsd"
	"gopkg.in/unrolled/render.v1"

	"github.com/alphagov/metadata-api/content_api"
	"github.com/alphagov/metadata-api/need_api"
	"github.com/alphagov/metadata-api/performance_platform"
	"github.com/alphagov/metadata-api/request"
)

var (
	appDomain    = getEnvDefault("GOVUK_APP_DOMAIN", "alphagov.co.uk")
	port         = getEnvDefault("HTTP_PORT", "3000")
	timeout      = getEnvDefault("HTTP_TIMEOUT", "2000")
	httpProtocol = getHttpProtocol(appDomain)

	contentAPI     = httpProtocol + "://contentapi." + appDomain
	needAPI        = httpProtocol + "://need-api." + appDomain
	performanceAPI = "https://www.performance.service.gov.uk"

	renderer = render.New(render.Options{})

	loggingMiddleware = negronilogrus.NewCustomMiddleware(
		logrus.InfoLevel, &logrus.JSONFormatter{}, "metadata-api")
	logging = loggingMiddleware.Logger

	statsdClient = newStatsDClient("localhost:8125", "metadata-api.")
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	renderer.JSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

func InfoHandler(contentAPI, needAPI, performanceAPI string, config *Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var needs []*need_api.Need = make([]*need_api.Need, 0)

		slug := r.URL.Path[len("/info"):]

		if len(slug) <= 1 || slug == "/" {
			renderError(w, http.StatusNotFound, "not found")
			return
		}

		artefactStart := time.Now()
		artefact, err := content_api.FetchArtefact(contentAPI, config.BearerTokenContentAPI, slug)
		statsDTiming("artefact", artefactStart, time.Now())
		if err != nil {
			if err == request.NotFoundError {
				renderError(w, http.StatusNotFound, err.Error())
				return
			}

			renderError(w, http.StatusInternalServerError, "Artefact: "+err.Error())
			return
		}

		needStart := time.Now()
		for _, needID := range artefact.Details.NeedIDs {
			need, err := need_api.FetchNeed(needAPI, config.BearerTokenNeedAPI, needID)
			if err != nil {
				renderError(w, http.StatusInternalServerError, "Need: "+err.Error())
				return
			}
			needs = append(needs, need)
		}
		statsDTiming("needs", needStart, time.Now())

		performanceStart := time.Now()
		ppClient := performanceclient.NewDataClient(performanceAPI, logging)
		is_multipart := (len(artefact.Details.Parts) != 0) || (artefact.Format == "smart-answer")
		performance, err := performance_platform.SlugStatistics(ppClient, slug, is_multipart)
		statsDTiming("performance", performanceStart, time.Now())
		if err != nil {
			renderError(w, http.StatusInternalServerError, "Performance: "+err.Error())
			return
		}

		metadata := &Metadata{
			Artefact:     artefact,
			Needs:        needs,
			Performance:  performance,
			ResponseInfo: &ResponseInfo{Status: "ok"},
		}

		renderer.JSON(w, http.StatusOK, metadata)
	}
}

func main() {
	config, err := ReadConfig("config.json")
	if err != nil {
		logging.Fatalln("Couldn't load configuration", err)
	}

	i, err := strconv.Atoi(timeout)
	if err != nil {
		logging.Fatalln("Timeout wasn't an integer of milliseconds")
	}

	http.DefaultClient.Timeout = time.Duration(i) * time.Millisecond
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/healthcheck", HealthCheckHandler)
	httpMux.HandleFunc("/info/", InfoHandler(
		contentAPI, needAPI, performanceAPI, config))

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

func getHttpProtocol(appDomain string) string {
	if appDomain == "dev.gov.uk" {
		return "http"
	}

	return "https"
}

func newStatsDClient(host, prefix string) *statsd.StatsdClient {
	statsdClient := statsd.NewStatsdClient(host, prefix)
	statsdClient.CreateSocket()

	return statsdClient
}

func statsDTiming(label string, start, end time.Time) {
	statsdClient.Timing("time."+label,
		int64(end.Sub(start)/time.Millisecond))
}
