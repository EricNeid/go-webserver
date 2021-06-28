package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

const timeFormat string = "Mon Jan 2 15:04:05 2006"

type ApplicationServer struct {
	Webserver *http.Server
	Logger    *log.Logger
}

func NewApplicationServer(logger *log.Logger, listenAddr string) ApplicationServer {
	// create webserver
	router := http.NewServeMux()

	// configure routes
	router.HandleFunc("/", logCall(logger, welcome))

	return ApplicationServer{
		Logger: logger,
		Webserver: &http.Server{
			Addr:         listenAddr,
			Handler:      router,
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}

func (srv ApplicationServer) GracefullShutdown(quit <-chan os.Signal, done chan<- bool) {
	<-quit
	server := srv.Webserver
	logger := srv.Logger

	logger.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(done)
}

func (srv ApplicationServer) ListenAndServe() error {
	return srv.Webserver.ListenAndServe()
}

func logCall(logger *log.Logger, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now()
		logger.Printf("%s - %s\n", timestamp.Format(timeFormat), r.URL.Path)
		handler(w, r)
	}
}
