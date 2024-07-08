package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validate struct {
	validate *validator.Validate
}

func New() *Validate {
	return &Validate{validate: validator.New()}
}

func (v *Validate) Validate(i interface{}) error {
	return v.validate.Struct(i)
}
