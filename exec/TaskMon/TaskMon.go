package TaskMon

import (
	"fmt"
	"io"
	"log"
//		"os"
	"encoding/json"
//	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	typ "../../common/types"
	"../docker"
)

//Stats structure is to populate docker stats
type StatsInfo struct {
	StatsTime    string  `json:"read"`
	Network      statnet `json:"network"`
	CStats       cstat   `json:"cpu_stats"`
	MStats       mstat   `json:"memory_stats"`
	BlockIOStats bstat   `json:"blockio_stats"`
}

type statnet struct {
	RxBytes   int64 `json:"rx_bytes"`
	RxPackets int64 `json:"rx_packets"`
	RxErrors  int   `json:"rx_errors"`
	RxDropped int   `json:"rx_dropped"`
	TxBytes   int64 `json:"tx_bytes"`
	TxPackets int64 `json:"tx_packets"`
	TxErrors  int   `json:"tx_errors"`
	TxDropped int   `json:"tx_dropped"`
}

type cstat struct {
	CpuUsage       usage      `json:"cpu_usage"`
	SCpuUsage      int64      `json:"system_cpu_usage"`
	ThrottlingData throttling `json:"throttling_data"`
}

type usage struct {
	TotalUsage        int    `json:"total_usage"`
	PerCpuUsage       string `json:"percpu_usage"`
	UsageInKernelMode int    `json:"usage_in_kernel_mode"`
	UsageInUserMode   int    `json:"usage_in_user_mode"`
}

type throttling struct {
	Periods          int `json:"periods"`
	ThrottledPeriods int `json:"throttled_periods"`
	ThrottledTime    int `json:"throttled_time"`
}

type mstat struct {
	Usage    int64         `json:"usage"`
	MaxUsage int64         `json:"max_usage"`
	Stat st `json:"stats"`
	FailCnt  int           `json:"failcnt"`
	Limit    int           `json:"limit"`
}

type bstat struct {
	IOServiceBytesRecursive []string `json:"io_service_bytes_recursive"`
	IOServiceRecursive      []string `json:"io_serviced_recursive"`
	IOQueueRecursive        []string `json:"io_queue_recursive"`
	IOServiceTimeRecursive  []string `json:"io_service_time_recursive"`
	IOWaitTimeRecursive     []string `json:"io_wait_time_recursive"`
	IOMergeTimeRecursive    []string `json:"io_merged_recursive"`
	IOTimeRecursive         []string `json:"io_time_recursive"`
	SectorRecursive         []string `json:"sectors_recursive"`
}
type st struct{}

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
	MS_Sync   bool      //Make this as master after sync
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

func (T *TaskMon) getStats() (StatsInfo,error) {

	var err error
	var data StatsInfo

	resp, err := http.Get("http://" + T.IP + fmt.Sprintf(":%d", T.Port) + "/containers/" + T.Container.ID + "/stats")
	if err != nil {
		log.Printf("docker container error\n", err)
	}
	defer resp.Body.Close()

	body := io.Reader(resp.Body)

	if err := json.NewDecoder(body).Decode(&data); err != nil {
		log.Printf("I am here Json Unmarshall error = %v", err)
	}
	fmt.Println(data)

	return data,nil
}

func (T *TaskMon) launchWorkload(isSlave bool, IP string, port string) bool {

	var err error
	if isSlave {
		err = T.Container.Run(T.P.ID, T.Image, []string{"server", fmt.Sprintf("--port %d", T.Port), fmt.Sprintf("--Slaveof %s %s", IP, port)}, int64(T.P.MemCap), T.P.ID+".log")
	} else {
		err = T.Container.Run(T.P.ID, T.Image, []string{"server", fmt.Sprintf("--port %d", T.Port)}, int64(T.P.MemCap), T.P.ID+".log")
//				err = T.Container.Run("test", T.Image, []string{}, int64(1), "test.log")
	}

	if err != nil {
		//Print some error
		return false
	}

	//hack otherwise its too quick to have the server receiving connections
	time.Sleep(time.Second)

	//get the connected client immediately after for monitoring and other functions
//	_,err = T.getStats()
//	if err != nil {
		//Print some error
//		return false
//	}

	return true
}

//UpdateStats Update the stats structure and flush it to the Store/DB
func (T *TaskMon) UpdateStats() bool {

        var workloadStats typ.Stats
        //var err error


	data,_  := T.getStatsInfo()

	fmt.Println(data.CStats.CpuUsage.TotalUsage)

        worklaodStats.RxBytes = data.Network.RxBytes
        worklaodStats.CpuTotalUsage =  data.CStats.CpuUsage.TotalUsage
        worklaodStats.MemoryUsage =  data.MStats.Usage 
        worklaodStats.BlockIOStats =   data.BlockIOStats.IOServiceBytesRecursive

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
/*	for !T.IsSyncComplete() {
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
				T.getStats()
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
	case T.P.Msg == "MASTER":
		T.MakeMaster()
	}
	//Once you have read the message delete the message.
	T.P.Msg = ""
	T.P.SyncMsg()

}
