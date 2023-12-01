package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UID             string             `bson:"uid,omitempty" json:"uid,omitempty"`
	Email           string             `bson:"email,omitempty" json:"email,omitempty"`
	Username        string             `bson:"username,omitempty" json:"username,omitempty"`
	Password        string             `bson:"password,omitempty" json:"password,omitempty"`
	ConfirmPassword string             `bson:"confirmpassword,omitempty" json:"confirmpassword,omitempty"`
	Role            string             `bson:"role,omitempty" json:"role,omitempty"`
}

type Credential struct {
	Status  bool   `bson:"status" json:"status"`
	Token   string `bson:"token,omitempty" json:"token,omitempty"`
	Message string `bson:"message,omitempty" json:"message,omitempty"`
	Data    []User `bson:"data,omitempty" json:"data,omitempty"`
}

type Todo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Deadline    string             `bson:"deadline,omitempty" json:"deadline,omitempty"`
	TimeStamp   TimeStamp          `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	IsDone      bool               `bson:"isdone,omitempty" json:"isdone,omitempty"`
	User        User               `bson:"user,omitempty" json:"user,omitempty"`
}

type TimeStamp struct {
	CreatedAt time.Time `bson:"createdat,omitempty" json:"createdat,omitempty"`
	UpdatedAt time.Time `bson:"updatedat,omitempty" json:"updatedat,omitempty"`
}

type TodoResponse struct {
	Status  bool   `bson:"status" json:"status"`
	Message string `bson:"message,omitempty" json:"message,omitempty"`
	Data    []Todo `bson:"data,omitempty" json:"data,omitempty"`
}
