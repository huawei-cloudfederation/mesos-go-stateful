package mesoslib

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/logs"
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/store/etcd"
	typ "github.com/huawei-cloudfederation/mesos-go-stateful/common/types"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
	"time"
)

//WorkloadScheduler scheudler struct
type WorkloadScheduler struct {
	executor *mesos.ExecutorInfo
}

//NewWorkloadScheduler Constructor
func NewWorkloadScheduler(exec *mesos.ExecutorInfo) *WorkloadScheduler {

	return &WorkloadScheduler{executor: exec}
}

//Registered Scheduler register call back initializes the timestamp and framework id
func (S *WorkloadScheduler) Registered(driver sched.SchedulerDriver, frameworkID *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	logs.Printf("Framework %s Registered %v", typ.Cfg.FrameworkName, frameworkID)

	FwIDKey := etcd.ETCD_CONFDIR + "/FrameworkID"
	typ.Gdb.Set(FwIDKey, frameworkID.GetValue())
	FwTstamp := etcd.ETCD_CONFDIR + "/RegisteredAt"
	typ.Gdb.Set(FwTstamp, time.Now().String())
}

//Reregistered re-register call back simply updates the timestamp
func (S *WorkloadScheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	logs.Printf("Famework %s Re-registered", typ.Cfg.FrameworkName)
	FwTstamp := etcd.ETCD_CONFDIR + "/RegisteredAt"
	typ.Gdb.Set(FwTstamp, time.Now().String())
}

//Disconnected Not implemented call back
func (S *WorkloadScheduler) Disconnected(sched.SchedulerDriver) {
	logs.Printf("Framework %s Disconnected", typ.Cfg.FrameworkName)
}

//ResourceOffers The moment we recive some offers we loop through the OfferList (container/list)
//see if we have any task that will fit this offers being sent
func (S *WorkloadScheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {

	//No work to do so reject all the offers we just received
	offerCount := typ.OfferList.Len()
	if offerCount <= 0 {
		//Reject the offers nothing to do now
		ids := make([]*mesos.OfferID, len(offers))
		for i, offer := range offers {
			ids[i] = offer.Id
		}
		driver.LaunchTasks(ids, []*mesos.TaskInfo{}, &mesos.Filters{})
		//logs.Printf("No task to peform reject all the offer")
		return
	}

	//We have some task and should consume the offers sent by the masters
	//Pick one task and check if any of the offer is suitable

	//Loop thought he offers
	for _, offer := range offers {

		cpuResources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
			return res.GetName() == "cpus"
		})
		cpus := 0.0
		for _, res := range cpuResources {
			cpus += res.GetScalar().GetValue()
		}

		memResources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
			return res.GetName() == "mem"
		})
		mems := 0.0
		for _, res := range memResources {
			mems += res.GetScalar().GetValue()
		}

		logs.Printf("Received Offer with CPU=%v MEM=%v OfferID=%v", cpus, mems, offer.Id.GetValue())
		var tasks []*mesos.TaskInfo

		//Loop through the tasks
		for tskEle := typ.OfferList.Front(); tskEle != nil; {

			tsk := tskEle.Value.(typ.Offer)
			tskCPUFloat := float64(tsk.Cpu)
			tskMemFloat := float64(tsk.Mem)

			var tmpData []byte

			if tsk.IsMaster {
				tmpData = []byte(fmt.Sprintf("%d Master", tsk.Mem))
			} else {
				tmpData = []byte(fmt.Sprintf("%d SlaveOf %s", tsk.Mem, tsk.MasterIpPort))
			}

			if cpus >= tskCPUFloat && mems >= tskMemFloat && typ.Agents.Canfit(offer.SlaveId.GetValue(), tsk.Name, tsk.DValue) {
				tskID := &mesos.TaskID{Value: proto.String(tsk.Taskname)}
				mesosTsk := &mesos.TaskInfo{
					Name:     proto.String(tsk.Taskname),
					TaskId:   tskID,
					SlaveId:  offer.SlaveId,
					Executor: S.executor,
					Resources: []*mesos.Resource{
						util.NewScalarResource("cpus", tskCPUFloat),
						util.NewScalarResource("mem", tskMemFloat),
					},
					Data: tmpData,
				}
				mems -= tskMemFloat
				cpus -= tskCPUFloat

				currentTask := tskEle
				tskEle = tskEle.Next()
				typ.OfferList.Remove(currentTask)
				tasks = append(tasks, mesosTsk)
				typ.Agents.Add(offer.SlaveId.GetValue(), tsk.Name, 1)

			} else {
				tskEle = tskEle.Next()
			}
			//Check if this task is suitable for this offer
		}
		driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, &mesos.Filters{})
		logs.Printf("Launched %d tasks from this offer", len(tasks))
	}
	logs.Printf("workload Receives offer")
}

//StatusUpdate Simply recives the update and passes it to the Maintainer goroutine
func (S *WorkloadScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {

	var ts typ.TaskUpdate
	ts.Name = status.GetTaskId().GetValue()
	ts.State = status.GetState().String()
	ts.Data = status.GetData()
	logs.Printf("workload Task Update received")
	logs.Printf("Status=%v", ts)

	//Send it across to the channel to maintainer
	//typ.Mchan <- &ts

}

//OfferRescinded Not implemented
func (S *WorkloadScheduler) OfferRescinded(_ sched.SchedulerDriver, oid *mesos.OfferID) {
	logs.Printf("offer rescinded: %v", oid)
}

//FrameworkMessage not implemented
func (S *WorkloadScheduler) FrameworkMessage(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, msg string) {
	logs.Printf("framework message from executor %q slave %q: %q", eid, sid, msg)
}

//SlaveLost Not implemented
func (S *WorkloadScheduler) SlaveLost(_ sched.SchedulerDriver, sid *mesos.SlaveID) {
	logs.Printf("slave lost: %v", sid)
}

//ExecutorLost Not implemented
func (S *WorkloadScheduler) ExecutorLost(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, code int) {
	logs.Printf("executor %q lost on slave %q code %d", eid, sid, code)
}

//Error Not implemeted
func (S *WorkloadScheduler) Error(_ sched.SchedulerDriver, err string) {
	logs.Printf("Scheduler received error:%v", err)
}
