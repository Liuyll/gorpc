package service

import "reflect"

type methodType struct {
	name  string
	args  reflect.Type
	reply reflect.Type
}

func (m *methodType) newArgs() reflect.Value {
	t := m.args

	if t.Kind() == reflect.Pointer {
		t = m.args.Elem()
	}

	return reflect.New(t).Elem()
}
