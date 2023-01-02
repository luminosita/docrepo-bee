package config

type ConfigLoader[T any] interface {
	ReadConfig(string, *T) error
	Validate(*T) []error
}
