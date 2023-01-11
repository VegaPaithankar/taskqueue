package main

import (
	"fmt"
	"math/rand"

	"github.com/VegaPaithankar/taskqueue"
)

func main() {
	tq := taskqueue.NewTaskQueue()
	go tq.Run()
	for i := 0; i < 30; i++ {
		var taskType string
		// flip a coin
		if rand.Intn(2) == 0 {
			taskType = taskqueue.ShortTaskType
		} else {
			taskType = taskqueue.LongTaskType
		}
		id, err := tq.AddTaskRequest(taskType, fmt.Sprintf("hello from task number %d", i))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Added task %d with type %s\n", id, taskType)
	}
	tq.Done() // stop when all tasks are done
}
