package activities

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/model"
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/pkg"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type RequestError struct{}

func (m *RequestError) Error() string {
	return "Request"
}

func AntiFraudChecks(ctx context.Context, req model.BaseRequest) (model.AntiFraudResponse, error) {
	return httpActivity[model.AntiFraudResponse](ctx, viper.GetString("anti_fraud_service_url")+"/checks", req)
}

func CreateUser(ctx context.Context, req model.BaseRequest) (model.UserResponse, error) {
	return httpActivity[model.UserResponse](ctx, viper.GetString("user_service_url")+"/users", req)
}

func CreateAccount(ctx context.Context, req model.BaseRequest) (model.AccountResponse, error) {
	return httpActivity[model.AccountResponse](ctx, viper.GetString("account_service_url")+"/accounts", req)
}

func CreateAgreement(ctx context.Context, req model.BaseRequest) (model.AgreementResponse, error) {
	return httpActivity[model.AgreementResponse](ctx, viper.GetString("agreement_service_url")+"/agreements", req)
}

func ValidateSignature(ctx context.Context, req model.BaseRequest) (model.SignatureResponse, error) {
	return httpActivity[model.SignatureResponse](ctx, viper.GetString("signature_service_url")+"/signatures", req)
}

func CreateCard(ctx context.Context, req model.BaseRequest) (model.CardResponse, error) {
	cardResponse, err := httpActivity[model.CardResponse](ctx, viper.GetString("card_service_url")+"/cards", req)
	if err == nil {
		cardResponse.Number = pkg.MaskPan(cardResponse.Number)
	}
	return cardResponse, err
}

var client = resty.New().SetDebug(true)

func makeRequest[T any](url string, req model.BaseRequest, res *T) error {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeaders(req.Headers).
		SetBody(req.Body).
		SetResult(res).
		Post(url)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return &RequestError{}
	}

	return nil
}

func httpActivity[T any](ctx context.Context, url string, req model.BaseRequest) (T, error) {
	var res T
	err := makeRequest(url, req, &res)
	return res, err
}
