package validators

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func numberValidation() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		result, _ := regexp.MatchString("^([0-9]{8,12})$", fl.Field().String())
		if !result {
			fmt.Println("Phone number is not valid!")
		}
		return result
	}
}
