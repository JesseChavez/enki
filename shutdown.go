package enki

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var Shutdown []func() error

var StatusCode = make(chan string)

func sigGracefulShutdown(ctx context.Context, server *http.Server, halt chan struct{}) {
	sigint := make(chan os.Signal, 1)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigint

	log.Println("IPC signal, shutdown server")

	gracefulShutdown(ctx, server, halt)
}

func apiGracefulShutdown(ctx context.Context, server *http.Server, halt chan struct{}) {
	<-StatusCode

	log.Println("API signal, shutdown server")

	gracefulShutdown(ctx, server, halt)
}

func gracefulShutdown(ctx context.Context, server *http.Server, halt chan struct{}) {
	log.Println("shutting down server gracefully")

	// stop receiving any request.
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("shutdown error", err)
	}

	// close any other things db, redis, etc.
	for i := range Shutdown {
		Shutdown[i]()
	}

	close(halt)
}
