package util

import (
	"crypto/sha256"
	"encoding/hex"
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

func CheckPassword(password, passHash string) error {
	inputPassHash := sha256.Sum256([]byte(password))
	if hex.EncodeToString(inputPassHash[:]) != passHash {
		return fmt.Errorf("password is wrong")
	}
	return nil
}
