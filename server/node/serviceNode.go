package node

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/google/uuid"
	"time"
)

type ServiceNode struct {
	center *clientv3.Client
	ttl int64
	name string
	address string
}

func NewServiceNode(address string, name string, ttl int64) (error, *ServiceNode) {
	conf := clientv3.Config{
		Endpoints: []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}

	serviceNode := new(ServiceNode)
	if client, err := clientv3.New(conf); err != nil {
		return err, nil
	} else {
		serviceNode.center = client
	}

	serviceNode.ttl = ttl
	serviceNode.name = name
	serviceNode.address = address

	if err := serviceNode.registerNode(); err != nil {
		return err, nil
	}
	return nil, serviceNode
}

func (this ServiceNode) registerNode() error {
	lease := clientv3.Lease(this.center)

	ctx := context.TODO()
	grantRes, err := lease.Grant(ctx, this.ttl)

	if err != nil {
		return err
	}

	alives, err := lease.KeepAlive(ctx, grantRes.ID)
	this.listenPing(alives)

	uid := uuid.New()
	this.put(uid.String(), grantRes)

	return nil
}

func (this ServiceNode) listenPing(prod <-chan *clientv3.LeaseKeepAliveResponse)  {

}

func (this ServiceNode) put(uid string, mGrant interface{}) {
	kv := clientv3.NewKV(this.center)
	registerName := fmt.Sprintf("%s/%s", this.name, uid)
	fmt.Println("registerName:", registerName, )

	if grantRes, ok := mGrant.(*clientv3.LeaseGrantResponse); ok {
		kv.Put(context.TODO(), registerName, this.address, clientv3.WithLease(grantRes.ID))
	} else {
		kv.Put(context.TODO(), registerName, this.address)
	}
}


