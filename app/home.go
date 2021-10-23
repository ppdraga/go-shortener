package app

import (
	"fmt"
	"github.com/ppdraga/go-shortener/prom"
	"net/http"
)

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prom.HomeHandlerCounter.Inc()
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "WELCOME! This is a simple shortener service!\n")
	}
}
