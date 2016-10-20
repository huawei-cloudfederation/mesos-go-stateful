package main

import (
	"./httplib"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
			log.Printf("Error marshalling the dummy config file. Exiting %v", err)
			return
		}
		fmt.Printf("%s\n", string(configBytes))
		return
	}

	cfgFile, err := ioutil.ReadFile(*cfgFileName)

	if err != nil {
		log.Printf("Error Reading the configration file. Resorting to default values")
	}
	err = json.Unmarshal(cfgFile, &Cfg)
	if err != nil {
		log.Fatalf("Error parsing the config file %v", err)
	}
	log.Printf("Configuration file is = %v", Cfg)

	//start http server
	httplib.Run(Cfg.HTTPPort)
}
