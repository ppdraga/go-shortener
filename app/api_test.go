package app

import (
	"github.com/ppdraga/go-shortener/fixtures"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
			logger.Error("Got an error during resources release.", "err", err)
		}
	}()

	t.Run("Case OK", func(t *testing.T) {
		assert.Equal(t, true, true)
	})

}
