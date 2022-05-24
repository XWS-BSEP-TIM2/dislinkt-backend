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
