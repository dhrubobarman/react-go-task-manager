package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Body      string             `json:"body"`
	Completed bool               `json:"completed"`
}

var collection *mongo.Collection

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading env file")
	}

	DB_URI := os.Getenv("DB_URI")
	clientOpt := options.Client().ApplyURI(DB_URI)
	client, err := mongo.Connect(context.Background(), clientOpt)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to database")

	collection = client.Database("task-manager-go").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")

	if port == "" {
		port = "4000"
	}

	log.Fatal(app.Listen(":" + port))

}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return c.Status(200).JSON(fiber.Map{"data": todos, "success": true})
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Body is required", "success": false})
	}

	inserResult, err := collection.InsertOne(context.Background(), todo)

	if err != nil {
		return err
	}
	todo.Id = inserResult.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(fiber.Map{"data": todo, "success": true})
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	queryValue := c.Query("success")

	isSuccess := true

	if queryValue == "false" {
		isSuccess = false
	}

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid id", "success": false})
	}
	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"completed": isSuccess}}
	var result bson.M
	err = collection.FindOneAndUpdate(context.Background(), filter, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&result)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "data": result})
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid id", "success": false})
	}
	filter := bson.M{"_id": objectId}
	result := collection.FindOneAndDelete(context.Background(), filter)
	if result.Err() != nil {
		return result.Err()
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "data": result})
}
