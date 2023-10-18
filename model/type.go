package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username string             `bson:"username,omitempty" json:"username,omitempty"`
	Password string             `bson:"password,omitempty" json:"password,omitempty"`
}

type Token struct {
	Token_String  string    `bson:"tokenstring,omitempty" json:"tokenstring,omitempty"`
	Expired_Token time.Time `bson:"expiredtoken,omitempty" json:"expiredtoken,omitempty"`
}
