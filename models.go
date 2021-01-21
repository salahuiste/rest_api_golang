package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Create Structs
type Friend struct {
	id   int    `json:"id" bson:"id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
}
type User struct {
	id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ID         string             `json:"ID,omitempty" bson:"ID,omitempty"`
	Password   string             `json:"password" bson:"password,omitempty"`
	IsActive   bool               `json:"isActive" bson:"isActive,omitempty"`
	Balance    string             `json:"balance" bson:"balance,omitempty"`
	Age        string             `json:"age" bson:"age,omitempty"`
	Name       string             `json:"name" bson:"name,omitempty"`
	Gender     string             `json:"gender" bson:"gender,omitempty"`
	Company    string             `json:"company" bson:"company,omitempty"`
	Phone      string             `json:"phone" bson:"phone,omitempty"`
	Address    string             `json:"address" bson:"address,omitempty"`
	About      string             `json:"about" bson:"about,omitempty"`
	Registered string             `json:"registered" bson:"registered,omitempty"`
	Latitude   float32            `json:"latitude" bson:"latitude,omitempty"`
	Longitude  float32            `json:"longitude" bson:"longitude,omitempty"`
	Tags       []string           `json:"tags" bson:"tags,omitempty"`
	Friends    []Friend           `json:"friends" bson:"friends,omitempty"`
	data       string             `json:"data" bson:"data,omitempty"`
}
