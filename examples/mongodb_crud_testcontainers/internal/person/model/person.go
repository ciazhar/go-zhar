package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Person struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" query:"id" json:"id"`
	Name    string             `bson:"name" query:"name" json:"name"`
	Email   string             `bson:"email" query:"email" json:"email"`
	Age     int                `bson:"age" query:"age" json:"age"`
	Address Address            `bson:"address" query:"address" json:"address"`
	Hobbies []string           `bson:"hobbies" query:"hobbies" json:"hobbies"`
}

type Address struct {
	Street  string `bson:"street" json:"street"`
	City    string `bson:"city" json:"city"`
	Country string `bson:"country" json:"country"`
}

type PersonMinified struct {
	Id      string `bson:"_id" json:"id"`
	Name    string `bson:"name" json:"name"`
	Country string `bson:"country" json:"country"`
}

type UpdatePersonForm struct {
	Name    string   `json:"name,omitempty"`
	Email   string   `json:"email,omitempty"`
	Age     int      `json:"age,omitempty"`
	Address Address  `json:"address,omitempty"`
	Hobbies []string `json:"hobbies,omitempty"`
}
