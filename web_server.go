package enki

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

func (ek *Enki) ListenAndServe() {
	port := fmt.Sprintf(":%v", webPort)

	server := &http.Server{
		Addr:         port,
		Handler:      ek.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	log.Println("Web Applications is starting...")
	log.Println("* Enki version:", ek.Version())
	log.Println("*    Time zone:", timeZone)
	log.Println("*   Go version:", runtime.Version())
	log.Println("*   Process ID:", os.Getpid())
	log.Println("*   Using port:", port)
	log.Println("Use Ctrl-C to stop")

	halt := make(chan struct{})

	ctx := context.Background()

	go gracefulShutdown(ctx, server, halt)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-halt
}
