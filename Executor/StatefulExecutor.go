package Executor

import (
	typ "github.com/huawei-cloudfederation/mesos-go-stateful/common/types"
)

//An interface for Custom Executor
type StatefulExecutor interface {
	Config (*typ.Task)
	TaskStarted (*typ.Task)
	Cleanup(*typ.Task)
	UpdateConfig(*typ.Task)
}
