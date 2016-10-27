package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"../common/wlogs"
	"./httplib"
)

type Config struct {
	HTTPPort string //Defaults to 8080 if otherwise specify explicitly
}

//NewDefaultConfig Default Constructor to create a config file
func NewDefaultConfig() Config {
	return Config{
		HTTPPort: "5055",
	}
}

func main() {

	cfgFileName := flag.String("config", "./config.json", "Supply the location of configuration file")
	dumpConfig := flag.Bool("DumpEmptyConfig", false, "Dump Empty Config file")
	flag.Parse()

	Cfg := NewDefaultConfig()

	if *dumpConfig == true {
		configBytes, err := json.MarshalIndent(Cfg, " ", "  ")
		if err != nil {
			wlogs.Info("Error marshalling the dummy config file. Exiting %v", err)
			return
		}
		wlogs.Info("%s\n", string(configBytes))
		return
	}

	cfgFile, err := ioutil.ReadFile(*cfgFileName)

	if err != nil {
		wlogs.Info("Error Reading the configration file. Resorting to default values")
	}
	err = json.Unmarshal(cfgFile, &Cfg)
	if err != nil {
		wlogs.Fatal("Error parsing the config file %v", err)
	}
	wlogs.Info("Configuration file is = %v", Cfg)

	//start http server
	httplib.Run(Cfg.HTTPPort)
}
