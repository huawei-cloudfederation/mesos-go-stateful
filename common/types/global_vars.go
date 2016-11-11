package types

import (
	"container/list"

	"../agentstate"
	"../store"
)

var (
	Gdb   store.DB //Gdb Golabal variables related to db connection/instace
	MemDb *InMem   //In memory store

	OfferList *list.List        //list for having offer
	Cchan     chan TaskCreate   //Channel for Creator
	Mchan     chan *TaskUpdate  //Channel for Maintainer
	Dchan     chan TaskMsg      //Channel for Destroyer
	Agents    *agentstate.State //A Global View of agents and the Instnaces book keeping
	Wconfig   *Config           //A Global View of Config
)

//Global constants for Instance Status
//CREATING/ACTIVE/DELETED/DISABLED
const (
	INST_STATUS_CREATING = "CREATING"
	INST_STATUS_RUNNING  = "RUNNING"
	INST_STATUS_DISABLED = "DISABLED"
	INST_STATUS_DELETED  = "DELETED"
)

//Const for instance type
const (
	INST_TYPE_SINGLE       = "S"  //A Single instance server
	INST_TYPE_MASTER_SLAVE = "MS" //A workload instance with master-slave
)

//const for type of the server
const (
	PROC_TYPE_MASTER = "M"
	PROC_TYPE_SLAVE  = "S"
)
