package config

type Conf struct {
	App      App           `json:"app"`      // 应用配置
	Server   Server        `json:"server"`   // 服务配置
	Log      Log           `json:"log"`      // 日志配置
	DB       map[string]DB `json:"db"`       // 数据库配置
	Redis    []Redis       `json:"redis"`    // redis配置
	MQ       MQ            `json:"mq"`       // mq配置
	Trace    Trace         `json:"trace"`    // 链路追踪
	Dingtalk Dingtalk      `json:"dingtalk"` // 钉钉配置
}

type App struct {
	Name         string `json:"name"` // 应用名称
	Env          string `json:"env"`  // 环境
	Key          string `json:"key"`
	ServerNumber int    `json:"server_number"` // 服务器编号
}

type Server struct {
	Http Network `json:"http"` // http配置
	Rpc  Network `json:"rpc"`  // rpc配置
}

type Log struct {
	Path string `json:"path"` // 日志路径
}

type Network struct {
	Addr string `json:"addr"` // 地址
	Mode string `json:"mode"` // 模式（etcd、直连）
}

type DB struct {
	Driver       string   `json:"driver"`        // 数据库驱动
	Host         string   `json:"host"`          // 地址
	Sources      []string `json:"sources"`       // 主库
	Replicas     []string `json:"replicas"`      // 从库
	Port         int      `json:"port"`          // 端口
	Username     string   `json:"username"`      // 用户名
	Password     string   `json:"password"`      // 密码
	AuthDatabase string   `json:"auth_database"` // 验证数据库（MongoDB）
	Database     string   `json:"database"`      // 数据库
	Alias        string   `json:"alias"`         // 别名
	Options      string   `json:"options"`       // 选项
}

type Redis struct {
	Host     string `json:"host"`     // 地址
	Port     int    `json:"port"`     // 端口
	Database int    `json:"database"` // 数据库
	Alias    string `json:"alias"`    // 别名
	UserName string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

type MQ struct {
	Endpoint  []string `json:"endpoint"`   // 地址
	AccessKey string   `json:"access_key"` // accessKey
	SecretKey string   `json:"secret_key"` // secretKey
	Namespace string   `json:"namespace"`  // namespace（instanceId）
	Env       string   `json:"env"`        // environment（几服）
}

type Etcd struct {
	Hosts       []string
	Key         string
	ID          int64  `json:",optional"`
	User        string `json:",optional"`
	Pass        string `json:",optional"`
	CertFile    string `json:",optional"`
	CertKeyFile string `json:",optional=CertFile"`
	CACertFile  string `json:",optional=CertFile"`
}

// Dingtalk 钉钉配置
type Dingtalk struct {
	Robots Robots `json:"robots"`
}

// Robots 钉钉机器人
type Robots struct {
	AlarmSecret string `json:"alarm_secret"` // 码商机器人警报秘钥
}

// Trace 链路追踪
type Trace struct {
	Endpoint string `json:"endpoint"`
	UrlPath  string `json:"url_path"`
}
