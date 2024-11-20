package main

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/temporal_trip_planner/internal/model"
	"github.com/ciazhar/go-start-small/examples/temporal_trip_planner/internal/workflows"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {
	// Create Temporal client
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	// Define workflow inputs
	destination := "Los Angeles"
	date := "2024-12-25"
	userPreferences := model.UserPreferences{
		Budget:   1000,
		Interest: []string{"parks", "historical sites"},
	}

	// Start the workflow
	options := client.StartWorkflowOptions{
		TaskQueue: model.TripPlannerWorkflowQueue,
	}
	we, err := c.ExecuteWorkflow(context.Background(), options, workflows.TripPlannerWorkflow, destination, date, userPreferences)
	if err != nil {
		log.Fatalln("Unable to execute workflow:", err)
	}

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalf("Error getting workflow result: %v", err)
	}

	// Output workflow results
	log.Printf("Started workflow with ID %s and Run ID %s\n", we.GetID(), we.GetRunID())
	log.Printf("Workflow result: %s", result)
}
