package main

import (
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/validator"
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
	user := User{
		//Username: "Joeybloggs",
		Tagline: "This tagline is way too long.",
	}

	userGokil := UserGokil{
		Username:   "ciazhar",
		Tagline:    "aman",
		Impression: "gak gokil",
	}

	//logger
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})

	//validate english
	validate := validator.New("en", log)
	err := validate.ValidateStruct(user)
	if err != nil {
		log.Infof("validateStruct : %v", err)
	}

	//validate indo
	validateIndo := validator.New("id", log)
	err = validateIndo.ValidateStruct(user)
	if err != nil {
		log.Infof("validateStruct : %v", err)
	}

	//validate custom
	err = validate.RegisterCustomValidation([]validator.CustomValidator{{"is_gokil", isGokil}})
	if err != nil {
		return
	}
	err = validate.ValidateStruct(userGokil)
	if err != nil {
		log.Infof("validateStruct : %v", err)
	}

	//validate override
	err = validate.OverrideTranslationFieldOnly("is_gokil", "{0} harusnya gokil!", "field")
	if err != nil {
		log.Infof("error override : %v", err)
	}
	err = validate.ValidateStruct(userGokil)
	if err != nil {
		log.Infof("validateStruct : %v", err)
	}

	err = validateIndo.OverrideTranslationFieldOnly("lt", "{0} harus kurang dari {1}", "field", "param")
	if err != nil {
		log.Infof("error override : %v", err)
	}
	err = validateIndo.ValidateStruct(user)
	if err != nil {
		log.Infof("validateStruct : %v", err)
	}

}

func isGokil(fl validator2.FieldLevel) bool {
	return fl.Field().String() == "gokil"
}
