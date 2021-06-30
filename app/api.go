package app

import (
	"fmt"
	"github.com/ppdraga/go-shortener/database"
	"net/http"
)

func APIHomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "WELCOME! This is an API interface of a simple shortener service!\n")
	}
}

func APIHandler(rsc *database.R) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "WELCOME! This will be an API interface of a simple shortener service!\n")
	}
}
