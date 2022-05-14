package center

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"sync"
	"time"
)

// prefix -> map[key_n]address
type ServiceEndpointCache map[string]map[string]string

type serviceWatchUpdateController struct {
	cancel func()
	updateTime time.Time
}

type Center struct {
	client *clientv3.Client
	serviceEndpointCache ServiceEndpointCache
	mu *sync.RWMutex
	watchTTL int64
	watchServiceUpdateMap sync.Map
}

func NewCenter(registerCenterAddress []string) (*Center, error) {
	config := clientv3.Config{
		Endpoints: registerCenterAddress,
		DialTimeout:  time.Second * 30,
		DialKeepAliveTimeout: time.Second * 30,
	}

	client, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}

	center := new(Center)
	center.client = client
	center.mu = new(sync.RWMutex)
	center.serviceEndpointCache = make(ServiceEndpointCache)

	return center, nil
}

func (c Center) findService(prefix string) (error) {
	ctx := context.TODO()
	res, err := c.client.Get(ctx, prefix, clientv3.WithPrefix())

	if err != nil {
		return err
	}

	for _, ev := range res.Kvs {
		if c.serviceEndpointCache[prefix] == nil {
			c.serviceEndpointCache[prefix] = make(map[string]string)
		}
		c.serviceEndpointCache[prefix][string(ev.Key)] = string(ev.Value)
	}

	c.WatchService(prefix)

	return nil
}

func (c Center) WatchService(prefix string) {
	v, ok := c.watchServiceUpdateMap.Load(prefix)
	if ok {
		if controller, ok := v.(*serviceWatchUpdateController); ok {
			controller.updateTime = time.Now()
			return
		}
	}

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		cb := func() {cancel()}
		c.watchServiceInterval(prefix, c.watchTTL, &cb)
		ch := c.client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())

		for v := range ch {
			for _, v := range v.Events {
				key := string(v.Kv.Key)
				val := string(v.Kv.Value)

				c.mu.Lock()
				switch v.Type {
				// put
				case mvccpb.PUT: {
					c.serviceEndpointCache[prefix][key] = val
					fmt.Printf("listen service: %s change \n", key)
				}
				case mvccpb.DELETE: {
					delete(c.serviceEndpointCache[prefix], key)
					fmt.Printf("listen service: %s delele \n", key)
				}
				}
				c.mu.Unlock()
			}
		}
	}()
}

func (c Center) GetServiceEndpoint(service string) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	sm := c.serviceEndpointCache[service]
	if sm == nil {
		err := c.findService(service)
		if err != nil {
			return nil, err
		}
		sm = c.serviceEndpointCache[service]
	}

	ret := make([]string, 0)

	for _, address := range sm {
		ret = append(ret, address)
	}

	return ret, nil
}

func (c Center) watchServiceInterval(service string, sleepTime int64, callback *func()) {
	v, ok := c.watchServiceUpdateMap.Load(service)
	if ok {
		if controller, ok := v.(*serviceWatchUpdateController); ok {
			controller.updateTime = time.Now()
			return
		}
	}

	go func() {
		time.Sleep(time.Duration(sleepTime)*time.Millisecond)
		if lastUpdateTime := c.readServiceWatchUpdateTime(service); lastUpdateTime != -1 {
			curTime := time.Now().UnixNano() / 1e6
			if curTime - lastUpdateTime < c.watchTTL {
				c.watchServiceInterval(service, sleepTime, callback)
			} else {
				if callback != nil {
					(*callback)()
				}
			}
		}
	}()
}

func (c Center) readServiceWatchUpdateTime(service string) int64 {
	v, ok := c.watchServiceUpdateMap.Load(service)
	if ok {
		if controller, ok := v.(*serviceWatchUpdateController); ok {
			return controller.updateTime.UnixNano() / 1e6
		}
	}
	return -1
}
