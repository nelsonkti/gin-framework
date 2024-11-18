package tool

import (
	"go-framework/internal/common/tool/dingtalk_tool"
	"go-framework/internal/data/common_data/tool_data"
)

type Container struct {
	DingtalkTool *dingtalk_tool.Dingtalk
}

func Register(svc *tool_data.SvcContext) *Container {
	return &Container{
		DingtalkTool: dingtalk_tool.NewDingtalkTool(svc.Conf),
	}
}
