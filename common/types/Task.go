package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/huawei-cloudfederation/mesos-go-stateful/common/logs"
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/store/etcd"
)

//Task A standalone task KV store is usually started in any slave (Linux) like below
//$./server -p <PORT> ..... {OPTIONS}
//This stand alone server will be an actual unix process bound to a particular port witha PID
//A workload Master Slave setup will have two such "server" processes running in either the same machine or two different machines
//The below structure "Task" is a representation of such a running 'server' process started via this framework
type Task struct {
	Name     string    //Name of the Task
	IName    string    //Instance
	ID       string    //UUID of this task
	Pid      int       //Unix Process id of this running instance
	State    string    //Current state of the process Active/Dead/Crashed etc.,
	Type     string    //Type of the PROC master/Slave etc.,
	SlaveOf  string    //Slave of which workload master
	Stats    TaskStats //All other statistics apart from Memory usage to be stored as a json/string
	ConfFile string    //Config file to be used for this task
	EID      string    //Executor ID of this PROC  .. Just in case we need to send a workload messsage
	SID      string    //Slave ID of this PROC .. Just in case we need to send a workload message
	Nodename string    //node name or path/node to write data (store)
}

//ProcJson Fields to be packed in a json when replied to a HTTP REST query.
type TaskStats struct {
	IName    string
	IP       string
	Port     string
	CmdInfo  string
	Capacity WLSpec
	Used     WLSpec
}

func (p *TaskStats) ToJson() string {
	rc, err := json.Marshal(p)
	if err != nil {
		logs.Printf("ProcStats: Json Marshall error %v", err)
		return "{}"
	}
	return string(rc)
}

func (p *TaskStats) FromJson(data string) error {
	err := json.Unmarshal([]byte(data), p)
	if err != nil {
		logs.Printf("ProcStats: Json Unmarshall error %v", err)
	}
	return err
}

//NewProc Constructor for a PROC struct, this does not sync anything to the DB
func NewTask(TskName string, Capacity WLSpec, Type string, SlaveOf string) *Task {

	var tmp Task
	iname, id := TaskSplitNames(TskName)
	tmp.IName = iname
	tmp.ID = id
	tmp.Name = TskName
	tmp.Stats.Capacity.Copy(Capacity)

	tmp.Nodename = etcd.ETCD_INSTDIR + "/" + tmp.IName + "/Procs/" + tmp.ID
	return &tmp
}

func LoadTask(TskName string) *Task {

	var tmp Task

	iname, id := TaskSplitNames(TskName)
	tmp.IName = iname
	tmp.ID = id
	tmp.Name = TskName
	tmp.Nodename = etcd.ETCD_INSTDIR + "/" + tmp.IName + "/Procs/" + tmp.ID

	tmp.Load()

	return &tmp
}

//SplitNames Will return InstanceName and TaskName seperately
func TaskSplitNames(TskName string) (IName, TName string) {

	Tids := strings.Split(TskName, "::")

	if len(Tids) != 2 {
		logs.Printf("Proc.Load() Wrong format Task Name %s", TskName)
		return "", ""
	}

	return Tids[0], Tids[1]
}

func (T *Task) SyncStats() bool {

	if Gdb.IsSetup() != true {
		return false
	}
	Gdb.CreateSection(T.Nodename)
	Gdb.Set(T.Nodename+"/Stats", T.Stats.ToJson())
	return true
}

func (T *Task) GetDBKey(Key string) string {

	str, err := Gdb.Get(T.Nodename + "/" + Key)
	if err != nil {
		logs.Printf("TASK: Unable to read Key %s from the Store err=%v", Key, err)
		return ""
	}
	return str
}

func (T *Task) Sync() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	Gdb.CreateSection(T.Nodename)

	Gdb.Set(T.Nodename+"/Instance", T.IName)
	Gdb.Set(T.Nodename+"/ID", T.ID)
	Gdb.Set(T.Nodename+"/Pid", fmt.Sprintf("%d", T.Pid))
	Gdb.Set(T.Nodename+"/Stats", T.Stats.ToJson())
	Gdb.Set(T.Nodename+"/State", T.State)
	Gdb.Set(T.Nodename+"/SlaveOf", T.SlaveOf)
	Gdb.Set(T.Nodename+"/EID", T.EID)
	Gdb.Set(T.Nodename+"/SID", T.SID)
	Gdb.Set(T.Nodename+"/Type", T.Type)

	return true

}

//Load the latest from KV store
func (T *Task) Load() bool {

	if Gdb.IsSetup() != true {
		return false
	}
	T.Name = T.GetDBKey("Instance")
	T.ID = T.GetDBKey("ID")
	T.State = T.GetDBKey("State")
	T.SlaveOf = T.GetDBKey("SlaveOf")
	T.EID = T.GetDBKey("EID")
	T.SID = T.GetDBKey("SID")
	T.Type = T.GetDBKey("Type")
	T.Stats.FromJson(T.GetDBKey("Stats"))

	return true
}
