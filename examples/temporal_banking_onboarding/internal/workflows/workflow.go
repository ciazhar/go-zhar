package workflows

import (
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/activities"
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/model"
	"github.com/ciazhar/go-start-small/examples/temporal_trip_planner/pkg/temporal"
	"go.temporal.io/sdk/workflow"
)

func Onboarding(ctx workflow.Context, request model.OnboardingRequest) (model.OnboardingResponse, error) {
	logger := workflow.GetLogger(ctx)

	workflowId := workflow.GetInfo(ctx).WorkflowExecution.ID

	logger.Info("Starting onboarding workflow", "WorkflowID", workflowId)

	options := temporal.GetDefaultActivityOptions()

	currentState := model.OnboardingResponse{State: model.ProcessingState}

	if err := workflow.SetQueryHandler(ctx, model.CurrentStateQuery, func() (model.OnboardingResponse, error) {
		return currentState, nil
	}); err != nil {
		currentState = model.OnboardingResponse{State: model.FailedState}
		return currentState, err
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	userInput := model.UserRequest{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		City:      request.City,
	}

	// 1. Execute antifraud checks
	logger.Info("Executing antifraud checks for user", "Email", userInput.Email)
	antiFraudChecksResult, err := ExecuteActivity[model.BaseRequest, model.AntiFraudResponse](ctx, activities.AntiFraudChecks, toBaseRequest(workflowId, userInput))

	if err != nil {
		currentState = model.OnboardingResponse{State: model.FailedState}
		return currentState, err
	}

	// If anti fraud checks failed, return fraud_not_passed state
	if !antiFraudChecksResult.Passed {
		logger.Warn("User did not pass antifraud checks", "Email", userInput.Email)
		currentState = model.OnboardingResponse{State: model.FraudNotPassedState}
		return currentState, nil
	}

	// 2. Create user
	logger.Info("Creating user", "Email", userInput.Email)
	createUserResult, err := ExecuteActivity[model.BaseRequest, model.UserResponse](ctx, activities.CreateUser, toBaseRequest(workflowId, userInput))
	if err != nil {
		logger.Error("Unable to create user", "Email", userInput.Email)
		currentState = model.OnboardingResponse{State: model.FailedState}
		return currentState, err
	}

	// 3. Create account
	logger.Info("Creating account for user", "UserID", createUserResult.ID)
	accountInput := model.AccountRequest{
		UserID:   createUserResult.ID,
		Type:     request.AccountType,
		Currency: request.Currency,
	}

	createAccountResult, err := ExecuteActivity[model.BaseRequest, model.AccountResponse](ctx, activities.CreateAccount, toBaseRequest(workflowId, accountInput))
	if err != nil {
		logger.Error("Unable to create account", "UserID", createUserResult.ID)
		currentState = model.OnboardingResponse{State: model.FailedState}
		return currentState, err
	}

	// 4. Create agreement
	logger.Info("Creating agreement for user", "UserID", createUserResult.ID, "AccountID", createAccountResult.ID)
	agreementInput := model.AgreementRequest{
		UserID:    createUserResult.ID,
		AccountID: createAccountResult.ID,
	}

	createAgreementResult, err := ExecuteActivity[model.BaseRequest, model.AgreementResponse](ctx, activities.CreateAgreement, toBaseRequest(workflowId, agreementInput))
	if err != nil {
		currentState = model.OnboardingResponse{State: model.FailedState}
		return currentState, err
	}

	// 5. Wait for signature
	logger.Info("Waiting for agreement signature", "AgreementID", createAgreementResult.ID)
	currentState = model.OnboardingResponse{
		State: model.WaitingForAgreementSignState,
		Data:  map[string]any{"link": createAgreementResult.Link},
	}

	var signatureSignal model.Signature

	signalChan := workflow.GetSignalChannel(ctx, model.SignatureSignal)
	signalChan.Receive(ctx, &signatureSignal)

	// Send signature for validation
	signatureInput := model.SignatureRequest{
		AgreementID: createAgreementResult.ID,
		Signature:   signatureSignal.Signature,
	}

	signatureResult, err := ExecuteActivity[model.BaseRequest, model.SignatureResponse](ctx, activities.ValidateSignature, toBaseRequest(workflowId, signatureInput))
	if err != nil {
		logger.Error("Unable to validate signature", "AgreementID", createAgreementResult.ID)
		currentState = model.OnboardingResponse{State: model.FailedState}
		return currentState, err
	}

	// If signature is not valid, return signature_not_valid state
	if !signatureResult.Valid {
		logger.Warn("Signature is not valid", "AgreementID", createAgreementResult.ID)
		currentState = model.OnboardingResponse{State: model.SignatureNotValidState}
		return currentState, nil
	}

	currentState = model.OnboardingResponse{State: model.ProcessingState}

	// 6. Create card
	logger.Info("Creating card for account", "AccountID", createAccountResult.ID)
	cardInput := model.CardRequest{
		AccountID: createAccountResult.ID,
	}

	createCardResult, err := ExecuteActivity[model.BaseRequest, model.CardResponse](ctx, activities.CreateCard, toBaseRequest(workflowId, cardInput))
	if err != nil {
		logger.Error("Unable to create card", "AccountID", createAccountResult.ID)
		currentState = model.OnboardingResponse{State: model.FailedState}
		return currentState, err
	}

	// 7. Onboarding completed
	logger.Info("Onboarding completed", "UserID", createUserResult.ID, "AccountID", createAccountResult.ID, "CardID", createCardResult.ID)
	currentState = toFinalState(createAccountResult, createCardResult)

	return currentState, nil
}

func ExecuteActivity[I any, R any](ctx workflow.Context, activityFunc interface{}, input I) (R, error) {
	var res R
	err := workflow.ExecuteActivity(ctx, activityFunc, input).Get(ctx, &res)
	return res, err
}

func toBaseRequest(workflowId string, body any) model.BaseRequest {
	return model.BaseRequest{
		Headers: map[string]string{"X-Request-Id": workflowId},
		Body:    body,
	}
}

func toFinalState(accountResult model.AccountResponse, cardResult model.CardResponse) model.OnboardingResponse {
	return model.OnboardingResponse{
		State: model.CompletedState,
		Data: map[string]any{
			"account": accountResult,
			"card":    cardResult,
		},
	}
}
