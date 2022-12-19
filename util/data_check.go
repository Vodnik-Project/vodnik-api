package util

import (
	"fmt"
	"reflect"
)

// return error if given fields of strust be empty
func CheckEmpty(s interface{}, args []string) error {
	v := reflect.ValueOf(s)
	for _, i := range args {
		if v.FieldByName(i).String() == "" {
			return fmt.Errorf("%v can't be empty", i)
		}
	}
	return nil
}
