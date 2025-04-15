package models

import (
	"time"
)

type Users struct {
	ID            string    `json:"id" bson:"id" validate:"required"`
	First_Name    string    `json:"first_name" bson:"first_name" validate:"required"`
	Last_name     string    `json:"last_name" bson:"last_name"  validate:"required"`
	Email         string    `json:"email" bson:"email" validate:"required"`
	Password      string    `json:"password" bson:"password" validate:"required"`
	Token         string    `bson:"token" json:"token" validate:"required"`
	Refresh_Token string    `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	Created_time  time.Time `json:"created_time,omitempty" bson:"created_time,omitempty"`
	Updated_time  time.Time `json:"updated_time,omitempty" bson:"updated_time,omitempty"`
}
