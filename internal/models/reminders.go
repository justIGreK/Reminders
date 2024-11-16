package models

import "time"

type Reminder struct {
	ID           string    `bson:"_id,omitempty"`
	UserID       string    `bson:"user_id"`
	Action       string    `bson:"action"`
	Time         time.Time `bson:"utc_time"`
	OriginalTime time.Time `bson:"time"`
	IsActive     bool      `bson:"is_active"`
}

type CreateRms struct{
	UserID string
	Action string
	Time string
	Date *string
}