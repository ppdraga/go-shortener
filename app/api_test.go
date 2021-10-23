package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ppdraga/go-shortener/fixtures"
	linkc "github.com/ppdraga/go-shortener/internal/shortener/link"
	"github.com/ppdraga/go-shortener/internal/shortener/link/datatype"
	linkwdb "github.com/ppdraga/go-shortener/internal/shortener/link/withdb"
	"github.com/ppdraga/go-shortener/prom"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAPIHandler(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	rsc, err := fixtures.InitTestSQLite(logger)
	if err != nil {
		logger.Errorf("Can't initialize resources. Err: %v", err)
	}
	defer func() {
		err := rsc.Release()
		if err != nil {
			logger.Errorf("Got an error during resources release. %v", err)
		}
	}()
	linkdb := linkwdb.New(rsc.DB)
	linkCtrl := linkc.NewController(linkdb)
	apiHandler := APIHandler(linkCtrl)
	redirectHandler := RedirectHandler(linkCtrl)

	t.Run("Case Add Link", func(t *testing.T) {
		url := "http://url.com/something"
		customName := "Custom"
		linkBody := datatype.Link{Resource: &url, CustomName: &customName}
		body, err := json.Marshal(linkBody)
		if err != nil {
			logger.Errorf("Can't marshal object. %v", err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/_api/link/", bytes.NewReader(body))
		rw := httptest.NewRecorder()
		prom.Init()
		apiHandler.ServeHTTP(rw, req)

		assert.Equal(t, rw.Code, http.StatusCreated)

		var bodyItem datatype.Link
		json.NewDecoder(rw.Body).Decode(&bodyItem)
		assert.Equal(t, *bodyItem.Resource, *linkBody.Resource)
		assert.Equal(t, *bodyItem.CustomName, *linkBody.CustomName)
	})

	t.Run("Case Get Link", func(t *testing.T) {
		// Add link
		url := "http://url.com/something-get"
		customName := "Custom-get"
		linkBody := datatype.Link{Resource: &url, CustomName: &customName}
		body, err := json.Marshal(linkBody)
		if err != nil {
			logger.Errorf("Can't marshal object. %v", err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/_api/link/", bytes.NewReader(body))
		rw := httptest.NewRecorder()
		apiHandler.ServeHTTP(rw, req)

		assert.Equal(t, rw.Code, http.StatusCreated)

		// Get link
		var bodyItem datatype.Link
		json.NewDecoder(rw.Body).Decode(&bodyItem)

		path := fmt.Sprintf("/_api/link/%d", *bodyItem.ID)
		req, _ = http.NewRequest(http.MethodGet, path, nil)
		rw = httptest.NewRecorder()
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"id": fmt.Sprint(*bodyItem.ID),
		}
		req = mux.SetURLVars(req, vars)

		apiHandler.ServeHTTP(rw, req)

		var respBodyItem datatype.Link
		json.NewDecoder(rw.Body).Decode(&respBodyItem)
		fmt.Println(respBodyItem)

		assert.Equal(t, rw.Code, http.StatusOK)
		assert.Equal(t, *respBodyItem.Resource, *linkBody.Resource)
		assert.Equal(t, *respBodyItem.CustomName, *linkBody.CustomName)
	})

	t.Run("Case Redirect", func(t *testing.T) {
		// Add link
		url := "http://url.com/something3"
		customName := "Custom3"
		linkBody := datatype.Link{Resource: &url, CustomName: &customName}
		body, err := json.Marshal(linkBody)
		if err != nil {
			logger.Errorf("Can't marshal object. %v", err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/_api/link/", bytes.NewReader(body))
		rw := httptest.NewRecorder()
		apiHandler.ServeHTTP(rw, req)

		assert.Equal(t, rw.Code, http.StatusCreated)

		var bodyItem datatype.Link
		json.NewDecoder(rw.Body).Decode(&bodyItem)

		// Try redirect
		path := "/" + *bodyItem.ShortLink
		req, _ = http.NewRequest(http.MethodGet, path, nil)
		rw = httptest.NewRecorder()

		redirectHandler.ServeHTTP(rw, req)
		assert.Equal(t, rw.Code, http.StatusFound)
		assert.Equal(t, rw.Result().Header["Location"][0], *linkBody.Resource)
	})

}
