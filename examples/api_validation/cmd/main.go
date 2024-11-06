package main

import (
	"context"

	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/validator"
	validator2 "github.com/go-playground/validator/v10"
)

type User struct {
	Username string `validate:"required"`
	Tagline  string `validate:"required,lt=10"`
}

type UserGokil struct {
	Username   string `validate:"required"`
	Tagline    string `validate:"required,lt=10"`
	Impression string `validate:"is_gokil"`
}

func main() {
	// Initialize logger
	logger.InitLogger(logger.LogConfig{
		ConsoleOutput: true,
	})

	user := User{
		// Username: "Joeybloggs",
		Tagline: "This tagline is way too long.",
	}

	userGokil := UserGokil{
		Username:   "ciazhar",
		Tagline:    "aman",
		Impression: "gak gokil",
	}

	// Validator with English settings
	validate := validator.New("en")
	err := validate.ValidateStruct(user)
	if err != nil {
		logger.LogError(context.Background(), err, "validateStruct", nil)
	}

	// Validator with Indonesian settings
	validateIndo := validator.New("id")
	err = validateIndo.ValidateStruct(user)
	if err != nil {
		logger.LogError(context.Background(), err, "validateStruct", nil)
	}

	// Validate with custom validation
	// Register custom "is_gokil" validation
	err = validate.RegisterCustomValidation([]validator.CustomValidator{{Tag: "is_gokil", Handler: isGokil}})
	if err != nil {
		logger.LogError(context.Background(), err, "error override", nil)
	}
	err = validate.ValidateStruct(userGokil)
	if err != nil {
		logger.LogError(context.Background(), err, "validateStruct", nil)
	}

	// Override the translation for "is_gokil"
	err = validate.OverrideTranslationFieldOnly("is_gokil", "{0} harusnya gokil!", "field")
	if err != nil {
		logger.LogError(context.Background(), err, "error override", nil)
	}
	err = validate.ValidateStruct(userGokil)
	if err != nil {
		logger.LogError(context.Background(), err, "validateStruct", nil)
	}

	// Override the translation for "is_gokil" and "lt"
	err = validateIndo.OverrideTranslationFieldOnly("lt", "{0} harus kurang dari {1}", "field", "param")
	if err != nil {
		logger.LogError(context.Background(), err, "error override", nil)
	}
	err = validateIndo.ValidateStruct(user)
	if err != nil {
		logger.LogError(context.Background(), err, "validateStruct", nil)
	}

}

func isGokil(fl validator2.FieldLevel) bool {
	return fl.Field().String() == "gokil"
}
