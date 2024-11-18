package rpc

type ClientConf struct {
	Etcd      Etcd   `json:"etcd"`
	Endpoint  string `json:"endpoint"`
	Insecure  bool   `json:"insecure"`
	Namespace string `json:"namespace"`
}

type Etcd struct {
	Hosts       []string
	Key         string
	ID          int64  `json:",optional"`
	User        string `json:",optional"`
	Pass        string `json:",optional"`
	Namespace   string `json:",optional"`
	CertFile    string `json:",optional"`
	CertKeyFile string `json:",optional=CertFile"`
	CACertFile  string `json:",optional=CertFile"`
}
