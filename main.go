package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	tempVars := "8080"
	engine := html.NewFileSystem(http.Dir("./views"), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/dashboard", HandleGettingUserDashboard)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Post("/update-datapoint", func(c *fiber.Ctx) error {
		fmt.Println(string(c.Body()))
		// make the DB change with the user id, the datapoint id, and the data for the slide value

		return nil
	})

	app.Static("/static", "./static")

	log.Fatal(app.Listen(":" + tempVars))
}

// middleware to authenticate, then use htmx to then show dashboard data

func HandleGettingUserDashboard(c *fiber.Ctx) error {
	// get all user and dashboard info.
	type DataPoint struct {
		Id              int
		Title           string
		IsComplete      bool
		LevelOfComplete int
	}

	type DataList struct {
		Data []DataPoint
	}

	type User struct {
		Id   int
		Name string
		Data DataList
	}

	var user User
	var userData DataList
	userData.Data = []DataPoint{}

	var dummyData1 = DataPoint{Title: "Going to the gym", IsComplete: true, LevelOfComplete: 100, Id: 0}
	var dummyData2 = DataPoint{Title: "Coding for more than 4 hours a day", IsComplete: false, LevelOfComplete: 100, Id: 1}
	var dummyData3 = DataPoint{Title: "Take a walk", IsComplete: false, LevelOfComplete: 50, Id: 2}
	userData.Data = append(userData.Data, dummyData1, dummyData2, dummyData3)

	user.Name = "Billy"
	user.Id = 3999
	user.Data = userData

	fmt.Println(user)

	return c.Render("dashboard", fiber.Map{"User": user, "Data": user.Data.Data})
}
