package enki

import (
	"reflect"

	"github.com/JesseChavez/enki/job"
)

func (ek *Enki) InitQueueing() IJobSupport {
	enqueuer := ek.JobSupport

	return enqueuer
}

func (ek *Enki) TypeOf(item any) job.ActiveJob {
	kind := reflect.TypeOf(item)

	return kind
}
