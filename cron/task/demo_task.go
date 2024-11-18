package task

import (
	"fmt"
	"reflect"
)

type DemoTask struct {
}

func (*DemoTask) Rule() string {
	return "*/1 * * * *"
}

func (*DemoTask) Run() {
	fmt.Println("DemoTask 2222222")
	panic("123123")
}

func (m *DemoTask) Name() string {
	return reflect.TypeOf(*m).String()
}
