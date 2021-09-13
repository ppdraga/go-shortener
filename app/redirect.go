package app

import (
	"fmt"
	"github.com/ppdraga/go-shortener/internal/restapi"
	linkc "github.com/ppdraga/go-shortener/internal/shortener/link"
	"net/http"
)

func RedirectHandler(linkCtrl *linkc.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortLink := r.URL.Path
		if shortLink == "" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "WELCOME! This is a Redirector of simple shortener service!\n")
			return
		}
		link, err := linkCtrl.FindLink(shortLink[1:])
		if err != nil {
			errMsg := fmt.Sprintf("Can't find link: %v", err)
			restapi.ResponseBadRequest(errMsg, w)
			return
		}
		http.Redirect(w, r, *link.Resource, http.StatusFound)
	}
}
