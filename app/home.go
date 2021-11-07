package app

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	linkc "github.com/ppdraga/go-shortener/internal/shortener/link"
	"net/http"
)

func HomeHandler(linkCtrl *linkc.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span, _ := opentracing.StartSpanFromContextWithTracer(r.Context(), linkCtrl.Tracer, "APIHomeHandler")
		defer span.Finish()

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "WELCOME! This is a simple shortener service!\n")
	}
}
