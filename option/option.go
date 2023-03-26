package option

type Option[T any] struct {
	Present bool
	Data    T
}

func Some[T any](data T) Option[T] {
	return Option[T]{
		Present: true,
		Data:    data,
	}
}

func None[T any]() Option[T] {
	return Option[T]{
		Present: false,
	}
}
