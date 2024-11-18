package job

import (
	"fmt"
	"go-framework/util/mq/queue"
	"reflect"
	"time"
)

var _ queue.Job = (*OrderJob)(nil)

type OrderJob struct {
}

func (o *OrderJob) Name() string {
	return reflect.TypeOf(*o).String()
}

func (o *OrderJob) Execute(bytes []byte) error {
	fmt.Println(string(bytes))
	time.Sleep(time.Second * 5)
	//fmt.Println("准备进行报错")
	//panic("not implemented")
	fmt.Println("执行完成")
	return nil
}
