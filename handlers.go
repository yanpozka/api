package main

import "net/http"

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
}

// func checkErr(err error, w http.ResponseWriter, message string) bool {
// 	if err != nil {
// 		log.Println(message, err)
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 	}
// 	return err != nil
// }
