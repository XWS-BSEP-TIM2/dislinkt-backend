package validators

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func emailValidation() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		result, _ := regexp.MatchString("\\S+@\\S+\\.\\S+", fl.Field().String())
		if !result {
			fmt.Println("Email is not valid!")
		}
		return result
	}
}
