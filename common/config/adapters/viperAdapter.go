package adapters

import (
	"github.com/luminosita/bee/common/validators"
	"github.com/luminosita/bee/common/validators/adapters"
	"github.com/spf13/viper"
)

type ViperAdapter[T any] struct {
	viper     *viper.Viper
	validator validators.Validator[T]
}

func NewViperAdapter[T any](viper *viper.Viper) *ViperAdapter[T] {
	return &ViperAdapter[T]{
		viper:     viper,
		validator: adapters.NewValidatorAdapter[T](),
	}
}

func (v *ViperAdapter[T]) ReadConfig(key string, config *T) error {
	return v.viper.UnmarshalKey(key, config)
}

func (v *ViperAdapter[T]) Validate(config *T) []error {
	return v.validator.Validate(config)
}
