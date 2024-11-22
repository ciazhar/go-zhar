package main

import (
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/activities"
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/model"
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/workflows"
	"github.com/ciazhar/go-start-small/pkg/config"
	"log/slog"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	config.InitConfig(
		config.Config{
			Source: "file",
			Type:   "json",
			File: config.FileConfig{
				FileName: "config.json",
				FilePath: "configs",
			},
		})

	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		slog.Error("Unable to create Temporal client:", slog.String("error", err.Error()))
		return
	}
	defer temporalClient.Close()

	w := worker.New(temporalClient, model.OnboardingTaskQueue, worker.Options{})

	w.RegisterWorkflow(workflows.Onboarding)
	w.RegisterActivity(activities.AntiFraudChecks)
	w.RegisterActivity(activities.CreateUser)
	w.RegisterActivity(activities.CreateAccount)
	w.RegisterActivity(activities.CreateAgreement)
	w.RegisterActivity(activities.ValidateSignature)
	w.RegisterActivity(activities.CreateCard)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		slog.Error("Unable to start worker", slog.String("error", err.Error()))
		return
	}
}
