package entities

import "time"

type ExercisePoints struct {
	Username string  `json:"username,omitempty"`
	Points   []Point `json:"points"`
}

type Point struct {
	Timestamp      time.Time `json:"timestamp"`
	Username       string    `json:"username,omitempty"`
	ExerciseNumber int       `json:"exerciseNumber"`
	Amount         int       `json:"amount"`
}
