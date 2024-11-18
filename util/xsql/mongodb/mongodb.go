package mongodb

import (
	"fmt"
	"go-framework/util/xsql/config"
	"go-framework/util/xsql/databese"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)
import "go.mongodb.org/mongo-driver/mongo"
import "context"

type MongoDB struct {
	c      map[string]config.DBConfig
	client map[string]*mongo.Database
}

func NewMongoDB() *MongoDB {
	return &MongoDB{
		client: make(map[string]*mongo.Database),
	}
}

func (m *MongoDB) Name() string {
	return "mongodb"
}

func (m *MongoDB) Connect(c map[string]config.DBConfig) {
	m.c = c
	for _, dbConfig := range m.c {
		err := m.connect(dbConfig)
		if err != nil {
			panic(fmt.Sprintf("The database %s connection failed， error: %s", dbConfig.Database, err))
		}
	}
}

func (m *MongoDB) connect(c config.DBConfig) error {
	// 创建客户端选项
	url := generateMongoDBURL(c)
	clientOptions := options.Client().ApplyURI(url)

	clientOptions.SetMaxPoolSize(50)
	clientOptions.SetMaxConnIdleTime(time.Hour)

	// 建立到 MongoDB 的连接
	var err error
	var client *mongo.Client
	if client, err = mongo.Connect(context.Background(), clientOptions); err != nil {
		return err
	}
	// 检查连接是否成功
	if err = client.Ping(context.Background(), nil); err != nil {
		return err
	}
	alias := c.Alias
	if alias == "" {
		alias = c.Database
	}

	m.client[alias] = client.Database(c.Database)
	return nil
}

func (m *MongoDB) ConnType(database string) bool {
	if database != "mongodb" {
		return false
	}
	return true
}

func (m *MongoDB) Result(c *databese.Engine) {
	c.Mongo = m.client
}

// generateMongoDBURL 生成 MongoDB 连接地址
func generateMongoDBURL(c config.DBConfig) string {
	urlMap := make(map[string]struct{})

	// 添加主机地址
	if c.Host != "" {
		urlMap[c.Host] = struct{}{}
	}

	// 添加数据源和副本地址
	addUniqueURLs(urlMap, c.Sources)
	addUniqueURLs(urlMap, c.Replicas)

	// 生成连接地址
	combinedURL := combineURLs(urlMap)
	mongoURL := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", c.Username, c.Password, combinedURL, c.Port, c.AuthDatabase)

	// 添加连接选项
	if c.Options != "" {
		mongoURL = appendOptions(mongoURL, c.Options)
	}

	return mongoURL
}

// addUniqueURLs 将唯一的 URL 地址添加到 map 中
func addUniqueURLs(urlMap map[string]struct{}, urls []string) {
	for _, url := range urls {
		if url != "" {
			urlMap[url] = struct{}{}
		}
	}
}

// combineURLs 将唯一的 URL 地址组合成一个字符串
func combineURLs(urlMap map[string]struct{}) string {
	var builder strings.Builder
	for url := range urlMap {
		builder.WriteString(url)
		builder.WriteString(",")
	}
	return strings.TrimSuffix(builder.String(), ",")
}

// appendOptions 将连接选项添加到连接字符串中
func appendOptions(mongoURL, options string) string {
	return fmt.Sprintf("%s?%s", mongoURL, options)
}
