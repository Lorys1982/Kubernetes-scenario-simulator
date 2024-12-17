package utils

// Option simple implementation from rust, if functions with multiple
// optional parameters are ever needed
type Option[T any] struct {
	none bool
	some T
}

func (o *Option[T]) IsNone() bool {
	return o.none
}
func (o *Option[T]) IsSome() bool {
	return !o.none
}
func (o *Option[T]) GetSome() T {
	return o.some
}

func Some[T any](data T) Option[T] {
	return Option[T]{
		none: false,
		some: data,
	}
}

func None[T any]() Option[T] {
	return Option[T]{
		none: true,
	}
}
