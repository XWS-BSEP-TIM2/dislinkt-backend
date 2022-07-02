package validators

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/go-playground/validator/v10"
)

func PasswordValidator(ctx context.Context, v *validator.Validate) {
	span := tracer.StartSpanFromContext(ctx, "PasswordValidator")
	defer span.Finish()

	_ = v.RegisterValidation("password_validation", validatePassword())
}

func UsernameValidator(ctx context.Context, v *validator.Validate) {
	span := tracer.StartSpanFromContext(ctx, "UsernameValidator")
	defer span.Finish()

	_ = v.RegisterValidation("username_validation", usernameValidation())
}

func EmailValidator(ctx context.Context, v *validator.Validate) {
	span := tracer.StartSpanFromContext(ctx, "EmailValidator")
	defer span.Finish()

	_ = v.RegisterValidation("email_validation", emailValidation())
}

func NameValidator(ctx context.Context, v *validator.Validate) {
	span := tracer.StartSpanFromContext(ctx, "NameValidator")
	defer span.Finish()

	_ = v.RegisterValidation("name_validation", nameValidation())
}

func NumberValidator(ctx context.Context, v *validator.Validate) {
	span := tracer.StartSpanFromContext(ctx, "NumberValidator")
	defer span.Finish()

	_ = v.RegisterValidation("number_validation", numberValidation())
}
