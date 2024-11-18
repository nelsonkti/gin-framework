package es

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"go-framework/util/helper"
)

type Client struct {
	client *elasticsearch.Client
}

func NewClient(conf interface{}) (*Client, error) {
	var configs elasticsearch.Config
	err := helper.UnMarshalWithInterface(conf, &configs)
	if err != nil {
		panic(err)
	}

	// 创建客户端连接
	client, err := elasticsearch.NewClient(configs)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch.NewClient failed, err:%v\n", err)
	}

	// 检查客户端是否可以连接到Elasticsearch集群
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("Error getting response: %v\n", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error: %v", res)
	}

	return &Client{client: client}, err
}

func (c *Client) Client() *elasticsearch.Client {
	return c.client
}
