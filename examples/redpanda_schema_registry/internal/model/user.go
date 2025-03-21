package model

// User Struct
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// AvroSchema Avro Schema
const AvroSchema = `{
	"type": "record",
	"name": "User",
	"fields": [
		{"name": "id", "type": "int"},
		{"name": "name", "type": "string"},
		{"name": "email", "type": "string"}
	]
}`
