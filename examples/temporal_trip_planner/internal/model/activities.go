package model

type FetchFlightInfoRequest struct {
	Destination string
	Date        string
}

type FetchFlightInfoResponse struct {
	From  string
	To    string
	Date  string
	Price int
}

type FetchPlacesInfoRequest struct {
	Destination string
}

type FetchPlacesInfoResponse struct {
	Name   string
	Rating float64
	Type   string
}

type FetchReviewsAndActivitiesRequest struct {
	Destination string
}

type FetchReviewsAndActivitiesResponse struct {
	Activity string
	Reviews  float64
}

type GenerateItineraryWithChatGPTRequest struct {
	UserPreferences UserPreferences
	FlightInfo      FetchFlightInfoResponse
	Places          []FetchPlacesInfoResponse
	Activities      []FetchReviewsAndActivitiesResponse
}

type UserPreferences struct {
	Budget   int
	Interest []string
}

type GenerateItineraryWithChatGPTResponse struct {
	Itinerary string
}
