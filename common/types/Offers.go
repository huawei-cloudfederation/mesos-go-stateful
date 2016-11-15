package types


//Offer Structure that is used between creator and Mesos Scheduler
type Offer struct {
	Name         string //Name of the instance
	Taskname     string //Name of the redis proc
	CmdInfo	     []string //What ever command line argument need to be supplied
	Spec 	     WLSpec //Workload specification
	DValue       int    //Distribution value of the workload
}

//NewOffer Returns a new offer which will be interpreted by the scheduler
func NewOffer(name string, tname string, C []string, dvalue int, S WLSpec) Offer {
	return Offer{Name: name, Taskname: tname, CmdInfo:C, DValue:dvalue, Spec:S}
}
