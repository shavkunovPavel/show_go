package model

import (
	"admigo/common"
	"fmt"
)

type Result struct {
	Message string            `json:"message"`
	Name    string            `json:"name,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// Get a result with errors
func GetErrorResult(errors map[string]string) (res *Result) {
	res = &Result{Message: "error", Errors: errors}
	return
}

func GetOkName(message_in string, name_in string) (res *Result) {
	res = &Result{Message: message_in, Name: name_in}
	return
}

// Get ok result
func GetOk(message_in string) (res *Result) {
	res = &Result{Message: message_in}
	return
}

// Checks a map of fields for empty
func Required(errs *map[string]string, fields map[string][]string) {
	for field, val := range fields {
		if len(val[0]) == 0 {
			(*errs)[field] = fmt.Sprintf(common.Mess().Required, val[1])
		}
	}
	return
}

// Checks two values for equality
func Confirmed(errs *map[string]string, values map[string][]string) {
	for field, val := range values {
		if val[0] != val[1] {
			(*errs)[field] = fmt.Sprintf(common.Mess().Confirmed, val[2])
		}
	}
}
