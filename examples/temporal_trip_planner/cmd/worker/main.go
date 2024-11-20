package main

import (
	"github.com/ciazhar/go-start-small/examples/temporal_trip_planner/internal/activities"
	"github.com/ciazhar/go-start-small/examples/temporal_trip_planner/internal/workflows"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create Temporal client
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	// Register worker
	w := worker.New(c, "TripPlannerTaskQueue", worker.Options{})
	w.RegisterWorkflow(workflows.TripPlannerWorkflow)
	w.RegisterActivity(activities.FetchFlightInfo)
	w.RegisterActivity(activities.FetchPlacesInfo)
	w.RegisterActivity(activities.FetchReviewsAndActivities)
	w.RegisterActivity(activities.GenerateItineraryWithChatGPT)

	// Start worker
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Unable to start worker:", err)
	}
}
