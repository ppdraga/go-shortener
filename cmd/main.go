package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/ppdraga/go-shortener/app"
	"github.com/ppdraga/go-shortener/database"
	linkc "github.com/ppdraga/go-shortener/internal/shortener/link"
	linkwdb "github.com/ppdraga/go-shortener/internal/shortener/link/withdb"
	"github.com/ppdraga/go-shortener/settings"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	logger.Info("Starting the application")

	if err := godotenv.Load(); err != nil {
		logger.Infof("No .env file found")
	}

	port := os.Getenv("PORT")
	settings.Config = settings.Settings{
		Port: port,
	}
	logger.Infof("Application access port: %s", port)

	// Database init
	rsc, err := database.InitDB(logger)
	if err != nil {
		//logger.Panic("Can't initialize resources.", "err", err)
		logger.Infof("Can't initialize resources. %v", err)
	}
	defer func() {
		err := rsc.Release()
		if err != nil {
			logger.Error("Got an error during resources release.", "err", err)
		}
	}()

	linkdb := linkwdb.New(rsc.DB)
	linkCtrl := linkc.NewController(linkdb)

	r := mux.NewRouter()
	r.HandleFunc("/", app.HomeHandler()).Methods("GET")
	r.HandleFunc("/_api", app.APIHomeHandler()).Methods("GET")
	r.HandleFunc("/_api/", app.APIHomeHandler()).Methods("GET")
	r.HandleFunc("/_api/link", app.APIHandler(linkCtrl))
	r.HandleFunc("/_api/link/", app.APIHandler(linkCtrl))
	r.HandleFunc("/_api/link/{id:[0-9]+}", app.APIHandler(linkCtrl))
	r.HandleFunc("/{.+}", app.RedirectHandler(linkCtrl)).Methods("GET")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	shutdown := make(chan error, 1)

	server := http.Server{
		Addr:    net.JoinHostPort("", settings.Config.Port),
		Handler: r,
	}

	go func() {
		err := server.ListenAndServe()
		logger.Errorf("Got an error during ListenAndServe: %v", err)
		shutdown <- err
	}()
	logger.Infof("The service is ready to listen and serve")

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

	logger.Infof("The service is stopping...")
	err = server.Shutdown(context.Background())
	if err != nil {
		logger.Infof("Got an error during service shutdown: %v", err)
	}
}
