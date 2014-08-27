package main

import (
	"fmt"
	"log"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {
	http.HandleFunc("/healthcheck", HealthCheckHandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
