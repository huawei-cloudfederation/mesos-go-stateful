package TaskMon

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	typ "../../common/types"
	"../docker"
)


//TaskMon This structure is used to implement a monitor thread/goroutine for a running task
//This structure should be extended only if more functionality is required on the Monitoring functionality
//A task objec is created within this and monitored hence forth
type TaskMon struct {
	P         *typ.Proc //The task structure that should be used
	Pid       int       //The Pid of the running task
	IP        string    //IP address the task instance should bind to
	Port      int       //The port number of this task instance to be started
	Ofile     io.Writer //Stdout log file to be re-directed to this io.writer
	Efile     io.Writer //stderr of the task instance should be re-directed to this file
	MonChan   chan int
	Container *docker.Dcontainer //A handle for the Container package
	Image     string             //Name of the Image that should be pulled
	L         *log.Logger        //to redirect log outputs to a file
}

//NewTaskMon Create a new monitor based on the Data sent along with the TaskInfo
//The data could have the following details
func NewTaskMon(tskName string, IP string, Port int, data string, L *log.Logger, Image string) *TaskMon {

	var T TaskMon
	var P *typ.Proc

	T.MonChan = make(chan int)
	T.Port = Port
	T.IP = IP

	//ToDo does this need error handling
	T.L = L

	T.L.Printf("Split data received is %v\n", data)

	splitData := strings.Split(data, " ")
	if len(splitData) < 1 || len(splitData) > 4 {
		//Print an error this is not suppose to happen
		T.L.Printf("TaskMon Splitdata error %v\n", splitData)
		return nil
	}

	Cap, _ := strconv.Atoi(splitData[0])

	switch splitData[1] {
	case "Master":
		P = typ.NewProc(tskName, Cap, "M", "")
		T.L.Printf("created proc for new MASTER\n")
		break
	case "SlaveOf":
		P = typ.NewProc(tskName, Cap, "S", splitData[2])
		break
	}
	T.P = P
	//ToDo each instance should be started with its own dir and specified config file
	//ToDo Stdout file to be tskname.stdout
	//ToDo stderere file to be tskname.stderr
	T.Container = &docker.Dcontainer{}
	T.Image = Image

	return &T
}


func (T *TaskMon) launchWorkload(isSlave bool, IP string, port string) bool {

	var err error
	if isSlave {
		err = T.Container.Run(T.P.ID, T.Image, []string{"server", fmt.Sprintf("--port %d", T.Port), fmt.Sprintf("--Slaveof %s %s", IP, port)}, int64(T.P.MemCap), T.P.ID+".log")
	} else {
		err = T.Container.Run(T.P.ID, T.Image, []string{"server", fmt.Sprintf("--port %d", T.Port)}, int64(T.P.MemCap), T.P.ID+".log")
	}

	if err != nil {
		//Print some error
		return false
	}

	//hack otherwise its too quick to have the server receiving connections
	time.Sleep(time.Second)


	return true
}

//UpdateStats Update the stats structure and flush it to the Store/DB
func (T *TaskMon) UpdateStats() bool {


	var workloadStats typ.Stats

	data,err  := T.Container.GetStats()

	if err != nil {
		log.Println("GetStats error",err)
                return false
	}


        workloadStats.NRxbytes = data.Network.RxBytes
        workloadStats.CpuUsage =  data.CStats.CpuUsage.TotalUsage
        workloadStats.Mem =  data.MStats.Usage
        workloadStats.BlkIOstats =   data.BlockIOStats.IOServiceBytesRecursive

        errSync := T.P.SyncStats(workloadStats)
        if !errSync {
                T.L.Printf("Error syncing stats to store")
                return false
        }
        return true
}



//Start the workload be it Master or Slave
func (T *TaskMon) Start() bool {

	if T.P.SlaveOf == "" {
		return T.StartMaster()
	}

	return T.StartSlave()
}

//StartMaster Start the workload as a master
func (T *TaskMon) StartMaster() bool {

	var ret = false
	//Command Line
	ret = T.launchWorkload(false, "", "")
	if ret != true {
		return ret
	}

	T.Pid = 0
	T.P.Pid = 0
	T.P.Port = fmt.Sprintf("%d", T.Port)
	T.P.IP = T.IP
	T.P.State = "Running"
	T.P.Sync()

	return true
}

//StartSlave start the workload as a slave, should be called to point to a valid master
func (T *TaskMon) StartSlave() bool {
	var ret = false
	//Command Line
	slaveof := strings.Split(T.P.SlaveOf, ":")
	if len(slaveof) != 2 {
		T.L.Printf("Unacceptable SlaveOf value %vn", T.P.SlaveOf)
		return false
	}

	//Command Line
	ret = T.launchWorkload(true, slaveof[0], slaveof[1])
	if ret != true {
		return ret
	}

	T.Pid = 0
	T.P.Pid = 0
	T.P.Port = fmt.Sprintf("%d", T.Port)
	T.P.IP = T.IP
	T.P.State = "Running"

	T.P.Sync()

	return true
}

//Monitor Primary monitor thread started for every workload
func (T *TaskMon) Monitor() bool {

	//wait for a second for the server to start
	//ToDo: is it needed

	CheckMsgCh := time.After(time.Second)
	UpdateStatsCh := time.After(2 * time.Second)

	for {
		if T.P.State == "Running" {
			select {

			case <-T.MonChan:
				//ToDo:update state if needed
				//signal to stop monitoring this
				T.L.Printf("Stopping TaskMon for %s %s", T.P.IP, T.P.Port)
				return false

			case <-CheckMsgCh:
				//this is to check communication from scheduler; mesos messages are not reliable
				T.CheckMsg()
				CheckMsgCh = time.After(time.Second)

			case <-UpdateStatsCh:
				T.UpdateStats()
				UpdateStatsCh = time.After(2 * time.Second)

			}
		} else {
			time.Sleep(time.Second)
		}

	}

}

//Stop we have been told to stop the worklaod
func (T *TaskMon) Stop() bool {

	//send kill command for a graceful exit of the worklaod
	//the server exited graceful will reflect at the task status FINISHED

		errMsg := T.Die()
		if !errMsg { //message should be read by scheduler
			T.L.Printf("Killing the worklaod also did not work for  IP:%s and port:%d", T.IP, T.Port)
			return false
		}

	return true

}

//Die Kill the workload
func (T *TaskMon) Die() bool {
	//err := nil
	err := T.Container.Kill()
	if err != nil {
		T.L.Printf("Unable to kill the process %v", err)
		return false
	}

	//either the shutdown or a kill will stop the monitor also
	return true
}

//CheckMsg constantly keep checking if there is a new message for this workload
func (T *TaskMon) CheckMsg() {
	//check message from scheduler
	//currently we do it to see if scheduler asks us to quit

	//ToDo better error handling needed
	err := T.P.LoadMsg()
	if !err {
		T.L.Printf("Failed While Loading msg for workload %v from node %v", T.P.ID, T.P.Nodename)
		return
	}

	switch {
	case T.P.Msg == "SHUTDOWN":
		err = T.Stop()
		if err {

			T.L.Printf("failed to stop the server")
		}
		//in any case lets stop monitoring
		T.MonChan <- 1
		return
	}
	//Once you have read the message delete the message.
	T.P.Msg = ""
	T.P.SyncMsg()

}
