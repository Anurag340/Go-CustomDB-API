package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// User struct to define user data
type User struct {
	Name    string `json:"name"`
	Age     string `json:"age"`
	Contact string `json:"contact"`
	Company string `json:"company"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	Zipcode string `json:"zipcode"`
}

// Driver struct to represent the database driver
type Driver struct {
	dir     string
	mutexes map[string]*sync.Mutex
}

// New initializes a new Driver (simulating a database setup)
func New(dir string) (*Driver, error) {
	dir = filepath.Clean(dir)
	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
	}

	// Create directory if it doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf("unable to create database directory: %v", err)
		}
	}

	return &driver, nil
}

// Write simulates writing a User record to a file
func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" || resource == "" {
		return fmt.Errorf("missing collection or resource")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(dir, resource+".json")
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// Read simulates reading a User record from a file
func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" || resource == "" {
		return fmt.Errorf("missing collection or resource")
	}

	filePath := filepath.Join(d.dir, collection, resource+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	return json.Unmarshal(data, v)
}

// Delete simulates deleting a User record from a file
func (d *Driver) Delete(collection, resource string) error {
	filePath := filepath.Join(d.dir, collection, resource+".json")
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("could not delete file: %v", err)
	}
	return nil
}

// getOrCreateMutex returns a mutex for the collection to ensure thread safety
func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutexes[collection] = &sync.Mutex{}
	return d.mutexes[collection]
}

// Controller for creating a user
func CreateUser(db *Driver) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Error parsing body: %v", err))
		}

		// Save to the database
		err := db.Write("users", user.Name, user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error saving user: %v", err))
		}

		return c.Status(fiber.StatusCreated).JSON(user)
	}
}

// Controller for getting a user by name
func GetUser(db *Driver) fiber.Handler {
	return func(c *fiber.Ctx) error {
		name := c.Params("name")

		var user User
		err := db.Read("users", name, &user)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("User not found: %v", err))
		}

		return c.JSON(user)
	}
}

// Controller for getting all users
func GetAllUsers(db *Driver) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dirPath := filepath.Join(db.dir, "users")
		files, err := os.ReadDir(dirPath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error reading users directory: %v", err))
		}

		var users []User
		for _, file := range files {
			if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
				continue
			}

			var user User
			name := file.Name()[:len(file.Name())-len(".json")]
			err := db.Read("users", name, &user)
			if err == nil {
				users = append(users, user)
			}
		}

		return c.JSON(users)
	}
}


// Controller for deleting a user by name
func DeleteUser(db *Driver) fiber.Handler {
	return func(c *fiber.Ctx) error {
		name := c.Params("name")

		err := db.Delete("users", name)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Error deleting user: %v", err))
		}

		return c.SendString(fmt.Sprintf("User '%s' deleted", name))
	}
}

func main() {
	// Initialize app
	app := fiber.New()

	// Set up database
	dir := "./data"
	db, err := New(dir)
	if err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}

	// Define routes
	app.Post("/users", CreateUser(db))
	app.Get("/users/:name", GetUser(db))
	app.Delete("/users/:name", DeleteUser(db))
	app.Get("/users", GetAllUsers(db)) 


	// Start the server
	err = app.Listen(":8001")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
