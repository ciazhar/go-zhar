package workflows

import (
	"github.com/ciazhar/go-start-small/examples/temporal_trip_planner/internal/activities"
	"github.com/ciazhar/go-start-small/examples/temporal_trip_planner/internal/model"
	"github.com/ciazhar/go-start-small/examples/temporal_trip_planner/pkg/temporal"
	"go.temporal.io/sdk/workflow"
)

// TripPlannerWorkflow orchestrates fetching data and generating an itinerary
func TripPlannerWorkflow(ctx workflow.Context, destination string, date string, userPreferences model.UserPreferences) (string, error) {
	ao := temporal.GetDefaultActivityOptions()
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Step 1: Fetch flight info
	flightInfo, err := ExecuteActivity[
		model.FetchFlightInfoRequest,
		model.FetchFlightInfoResponse,
	](ctx, activities.FetchFlightInfo, model.FetchFlightInfoRequest{
		Destination: destination,
		Date:        date,
	})
	if err != nil {
		return "", err
	}

	// Step 2: Fetch places to visit
	places, err := ExecuteActivity[
		model.FetchPlacesInfoRequest,
		[]model.FetchPlacesInfoResponse,
	](ctx, activities.FetchPlacesInfo, model.FetchPlacesInfoRequest{
		Destination: destination,
	})
	if err != nil {
		return "", err
	}

	// Step 3: Fetch reviews and activities
	activitiesList, err := ExecuteActivity[
		model.FetchReviewsAndActivitiesRequest,
		[]model.FetchReviewsAndActivitiesResponse,
	](ctx, activities.FetchReviewsAndActivities, model.FetchReviewsAndActivitiesRequest{
		Destination: destination,
	})
	if err != nil {
		return "", err
	}

	// Step 4: Generate itinerary using ChatGPT
	itinerary, err := ExecuteActivity[
		model.GenerateItineraryWithChatGPTRequest,
		model.GenerateItineraryWithChatGPTResponse,
	](ctx, activities.GenerateItineraryWithChatGPT, model.GenerateItineraryWithChatGPTRequest{
		UserPreferences: userPreferences,
		FlightInfo:      flightInfo,
		Places:          places,
		Activities:      activitiesList,
	})
	if err != nil {
		return "", err
	}

	// Return final itinerary
	return itinerary.Itinerary, nil
}

func ExecuteActivity[I any, R any](ctx workflow.Context, activityFunc interface{}, input I) (R, error) {
	var res R
	err := workflow.ExecuteActivity(ctx, activityFunc, input).Get(ctx, &res)
	return res, err
}
