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

func gracefulShutdown(ctx context.Context, server *http.Server, halt chan struct{}) {
	sigint := make(chan os.Signal, 1)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

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
