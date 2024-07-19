package main

import (
	"gorm.io/gorm"
)

var db *gorm.DB

// User model
type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Habits   []Habit
}

// Habit model
type Habit struct {
	gorm.Model
	UserID          uint
	Title           string
	IsComplete      bool
	LevelOfComplete int
}
