package utils

import (
	"reflect"
)

func SetInterfacePtr(p interface{}, v interface{}) {
	rv := reflect.ValueOf(p)
	rt := reflect.TypeOf(p)

	if rt.Kind() != reflect.Ptr {
		panic("only support pointer type")
	}

	rv.Elem().Set(reflect.ValueOf(v))
}