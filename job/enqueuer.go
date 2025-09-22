package job

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/go-rel/rel"
	"github.com/google/uuid"
)

type ActiveJob = reflect.Type

type Args = map[string]string

type Enqueuer struct {
	env     string
	db      rel.Repository
	JobList map[string]ActiveJob
}

type JobDefaults struct {
	Queue    string
	Priority int
	runAt    time.Time
}

func New(env string, db rel.Repository) *Enqueuer {
	enqueuer := Enqueuer{
		env:     env,
		db:      db,
		JobList: make(map[string]ActiveJob),
	}

	return &enqueuer
}

func (enq *Enqueuer) Register(queue string, job ActiveJob) {
	log.Println("====>", job.Name())

	switch queue {
	case "default":
		enq.JobList[job.Name()] = job
	default:
		log.Println("Unknown queue:", queue)
	}
}

func (enq *Enqueuer) PerformNow(jobName string, args Args) (string, error) {
	// log.Println("registered queues:", enq.JobList)
	log.Println("performing job:", jobName)

	jobModel := enq.JobList[jobName]

	if jobModel == nil {
		log.Println("Job not defined:", jobModel)
		return "", errors.New("Undefined Job")
	}
	
	job := reflect.New(jobModel)

	// log.Println("jxxx:", job)

	method := job.MethodByName("Perform")

	if !method.IsValid() {
		log.Println("Perform not defined")
		return "", errors.New("Undefined Perform menthod")
	}


	id := uuid.New().String()
	// Prepare arguments for the method call
	// args := []reflect.Value{reflect.ValueOf("my args?")}
	mArgs := []reflect.Value{reflect.ValueOf(args)}

	// Call the method
	wrappedFailure := method.Call(mArgs)

	failure := wrappedFailure[0].Interface().([]error)

	if len(failure) > 0 {
		log.Println("Error:", failure)
		return id, failure[0]
	}

	return id, nil
}

func (enq *Enqueuer) Enqueue(task *Task, args Args) (string, error) {
	var err error

	jobName := task.name

	// log.Println("registered queues:", enq.JobList)

	jobModel := enq.JobList[jobName]

	if jobModel == nil {
		log.Println("Job not defined:", jobName)
		return "", errors.New("Undefined Job")
	}

	sArgs, err := json.Marshal(args)

	if err != nil {
		log.Println("Failure serializing Args:", args)
		return "", err
	}


	job := reflect.New(jobModel)
	
	// log.Println("xxx:", job)

	values := enq.configuredValues(task, job)

	id := uuid.New().String()

	record := QueuedJob{}

	record.Queue    = values.Queue
	record.Handler  = "JobExecutor"
	record.JobClass = jobName
	record.JobId    = id
	record.Priority = values.Priority
	record.Attempts = 0
	record.State    = "scheduled"
	record.Args     = string(sArgs)
	record.RunAt    = values.runAt

	ctx := context.Background()

	err = enq.db.Insert(ctx, &record)

	return id , err
}

func (enq *Enqueuer) configuredValues(task *Task, job reflect.Value) *JobDefaults {
	queue := "default"
	priority := 1
	runAt := time.Now()

	if task.queue != nil {
		queue = *task.queue
	}

	if task.priority != nil {
		priority = *task.priority
	}

	if task.runAt != nil {
		runAt = *task.runAt
	}

	values := JobDefaults{
		Queue: queue,
		Priority: priority,
		runAt: runAt,
	}

	method := job.MethodByName("Init")

	if !method.IsValid() {
		log.Println("Init not defined")
		return &values
	}

	mArgs := []reflect.Value{}

	// initialize job configuration
	_ = method.Call(mArgs)

	priorityField := reflect.Indirect(job).FieldByName("Priority")


	if enq.fieldPresent(priorityField, values.Priority) {
		values.Priority = int(priorityField.Int())
	}

	queueField := reflect.Indirect(job).FieldByName("Queue")


	if enq.fieldPresent(queueField, values.Queue) {
		values.Queue = queueField.String()
	}

	return &values
}

func (enc *Enqueuer) fieldPresent(field reflect.Value, value any) bool {
	return field.IsValid() && field.CanConvert(reflect.TypeOf(value))
}
