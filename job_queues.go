package enki

import (
	"log"
	"reflect"
)


type ActiveJob = reflect.Type

type Queues struct {
	List  map[string]ActiveJob
}


func (q *Queues) Register(queue string, job ActiveJob) {
	q.List = make(map[string]ActiveJob)

	log.Println("====>", job.Name())
	switch queue {
	case "default":
		q.List[job.Name()] = job
	default:
		log.Println("Unknown queue:", queue)
	}
}


func (ek *Enki) InitQueueing() *Queues {
	Queues := &Queues{}

	ek.Queues = Queues

	return Queues
}

func (ek *Enki) TypeOf(item any) ActiveJob {
	kind := reflect.TypeOf(item)

	return kind
}
