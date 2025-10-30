package validator

import (
	"fmt"
	"regexp"

	validators "github.com/go-playground/validator/v10"
)

// Validator is an interface for validating structs.
type Validator interface {
	ValidateStruct(inf interface{}) error
	TranslateError(err error) []string
}

type validator struct {
	validator *validators.Validate
}

// New creates a new Validator instance.
func New() Validator {
	v := validators.New()

	err := v.RegisterValidation("password", func(fl validators.FieldLevel) bool {
		password := fl.Field().String()
		return isValidPassword(password)
	})

	if err != nil {
		panic(err)
	}

	return &validator{
		validator: v,
	}
}

// ValidateStruct validates a struct using the validator.
func (v *validator) ValidateStruct(inf interface{}) error {

	return v.validator.Struct(inf)
}

func (v *validator) TranslateError(err error) []string {
	if err == nil {
		return nil
	}

	var errs []string
	for _, err := range err.(validators.ValidationErrors) {
		var msg string

		switch err.Tag() {
		case "password":
			msg = "Password must be at least 8 characters long, contain uppercase, lowercase, number, and special character"
		default:
			msg = fmt.Sprintf("%s is not valid", err.Field())
		}

		errs = append(errs, msg)
	}

	return errs
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return hasLower && hasUpper && hasNumber && hasSpecial
}
