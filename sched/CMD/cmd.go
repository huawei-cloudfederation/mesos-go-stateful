package CMD

import (
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/id"
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/logs"
	typ "github.com/huawei-cloudfederation/mesos-go-stateful/common/types"
	"time"
)

type CMD struct {
	CB     typ.StateFul
	TaskCh chan bool
}

func NewCMD(C typ.StateFul) *CMD {
	var cmd CMD
	cmd.CB = C
	cmd.TaskCh = make(chan bool)

	return &cmd
}

func (C *CMD) Run() {

	//Start all the go-routine
	go C.Creator()
	go C.Maintainer()
	go C.Destroy()
	go typ.TaskList.EventHandler(C.TaskListQueued, C.TaskListEmpty, time.Second)
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
					//Loop and create as many slaves
					for i := I.Slaves; i < I.ExpSlaves; i++ {
						typ.OfferList.EnQ(typ.NewOffer(IRequest.Name, IRequest.Name+"-"+id.NewUIIDstr(), CmdInfo, I.DValue, IRequest.Spec))
					}
				} else if I.Masters == 0 {
					//We have to re-create Master
					CmdInfo = C.CB.Config(I, true)
					typ.OfferList.EnQ(typ.NewOffer(IRequest.Name, IRequest.Name+"-"+id.NewUIIDstr(), CmdInfo, I.DValue, IRequest.Spec))
				}
			} else {
				//Sent by HTTP module to create a new instance
				I = typ.NewInstance(IRequest.Name, 1, IRequest.NSlaves, IRequest.Spec)
				logs.Printf("CREATOR: Recived %v from HTTP", IRequest)
				//Just create Master instance offer and get out
				CmdInfo = C.CB.Config(I, true)
				typ.MemDb.Add(I.Name, I)
				typ.OfferList.EnQ(typ.NewOffer(IRequest.Name, IRequest.Name+"::"+id.NewUIIDstr(), CmdInfo, I.DValue, IRequest.Spec))
			}

		}
	}
}

func (C *CMD) TaskListQueued() bool {

	logs.Printf("TaskList is Queued")
	go func() {
		C.TaskCh <- true
	}()
	return true
}

func (C *CMD) TaskListEmpty() bool {

	logs.Printf("TaskList is Empty")
	return true
}

func (C *CMD) Maintainer() {

	var err error

	logs.Printf("MAINTAINER: Started")
	for {
		select {
		case <-C.TaskCh:

			for tEle := typ.TaskList.Front(); tEle != nil; {

				tskUpdate := tEle.Value.(typ.TaskUpdate)
				iname, id := typ.TaskSplitNames(tskUpdate.Name)
				if !typ.MemDb.IsValid(iname) {
					logs.Printf("MAINTAINOR: Recived a task update of a Non-existing Instnace %v. Ignoring...", iname)
					typ.TaskList.Delete(tEle)
					continue
				}
				I := typ.MemDb.Get(iname)
				tsk, isvalid := I.Procs[id]
				if isvalid == false {
					logs.Printf("MAINTAINOR: Recived an Update of Non-Existant TASK Instance = %v Task = %v", iname, id)
					typ.TaskList.Delete(tEle)
					continue
				}
				switch tskUpdate.State {
				case "TASK_STAGING":
				case "TASK_STARTING":
				case "TASK_RUNNING":
					//Invoke the Call back
					if tsk.Type == "M" { //Should have some endpoint
						//Invoke the master CalBack
						err = C.CB.MasterRunning(I)
					} else {

						//Invoke the slave call back
						err = C.CB.SlaveRunning(I)
					}
					if err != nil {
						logs.Printf("Error occured CallBack Invokating of Master/Slave")
					}
				case "TASK_FINISHED":
					//When task finish execution themselves
					if I.Status != typ.INST_STATUS_DELETED {
						//Something wrong we did not initiate the SHUTDOWN signal, treate it like any other CRASH/LOST signal
						if tsk.Type == "M" {
							err = C.CB.MasterLost(I)
						} else {
							err = C.CB.SlaveLost(I)
						}
						logs.Printf("MAINTAINOR: TaskLost Call Back Invocation Error:%v", err)
					}

				case "TASK_ERROR", "TASK_FAILED", "TASK_KILLED", "TASK_LOST":

					if tsk.Type == "M" {
						err = C.CB.MasterLost(I)
					} else {
						err = C.CB.SlaveLost(I)
					}
					logs.Printf("MAINTAINOR: TaskLost Call Back Invocation Error:%v", err)

				}

			}
		}
	}
	logs.Printf("MAINTAINER: Terminated")
}

func (C *CMD) Destroy() {

	logs.Printf("Starting Destroyer")
}
