// Copyright © 2025 Ping Identity Corporation

package utils

func Pointer[T any](value T) *T {
	return &value
}
