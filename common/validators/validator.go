package validations

import (
	"github.com/luminosita/bee/common/errors"
)

type Validator[T any] interface {
	Validate(T) []*errors.Error
}
