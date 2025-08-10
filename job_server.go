package enki

import (
	"context"
	"log"
	"os"
	"reflect"
	"runtime"

	"github.com/JesseChavez/enki/job"
)

func (ek *Enki) StartAndProcess() {
	server := &job.JobExecutor{
		Env:  ek.Env,
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

	jobString := initArgs[0]

	log.Println("Job:", jobString)

	jobType := ek.Queues.List[jobString]

	if jobType == nil {
		log.Println("Job not defined:", jobString)
		return
	}

	job := reflect.New(jobType)

	// log.Println("Registered Job:", job)

	method := job.MethodByName("Perform")

	if !method.IsValid() {
		log.Println("Perform not defined")
		return
	}

	// Prepare arguments for the method call
	// args := []reflect.Value{reflect.ValueOf("my args?")}
	args := []reflect.Value{}

	// Call the method
	method.Call(args)
}
