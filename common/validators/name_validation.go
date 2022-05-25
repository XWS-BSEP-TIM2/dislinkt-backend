package validators

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func nameValidation() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		result, _ := regexp.MatchString("[A-Za-z ,.'-]+", fl.Field().String())
		if !result {
			fmt.Println("Name is not valid!")
		}
		return result
	}
}
