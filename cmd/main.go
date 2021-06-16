package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/ppdraga/go-shortener/app"
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
	settings.Config = settings.Settings{
		Port: os.Getenv("PORT"),
	}
	// TODO get port from settings file or environment
	//port := "8080"

	r := mux.NewRouter()
	r.HandleFunc("/", app.HomeHandler())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	shutdown := make(chan error, 1)

	server := http.Server{
		Addr:    net.JoinHostPort("", settings.Config.Port),
		Handler: r,
	}

	go func() {
		err := server.ListenAndServe()
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
	err := server.Shutdown(context.Background())
	if err != nil {
		logger.Infof("Got an error during service shutdown: %v", err)
	}
}
