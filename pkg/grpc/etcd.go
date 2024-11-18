package grpc

import (
	"go-framework/pkg/registry/etcd"
	"go-framework/pkg/rpc"
	etcdclient "go.etcd.io/etcd/client/v3"
)

func RegistryEtcd(conf rpc.Etcd, opts ...etcd.Option) (*etcd.Registry, error) {
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: conf.Hosts,
	})
	if err != nil {
		return nil, err
	}
	return etcd.New(client, opts...), nil
}
