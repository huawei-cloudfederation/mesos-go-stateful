package main

import (
	"os"
	"flag"
	"io/ioutil"
	"encoding/json"
	typ "../common/types"
	"../common/logs"
	"./httplib"
	"./mesoslib"
)

//Declare all the Constants to be used in this file
const (
	//HTTP_SERVER_PORT Rest server of the scheduler by default
	HTTP_SERVER_PORT = "8080"
)

func main() {

	cfgFileName := flag.String("config", "./config.json", "Supply the location of configuration file")
	dumpConfig := flag.Bool("DumpEmptyConfig", false, "Dump Empty Config file")
	flag.Parse()

	cfg := typ.NewDefaultConfig()

	if *dumpConfig == true {
		configBytes, err := json.MarshalIndent(cfg, " ", "  ")
		if err != nil {
			logs.Printf("Error marshalling the dummy config file. Exiting %v", err)
			return
		}
		logs.Printf("%s\n", string(configBytes))
		return
	}

	cfgFile, err := ioutil.ReadFile(*cfgFileName)

	if err != nil {
		logs.Printf("Error Reading the configration file. Resorting to default values")
	}
	err = json.Unmarshal(cfgFile, &cfg)
	if err != nil {
		logs.FatalInfo("Error parsing the config file %v", err)
	}
	logs.Printf("Configuration file is = %v", cfg)

	logs.Printf("*****************************************************************")
	logs.Printf("*********************Starting Workload-Scheduler******************")
	logs.Printf("*****************************************************************")
	//Command line argument parsing

	//Facility to overwrite the etcd endpoint for scheduler if its running in the same docker container and expose a different one for executors

	dbEndpoint := os.Getenv("ETCD_LOCAL_ENDPOINT")

	if dbEndpoint == "" {
		dbEndpoint = cfg.DBEndPoint
	}

	//Initalize the common entities like store, store configuration etc.
	isInit, err := typ.Initialize(cfg.DBType, dbEndpoint)
	if err != nil || isInit != true {
		logs.FatalInfo("Failed to intialize Error:%v return %v", err, isInit)
	}

	logs.Printf("Configuration file is = %v", cfg.WInfo.Image)
	//Start the Mesos library
	go mesoslib.Run(cfg.Master, cfg.ArtifactIP, cfg.ArtifactPort, cfg.ExecutorPath, cfg.WInfo.Image, cfg.DBType, cfg.DBEndPoint, cfg.FrameworkName, cfg.UserName)



	//start http server
	httplib.Run(cfg.HTTPPort)

	logs.Printf("*****************************************************************")
	logs.Printf("*********************Finished Workload-Scheduler******************")
	logs.Printf("*****************************************************************")

}
