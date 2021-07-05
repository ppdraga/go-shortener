package app

import (
	"bytes"
	"encoding/json"
	"github.com/ppdraga/go-shortener/fixtures"
	linkc "github.com/ppdraga/go-shortener/internal/shortener/link"
	"github.com/ppdraga/go-shortener/internal/shortener/link/datatype"
	linkwdb "github.com/ppdraga/go-shortener/internal/shortener/link/withdb"
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
		apiHandler.ServeHTTP(rw, req)

		assert.Equal(t, rw.Code, http.StatusCreated)

		var bodyItem datatype.Link
		json.NewDecoder(rw.Body).Decode(&bodyItem)
		assert.Equal(t, bodyItem.Resource, linkBody.Resource)
		assert.Equal(t, bodyItem.CustomName, linkBody.CustomName)
	})

}
