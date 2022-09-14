package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"ipchecker/ipchecker"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	EnvCIDRListFile = "CIDR_LIST_FILE"
	EnvHost         = "HOST"
	EnvPort         = "PORT"
)

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(contentTypeHeader, contentTypeJSON)
		next.ServeHTTP(w, r)
	})
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// get env variables
	host := getEnv(EnvHost, "0.0.0.0")
	port := getEnv(EnvPort, "8000")
	cidrFilename := os.Getenv(EnvCIDRListFile)

	var err error

	// create api
	ch, err := ipchecker.CreateIpChecker(cidrFilename)
	if err != nil {
		log.Panicf("create checker failed: %v", err)
		return
	}

	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	router.UseEncodedPath()
	router.HandleFunc("/check", ch.RequestHandler).Methods("GET").Queries("ip", "{ip}")

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%s", host, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panicf("listen and serve %v", err)
		return
	}
}
