package TaskMon

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
	"os/exec"

	typ "../../common/types"
	"../docker"
)

//TaskMon This structure is used to implement a monitor thread/goroutine for a running task
//This structure should be extended only if more functionality is required on the Monitoring functionality
//A task objec is created within this and monitored hence forth
type TaskMon struct {
	P       *typ.Proc //The task structure that should be used
	Pid     int       //The Pid of the running task
	IP      string    //IP address the task instance should bind to
	Port    int       //The port number of this task instance to be started
	Ofile   io.Writer //Stdout log file to be re-directed to this io.writer
	Efile   io.Writer //stderr of the task instance should be re-directed to this file
	 MS_Sync bool      //Make this as master after sync
	MonChan chan int
	Container *docker.Dcontainer  //A handle for the Container package
	Image     string              //Name of the Image that should be pulled
	L         *log.Logger         //to redirect log outputs to a file
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

func (T *TaskMon) gethealthCheck() {
	
	//return client
}

func (T *TaskMon) launchWorkload(isSlave bool,IP string, port string) bool {

	var err error
        if isSlave {
                err = T.Container.Run(T.P.ID, T.Image, []string{"server", fmt.Sprintf("--port %d", T.Port), fmt.Sprintf("--Slaveof %s %s", IP, port)}, int64(T.P.MemCap), T.P.ID+".log")
        } else {
                err = T.Container.Run(T.P.ID, T.Image, []string{"server", fmt.Sprintf("--port %d", T.Port)}, int64(T.P.MemCap), T.P.ID+".log")
        //        err = T.Container.Run("test",T.Image, []string{},int64(1), "test.log")
        }

	fmt.Println(T.Container.ID)
        if err != nil {
                //Print some error
                return false
        }

        //hack otherwise its too quick to have the server receiving connections
        time.Sleep(time.Second)

	 //get the connected client immediately after for monitoring and other functions
       // T.Client = T.gethealthCheck()
        T.gethealthCheck()

	return true
}

//Start the workload be it Master or Slave
func (T *TaskMon) Start() bool {

	if T.P.SlaveOf == "" {
		return T.StartMaster()
	}

	if !T.MS_Sync {
		return T.StartSlave()
	}
	//Posibly a scale request so start it as a slave, sync then make as master
	return T.StartSlaveAndMakeMaster()

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

	//Monitor the worklaod to check if the sync is complete
	/*for !T.IsSyncComplete() {
		time.Sleep(time.Second)
	}*/
	T.Pid = 0
	T.P.Pid = 0
	T.P.Port = fmt.Sprintf("%d", T.Port)
	T.P.IP = T.IP
	T.P.State = "Running"

	T.P.Sync()

	return true
}

//StartSlaveAndMakeMaster Start is as a slave and make it as a master, should be done for replication or adding a new slave
func (T *TaskMon) StartSlaveAndMakeMaster() bool {
	var ret = false
	//Command Line
	slaveof := strings.Split(T.P.SlaveOf, ":")
	if len(slaveof) != 2 {
		T.L.Printf("Unacceptable SlaveOf value %vn", T.P.SlaveOf)
		return false
	}

	ret = T.launchWorkload(true, slaveof[0], slaveof[1])
	if ret != true {
		return ret
	}

	T.Pid = 0

	//Monitor the workload to check if the sync is complete
	/*for !T.IsSyncComplete() {
		time.Sleep(time.Second)
	}*/
	//Make this workload as master
	T.MakeMaster()

	T.Pid = 0
	T.P.Pid = 0
	T.P.Port = fmt.Sprintf("%d", T.Port)
	T.P.IP = T.IP
	T.P.State = "Running"
	T.P.Sync()

	return true
}

func fetchSubSection(value string, SubSection string) string {
	arr := strings.Split(value, "\r\n")

	for _, key := range arr {
		if strings.Contains(key, SubSection) {
			subArr := strings.Split(key, ":")
			if len(subArr) != 2 {
				return ""
			}
			return subArr[1]
		}
	}
	return ""
}

//GetWorkloadInfo Connect to the  Proc and collect info we need
/*func (T *TaskMon) GetWorkloadInfo(Section string, Subsection string) string {

	value, err := T.Client.Info(Section).Result()
	if err != nil {
		T.L.Printf("STATS collection returned error on IP:%s and PORT:%d Err:%v for section %s subsection %s", T.IP, T.Port, err, Section, Subsection)
		return ""
	}
	return fetchSubSection(value, Subsection)
}*/

//UpdateStats Update the stats structure and flush it to the Store/DB
/*func (T *TaskMon) UpdateStats() bool {

	var workloadStats typ.Stats
	var txt string
	var err error

	txt = T.GetWorklaodInfo("Memory", "used_memory")
	workloadStats.Mem, err = strconv.ParseInt(txt, 10, 64)
	if err != nil {
		T.L.Printf("UpdateStats(Mem) Unable to convert %s to number %v", txt, err)
	}

	txt = T.GetWorkloadInfo("Server", "uptime_in_seconds")
	worklaodStats.Uptime, err = strconv.ParseInt(txt, 10, 64)
	if err != nil {
		T.L.Printf("UpdateStats(Uptime) Uptime Unable to convert %s to number %v", txt, err)
	}

	txt = T.GetWorklaodInfo("Clients", "connected_clients")
	workloadStats.Clients, err = strconv.Atoi(txt)
	if err != nil {
		T.L.Printf("UpdateStats(Clients) Unable to convert %s to number %v", txt, err)
	}

	txt = T.GetWorklaodInfo("Replication", "master_last_io_seconds_ago")
	worklaodStats.LastSyced, err = strconv.Atoi(txt)
	if err != nil && txt != "" {
		T.L.Printf("UpdateStats(master_last_io) Unable to convert %s to number %v", txt, err)
	}

	txt = T.GetWorklaodInfo("Replication", "slave_repl_offset")
	worklaodStats.SlaveOffset, err = strconv.ParseInt(txt, 10, 64)
	if err != nil && txt != "" {
		T.L.Printf("UpdateStats(slave_repl_offset) Unable to convert %s to number %v", txt, err)
	}

	txt = T.GetWorklaodInfo("Replication", "slave_priority")
	worklaodStats.SlavePriority, err = strconv.Atoi(txt)
	if err != nil && txt != "" {
		T.L.Printf("UpdateStats(slave_priority) Unable to convert %s to number %v", txt, err)
	}

	errSync := T.P.SyncStats(workloadStats)
	if !errSync {
		T.L.Printf("Error syncing stats to store")
		return false
	}
	return true
}*/

//Monitor Primary monitor thread started for every workload
func (T *TaskMon) Monitor() bool {

	//wait for a second for the server to start
	//ToDo: is it needed

	CheckMsgCh := time.After(time.Second)
	//UpdateStatsCh := time.After(2 * time.Second)

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

/*			case <-UpdateStatsCh:
				T.UpdateStats()
				UpdateStatsCh = time.After(2 * time.Second)
*/			}
		} else {
			time.Sleep(time.Second)
		}

	}

}

//Stop we have been told to stop the worklaod 
/*func (T *TaskMon) Stop() bool {

	//send SHUTDOWN command for a graceful exit of the worklaod 
	//the server exited graceful will reflect at the task status FINISHED
	_, err := T.Client.Shutdown().Result()
	if err != nil {
		T.L.Printf("problem shutting down the worklaod at IP:%s and port:%d with error %v", T.IP, T.Port, err)

		//in this error case the scheduler will get a task killed notification
		//but will also see that the status it updated was SHUTDOWN, thus will handle it as OK

		errMsg := T.Die()
		if !errMsg { //message should be read by scheduler
			T.L.Printf("Killing the worklaod also did not work for  IP:%s and port:%d", T.IP, T.Port)
		}
		return false
	}

	return true

}*/

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
		/*err = T.Stop()
		if err {

			T.L.Printf("failed to stop the server")
		}*/
		//in any case lets stop monitoring
		T.MonChan <- 1
		return
	case T.P.Msg == "MASTER":
		T.MakeMaster()
	case strings.Contains(T.P.Msg, "SLAVEOF"):
		T.TargetNewMaster(T.P.Msg)
		//If this is the message then this particular workload will become slave of a different master
	}
	//Once you have read the message delete the message.
	T.P.Msg = ""
	T.P.SyncMsg()

}

//IsSyncComplete Should be called by the Monitors on Slave worklaod, this gives the boolien answer if the sync is completed or not
/*func (T *TaskMon) IsSyncComplete() bool {

	//time.Sleep(1 * time.Second)

	if T.Client == nil {
		return false
	}

	respStr, err := T.Client.Info("replication").Result()
	if err != nil {
		T.L.Printf("getting the repication stats from server at IP:%s and port:%d", T.IP, T.Port)
		//dont return but try next time in another second/.1 second
	}

	respArr := strings.Split(respStr, "\n")
	for _, resp := range respArr {
		T.L.Printf("resp = %v", resp)
		r := strings.Split(resp, ":")
		switch r[0] {
		case "role":
			if !strings.Contains(r[1], "slave") {
				T.L.Printf("Trying to call is sync, but this server is not really a slave IP:%s, port:%d", T.IP, T.Port)
				return false
			}
			continue
		case "master_sync_in_progress":
			if !strings.Contains(r[1], "0") {
				T.L.Printf("Sync not complete yet in slave IP:%s, port:%d", T.IP, T.Port)
				return false
			}
			return true
		case "master_sync_last_io_seconds_ago":
			//If the sync is completed then return true
			return true
		default:
			continue
		}

	}

	//if we did not find a master_sync_in_progress or slave in return at all, then some other problem, try again
	return false
}
*/
//MakeMaster Make a worklaod as a master (ie: supply the command "slaveof no on" to the worklaod 
func (T *TaskMon) MakeMaster() bool {

	//send the slaveof no one command to server
	/*_, err := T.Client.SlaveOf("no", "one").Result()
	if err != nil {
		T.L.Printf("Error turning the slave to Master at IP:%s and port:%d", T.IP, T.Port)
		return false
	}*/

	T.L.Printf("Slave of NO ONE worked")
	return true
}

//TargetNewMaster Make this worklaod now target a new master, should be done when a new slave is promoted
func (T *TaskMon) TargetNewMaster(Msg string) bool {

	SlaveofArry := strings.Split(Msg, " ") //Split it with space as while we are sending from the sheduler we send it of the format SLAVEOF<SPACE>IP<SPACE>PORT
	if len(SlaveofArry) != 3 {             //This should have three elements otherwise its an error

		T.L.Printf("Writing SLAVE of COMMAND %s", Msg)
		return false

	}

	//send the slaveof IP (Arry[1]) and PORT (Array[1])
	/*_, err := T.Client.SlaveOf(SlaveofArry[1], SlaveofArry[2]).Result()
	if err != nil {
		T.L.Printf("Error turning the slave to Master at IP:%s and port:%d", T.IP, T.Port)
		return false
	}*/

	T.L.Printf("Slave of %s %s worked", SlaveofArry[1], SlaveofArry[2])
	return true
}
