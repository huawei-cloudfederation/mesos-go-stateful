package types

//An interface which will be used as a base to refer
type StateFul interface {

	//Configure  the Instance if you want and return update command argument to be supplied with th task
	Config(I *Instance) string

	//Start the task generally
	Start(I *Instance) error

	//Start a Master Specifically
	StartMaster(I *Instance) error

	//Start a slave
	StartSlave(I *Instance) error

	//Handle status update of Master Running
	MasterRunning(I *Instance) error

	//Handle status update of slave runing
	SlaveRunning(I *Instance) error

	//Master is dead deal with it
	MasterLost(I *Instance) error

	//Slave is dead probably start a new one
	SlaveLost(I *Instance) error
}
