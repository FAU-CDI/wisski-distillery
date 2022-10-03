package slicesx

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Any returns true if test returns true for any of values.
func Any[T any](values []T, test func(T) bool) bool {
	for _, v := range values {
		if test(v) {
			return true
		}
	}
	return false
}

// Filter filters values in place
func Filter[T any](values []T, filter func(T) bool) []T {
	results := values[:0]
	for _, value := range values {
		if filter(value) {
			results = append(results, value)
		}
	}
	return results
}

// Partition partitions values in T by the given functions.
func Partition[T any, P comparable](values []T, partition func(value T) P) map[P][]T {
	result := make(map[P][]T)
	for _, v := range values {
		part := partition(v)
		result[part] = append(result[part], v)
	}
	return result
}

// FilterClone is like [Filter], but creates a new slice
func FilterClone[T any](values []T, filter func(T) bool) (results []T) {
	for _, value := range values {
		if filter(value) {
			results = append(results, value)
		}
	}
	return
}

// NonSequential sorts values, and then removes elements for which test() returns true.
// NonSequential does not re-allocate, but uses the existing slice.
func NonSequential[T constraints.Ordered](values []T, test func(prev, current T) bool) []T {
	if len(values) < 2 {
		return values
	}

	// sort the values and make a results array
	slices.Sort(values)
	results := values[:1]

	// do the filter loop
	prev := results[0]
	for _, current := range values[1:] {
		if !test(prev, current) {
			results = append(results, current)
		}
		prev = current
	}

	return results
}
