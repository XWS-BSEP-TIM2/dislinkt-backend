package validators

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func usernameValidation() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		if len(fl.Field().String()) < 4 || len(fl.Field().String()) > 16 {
			fmt.Println("Your username must be between 4 and 16 characters long.")
			return false
		}
		result, _ := regexp.MatchString("(.*[\"'!@#$%^&*(){}\\[:;\\]<>,\\.?~_+\\-\\\\=|/].*)", fl.Field().String())
		if result {
			fmt.Println("Description contains special characters that are not allowed!")
		}
		return !result
	}
}
