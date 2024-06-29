package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("hello world")
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{
		{Id: 1, Body: "Hello", Completed: true},
		{Id: 2, Body: "World", Completed: false},
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Hello World!",
		})
	})

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"data":    &todos,
		})
	})

	// Create a todo
	app.Post("api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Body is required",
			})
		}
		todo.Id = len(todos) + 1
		todos = append(todos, *todo)
		return c.Status(201).JSON(fiber.Map{
			"success": true,
			"data":    &todo,
		})
	})

	// Update a todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos[i].Completed = !todos[i].Completed
				return c.Status(200).JSON(fiber.Map{
					"success": true,
					"data":    &todos[i],
				})
			}
		}
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Todo not found",
		})
	})

	// Delete a todo
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{
					"success": true,
					"data":    &todo,
				})
			}
		}
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Todo not found",
		})
	})

	log.Fatal(app.Listen(":" + PORT))
}

type Todo struct {
	Id        int    `json:"id"`
	Body      string `json:"body"`
	Completed bool   `json:"completed"`
}
