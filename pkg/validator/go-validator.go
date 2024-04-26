package validator

import (
	"errors"
	"fmt"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	id_translations "github.com/go-playground/validator/v10/translations/id"
	"strings"
)

type Validator struct {
	validate   *validator.Validate
	translator ut.Translator
}

func New(language string, logger logger.Logger) Validator {

	validate := validator.New()
	translator, err := getTranslator(validate, language)
	if err != nil {
		logger.Fatalf("Failed to get translator: %v", err)
	}

	return Validator{
		validate:   validate,
		translator: translator,
	}
}

func getTranslator(validate *validator.Validate, language string) (ut.Translator, error) {
	if language == "en" {
		enT := en.New()
		uni := ut.New(enT, enT)
		trans, _ := uni.GetTranslator("en")
		err := en_translations.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			return nil, err
		}
		return trans, nil
	} else if language == "id" {
		id := id.New()
		uni := ut.New(id, id)
		trans, _ := uni.GetTranslator("id")
		err := id_translations.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			return nil, err
		}
		return trans, nil
	}
	return nil, errors.New("invalid language")

}

func (v Validator) ValidateStruct(s interface{}) error {
	if err := v.validate.Struct(s); err != nil {

		var errs validator.ValidationErrors
		errors.As(err, &errs)

		var concatenatedErr strings.Builder
		for i, message := range errs {
			if i+1 < len(errs) {
				concatenatedErr.WriteString(fmt.Sprintf("%s, ", message.Translate(v.translator)))
			} else {
				concatenatedErr.WriteString(fmt.Sprintf("%s. ", message.Translate(v.translator)))
			}
		}

		return fmt.Errorf(concatenatedErr.String())
	}
	return v.validate.Struct(s)
}

type CustomValidator struct {
	Tag     string
	Handler validator.Func
}

func (v Validator) RegisterCustomValidation(validators []CustomValidator) error {

	for _, vv := range validators {
		err := v.validate.RegisterValidation(vv.Tag, vv.Handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v Validator) OverrideTranslationFieldOnly(tag string, message string, params ...string) error {
	err := v.validate.RegisterTranslation(tag, v.translator, func(ut ut.Translator) error {
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
	return err
}
