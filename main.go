package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ElGraam/Todos/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var queries *db.Queries

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection settings
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")

	// Print environment variables to verify
	fmt.Println("DB_USER:", user)
	fmt.Println("DB_PASSWORD:", password)
	fmt.Println("DB_HOST:", host)
	fmt.Println("DB_NAME:", dbname)

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbname)

	// Open database connection
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlDB.Close()

	// Initialize sqlc queries
	queries = db.New(sqlDB)

	// Initialize Fiber app
	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Routes
	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	// Start server
	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}

func getTodos(c *fiber.Ctx) error {
	log.Println("Handling GET /api/todos request")
	todos, err := queries.ListTodos(c.Context())
	if err != nil {
		log.Printf("Error in getTodos: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch todos: %v", err),
		})
	}

	// Convert the todos to a JSON-friendly format
	jsonTodos := make([]map[string]interface{}, len(todos))
	for i, todo := range todos {
		jsonTodo := map[string]interface{}{
			"id":        todo.ID,
			"body":      todo.Body,
			"completed": todo.Completed,
		}
		if !todo.CreatedAt.IsZero() { // Check if CreatedAt is not the zero value
			jsonTodo["created_at"] = todo.CreatedAt
		} else {
			jsonTodo["created_at"] = nil
		}
		jsonTodos[i] = jsonTodo
	}

	log.Printf("Successfully retrieved %d todos", len(todos))
	return c.Status(200).JSON(jsonTodos)
}

func createTodo(c *fiber.Ctx) error {
	log.Println("Received createTodo request")
	var input db.CreateTodoParams
	if err := c.BodyParser(&input); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Cannot parse JSON: %v", err)})
	}
	log.Printf("Parsed input: %+v", input)

	if input.Body == "" {
		log.Println("Empty todo body received")
		return c.Status(400).JSON(fiber.Map{"error": "Body is required"})
	}

	// Set default values
	input.Completed = false

	log.Println("Attempting to create todo in database")
	todo, err := queries.CreateTodo(c.Context(), input)
	if err != nil {
		log.Printf("Error creating todo: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create todo: %v", err)})
	}

	log.Printf("Todo created successfully: %+v", todo)
	return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var updateParams db.UpdateTodoParams
	updateParams.ID = int32(id)

	existingTodo, err := queries.GetTodoByID(c.Context(), updateParams.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch existing todo"})
	}
	updateParams.Completed = !existingTodo.Completed

	if err := c.BodyParser(&updateParams); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	updatedTodo, err := queries.UpdateTodo(c.Context(), updateParams)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update todo"})
	}

	return c.Status(200).JSON(updatedTodo)
}

func deleteTodo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = queries.DeleteTodo(c.Context(), int32(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete todo"})
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}
