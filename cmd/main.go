package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/ppdraga/go-shortener/app"
	"github.com/ppdraga/go-shortener/database"
	linkc "github.com/ppdraga/go-shortener/internal/shortener/link"
	linkwdb "github.com/ppdraga/go-shortener/internal/shortener/link/withdb"
	"github.com/ppdraga/go-shortener/settings"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

type zapWrapper struct {
	logger *zap.Logger
}

func (w *zapWrapper) Error(msg string) {
	w.logger.Error(msg)
}

func (w *zapWrapper) Infof(msg string, args ...interface{}) {
	w.logger.Sugar().Infof(msg, args...)
}

func initJaeger(service string, logger *zap.Logger) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(&zapWrapper{logger: logger}))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return tracer, closer
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()

	// Tracing service
	tracer, closer := initJaeger("shortener", logger)
	defer closer.Close()

	logger.Info("Starting the application")

	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found")
	}

	port := os.Getenv("PORT")
	settings.Config = settings.Settings{
		Port: port,
	}
	logger.Info(fmt.Sprintf("Application access port: %s", port))

	// Database init
	rsc, err := database.InitDB(logger)
	if err != nil {
		//logger.Panic("Can't initialize resources.", "err", err)
		logger.Info("Can't initialize resources.", zap.Error(err))
	}
	defer func() {
		err := rsc.Release()
		if err != nil {
			logger.Error("Got an error during resources release.", zap.Error(err))
		}
	}()

	linkdb := linkwdb.New(rsc.DB)
	linkCtrl := linkc.NewController(linkdb, logger, tracer)

	r := mux.NewRouter()
	r.HandleFunc("/", app.HomeHandler(linkCtrl)).Methods("GET")
	r.HandleFunc("/_api", app.APIHomeHandler(linkCtrl)).Methods("GET")
	r.HandleFunc("/_api/", app.APIHomeHandler(linkCtrl)).Methods("GET")
	r.HandleFunc("/_api/link", app.APIHandler(linkCtrl))
	r.HandleFunc("/_api/link/", app.APIHandler(linkCtrl))
	r.HandleFunc("/_api/link/{id:[0-9]+}", app.APIHandler(linkCtrl))
	r.HandleFunc("/{.+}", app.RedirectHandler(linkCtrl)).Methods("GET")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	shutdown := make(chan error, 1)

	server := http.Server{
		Addr:    net.JoinHostPort("", settings.Config.Port),
		Handler: nethttp.Middleware(tracer, r),
	}

	go func() {
		err := server.ListenAndServe()
		logger.Error("Got an error during ListenAndServe", zap.Error(err))
		shutdown <- err
	}()
	logger.Info("The service is ready to listen and serve")

	select {
	case killSignal := <-interrupt:
		switch killSignal {
		case os.Interrupt:
			logger.Info("Got SIGINT...")
		case syscall.SIGTERM:
			logger.Info("Got SIGTERM...")
		}
	case <-shutdown:
		logger.Info("Got an error...")
	}

	logger.Info("The service is stopping...")
	err = server.Shutdown(context.Background())
	if err != nil {
		logger.Info("Got an error during service shutdown", zap.Error(err))
	}
}
