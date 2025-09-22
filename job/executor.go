package job

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/sort"
	"github.com/go-rel/rel/where"
)

// initial delay for the first retry (in seconds)
const INITIAL_RETRY_DELAY float64 = 10

// exponential growth of retry delays (binary exponential backoff algorithm)
const EXPONENTIAL_BACKOFF float64 = 2

type JobExecutor struct {
	Env        string
	DB         rel.Repository
	Cancel     context.CancelFunc
	JobSupport IJobSupport
}


func (ex *JobExecutor) StartAndProcess() error {

	launcher := Launcher{
		Env: ex.Env,
	}

	err := launcher.Start()

	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	ex.Cancel = cancel

	go ex.TimerTask(ctx)

	return  nil
}

func (ex *JobExecutor) Shutdown(ctx context.Context) error {

	ex.Cancel()

	ctx.Done()

	return nil
}

// Runs some task every 1 second.
// If canceled goroutine should be stopped.
func (ex *JobExecutor) TimerTask(ctx context.Context) {
  // Create a new ticker with a period of 2 second.
  ticker := time.NewTicker(2 * time.Second)

  for {
    select {
    case <-ticker.C:
      ex.ProcessJob(ctx)
    case <-ctx.Done():
      fmt.Println("stopping job processing")
      ticker.Stop()
      return
    }
  }
}

func (ex *JobExecutor) ProcessJob(ctx context.Context) {
	// fmt.Println("performing job")

	jobs := []QueuedJob{}

	maxAttempts := 15
	rightNow    := time.Now()
	maxPerRun   := 2

	cond := where.Lte("attempts", maxAttempts).And(where.Lte("run_at", rightNow))

	err := ex.DB.FindAll(
		ctx, &jobs, rel.Where(cond).Limit(maxPerRun), sort.Asc("priority"), sort.Asc("id"),
	)

	if err != nil {
		fmt.Println("Error fetching jobs:", err)
		return
	}

	for _, job := range jobs {
		failure := ex.executeJob(ctx, job)
		if failure { break }
	}
}

func (ex *JobExecutor) executeJob(ctx context.Context, job QueuedJob) bool {
	var err error

	rightNow := time.Now()

	err = ex.lockJob(ctx, rightNow, job)

	if err != nil {
		fmt.Println("Error locking jobs:", err)
		return true
	}

	err = ex.perform(job)

	if err != nil {
		fmt.Println("Error performing jobs:", err)
		err = ex.unlockJob(ctx, rightNow, job, err)

		if err != nil {
			fmt.Println("Error unlocking jobs:", err)
		}

		return true
	}

	err = ex.DB.Delete(ctx, &job)

	if err != nil {
		fmt.Println("Error deleting job:", err)
		return true
	}

	return false
}

func (ex *JobExecutor) perform(job QueuedJob) error {
	var err error

	jobName := job.JobClass

	args := Args{}

	err = json.Unmarshal([]byte(job.Args), &args)

	if err != nil {
		return err
	}

	_, err = ex.JobSupport.PerformNow(jobName, args)

	return err
}

func (ex *JobExecutor) lockJob(ctx context.Context, rightNow time.Time, job QueuedJob) error {
	var err error

	lockedAt := &rightNow
	
	// changed := rel.NewChangeset(&job)

	job.LockedBy = strconv.Itoa(os.Getpid())
	job.LockedAt = lockedAt
	job.State    = ex.detectState(job.Attempts)


	// err = ex.DB.Update(ctx, &job, changed)
	err = ex.DB.Update(ctx, &job)

	return err
}

func (ex *JobExecutor) unlockJob(ctx context.Context, rightNow time.Time, job QueuedJob, failure error) error {
	var err error

	failedAt := &rightNow
	attempts := job.Attempts + 1

	// changed := rel.NewChangeset(&job)

	job.LockedBy  = ""
	job.LockedAt  = nil
	job.State     = "retried"
	job.LastError = failure.Error()
	job.Attempts  = attempts
	job.RunAt     = ex.rescheduleAt(rightNow, attempts)
	job.FailedAt  = failedAt

	// err = ex.DB.Update(ctx, &job, changed)
	err = ex.DB.Update(ctx, &job)

	return err
}

func (ex *JobExecutor) rescheduleAt(right_now time.Time, attempts int) time.Time {
	first_delay := INITIAL_RETRY_DELAY
	exp_growth  := EXPONENTIAL_BACKOFF

	// period = first_delay * (exp_growth**(attempts - 1))
	period := first_delay * math.Pow(exp_growth, float64(attempts) - 1)

	duration := time.Duration(period)

	return right_now.Add(duration * time.Second)
}

func (ex *JobExecutor) detectState(attempts int) string {
	if attempts < 1 {
		return "running"
	} else {
		return "retrying"
	}
}
