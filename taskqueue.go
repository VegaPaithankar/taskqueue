package taskqueue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

const capacity = 10

type TaskQueue struct {
	taskCh   chan Task
	taskDone chan int64
	queue    []Task
	running  map[int64]Task
	mu       sync.Mutex
	nextId   int64
	done     bool
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		taskCh:   make(chan Task, capacity),
		taskDone: make(chan int64),
		running:  make(map[int64]Task),
		queue:    make([]Task, 0),
	}
}

func (t *TaskQueue) AddTaskRequest(taskType string, payload interface{}) (int64, error) {
	var task Task
	ctx := context.Background()
	switch taskType {
	case ShortTaskType:
		task = &ShortTask{
			ctx:     ctx,
			payload: payload.(string),
		}
	case LongTaskType:
		task = &LongTask{
			ctx:     ctx,
			payload: payload.(string),
		}
	default:
		return 0, fmt.Errorf("unknown task type: %s", taskType)
	}
	task.SetId(t.nextId)
	t.nextId++
	t.AddTask(task)
	return task.Id(), nil
}

func (t *TaskQueue) AddTask(task Task) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.running) < capacity {
		t.running[task.Id()] = task
		t.taskCh <- task
	} else {
		t.queue = append(t.queue, task)
	}
}

func (t *TaskQueue) AbortTask(id int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.running[id]; ok {
		// send signal to task to abort
		t.running[id].Abort()
	} else {
		return fmt.Errorf("task %d not found", id)
	}
	return nil
}

func (t *TaskQueue) Run() {
	wg := sync.WaitGroup{}
	for i := 0; i < capacity; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(i, t.taskCh, t.taskDone)
		}()
	}

	for !t.done {
		select {
		case id := <-t.taskDone:
			log.Printf("task %d (%s) done", id, t.running[id].toString())
			t.mu.Lock()
			delete(t.running, id)
			// dispatch next task
			if len(t.queue) > 0 {
				task := t.queue[0]
				t.queue = t.queue[1:]
				t.running[task.Id()] = task
				t.taskCh <- task
			}
			t.mu.Unlock()
		}
	}
	wg.Wait()
}

func (t *TaskQueue) Done() {
	for {
		// check if all tasks are done, if so, exit. Otherwise, sleep for 1 second
		n := t.Running()
		log.Printf("waiting for %d tasks to finish\n", n)
		if t.Running() == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	t.done = true
	close(t.taskCh)
}

func (t *TaskQueue) Running() int {
	return len(t.running)
}

func worker(workerId int, taskch chan Task, done chan int64) {
	for {
		select {
		case task, more := <-taskch:
			if !more {
				log.Printf("worker %d exiting\n", workerId)
				return
			}
			task.Run()
			done <- task.Id()
		}
	}
}
