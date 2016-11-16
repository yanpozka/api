package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var tserver *httptest.Server

func TestMain(m *testing.M) {
	tserver = httptest.NewServer(createRouter())
	defer tserver.Close()

	m.Run()
}

func TestHealthOK(t *testing.T) {
	r, err := http.Get(tserver.URL + "/health")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Fatalf("Expected: %q we got: %q", http.StatusText(http.StatusOK), http.StatusText(r.StatusCode))
	}

	dec := json.NewDecoder(r.Body)
	var status bool

	if err := dec.Decode(&status); err != nil {
		t.Fatal(err)
	}
	if status == false {
		t.Fatal("Wrong status:", status)
	}
}
