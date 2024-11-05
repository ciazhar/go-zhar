package model

type BloodAvailability struct {
	ApplicationID string `json:"application_id"`
	ID            int    `json:"id"`
	A             int    `json:"a"`
	B             int    `json:"b"`
	AB            int    `json:"ab"`
	O             int    `json:"o"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
