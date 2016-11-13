package sched

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/huawei-cloudfederation/mesos-go-stateful/common/logs"
	typ "github.com/huawei-cloudfederation/mesos-go-stateful/common/types"
	"github.com/huawei-cloudfederation/mesos-go-stateful/sched/httplib"
	"github.com/huawei-cloudfederation/mesos-go-stateful/sched/mesoslib"
)

//Declare all the Constants to be used in this file
const (
	//HTTP_SERVER_PORT Rest server of the scheduler by default
	HTTP_SERVER_PORT = "8080"
)

func ParseConfig(cfgFileName string) error {

	typ.Cfg = typ.NewDefaultConfig()

	cfgFile, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		logs.Printf("Error Reading the configration file. Resorting to default values")
		return err
	}
	err = json.Unmarshal(cfgFile, typ.Cfg)
	if err != nil {
		logs.FatalInfo("Error parsing the config file %v", err)
		return err
	}

	logs.Printf("Configuration file is = %v", *typ.Cfg)

	return nil
}

func Register(S typ.StateFul) {
	//Assign the new Statful custom scheduler to global variable so that we can call its function when we need it
	typ.CustomFW = S
}

func Init(confName string) error {

	//ParseConfig File
	err := ParseConfig(confName)
	if err != nil {

		logs.Printf("Parse error terminating scheduler %v", err)
		return err
	}

	logs.Printf("*****************************************************************")
	logs.Printf("*********************Starting Scheduler******************")
	logs.Printf("*****************************************************************")

	//Facility to overwrite the etcd endpoint for scheduler if its running in the same docker container and expose a different one for executors

	dbEndpoint := os.Getenv("ETCD_LOCAL_ENDPOINT")

	if dbEndpoint == "" {
		dbEndpoint = typ.Cfg.DBEndPoint
	}

	//Initalize the common entities like store, store configuration etc.
	isInit, err := typ.Initialize(typ.Cfg.DBType, dbEndpoint)
	if err != nil || isInit != true {
		logs.FatalInfo("Failed to intialize Error:%v return %v", err, isInit)
		return err
	}

	logs.Printf("Configuration file is = %v", *typ.Cfg)
	//Start the Mesos library
	go mesoslib.Run()

	//start http server
	httplib.Run(typ.Cfg.HTTPPort)

	logs.Printf("*****************************************************************")
	logs.Printf("*********************Finished Workload-Scheduler******************")
	logs.Printf("*****************************************************************")

	return nil

}
