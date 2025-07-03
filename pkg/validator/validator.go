package validator

import (
	"errors"
	"fmt"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/response"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	idtranslations "github.com/go-playground/validator/v10/translations/id"
)

type Validator struct {
	validate   *validator.Validate
	translator ut.Translator
}

// New initializes a new validator with translations
func New(language string) Validator {

	validate := validator.New()
	translator, err := getTranslator(validate, language)
	if err != nil {
		logger.LogFatal(err).Msg("failed to get translator")
	}

	return Validator{
		validate:   validate,
		translator: translator,
	}
}

// getTranslator sets up translations based on the language code
func getTranslator(validate *validator.Validate, language string) (ut.Translator, error) {
	if language == "en" {
		enT := en.New()
		uni := ut.New(enT, enT)
		trans, _ := uni.GetTranslator("en")
		err := entranslations.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			return nil, err
		}
		return trans, nil
	} else if language == "id" {
		id := id.New()
		uni := ut.New(id, id)
		trans, _ := uni.GetTranslator("id")
		err := idtranslations.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			return nil, err
		}
		return trans, nil
	}
	return nil, errors.New("invalid language")

}

func (v Validator) ValidateStruct(s interface{}) ([]response.ValidationError, error) {
	if err := v.validate.Struct(s); err != nil {

		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			var out []response.ValidationError
			for _, e := range errs {
				out = append(out, response.ValidationError{
					Field:   e.Field(),
					Message: e.Translate(v.translator),
				})
			}
			return out, fmt.Errorf("validation failed")
		}

		// Non-validation error
		return nil, err
	}

	return nil, nil
}

type CustomValidator struct {
	Tag     string
	Handler validator.Func
}

func (v Validator) RegisterCustomValidation(customValidators []CustomValidator) error {

	for _, vv := range customValidators {
		err := v.validate.RegisterValidation(vv.Tag, vv.Handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v Validator) OverrideTranslationFieldOnly(tag string, message string, params ...string) error {
	return v.validate.RegisterTranslation(tag, v.translator, func(ut ut.Translator) error {
		return ut.Add(tag, message, true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {

		for i := range params {
			switch params[i] {
			case "field":
				params[i] = fe.Field()
			case "param":
				params[i] = fe.Param()
			case "tag":
				params[i] = fe.Tag()
			}
		}
		t, _ := ut.T(tag, params...)
		return t
	})
}
