package service

import "reflect"

type service struct {
	Name       string
	ServiceMap map[string]*methodType
}

func (s service) NewService(svr interface{}) {
	s := new(service)

	rt := reflect.TypeOf(svr)
	rv := reflect.ValueOf(svr)
	s.Name = rt.Name()

	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		mt := m.Type()

		if mt.NumIn() > 3 {

		}

		if mt.NumOut() > 2 {

		}
	}
}
