package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ppdraga/go-shortener/internal/restapi"
	linkc "github.com/ppdraga/go-shortener/internal/shortener/link"
	"github.com/ppdraga/go-shortener/internal/shortener/link/datatype"
	"net/http"
)

func APIHomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "WELCOME! This is an API interface of a simple shortener service!\n")
	}
}

func APIHandler(linkCtrl *linkc.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		_ = vars
		if r.Method == "POST" {
			linkItem := new(datatype.Link)
			err := json.NewDecoder(r.Body).Decode(linkItem)
			if err != nil {
				restapi.ResponseBadRequest("Couldn't parse request body", w)
				return
			}
			linkCtrl.AddLink(linkItem)
		}

		if r.Method == "GET" {
			id, ok := vars["id"]
			fmt.Println(id, ok)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "WELCOME! This will be an API interface of a simple shortener service!\n")
	}
}

func APILinkHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "WELCOME! This is an API Link handler!\n")
	}
}
