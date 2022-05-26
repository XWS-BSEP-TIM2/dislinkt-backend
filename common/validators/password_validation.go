package validators

import (
	"bufio"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
	"path/filepath"
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

func IsPasswordCracked(password string) bool {

	p := filepath.FromSlash("./security/cracked_password.txt")
	println(p)

	f, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// do something with a line
		if password == scanner.Text() {
			return false
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return true
}
