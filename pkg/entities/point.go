package entities

import "time"

type ExercisePoints struct {
	Username string  `json:"username,omitempty"`
	Points   []Point `json:"points"`
}

type Point struct {
	Timestamp      time.Time `json:"timestamp,omitempty"`
	Username       string    `json:"username"`
	ExerciseNumber int       `json:"exerciseNumber"`
	Amount         int       `json:"amount"`
}
