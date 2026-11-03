package internal

import "github.com/go-playground/validator/v10"

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (cv *Validator) Validate(i any) error {
	return cv.validator.Struct(i)
}
