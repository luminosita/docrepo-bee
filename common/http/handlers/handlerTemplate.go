package handlers

import (
	"github.com/luminosita/bee/common/http"
	"github.com/luminosita/bee/common/validators"
	"github.com/luminosita/bee/common/validators/adapters"
)

type HandlerTemplate[T any] struct {
	validator validators.Validator[T]
	Handler[T]
}

func NewHandlerTemplate[T any](h Handler[T]) *HandlerTemplate[T] {
	return &HandlerTemplate[T]{
		validator: adapters.NewValidatorAdapter[T](),
		Handler:   h,
	}
}

func (bh *HandlerTemplate[T]) Process(req *http.HttpRequest) *http.HttpResponse {
	res := bh.Model(req)
	if res != nil {
		err := bh.validator.Validate(res)
		if err != nil {
			return http.BadValidation(err)
		}
	}

	bh.Handle(req)

	return http.Ok("")
}
