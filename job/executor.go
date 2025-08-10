package job

import (
	"context"
	"fmt"
	"time"
)

type JobExecutor struct {
	Env        string
	Cancel     context.CancelFunc
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
      ex.ProcessJob()
    case <-ctx.Done():
      fmt.Println("stopping job processing")
      ticker.Stop()
      return
    }
  }
}

func (ex *JobExecutor) ProcessJob() {
	fmt.Println("performing job:", "My Job")
}
