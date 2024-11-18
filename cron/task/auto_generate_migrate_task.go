package task

import (
	"fmt"
	"reflect"
)

type AutoGenerateMigrateTask struct {
}

func (*AutoGenerateMigrateTask) Rule() string {
	return "*/1 * * * *"
}

func (*AutoGenerateMigrateTask) Run() {
	fmt.Println("2222222")
}

func (m *AutoGenerateMigrateTask) Name() string {
	return reflect.TypeOf(*m).String()
}
