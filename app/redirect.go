package app

import (
	"fmt"
	"net/http"
)

func RedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "WELCOME! This is a Redirector of simple shortener service!\n")
	}
}
