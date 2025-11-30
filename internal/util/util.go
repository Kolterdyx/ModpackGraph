package util

func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func ValueOrDefault[T any](value *T, defaultValue T) T {
	if value != nil {
		return *value
	}
	return defaultValue
}

func ValueOrZero[T any](value *T) T {
	var zero T
	return ValueOrDefault(value, zero)
}

func Ptr[T any](value T) *T {
	return &value
}
