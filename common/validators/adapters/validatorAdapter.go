package adapters

import (
	"github.com/go-playground/validator/v10"
	"github.com/luminosita/bee/common/http"
)

type ValidatorAdapter struct {
	validate *validator.Validate
}

func NewValidatorAdapter() *ValidatorAdapter {
	return &ValidatorAdapter{
		validate: validator.New(),
	}
}

func (v *ValidatorAdapter) Validate(request *http.HttpRequest) []*http.ErrorResponse {
	var errors []*http.ErrorResponse
	err := v.validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element := http.ErrorResponse{
				FailedField: err.StructNamespace(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			}

			errors = append(errors, &element)
		}
	}
	return errors
}
