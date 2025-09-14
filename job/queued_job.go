package job

import (
	"context"
	"time"
)

// == Schema Information
//
// Table name: queued_jobs
//
//  id             :bigint           not null, primary key
//  queue          :string           not null
//  handler        :string           not null
//  job_class      :string(100)      not null
//  job_id         :string(36)       not null
//  priority       :integer          default(1), not null
//  attempts       :integer          default(0), not null
//  state          :string(20)       not null
//  run_at         :datetime
//  args           :text(2147483647)
//  last_error     :text(2147483647)
//  last_backtrace :text(2147483647)
//  failed_at      :datetime
//  locked_by      :string
//  locked_at      :datetime
//  on_hold        :boolean          default(FALSE), not null
//  created_at     :datetime         not null
//  updated_at     :datetime         not null
//

type QueuedJob struct {
	ID                int
	Queue             string
	Handler           string
	JobClass          string
	JobId             string
	Priority          int
	Attempts          int
	RunAt             time.Time
	State             string
	FailedAt          *time.Time
	LastError         string
	LastBacktrace     string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}


type QueuedJobEntity struct {
	Entity []QueuedJob
}

func (list *QueuedJobEntity) Select(ctx context.Context, page string) (err error) {

	return err
}
