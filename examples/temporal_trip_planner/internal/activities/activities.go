package activities

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/temporal_trip_planner/internal/model"
)

// FetchFlightInfo fetches flight details from a simulated API
func FetchFlightInfo(ctx context.Context, request model.FetchFlightInfoRequest) (model.FetchFlightInfoResponse, error) {
	flightData := model.FetchFlightInfoResponse{
		From:  "JFK",
		To:    request.Destination,
		Date:  request.Date,
		Price: 450,
	}
	return flightData, nil
}

// FetchPlacesInfo fetches places to visit from a simulated API
func FetchPlacesInfo(ctx context.Context, request model.FetchPlacesInfoRequest) ([]model.FetchPlacesInfoResponse, error) {
	places := []model.FetchPlacesInfoResponse{
		{
			Name:   "Central Park",
			Rating: 4.8,
			Type:   "park",
		},
		{
			Name:   "Statue of Liberty",
			Rating: 4.7,
			Type:   "monument",
		},
	}
	return places, nil
}

// FetchReviewsAndActivities fetches reviews and activities from a simulated API
func FetchReviewsAndActivities(ctx context.Context, request model.FetchReviewsAndActivitiesRequest) ([]model.FetchReviewsAndActivitiesResponse, error) {
	activities := []model.FetchReviewsAndActivitiesResponse{
		{
			Activity: "Boat Tour",
			Reviews:  4.5,
		},
		{
			Activity: "Museum Visit",
			Reviews:  4.6,
		},
	}
	return activities, nil
}

// GenerateItineraryWithChatGPT generates an itinerary based on input data
func GenerateItineraryWithChatGPT(ctx context.Context, request model.GenerateItineraryWithChatGPTRequest) (model.GenerateItineraryWithChatGPTResponse, error) {
	itinerary := fmt.Sprintf("Day 1: Fly from %s to %s. Visit %s, rated %.1f stars.\n",
		request.FlightInfo.From, request.FlightInfo.To, request.Places[0].Name, request.Places[0].Rating)
	itinerary += fmt.Sprintf("Day 2: Enjoy a %s with %.1f stars reviews.",
		request.Activities[0].Activity, request.Activities[0].Reviews)
	return model.GenerateItineraryWithChatGPTResponse{
		Itinerary: itinerary,
	}, nil
}
