package xhttp

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// RespJsonModifier 用于封装捕获和修改响应数据的逻辑
type RespJsonModifier struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// NewRespJsonModifier 创建一个新的 RespJsonModifier 实例
func NewRespJsonModifier(c *gin.Context) *RespJsonModifier {
	return &RespJsonModifier{
		ResponseWriter: c.Writer,
		body:           bytes.NewBufferString(""),
	}
}

func (rm *RespJsonModifier) Body() []byte {
	return rm.body.Bytes()
}

// Write 捕获响应体数据
func (rm *RespJsonModifier) Write(b []byte) (int, error) {
	return rm.body.Write(b)
}

func (rm *RespJsonModifier) WriteMeta(c *gin.Context, value interface{}) {
	rm.WriteResponse(c, "meta", value)
}

// WriteResponse 修改响应数据并写回响应
func (rm *RespJsonModifier) WriteResponse(c *gin.Context, key string, value interface{}) {
	responseBody := rm.Body()
	var response map[string]interface{}
	if err := json.Unmarshal(responseBody, &response); err == nil {
		newData, ok := response[key].(map[string]interface{})
		if !ok {
			newData = make(map[string]interface{})
		}

		// 设置新的键值对
		for k, v := range value.(map[string]interface{}) {
			newData[k] = v
		}
		if key != "meta" && newData != nil {
			response[key] = newData
		}
		rm.write(c, response)
	}
}

// WriteResponse 修改响应数据并写回响应
func (rm *RespJsonModifier) write(c *gin.Context, response map[string]interface{}) {
	// 获取控制器返回的响应体
	newResponseBody, _ := json.Marshal(response)

	// 清空原始响应并写入修改后的响应
	c.Writer = rm.ResponseWriter
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(newResponseBody)
}
