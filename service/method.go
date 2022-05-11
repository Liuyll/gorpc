package service

import (
	"errors"
	"fmt"
	"reflect"
)

type MethodType struct {
	value reflect.Value
	args  reflect.Type
	reply reflect.Type
}

func (m *MethodType) NewArgs() interface{} {
	t := m.args

	if t.Kind() == reflect.Ptr {
		t = m.args.Elem()
	}

	return reflect.New(t).Elem().Interface()
}

func (m *MethodType) NewReply() interface{} {
	t := m.reply

	if t.Kind() == reflect.Ptr {
		t = m.reply.Elem()
	}

	return reflect.New(t).Interface()
}

func (m *MethodType) Call(args interface{}, reply interface{}) {
	argsV := reflect.ValueOf(args)
	replyV := reflect.ValueOf(reply)

	params := []reflect.Value{argsV}
	ret := m.value.Call(params)[0]

	if replyV.Type().Kind() != reflect.Ptr {
		fmt.Println(errors.New("error: call must send pointer type reply"))
		return
	}

	replyV.Elem().Set(ret)
}
