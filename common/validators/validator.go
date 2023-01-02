package validators

type Validator[T any] interface {
	Validate(*T) []error
}
