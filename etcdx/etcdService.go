package etcdx

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tanenking/svrframe/config"
	"github.com/tanenking/svrframe/constants"
	"github.com/tanenking/svrframe/logx"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdRegister struct {
	leaseID       clientv3.LeaseID
	client        *clientv3.Client
	lease         int64
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	nd            EtcdNode
}

func InitEtcdRegister(watch_cb func(mvccpb.Event_EventType, string, EtcdNode)) (err error) {
	serviceCfg := config.GetServiceInfo()
	etcdCfg := config.GetEtcdInfo()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdCfg.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logx.ErrorF("%v", err)
		return
	}

	etcd := &etcdRegister{
		client: cli,
		lease:  ttl,
		nd: EtcdNode{
			ServiceType: constants.Service_Type,
			ServiceID:   serviceCfg.ServiceID,
			ZoneID:      serviceCfg.ZoneID,
			EtcdName:    etcdCfg.EtcdName,
			RpcAddr:     fmt.Sprintf("%s:%d", serviceCfg.ServiceHost, serviceCfg.RpcPort),
		},
	}

	watch_callback = watch_cb

	if len(etcdCfg.WatchList) > 0 {
		for _, ws := range etcdCfg.WatchList {
			newEtcdWatcher(etcdCfg.Endpoints, ws)
		}
	}

	err = etcd.startRunner()
	if err != nil {
		logx.ErrorF("etcd 启动失败")
		return
	}

	go etcd.update()

	return
}

func (etcd *etcdRegister) startRunner() (err error) {
	ctx1, cancel1 := context.WithTimeout(context.Background(), etcdTimeout)
	defer cancel1()
	resp, err := etcd.client.Grant(ctx1, etcd.lease)
	if err != nil {
		logx.ErrorF("%v", err)
		return err
	}
	ndJson, err := json.Marshal(etcd.nd)
	if err != nil {
		logx.ErrorF("%v", err)
		return err
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), etcdTimeout)
	defer cancel2()
	_, err = etcd.client.Put(ctx2, etcd.nd.EtcdName, string(ndJson), clientv3.WithLease(resp.ID))
	if err != nil {
		logx.ErrorF("%v", err)
		return err
	}
	leaseRespChan, err := etcd.client.KeepAlive(context.Background(), resp.ID)
	etcd.leaseID = resp.ID
	etcd.keepAliveChan = leaseRespChan
	logx.InfoF("etcd register success -> leaseID = %d, %v", etcd.leaseID, etcd.nd)

	return
}
func (etcd *etcdRegister) update() {
	defer func() {
		etcd.client.Revoke(context.Background(), etcd.leaseID)
		etcd.client.Close()
	}()
	exit_ch := constants.GetServiceStopListener().AddListener()
	for {
		select {
		case <-exit_ch.Done():
			return
		case success := <-etcd.keepAliveChan:
			if success != nil {
				//logx.DebugF("续约结果 -> %v", success)
			} else {
				etcd.startRunner()
			}
		default:
			time.Sleep(time.Second * 3)
		}
	}
}
