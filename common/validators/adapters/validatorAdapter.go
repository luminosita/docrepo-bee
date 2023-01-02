package adapters

import (
	"github.com/go-playground/validator/v10"
	"github.com/luminosita/bee/common/errors"
	"github.com/luminosita/bee/common/log"
)

type ValidatorAdapter[T any] struct {
	validator *validator.Validate
}

func NewValidatorAdapter[T any]() *ValidatorAdapter[T] {
	return &ValidatorAdapter[T]{
		validator: validator.New(),
	}
}

func (v *ValidatorAdapter[T]) Validate(obj *T) []error {
	var e []error

	log.GetLogger().Infof("Validation %+v", obj)

	err := v.validator.Struct(obj)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			element := errors.ValidationError{
				FailedField: err.StructNamespace(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			}

			e = append(e, &element)
		}
	}
	return e
}
