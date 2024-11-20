package model

// TripPlannerWorkflowRequest holds all the input data for the TripPlannerWorkflow
type TripPlannerWorkflowRequest struct {
	Destination     string
	Date            string
	UserPreferences UserPreferences
}

// TripPlannerWorkflowResponse holds the output data for the TripPlannerWorkflow
type TripPlannerWorkflowResponse struct {
	Itinerary string
}
