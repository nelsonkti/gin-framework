package config

import "go-framework/pkg/rpc"

type Config struct {
	App    App      `json:"app"`
	Server Server   `json:"server"`
	Etcd   rpc.Etcd `json:"etcd"`
}

type App struct {
	Name string `json:"name"`
}

type Server struct {
	Http Network `json:"http"`
	Rpc  Network `json:"rpc"`
}

type Network struct {
	Addr string `json:"addr"`
	Mode string `json:"mode"`
}
