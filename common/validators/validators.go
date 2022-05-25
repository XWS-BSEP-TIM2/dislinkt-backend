package validators

import (
	"context"
	"github.com/go-playground/validator/v10"
)

func PasswordValidator(ctx context.Context, v *validator.Validate) {
	_ = v.RegisterValidation("password_validation", validatePassword())
}

func UsernameValidator(ctx context.Context, v *validator.Validate) {
	_ = v.RegisterValidation("username_validation", usernameValidation())
}

func EmailValidator(ctx context.Context, v *validator.Validate) {
	_ = v.RegisterValidation("email_validation", emailValidation())
}

func NameValidator(ctx context.Context, v *validator.Validate) {
	_ = v.RegisterValidation("name_validation", nameValidation())
}

func NumberValidator(ctx context.Context, v *validator.Validate) {
	_ = v.RegisterValidation("number_validation", numberValidation())
}
