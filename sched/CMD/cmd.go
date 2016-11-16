package CMD

import (
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/logs"
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/id"
	typ "github.com/huawei-cloudfederation/mesos-go-stateful/common/types"
)

type CMD struct {
	CB typ.StateFul
}

func NewCMD(C typ.StateFul) *CMD {
	return &CMD{CB: C}
}

func (C *CMD) Run() {

	//Start all the go-routine
	go C.Creator()
	go C.Maintainer()
	go C.Destroy()
	//go typ.OfferList.EventHandler(JobListisQueued, JobListisEmpty, time.Second)
}



func (C *CMD) Creator() {

	var I *typ.Instance
	var CmdInfo []string

	if !typ.Gdb.IsSetup() {
		//If DB is not setup then return
		logs.Printf("CREATOR: DB is not setup")
		return
	}
	logs.Printf("Starting Creator")
	//Start an undefined Loop
	for {
		select {
		case IRequest := <-typ.Cchan:
			//Received a Instance Creation Request porcess it
			if typ.MemDb.IsValid(IRequest.Name) {
				//Sent by Maintainer (as we want to create something)
				I = typ.MemDb.Get(IRequest.Name)
				logs.Printf("CREATOR: Recived %v from Maintainer", IRequest)
				if I.Slaves < I.ExpSlaves {
					//We have to create slaves
					CmdInfo = C.CB.Config(I, false)
				} else if I.Masters == 0 {
					//We have to re-create Master
					CmdInfo = C.CB.Config(I, true)
				}
			} else {
				//Sent by HTTP module to create a new instance
				I = typ.NewInstance(IRequest.Name, 1, IRequest.NSlaves, IRequest.Spec)
				logs.Printf("CREATOR: Recived %v from HTTP", IRequest)
				//Just create Master instance offer and get out
				CmdInfo = C.CB.Config(I, true)
				typ.MemDb.Add(I.Name, I)
			}
			typ.OfferList.EnQ(typ.NewOffer(IRequest.Name, IRequest.Name+"-"+id.NewUIIDstr(), CmdInfo, I.DValue, IRequest.Spec))
		}
	}
}

func (C *CMD) Maintainer() {

	logs.Printf("Starting Maintainer")
	for {
		select {}
	}
}

func (C *CMD) Destroy() {

	logs.Printf("Starting Destroyer")
}
