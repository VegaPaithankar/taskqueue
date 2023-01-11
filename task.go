package taskqueue

import (
	"context"
	"log"
	"time"
)

type Task interface {
	Run()
	SetId(id int64)
	Id() int64
	toString() string
	Abort()
}

//////////////// ShortTask
const ShortTaskType = "shortTask"

type ShortTask struct {
	// task id
	id int64

	// payload
	payload string

	// context
	ctx context.Context

	// cancel function
	cancel context.CancelFunc
}

func (t *ShortTask) toString() string {
	return t.payload
}

// run function
func (t *ShortTask) Run() {
	log.Printf("Starting short task (3 seconds) with id %d", t.id)
	for i := 0; i < 3; i++ {
		select {
		case <-t.ctx.Done():
			log.Printf("task %d aborted", t.id)
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (t *ShortTask) SetId(id int64) {
	t.id = id
}

func (t *ShortTask) Id() int64 {
	return t.id
}

func (t *ShortTask) Abort() {
	// cancel context
	t.cancel()
}

//////////////// LongTask
const LongTaskType = "longTask"

type LongTask struct {
	// task id
	id int64
	// payload
	payload string

	// context
	ctx context.Context

	// cancel function
	cancel context.CancelFunc
}

func (t *LongTask) toString() string {
	return t.payload
}

// run function
func (t *LongTask) Run() {
	log.Printf("Starting long task (10 seconds) with id %d", t.id)
	for i := 0; i < 3; i++ {
		select {
		case <-t.ctx.Done():
			log.Printf("task %d aborted", t.id)
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (t *LongTask) SetId(id int64) {
	t.id = id
}

func (t *LongTask) Id() int64 {
	return t.id
}

func (t *LongTask) Abort() {
	// cancel context
	t.cancel()
}
