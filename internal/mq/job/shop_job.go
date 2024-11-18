package job

import (
	"fmt"
	"go-framework/util/mq/queue"
	"reflect"
)

var _ queue.Job = (*ShopJob)(nil)

type ShopJob struct {
}

func (o *ShopJob) Name() string {
	return reflect.TypeOf(*o).String()
}

func (o *ShopJob) Execute(bytes []byte) error {
	fmt.Println(string(bytes))
	return nil
}
