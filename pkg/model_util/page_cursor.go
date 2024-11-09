package model_util

type Page struct {
	Data      interface{} `json:"data"`
	TotalData int         `json:"total_data"`
	TotalPage int         `json:"total_page"`
}

type PageCursor struct {
	Data        interface{} `json:"data"`
	TotalData   int         `json:"total_data"`
	CurrentPage int         `json:"current_page"`
	TotalPage   int         `json:"total_page"`
	NextCursor  string      `json:"next_cursor"`
	PrevCursor  string      `json:"prev_cursor"`
}
