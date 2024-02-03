package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UID             string             `bson:"uid,omitempty" json:"uid,omitempty"`
	Email           string             `bson:"email,omitempty" json:"email,omitempty"`
	Phonenumber     string             `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
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
	Time        string             `bson:"time,omitempty" json:"time,omitempty"`
	File        GridFSFile         `bson:"file,omitempty" json:"file,omitempty"`
	Tags        Categories         `bson:"tags,omitempty" json:"tags,omitempty"`
	TimeStamps  TimeStamps         `bson:"timestamps,omitempty" json:"timestamps,omitempty"`
	User        User               `bson:"user,omitempty" json:"user,omitempty"`
}

type GridFSFile struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FileName    string             `bson:"filename,omitempty" json:"filename,omitempty"`
	FileSize    int64              `bson:"filesize,omitempty" json:"filesize,omitempty"`
	ContentType string             `bson:"contenttype,omitempty" json:"contenttype,omitempty"`
}

type TimeStamps struct {
	CreatedAt int64 `bson:"createdat,omitempty" json:"createdat,omitempty"`
	UpdatedAt int64 `bson:"updatedat,omitempty" json:"updatedat,omitempty"`
}

type Categories struct {
	Category string `bson:"category,omitempty" json:"category,omitempty"`
}

type TodoResponse struct {
	Status   bool         `bson:"status" json:"status"`
	Message  string       `bson:"message,omitempty" json:"message,omitempty"`
	Data     []Todo       `bson:"data,omitempty" json:"data,omitempty"`
	DataTags []Categories `bson:"datatags,omitempty" json:"datatags,omitempty"`
}

type TodoClear struct {
	IsDone    bool  `bson:"isdone,omitempty" json:"isdone,omitempty"`
	TimeClear int64 `bson:"timeclear,omitempty" json:"timeclear,omitempty"`
	Todo      Todo  `bson:"todo,omitempty" json:"todo,omitempty"`
}

type TodoClearResponse struct {
	Status  bool        `bson:"status" json:"status"`
	Message string      `bson:"message,omitempty" json:"message,omitempty"`
	Data    []TodoClear `bson:"data,omitempty" json:"data,omitempty"`
}

type LogTodo struct {
	TimeStamp int64                    `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	Action    string                   `bson:"action,omitempty" json:"action,omitempty"`
	TodoID    string                   `bson:"todoid,omitempty" json:"todoid,omitempty"`
	UserID    string                   `bson:"userid,omitempty" json:"userid,omitempty"`
	Change    []map[string]interface{} `bson:"change,omitempty" json:"change,omitempty"`
}

type LogTodoResponse struct {
	Status  bool      `bson:"status" json:"status"`
	Message string    `bson:"message,omitempty" json:"message,omitempty"`
	Data    []LogTodo `bson:"data,omitempty" json:"data,omitempty"`
}

type LogUser struct {
	TimeStamp int64                    `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	Action    string                   `bson:"action,omitempty" json:"action,omitempty"`
	ID        string                   `bson:"id,omitempty" json:"id,omitempty"`
	UserID    string                   `bson:"userid,omitempty" json:"userid,omitempty"`
	Change    []map[string]interface{} `bson:"change,omitempty" json:"change,omitempty"`
}
