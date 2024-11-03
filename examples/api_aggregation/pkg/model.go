package pkg

type GetUserProfileResponse struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type GetProductRecommendationsResponse struct {
	ProductId string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
}

type GetUserOrdersResponse struct {
	OrderID string   `json:"order_id"`
	Items   []string `json:"items"`
	Status  string   `json:"status"`
}

type GetDashboardDataResponse struct {
	Orders   []GetUserOrdersResponse             `json:"orders"`
	Products []GetProductRecommendationsResponse `json:"products"`
	User     GetUserProfileResponse              `json:"user"`
}
