package enki

import (
	"context"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/JesseChavez/enki/job"
)

func (ek *Enki) StartAndProcess() {
	server := &job.JobExecutor{
		Env:  ek.Env,
		DB: ek.DB,
		JobSupport: ek.JobSupport,
	}

	log.Println("Job Application is starting...")
	log.Println("* Enki version:", ek.Version())
	log.Println("*  Environment:", ek.Env)
	log.Println("*    Time zone:", timeZone)
	log.Println("*   Go version:", runtime.Version())
	log.Println("*   Process ID:", os.Getpid())
	log.Println("Use Ctrl-C to stop")
	log.Println("Starting BG jobs ...")

	halt := make(chan struct{})

	ctx := context.Background()

	go sigGracefulShutdown(ctx, server, halt)
	// go apiGracefulShutdown(ctx, server, halt)

	if err := server.StartAndProcess(); err != nil {
		// Error starting or closing listener:
		log.Fatalf("Job server StartAndProcess: %v", err)
	}

	<-halt
}

func (ek *Enki) StartAndWork() {
	initArgs := os.Args[2:]

	log.Println("Work args:", initArgs)

	if len(initArgs) < 1 {
		log.Println("Job argument is required")
		return
	}

	jobString := initArgs[0]

	args := Args{}

	jobArgs :=initArgs[1:]

	for _, jobArg := range jobArgs {
		delimiter := "="
		pair := strings.Split(jobArg, delimiter)
		if len(pair) == 2 {
			args[pair[0]] = pair[1]
		}

	}

	ek.JobSupport.PerformNow(jobString, args)
}
