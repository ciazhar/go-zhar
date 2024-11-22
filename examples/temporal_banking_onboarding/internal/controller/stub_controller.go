package controller

import (
	"fmt"
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func AntiFraudChecksHandler(c *fiber.Ctx) error {
	var request model.UserRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if strings.Contains(strings.ToLower(request.FirstName), "fraud") ||
		strings.Contains(strings.ToLower(request.LastName), "fraud") {

		return c.Status(http.StatusOK).JSON(model.AntiFraudResponse{
			Passed:  false,
			Comment: "User is fraud",
		})
	}

	return c.Status(http.StatusOK).JSON(model.AntiFraudResponse{
		Passed:  true,
		Comment: "ok",
	})

}

func UsersHandler(c *fiber.Ctx) error {
	var request model.UserRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user := model.UserResponse{
		ID:        uuid.New(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		City:      request.City,
	}

	return c.Status(http.StatusOK).JSON(user)
}

func AgreementsHandler(c *fiber.Ctx) error {
	var request model.AgreementRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id := uuid.New()
	agreement := model.AgreementResponse{
		ID:   id,
		Link: "https://perfect-bank.ua/agreements/" + id.String(),
	}

	return c.Status(http.StatusOK).JSON(agreement)
}

func SignaturesHandler(c *fiber.Ctx) error {
	var request model.SignatureRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if strings.Contains(request.Signature, "fraud") {
		return c.Status(http.StatusOK).JSON(model.SignatureResponse{
			ID:      uuid.New(),
			Valid:   false,
			Comment: "Signature is fraud",
		})
	}

	return c.Status(http.StatusOK).JSON(model.SignatureResponse{
		ID:      uuid.New(),
		Valid:   true,
		Comment: "ok",
	})
}

func AccountsHandler(c *fiber.Ctx) error {
	var request model.AccountRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	account := model.AccountResponse{
		ID:       uuid.New(),
		UserID:   request.UserID,
		Currency: request.Currency,
		Type:     request.Type,
		Iban:     fmt.Sprintf("UA8937040044%016d", rand.Int63n(9999999999999999-1000000000000000)+1000000000000000),
		Balance:  0,
	}

	return c.Status(http.StatusOK).JSON(account)
}

func CardsHandler(c *fiber.Ctx) error {
	var request model.CardRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	expiryDate := time.Now().AddDate(2, 0, 0)
	expiry := fmt.Sprintf("%02d/%02d", expiryDate.Month(), expiryDate.Year()%100)

	card := model.CardResponse{
		ID:        uuid.New(),
		AccountID: request.AccountID,
		Number:    fmt.Sprintf("%016d", rand.Int63n(9999999999999999-1000000000000000)+1000000000000000),
		Expire:    expiry,
		Cvv:       fmt.Sprintf("%03d", rand.Intn(1000)),
	}

	return c.Status(http.StatusOK).JSON(card)
}
