package center

type ServiceMap map[string][]string
type Center struct {
	serviceMap ServiceMap
}

func (c Center) RegisterService(serviceName string, address string) {
	if c.serviceMap[serviceName] == nil {
		c.serviceMap[serviceName] = make([]string, 1)
	}

	c.serviceMap[serviceName] = append(c.serviceMap[serviceName], address)
}