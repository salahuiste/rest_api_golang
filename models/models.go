package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Create Struct
type Friend struct {
	ID   int32  `json:"id" bson:"id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
}
type User struct {
	_id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ID         string             `json:"id" bson:"id,omitempty"`
	Password   string             `json:"password" bson:"password,omitempty"`
	IsActive   bool               `json:"isActive" bson:"isActive,omitempty"`
	Balance    string             `json:"balance" bson:"balance,omitempty"`
	Age        int32              `json:"age" bson:"age,omitempty"`
	Name       string             `json:"name" bson:"name,omitempty"`
	Gender     string             `json:"gender" bson:"gender,omitempty"`
	Company    string             `json:"company" bson:"company,omitempty"`
	Email      string             `json:"email" bson:"email,omitempty"`
	Phone      string             `json:"phone" bson:"phone,omitempty"`
	Address    string             `json:"address" bson:"address,omitempty"`
	About      string             `json:"about" bson:"about,omitempty"`
	Registered string             `json:"registered" bson:"registered,omitempty"`
	Latitude   float64            `json:"latitude" bson:"latitude,omitempty"`
	Longitude  float64            `json:"longitude" bson:"longitude,omitempty"`
	Tags       []string           `json:"tags" bson:"tags,omitempty"`
	Friends    []Friend           `json:"friends" bson:"friends,omitempty"`
	Data       string             `json:"data" bson:"data,omitempty"`
}

//login
type Login struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}
