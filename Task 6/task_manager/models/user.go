package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserRole string

const (
	RoleAdmin UserRole = "Admin"
	RoleUser  UserRole = "User"
)

func (role UserRole) IsValid() bool {
	switch role {
	case RoleAdmin, RoleUser:
		return true
	}
	return false
}

type User struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"-" bson:"password"`
	Role     UserRole           `json:"role" bson:"role"`
}

type UserRegisterLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
