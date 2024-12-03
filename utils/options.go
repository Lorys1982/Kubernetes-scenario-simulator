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
func (o *Option[T]) None() { o.none = true } // None sets the Option object to none
func (o *Option[T]) Some(data T) {
	o.none = false
	o.some = data
} // Some sets the Option object to some, filling it with data
func (o *Option[T]) GetSome() T {
	return o.some
}

func Some(data any) Option[any] {
	return Option[any]{
		none: false,
		some: data,
	}
}

func None() Option[any] {
	return Option[any]{
		none: true,
		some: nil,
	}
}
