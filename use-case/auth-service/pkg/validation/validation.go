package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var validate *validator.Validate

func InitValidation() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			// Format the error messages
			validationErrors = append(validationErrors,
				fmt.Sprintf("Field '%s' is invalid: %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf(strings.Join(validationErrors, ", "))
	}
	return nil
}
