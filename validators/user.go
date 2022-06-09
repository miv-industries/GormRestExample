package validators

import (
	"github.com/go-playground/validator/v10"
)

type User struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
}

var validateUser = validator.New()
