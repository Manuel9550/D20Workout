package service

import "time"

// Specific requests for certain endpoints

type PointsRequest struct {
	UserName  string    `json:"username"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
