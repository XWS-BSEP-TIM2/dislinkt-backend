package validators

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func validatePassword() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		if len(fl.Field().String()) < 8 {
			fmt.Println("Password must contain 8 characters or more!")
			return false
		}
		result, _ := regexp.MatchString("(.*[a-z].*)", fl.Field().String())
		if !result {
			fmt.Println("Password must contain lower case characters!")
		}
		result, _ = regexp.MatchString("(.*[A-Z].*)", fl.Field().String())
		if !result {
			fmt.Println("Password must contain upper case characters!")
		}
		result, _ = regexp.MatchString("(.*[0-9].*)", fl.Field().String())
		if !result {
			fmt.Println("Password must contain numbers!")
		}

		result, _ = regexp.MatchString("(.*[!@#$%^&*(){}\\[:;\\]<>,\\.?~_+\\-\\\\=|/].*)", fl.Field().String())
		if !result {
			fmt.Println("Password must contain special characters!")
		}
		return result
	}
}
