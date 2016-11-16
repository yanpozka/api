package main

import (
	"log"
	"net/http"
	"time"
)

func loggerMW(inner http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf("%s %q %s Time consumed: %s", r.Method, r.URL, r.RemoteAddr, time.Since(start))
	})
}

func headersMW(inner http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Server", "MioServer")

		w.Header().Set("Content-Type", "application/json")

		inner.ServeHTTP(w, r)
	})
}

// func empty(inner http.Handler) http.Handler {

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		inner.ServeHTTP(w, r)
// 	})
// }
