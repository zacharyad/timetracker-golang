package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// User struct
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

// Habit struct
type Habit struct {
	ID              int    `json:"id,omitempty"`
	UserID          int    `json:"user_id,omitempty"`
	Title           string `json:"title"`
	IsComplete      bool   `json:"is_complete"`
	LevelOfComplete int    `json:"level_of_complete"`
}

var db *sql.DB

func main() {
	var err error
	// Initialize the database
	db, err = InitDB("./database.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Seed the database
	err = SeedData(db)
	if err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	tempVars := "8080"
	engine := html.NewFileSystem(http.Dir("./views"), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/static", "./static")

	app.Get("/", HandleRenderLandingPage)
	app.Get("/dashboard", HandleGettingUserDashboard)
	app.Get("/habit/:id", HandleGetSingleHabit)
	app.Post("/update-datapoint", HandleCreateHabit)

	log.Fatal(app.Listen(":" + tempVars))
}

func HandleRenderLandingPage(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func HandleGetSingleHabit(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid habit ID"})
	}

	var habit Habit
	err = db.QueryRow("SELECT id, user_id, title, is_complete, level_of_complete FROM habits WHERE id = ?", id).
		Scan(&habit.ID, &habit.UserID, &habit.Title, &habit.IsComplete, &habit.LevelOfComplete)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Habit not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch habit"})
	}

	return c.JSON(habit)
}

func HandleCreateHabit(c *fiber.Ctx) error {
	var habit Habit
	if err := c.BodyParser(&habit); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// For simplicity, we're using a fixed user ID. In a real app, you'd get this from the authenticated user's session.
	const userID = 1

	log.Printf("Received new habit: Title=%s, IsComplete=%v, LevelOfComplete=%d",
		habit.Title, habit.IsComplete, habit.LevelOfComplete)

	result, err := db.Exec("INSERT INTO habits (user_id, title, is_complete, level_of_complete) VALUES (?, ?, ?, ?)",
		userID, habit.Title, habit.IsComplete, habit.LevelOfComplete)
	if err != nil {
		log.Printf("Error executing insert query: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create habit"})
	}

	newID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get new habit ID"})
	}

	habit.ID = int(newID)
	habit.UserID = userID

	log.Printf("Successfully created new habit: ID=%d", habit.ID)

	return c.Status(fiber.StatusCreated).JSON(habit)
}

func HandleGettingUserDashboard(c *fiber.Ctx) error {
	// For simplicity, we'll fetch the first user. In a real app, you'd get the user based on authentication.
	var user User
	err := db.QueryRow("SELECT id, name, email FROM users LIMIT 1").Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return c.Status(404).SendString("User not found")
	}

	rows, err := db.Query("SELECT id, title, is_complete, level_of_complete FROM habits WHERE user_id = ?", user.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var habits []Habit
	for rows.Next() {
		var h Habit
		err := rows.Scan(&h.ID, &h.Title, &h.IsComplete, &h.LevelOfComplete)
		if err != nil {
			return err
		}
		habits = append(habits, h)
	}

	fmt.Println(habits, user)

	return c.Render("dashboard", fiber.Map{
		"User": user,
		"Data": habits,
	})
}
