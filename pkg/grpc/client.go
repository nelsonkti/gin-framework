package grpc

import (
	"context"
	"go-framework/pkg/registry"
	"go-framework/pkg/registry/etcd"
	"go-framework/pkg/rpc"
	"go-framework/pkg/transport/grpc"
	googleGrpc "google.golang.org/grpc"
)

const EtcdEndpointPrefix = "discovery:///"

type Client struct {
	conf rpc.ClientConf
}

func NewClient(c *rpc.ClientConf, ctx context.Context) *googleGrpc.ClientConn {
	if c == nil {
		return nil
	}
	var err error
	var dis registry.Discovery
	endpoint := c.Endpoint
	if c.Etcd.Hosts != nil {
		dis, err = RegistryEtcd(c.Etcd, etcd.Namespace(c.Namespace))
		if err != nil {
			panic(err)
		}
		if c.Etcd.Key != "" {
			endpoint = c.Etcd.Key
		}
		endpoint = EtcdEndpointPrefix + endpoint
	}

	conn, err := grpc.DialWithInsecure(ctx, c.Insecure, grpc.WithEndpoint(endpoint), grpc.WithDiscovery(dis))
	if err != nil {
		panic(err)
	}
	return conn
}
