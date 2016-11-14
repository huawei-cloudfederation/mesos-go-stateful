package types

import (
	"container/list"

	"github.com/huawei-cloudfederation/mesos-go-stateful/common/agentstate"
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/logs"
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/store/etcd"
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/store/zookeeper"
)

//Initialize Initialize all the data strucutres in common package, should be called by the main program only and should be called only once per program
func Initialize(dbtype string, config string) (bool, error) {

	//Initalize all the communication channels
	OfferList = list.New()
	OfferList.Init()
	Cchan = make(chan HttpToCR)
	Mchan = make(chan *TaskUpdate) //Channel for Maintainer
	Dchan = make(chan TaskMsg)     //Channel for Destroyer

	Agents = agentstate.NewState()

	//Initalize the Internal in-memory storage
	MemDb = NewInMem()

	//Initalize the store db
	switch dbtype {
	case "etcd":
		Gdb = etcd.New()
		err := Gdb.Setup(config)
		if err != nil {
			logs.FatalInfo("Failed to setup etcd database error:%v", err)
		}
		return Gdb.IsSetup(), nil
	case "zookeeper":
		Gdb = zookeeper.New()
		err := Gdb.Setup(config)
		if err != nil {
			logs.FatalInfo("Failed to setup zookeeper database error:%v", err)
		}
		return Gdb.IsSetup(), nil
	}

	return true, nil
}
