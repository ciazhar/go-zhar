package service

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/model"
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/workflows"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"log/slog"
)

type OnboardingService struct {
	temporalClient client.Client
}

func NewOnboardingService(temporalClient client.Client) *OnboardingService {
	return &OnboardingService{temporalClient: temporalClient}
}

func (s *OnboardingService) CreateOnboarding(ctx context.Context, req *model.OnboardingRequest) (model.OnboardingStatusResponse, error) {
	workflowId := uuid.New()

	options := client.StartWorkflowOptions{
		ID:        workflowId.String(),
		TaskQueue: model.OnboardingTaskQueue,
	}

	_, err := s.temporalClient.ExecuteWorkflow(ctx, options, workflows.Onboarding, model.OnboardingRequest{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		City:        req.City,
		AccountType: req.AccountType,
		Currency:    req.Currency,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Unable to start the Workflow:", slog.String("Error", err.Error()))

		return model.OnboardingStatusResponse{
			ID:    workflowId,
			State: model.FailedState,
		}, nil

	}

	return model.OnboardingStatusResponse{
		ID:    workflowId,
		State: model.ProcessingState,
	}, nil
}

func (s *OnboardingService) GetOnboarding(ctx context.Context, id uuid.UUID) (model.OnboardingStatusResponse, error) {
	currentState, err := s.getCurrentState(ctx, id)

	if err != nil {
		return model.OnboardingStatusResponse{}, err
	}

	return model.OnboardingStatusResponse{
		ID:    id,
		State: currentState.State,
		Data:  currentState.Data,
	}, nil
}

func (s *OnboardingService) SignAgreement(ctx context.Context, id uuid.UUID, req *model.SignatureRequest) (model.OnboardingStatusResponse, error) {
	currentState, err := s.getCurrentState(ctx, id)

	if err != nil {
		return model.OnboardingStatusResponse{}, err
	}

	if currentState.State != model.WaitingForAgreementSignState {
		slog.WarnContext(ctx, "Invalid state for agreement signature")
		return model.OnboardingStatusResponse{
			ID:    id,
			State: currentState.State,
			Data:  currentState.Data,
		}, nil
	}

	signatureSignal := model.Signature{
		Signature: req.Signature,
	}

	err = s.temporalClient.SignalWorkflow(ctx, id.String(), "", model.SignatureSignal, signatureSignal)
	if err != nil {
		return model.OnboardingStatusResponse{}, err
	}

	return model.OnboardingStatusResponse{
		ID:    id,
		State: model.ProcessingState,
	}, nil
}

func (s *OnboardingService) getCurrentState(ctx context.Context, id uuid.UUID) (model.OnboardingResponse, error) {
	response, err := s.temporalClient.QueryWorkflow(ctx, id.String(), "", model.CurrentStateQuery)
	if err != nil {
		return model.OnboardingResponse{}, err
	}

	var currentState model.OnboardingResponse
	err = response.Get(&currentState)
	if err != nil {
		return model.OnboardingResponse{}, err
	}

	return currentState, nil
}
