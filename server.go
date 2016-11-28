package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

const (
	defaultLogFile    = "/tmp/apiserver.log"
	defaultListenPort = ":9090"
)

type paramsCtxType struct{}

var paramsCtxKey = paramsCtxType{}

func main() {
	{
		f, err := os.OpenFile(getOrDefault("LOG_FILE", defaultLogFile), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("error:", err)
		}

		log.SetOutput(io.MultiWriter(f, os.Stdout))
	}

	lp := getOrDefault("LISTEN_PORT", defaultListenPort)
	log.Printf("Starting listening in %q\n", lp)

	err := http.ListenAndServe(lp, createRouter())
	if err != nil {
		log.Fatal("Error ListenAndServe: ", err)
	}
}

func createRouter() http.Handler {
	router := httprouter.New()

	// panic recover
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		log.Printf("Recovering: %+v\nrequest: %s %q %s", v, r.Method, r.URL, r.RemoteAddr)
		debug.PrintStack()
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// not found handler
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Detected Not Found: %s %q %s", r.Method, r.URL, r.RemoteAddr)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Detected Method Now Allowed: %s %q %s", r.Method, r.URL, r.RemoteAddr)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	})

	commonChain := alice.New(loggerMW, headersMW)

	// routers:
	{
		router.Handler(http.MethodGet, "/health", commonChain.ThenFunc(health))

		router.GET("/resource/:resourceID", wrap(commonChain.ThenFunc(resource)))
	}

	return router
}

func wrap(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), paramsCtxKey, p)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

func getOrDefault(varName, defaultVal string) string {
	val := os.Getenv(varName)
	if val == "" {
		val = defaultVal
	}
	return val
}
