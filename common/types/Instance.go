package types

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/huawei-cloudfederation/mesos-go-stateful/common/logs"
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/store/etcd"
)

//A Instance structure that will be able to store a tree of data, Everything related to a workload intance
type Instance struct {
	Name       string           //Name of the instance
	Type       string           //Type of the instance "Single Instance = S; Master-Slave  = MS; Cluster = C"
	Capacity   int              //Capacity of the Instance in MB
	Masters    int              //Number of masters in this Instance
	Slaves     int              //Number of slaves in this Instance
	ExpMasters int              //Expected number of Masters
	ExpSlaves  int              //Expected number of Slaves
	Status     string           //Status of this instance "CREATING/RUNNING/DISABLED"
	Mname      string           //Name / task id of the master workload proc
	Snames     []string         //Name of the slave
	Spec       WLSpec           //Workload Specification (of each runnign process beloging to this Instance)
	DValue     int              //Distribution Value of WorkLoads (How to spread the workload)
	Procs      map[string]*Proc //An array of workload procs to be filled later
}

// NewInstance Creates a new instance variable
// Fills up the structure and updates the central store
// Returns an instance pointer
// Returns nil if the instance already exists
func NewInstance(Name string, Masters int, Slaves int, S WLSpec) *Instance {

	p := &Instance{Name: Name, ExpMasters: Masters, ExpSlaves: Slaves, Spec: S, DValue:2}
	return p
}

// LoadInstance Load an instance from the store using Instance Name from the store
// if the instance is unavailable then return nil
func LoadInstance(Name string) *Instance {

	if Gdb.IsSetup() != true {
		return nil
	}

	nodeName := etcd.ETCD_INSTDIR + "/" + Name

	if ok, _ := Gdb.IsKey(nodeName); !ok {
		return nil
	}

	I := &Instance{Name: Name}

	I.Load()

	return I

}

// Load Loads up the datastructure for the given Service Name to the struture
// If the Instance cannot be loaded the it returns an error
func (I *Instance) Load() bool {

	var err error
	var tmpStr string
	var SnamesKey []string

	if Gdb.IsSetup() != true {
		return false
	}

	nodeName := etcd.ETCD_INSTDIR + "/" + I.Name + "/"
	I.Type, err = Gdb.Get(nodeName + "Type")
	tmpStr, err = Gdb.Get(nodeName + "Capacity")
	I.Capacity, err = strconv.Atoi(tmpStr)
	tmpStr, err = Gdb.Get(nodeName + "Masters")
	I.Masters, err = strconv.Atoi(tmpStr)
	tmpStr, err = Gdb.Get(nodeName + "Slaves")
	I.Slaves, err = strconv.Atoi(tmpStr)
	tmpStr, err = Gdb.Get(nodeName + "ExpMasters")
	I.ExpMasters, err = strconv.Atoi(tmpStr)
	tmpStr, err = Gdb.Get(nodeName + "ExpSlaves")
	I.ExpSlaves, err = strconv.Atoi(tmpStr)
	I.Status, err = Gdb.Get(nodeName + "Status")
	I.Mname, err = Gdb.Get(nodeName + "Mname")
	SpecStr, err := Gdb.Get(nodeName + "Spec")
	I.Spec.FromJson(SpecStr)
	nodeNameSlaves := nodeName + "Snames/"
	SnamesKey, err = Gdb.ListSection(nodeNameSlaves, false)
	if err != nil {
		logs.Printf("The error value is %v", err)
	}

	for _, snamekey := range SnamesKey {
		_, sname := filepath.Split(snamekey)
		I.Snames = append(I.Snames, sname)
	}

	I.LoadProcs()

	return true
}

//Sync Writes the entier content of an instance into store, an instance could have many keys to be updated this is a write intensive function should be used carefully, do not call this if you are planning to update only a single attribute of an instance
func (I *Instance) Sync() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	nodeName := etcd.ETCD_INSTDIR + "/" + I.Name + "/"

	Gdb.CreateSection(nodeName)
	Gdb.Set(nodeName+"Type", I.Type)
	Gdb.Set(nodeName+"Masters", fmt.Sprintf("%d", I.Masters))
	Gdb.Set(nodeName+"Slaves", fmt.Sprintf("%d", I.Slaves))
	Gdb.Set(nodeName+"Capacity", fmt.Sprintf("%d", I.Capacity))
	Gdb.Set(nodeName+"ExpMasters", fmt.Sprintf("%d", I.ExpMasters))
	Gdb.Set(nodeName+"ExpSlaves", fmt.Sprintf("%d", I.ExpSlaves))
	Gdb.Set(nodeName+"Status", I.Status)
	Gdb.Set(nodeName+"Mname", I.Mname)
	Gdb.Set(nodeName+"Spec", I.Spec.ToJson())

	//Create Section for Slaves and Procs
	nodeNameSlaves := nodeName + "Snames/"

	Gdb.CreateSection(nodeNameSlaves)
	for _, sname := range I.Snames {
		Gdb.Set(nodeNameSlaves+sname, sname)
	}

	nodeNameProcs := nodeName + "Procs/"
	Gdb.CreateSection(nodeNameProcs)

	//for _, p := range I.Procs {
	//p.Sync()
	//}
	return true
}

//SyncType Write only the TYPE attribute to the DB/store
func (I *Instance) SyncType(string) bool {

	if Gdb.IsSetup() != true {
		return false
	}

	nodeName := etcd.ETCD_INSTDIR + "/" + I.Name + "/"
	Gdb.Set(nodeName+"Type", I.Type)
	return true
}

//SyncStatus Flushes only the status attribute to the DB
func (I *Instance) SyncStatus() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	nodeName := etcd.ETCD_INSTDIR + "/" + I.Name + "/"
	Gdb.Set(nodeName+"Status", I.Status)
	return true
}

//SyncSlaves Flushes only the Slaves attribute to the DB, used when a Slave died or promoted as a master
func (I *Instance) SyncSlaves() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	nodeName := etcd.ETCD_INSTDIR + "/" + I.Name + "/"
	Gdb.Set(nodeName+"Slaves", fmt.Sprintf("%d", I.Slaves))
	//Create Section for Slaves and Procs
	nodeNameSlaves := nodeName + "Snames/"

	Gdb.CreateSection(nodeNameSlaves)
	for _, sname := range I.Snames {
		Gdb.Set(nodeNameSlaves+sname, sname)
	}
	return true
}

//SyncMasters Flushes only the master attribute to the DB, used when a new workload master is choose.
func (I *Instance) SyncMasters() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	nodeName := etcd.ETCD_INSTDIR + "/" + I.Name + "/"
	Gdb.Set(nodeName+"Masters", fmt.Sprintf("%d", I.Masters))
	Gdb.Set(nodeName+"Mname", I.Mname)
	return true
}

//LoadProcs Should be called when all the PROCs need to be loaded to the lateest value, PS High DISK intensive function, should be used carefully
func (I *Instance) LoadProcs() bool {

	if I.Procs == nil {
		I.Procs = make(map[string]*Proc)
	}

	I.Procs[I.Mname] = LoadProc(I.Name + "::" + I.Mname)

	for _, n := range I.Snames {
		logs.Printf("Laoding proc key=%v ", n)
		I.Procs[n] = LoadProc(I.Name + "::" + n)
	}

	return true

}

//Instance_Json  Filtered elementes of an Instnace that will be sent as an HTTP response
type Instance_Json struct {
	Name     string
	Type     string
	Status   string
	Capacity int
	Master   *ProcJson
	Slaves   []*ProcJson
}

//ToJson_Obj Filtered elementes of an Instnace that will be sent as an HTTP response
func (I *Instance) ToJson_Obj() Instance_Json {

	var res Instance_Json
	res.Name = I.Name
	res.Type = I.Type
	res.Capacity = I.Capacity
	res.Status = I.Status

	if I.Status == INST_STATUS_RUNNING {
		var p *Proc
		p = I.Procs[I.Mname]
		res.Master = p.ToJson()
		for _, sname := range I.Snames {
			p = I.Procs[sname]
			res.Slaves = append(res.Slaves, p.ToJson())
		}
	}

	return res
}

//ToJson Marshall the Instane to a JSON
func (I *Instance) ToJson() string {

	var res Instance_Json
	res.Name = I.Name
	res.Type = I.Type
	res.Capacity = I.Capacity
	res.Status = I.Status

	if I.Status == INST_STATUS_RUNNING {
		var p *Proc
		p = I.Procs[I.Mname]
		res.Master = p.ToJson()
		res.Master.Port = p.Port
		for _, sname := range I.Snames {
			p = I.Procs[sname]
			res.Slaves = append(res.Slaves, p.ToJson())
		}
	}

	b, err := json.Marshal(res)

	if err != nil {
		return "Marshaling error"
	}

	return string(b)
}
