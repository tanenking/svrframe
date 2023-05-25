package etcdx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tanenking/svrframe/config"
	"github.com/tanenking/svrframe/constants"
	"github.com/tanenking/svrframe/logx"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdWatcher struct {
	client      *clientv3.Client
	watchPrefix string
	ndMap       map[string]EtcdNode
}

func newEtcdWatcher(endpoints []string, watchPrefix string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logx.ErrorF("%v", err)
		return
	}

	watcher := &etcdWatcher{
		client:      cli,
		watchPrefix: watchPrefix,
		ndMap:       make(map[string]EtcdNode),
	}
	go watcher.startRunner()
}
func (s *etcdWatcher) startRunner() error {
	//根据前缀获取现有的key
	ctx1, cancel1 := context.WithTimeout(context.Background(), etcdTimeout)
	defer cancel1()
	resp, err := s.client.Get(ctx1, s.watchPrefix, clientv3.WithPrefix())
	if err != nil {
		logx.ErrorF("%v", err)
		return err
	}

	for _, ev := range resp.Kvs {
		var nd = EtcdNode{}
		err := json.Unmarshal(ev.Value, &nd)
		if err == nil {
			watch_callback(mvccpb.PUT, string(ev.Key), nd)
		}
	}

	//监视前缀，修改变更的server
	go s.update()
	return nil
}
func (s *etcdWatcher) update() {
	rch := s.client.Watch(context.Background(), s.watchPrefix, clientv3.WithPrefix())
	logx.DebugF("watching prefix:%s now...", s.watchPrefix)

	defer func() {
		s.client.Close()
	}()
	exit_ch := constants.GetServiceStopListener().AddListener()
	for {
		select {
		case <-exit_ch.Done():
			return
		case wresp := <-rch:
			for _, ev := range wresp.Events {
				var ok bool = false
				key := string(ev.Kv.Key)
				var nd = EtcdNode{}
				switch ev.Type {
				case mvccpb.PUT: //修改或者新增
					err := json.Unmarshal(ev.Kv.Value, &nd)
					if err != nil {
						logx.ErrorF("%v, s = %s", err, ev.Kv.Value)
					} else {
						ok = true
						s.ndMap[key] = nd
					}
				case mvccpb.DELETE: //删除
					nd, ok = s.ndMap[key]
					if ok {
						delete(s.ndMap, key)
					} else {
						logx.ErrorF("not found key -> %s", key)
					}
				}
				if ok && watch_callback != nil {
					// not me
					if nd.EtcdName != config.GetEtcdInfo().EtcdName {
						watch_callback(ev.Type, key, nd)
					}
				}
			}
		default:
			time.Sleep(time.Second * 3)
		}
	}
}
