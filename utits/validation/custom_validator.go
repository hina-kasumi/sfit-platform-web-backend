package validation

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{}
}

func (cv *CustomValidator) PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasSpecial && hasNumber
}
