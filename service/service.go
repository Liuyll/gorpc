package service

import (
	"fmt"
	"reflect"
	"strings"
)

type Service struct {
	name       string
	serviceMap map[string]*MethodType
	self reflect.Value
}

func (s Service) GetMethod(methodName string) *MethodType {
	return s.serviceMap[methodName]
}

func (s Service) GetSelf() reflect.Value {
	return s.self
}

func NewService(svr interface{}) *Service {
	s := new(Service)
	s.serviceMap = make(map[string]*MethodType)

	rt := reflect.TypeOf(svr)
	rv := reflect.ValueOf(svr)
	s.name = rt.Name()
	s.self = reflect.ValueOf(svr)

	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		mt := m.Type
		mName := m.Name

		if strings.HasPrefix(mName, "Unmarshal") {
			continue
		}
		fmt.Println("register method:", mName, mt.Kind())

		if mt.NumIn() != 2 {
			fmt.Println(fmt.Sprintf("method %s arguments length not equal one", mName))
			continue
		}

		if mt.NumOut() != 1 {
			fmt.Println(fmt.Sprintf("method %s return length not equal one", mName))
			continue
		}

		if _, ok := rt.MethodByName("Unmarshal" + mName); !ok {
			fmt.Printf("method: %s must unmarshal way \n", mName)
			continue
		}

		s.serviceMap[mName] = &MethodType{
			rv.Method(i),
			mt.In(1),
			mt.Out(0),
			rv.MethodByName("Unmarshal" + mName),
		}
	}

	return s
}

