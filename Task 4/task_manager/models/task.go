package models

import "time"

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
	Id          int        `json:"id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     time.Time  `json:"duedate" binding:"required"`
	Status      TaskStatus `json:"status" binding:"required"`
}
