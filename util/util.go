package util

import (
	"cmp"
	"errors"
)

func Find[T any](slice []T, search func(T) bool) (T, error) {
	var zeroValue T

	for _, e := range slice {
		if search(e) {
			return e, nil
		}
	}
	return zeroValue, errors.New("not found")
}

func IndexOf[T comparable](slice []T, search T) (int, error) {
	for i, e := range slice {
		if e == search {
			return i, nil
		}
	}
	return -1, errors.New("not found")
}

func Clamp[T cmp.Ordered](value, low, high T) T {
    if value > high {
        return high
    }
    if value < low {
        return low
    }
    return value
}
