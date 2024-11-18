package nacos

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v9"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/file"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go-framework/util/xconfig/format"
	"os"
)

type ConfigNacos struct {
	Host        string `env:"CONFIG_NACOS_HOST"`
	UserName    string `env:"CONFIG_NACOS_USERNAME"`
	PassWd      string `env:"CONFIG_NACOS_PASSWD"`
	LogDir      string `env:"CONFIG_NACOS_LOG_DIR"`
	CacheDir    string `env:"CONFIG_NACOS_CACHE_DIR"`
	Port        uint64 `env:"CONFIG_NACOS_PORT" envDefault:"8848"`
	Group       string `env:"CONFIG_NACOS_GROUP" envDefault:"DEFAULT_GROUP"`
	NamespaceId string `env:"CONFIG_NACOS_NAMESPACEID"`
	DataId      string `env:"CONFIG_NACOS_DATAID"`
}

type Config struct {
	dataId      string
	formatName  string
	nacos       *ConfigNacos
	connCfg     *vo.NacosClientParam
	fileLoader  *format.Format
	client      config_client.IConfigClient
	ConfigCache map[string]interface{}
}

func NewConfig(formatName string) *Config {
	return &Config{
		formatName: formatName,
		fileLoader: format.NewFileFormat(),
	}
}

func (n *Config) Load() (map[string]interface{}, error) {
	n.nacos = new(ConfigNacos)
	if err := env.Parse(n.nacos); err != nil {
		panic(err)
	}

	if n.nacos.Host == "" || n.nacos.Port == 0 || n.nacos.UserName == "" || n.nacos.PassWd == "" {
		panic("nacos 配置文件信息缺失：" + fmt.Sprintf("nacos 配置信息：%+v", n.nacos))
	}

	// 初始化配置项
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(n.nacos.Host, n.nacos.Port),
	}

	separator := string(os.PathSeparator)
	currentPath := file.GetCurrentPath() + separator
	if n.nacos.CacheDir == "" {
		n.nacos.CacheDir = currentPath + "public" + separator + "nacos" + separator + "cache"
	}
	if n.nacos.LogDir == "" {
		n.nacos.LogDir = currentPath + "log"
	}

	cc := constant.NewClientConfig(
		constant.WithNamespaceId(n.nacos.NamespaceId),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(n.nacos.LogDir),
		constant.WithCacheDir(n.nacos.CacheDir),
		constant.WithUsername(n.nacos.UserName),
		constant.WithPassword(n.nacos.PassWd),
	)

	n.connCfg = &vo.NacosClientParam{
		ServerConfigs: sc,
		ClientConfig:  cc,
	}

	var err error
	n.client, err = clients.NewConfigClient(*n.connCfg)
	if err != nil {
		return nil, err
	}

	content, err := n.client.GetConfig(vo.ConfigParam{
		DataId: n.nacos.DataId,
		Group:  n.nacos.Group,
	})

	if err != nil {
		panic(err)
	}

	if content == "" {
		return nil, errors.New("empty content")
	}

	// 获取并缓存初始配置
	n.parseAndCacheConfig(content)

	// 初始化配置，并开始监听变更
	// n.initAndListenConfig()

	return n.ConfigCache, nil
}

func (n *Config) initAndListenConfig() {
	// 获取并缓存初始配置
	config, err := n.client.GetConfig(vo.ConfigParam{
		DataId: n.dataId,
		Group:  n.nacos.Group,
	})
	if err == nil {
		// 解析并缓存配置
		n.parseAndCacheConfig(config)
	}

	// 添加监听器
	_ = n.client.ListenConfig(vo.ConfigParam{
		DataId: n.dataId,
		Group:  n.nacos.Group,
		OnChange: func(namespace, group, dataId, data string) {
			// 配置发生变更时的处理逻辑
			n.parseAndCacheConfig(data)
		},
	})
}

func (n *Config) parseAndCacheConfig(config string) {
	// 解析配置字符串为map，这里假设配置是JSON格式
	var newConfig map[string]interface{}
	err := n.fileLoader.FileFormat[n.formatName].Load([]byte(config), &newConfig)
	if err == nil {
		n.ConfigCache = newConfig
	}
}
