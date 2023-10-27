package store

type Cursor[T any] interface {
	Route() T
	Cursor() string
}
