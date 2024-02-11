package util

type Null[T any] struct {
	Valid bool
	Data  T
}
