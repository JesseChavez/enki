package job

import "time"

type IJobSupport interface {
	Register(queue string, job ActiveJob)
	PerformNow(job string, args Args) (string, error)
	PerformLater(job string, args Args) (string, error)
	Wait(duration time.Duration) *Task
	WaitUntil(runAt time.Time) *Task
}

type Task struct {
	enqueuer *Enqueuer
	name     string
	args     Args
	queue    *string
	priority *int
	runAt    *time.Time
}

// Wait enqueues the job with the specified delay
func (enq *Enqueuer) Wait(duration time.Duration) *Task {
	runAt := time.Now().Add(duration)

	task := Task{
		enqueuer: enq,
		runAt: &runAt,
	}

	return &task
}

// WaitUntil enqueues the job at the time specified
func (enq *Enqueuer) WaitUntil(runAt time.Time) *Task {
	task := Task{
		enqueuer: enq,
		runAt: &runAt,
	}

	return &task
}

// Priority enqueues the job with the specified priority
func (enq *Enqueuer) Priority(priority int) *Task {
	task := Task{
		enqueuer: enq,
		priority: &priority,
	}

	return &task
}

// Queue enqueues the job on the specified queue
func (enq *Enqueuer) Queue(queue string) *Task {
	task := Task{
		enqueuer: enq,
		queue: &queue,
	}

	return &task
}

func (enq *Enqueuer) PerformLater(jobName string, args Args) (string, error) {
	task := Task{
		name: jobName,
		args: args,
	}

	id, err := enq.Enqueue(&task, args)

	return id, err
}

func (task *Task) PerformLater(jobName string, args Args) (string, error) {
	task.name = jobName
	task.args = args

	id, err := task.enqueuer.Enqueue(task, args)

	return id, err
}
