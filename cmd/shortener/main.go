package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var storage map[string]string

func GetHandler(w http.ResponseWriter, r *http.Request) {
	paramsArray := strings.Split(r.URL.Path, "/")
	if len(paramsArray) == 2 {

		url := storage[paramsArray[1]]
		if len(url) > 0 {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	} else {
		http.Error(w, "Incorrect parameters", http.StatusBadRequest)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	responseData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Incorrect parameters", http.StatusBadRequest)
	}

	responseString := string(responseData)
	if len(responseString) > 0 {
		key := RandString(8)
		storage[key] = responseString
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + key))

		return
	}

	http.Error(w, "Incorrect parameters", http.StatusBadRequest)
}

func RouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetHandler(w, r)
	} else if r.Method == http.MethodPost {
		PostHandler(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	storage = make(map[string]string)
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, RouteHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
