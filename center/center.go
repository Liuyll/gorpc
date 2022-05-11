package center

type ServiceMap map[string][]string
type Center struct {
	serviceMap ServiceMap
}

func (c Center) RegisterService(serviceName string, address string) {

}