package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatus string

const (
	Pending    TaskStatus = "Pending"
	InProgress TaskStatus = "In progress"
	Done       TaskStatus = "Done"
)

func (status TaskStatus) IsValid() bool {
	switch status {
	case Pending, InProgress, Done:
		return true
	}
	return false
}

type Task struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" binding:"required" bson:"title"`
	Description string             `json:"description" bson:"description"`
	DueDate     time.Time          `json:"duedate" binding:"required" bson:"duedate"`
	Status      TaskStatus         `json:"status" binding:"required" bson:"status"`
}
