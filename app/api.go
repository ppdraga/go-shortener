package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/ppdraga/go-shortener/internal/restapi"
	linkc "github.com/ppdraga/go-shortener/internal/shortener/link"
	"github.com/ppdraga/go-shortener/internal/shortener/link/datatype"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"strconv"
)

func APIHomeHandler(linkCtrl *linkc.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span, _ := opentracing.StartSpanFromContextWithTracer(r.Context(), linkCtrl.Tracer, "APIHomeHandler")
		defer span.Finish()

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "WELCOME! This is an API interface of a simple shortener service!\n")
	}
}

func APIHandler(linkCtrl *linkc.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span, _ := opentracing.StartSpanFromContextWithTracer(r.Context(), linkCtrl.Tracer, "APIHandler")
		defer span.Finish()

		linkCtrl.Logger.Info("APIHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})
		vars := mux.Vars(r)
		_ = vars
		if r.Method == "POST" {
			linkItem := new(datatype.Link)
			err := json.NewDecoder(r.Body).Decode(linkItem)
			if err != nil {
				restapi.ResponseBadRequest("Couldn't parse request body", w)
				return
			}
			err, linkID := linkCtrl.AddLink(linkItem)
			if err != nil {
				errMsg := fmt.Sprintf("Error while adding a link: %v", err)
				restapi.ResponseBadRequest(errMsg, w)
				return
			}
			item, err := linkCtrl.GetLink(linkID)
			if err != nil {
				errMsg := fmt.Sprintf("Error: %v", err)
				restapi.ResponseBadRequest(errMsg, w)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(*item)
			return
		}

		if r.Method == "GET" {
			id, ok := vars["id"]
			if !ok {
				restapi.ResponseBadRequest("Couldn't parse id param", w)
				return
			}
			idInt, err := strconv.Atoi(id)
			if err != nil {
				errMsg := fmt.Sprintf("Error %v", err)
				restapi.ResponseBadRequest(errMsg, w)
				return
			}
			item, err := linkCtrl.GetLink(int64(idInt))
			if err != nil {
				errMsg := fmt.Sprintf("Error: %v", err)
				restapi.ResponseBadRequest(errMsg, w)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(*item)
			return
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
