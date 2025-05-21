package model

type Product struct {
	ID    int32   `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	Image string  `json:"image"`
}
