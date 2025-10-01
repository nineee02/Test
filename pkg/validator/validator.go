package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		pattern := fl.Param()
		re := regexp.MustCompile(pattern)
		return re.MatchString(fl.Field().String())
	})

	return &Validator{validate: v}
}

func (v *Validator) Validate(s interface{}) error {
	return v.validate.Struct(s)
}

func IsValidationErrors(err error) (validator.ValidationErrors, bool) {
	vErrs, ok := err.(validator.ValidationErrors)
	return vErrs, ok
}
