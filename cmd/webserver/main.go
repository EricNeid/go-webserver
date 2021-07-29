package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/EricNeid/go-webserver/server"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	listenAddr string = ":5000"
)

func init() {
	value, isSet := os.LookupEnv("LISTEN_ADDR")
	if isSet {
		listenAddr = value
	}

	flag.StringVar(&listenAddr, "listen-addr", listenAddr, "server listen address")
	flag.Parse()

	log.SetFlags(log.LstdFlags)
	log.SetPrefix("http: ")
	log.SetOutput(os.Stdout)

	log.SetOutput(
		io.MultiWriter(
			os.Stdout,
			&lumberjack.Logger{
				Filename:   "log/webserver.log",
				MaxSize:    500, // megabytes
				MaxBackups: 3,
				MaxAge:     28, //days
			},
		),
	)
}

func main() {
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := server.NewApplicationServer(log.Default(), listenAddr)
	go server.GracefullShutdown(quit, done)

	log.Println("Server is ready to handle requests at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	log.Println("Server stopped")
}
