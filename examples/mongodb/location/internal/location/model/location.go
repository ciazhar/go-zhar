package model

type Location struct {
	Name       string    `json:"name"`
	Coordinate []float64 `json:"coordinate"` // Coordinates in [longitude, latitude] format
	Distance   float64   `json:"distance" bson:"distance"`
}

type InsertLocationForm struct {
	Name       string    `json:"name"`
	Coordinate []float64 `json:"coordinate"` // Coordinates in [longitude, latitude] format
}

type NearestLocationForm struct {
	Longitude   float64 `query:"longitude"`
	Latitude    float64 `query:"latitude"`
	MaxDistance int     `query:"max_distance"`
	Limit       int     `query:"limit"`
}
