package xconfig

import (
	"encoding/json"
	"fmt"
	"go-framework/util/xconfig/file"
	"go-framework/util/xconfig/nacos"
)

type ConfigReader interface {
	Load() (map[string]interface{}, error)
}

func New(c interface{}, confFile string) {
	var reader ConfigReader
	if confFile != "" {
		reader = file.NewConfig(confFile)
	} else {
		reader = nacos.NewConfig("yaml")
	}
	load(&c, reader)
}

func load(c interface{}, reader ConfigReader) {
	rawConfig, err := reader.Load()
	if err != nil {
		panic(err)
	}

	if rawConfig == nil {
		panic("config load error：must provide a config content")
	}

	configBytes, err := json.Marshal(rawConfig)
	if err != nil {
		panic(fmt.Errorf("failed to marshal config: %w", err))
	}

	if err := json.Unmarshal(configBytes, &c); err != nil {
		panic(fmt.Errorf("failed to unmarshal config into struct: %w", err))
	}
}
