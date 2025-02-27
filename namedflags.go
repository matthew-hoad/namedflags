// namedflags provides utils for creating and serializing bool-only structs in Go
package namedflags

import (
	"fmt"
	"reflect"
)

var uintSize int = int(reflect.TypeOf(uint(0)).Size()) * 8

// isValid check if the reflection value of the interface is valid
func isValid(reflection reflect.Value) error {
	// must be struct
	if reflection.Kind() != reflect.Struct {
		return fmt.Errorf(
			"expected struct, got %v",
			reflection.Kind(),
		)
	}
	// can't have more bits than uint
	if reflection.NumField() > uintSize {
		return fmt.Errorf(
			"struct can not have more than %d fields: has %d",
			uintSize,
			reflection.NumField(),
		)
	}
	// all fields must be bool
	for i := 0; i < reflection.NumField(); i++ {
		if reflection.Field(i).Type() != reflect.TypeOf(true) {
			return fmt.Errorf(
				"struct can only have bool fields: struct.%s is type %v",
				reflection.Type().Field(i).Name,
				reflection.Field(i).Type(),
			)
		}
	}
	return nil
}

// FromInt create a struct of bools from a uint
func FromInt[NF comparable](val uint) (res NF, err error) {
	res = *new(NF)
	reflection := reflect.ValueOf(&res).Elem()

	err = isValid(reflection)
	if err != nil {
		return
	}

	for i := 0; i < reflection.NumField(); i++ {
		reflection.Field(i).SetBool(false)
		if val&(1<<i) > 0 {
			reflection.Field(i).SetBool(true)
		} else {
			reflection.Field(i).SetBool(false)
		}
	}
	return
}

// ToInt create a uint from a struct of bools
func ToInt[NF comparable](instance NF) (res uint, err error) {
	res = uint(0)
	reflection := reflect.ValueOf(&instance).Elem()

	err = isValid(reflection)
	if err != nil {
		return
	}

	for i := 0; i < reflection.NumField(); i++ {
		if reflection.Field(i).Bool() {
			res += 1 << i
		}
	}
	return
}
