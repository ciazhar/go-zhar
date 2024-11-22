package model

type OnboardingState string

const (
	ProcessingState              OnboardingState = "processing"
	FailedState                  OnboardingState = "failed"
	FraudNotPassedState          OnboardingState = "fraud_not_passed"
	SignatureNotValidState       OnboardingState = "signature_not_valid"
	WaitingForAgreementSignState OnboardingState = "waiting_for_agreement_signature"
	CompletedState               OnboardingState = "completed"
)
