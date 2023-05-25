package etcdx

import (
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

const (
	ttl         = 10
	etcdTimeout = time.Second * 3
)

type EtcdNode struct {
	//服务大区ID
	ZoneID uint16
	//服务ID
	ServiceID uint16
	//服务类型
	ServiceType string
	EtcdName    string
	RpcAddr     string
}

var watch_callback func(mvccpb.Event_EventType, string, EtcdNode)

func init() {
	watch_callback = nil
}
