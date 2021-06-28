package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/EricNeid/go-webserver/server"
)

var (
	listenAddr string = ":5000"
)

func readEnvironmentVariables() {
	value, isSet := os.LookupEnv("LISTEN_ADDR")
	if isSet {
		listenAddr = value
	}
}

func readCli() {
	flag.StringVar(&listenAddr, "listen-addr", listenAddr, "server listen address")
	flag.Parse()
}

func main() {
	readEnvironmentVariables()
	readCli()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := server.NewApplicationServer(logger, listenAddr)
	go server.GracefullShutdown(quit, done)

	logger.Println("Server is ready to handle requests at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}
