package ex

import (
	"errors"
	"slices"
	"strings"
)

func AnyVales(source *[]string, values *[]string) bool {
	for _, item := range *source {
		if slices.Contains(*values, item) {
			return true
		}
	}
	return false
}

func AggregateErr(errs []error) error {
	e := make([]string, len(errs))
	for _, err := range errs {
		e = append(e, err.Error())
	}
	if len(e) == 0 {
		return nil
	}
	return errors.New(strings.Join(e, ", "))
}

func ToMap[TKey comparable, TSource any](source []TSource, selector func(TSource) TKey) map[TKey]TSource {
	result := make(map[TKey]TSource, len(source))
	for _, s := range source {
		result[selector(s)] = s
	}
	return result
}

func Map[TSource any, TResult any](source []TSource, selector func(TSource) TResult) []TResult {
	result := make([]TResult, len(source))
	for i := range source {
		result[i] = selector(source[i])
	}
	return result
}

func Add[TValue any](slice []TValue, item TValue) []TValue {
	if len(slice) == 0 {
		slice = make([]TValue, 0)
	}
	slice = append(slice, item)

	return slice
}

func Where[TValue any](source []TValue, selector func(TValue) bool) []TValue {
	if len(source) == 0 {
		return make([]TValue, 0)
	}
	res := make([]TValue, 0)
	for _, s := range source {
		if selector(s) {
			res = append(res, s)
		}
	}

	return res
}

func Include[T comparable](find T, slice []T) bool {
	for _, s := range slice {
		if s == find {
			return true
		}
	}
	return false
}
