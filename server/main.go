package main

import (
	"log"
	"taskqueue"

	"github.com/gofiber/fiber/v2"
)

func main() {
	tq := taskqueue.NewTaskQueue()
	go tq.Run()

	// Instantiate fiber routes.
	app := fiber.New()

	app.Get("/add-new-task", func(c *fiber.Ctx) error {
		id, err := tq.AddTaskRequest("FIXME", "FIXME")
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("Added task %d with type %s\n", id, "FIXME")
		// FIXME: send back taskid or error.
		return nil
	})

	app.Get("/task-abort", func(c *fiber.Ctx) error {
		var id int64 = -1 // FIXME: get id from request
		err := tq.AbortTask(id)
		if err != nil {
			log.Println(err)
			return err
		}
		// FIXME: send back success code
		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
