package helper

import "github.com/danangkonang/validation"

func Validation(data interface{}) ([]*validation.ValidationErrors, error) {
	a := validation.New()
	return a.MustValid(data)
}
