package main

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/temporal_trip_planner/internal/model"
	"github.com/ciazhar/go-start-small/examples/temporal_trip_planner/internal/workflows"
	"github.com/google/uuid"
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
	input := model.TripPlannerWorkflowRequest{
		Destination:     destination,
		Date:            date,
		UserPreferences: userPreferences,
	}

	workflowId := uuid.New().String()

	// Start the workflow
	options := client.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: model.TripPlannerWorkflowQueue,
	}
	result, err := ExecuteWorkflowWithResult[
		model.TripPlannerWorkflowRequest,
		model.TripPlannerWorkflowResponse,
	](context.Background(), c, options, workflows.TripPlannerWorkflow, input)
	if err != nil {
		log.Fatalln("Unable to execute workflow:", err)
	}

	// Output workflow results
	//log.Printf("Started workflow with ID %s and Run ID %s\n", we.GetID(), we.GetRunID())
	log.Printf("Workflow result: %s", result.Itinerary)
}

// ExecuteWorkflowWithResult executes a workflow with a unified input struct and returns the result
func ExecuteWorkflowWithResult[I any, R any](
	ctx context.Context,
	c client.Client,
	options client.StartWorkflowOptions,
	workflowFunc interface{},
	input I,
) (R, error) {
	we, err := c.ExecuteWorkflow(ctx, options, workflowFunc, input)
	if err != nil {
		return *new(R), fmt.Errorf("failed to start workflow: %w", err)
	}

	var result R
	err = we.Get(ctx, &result)
	if err != nil {
		return *new(R), fmt.Errorf("failed to get workflow result: %w", err)
	}
	return result, nil
}
