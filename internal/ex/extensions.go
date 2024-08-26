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
