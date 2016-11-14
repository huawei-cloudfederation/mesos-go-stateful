package types

import (
	"encoding/json"
)

//Specification of the worklaod
type WLSpec struct {
	CPU     float64 //Number of CPU each workload will require NOTE:an instance may have more than one workload
	Mem     float64 //Memory requirement of each workload  in MB
	Disk    float64 //Disk requirement of each workload in GB
	Network string  //Most preffered networkign layer HOST or DOckerBridge
	Image   string  // Docker Image of this workload
}

//ToJson convert to JSON string
func (W WLSpec) ToJson() string {

	b, e := json.Marshal(&W)
	if e != nil {
		return ""
	}
	return string(b)
}

//Read from a JSON structure
func (W WLSpec) FromJson(b string) error {

	e := json.Unmarshal([]byte(b), &W)
	if e != nil {
		return e
	}
	return nil
}

func (W *WLSpec) Copy(S WLSpec) {
	W.CPU = S.CPU
	W.Mem = S.Mem
	W.Disk = S.Disk
	W.Image = S.Image
	W.Network = S.Network
}

func (W *WLSpec) Default() {
	W.CPU = 1.0
	W.Mem = 100.00
	W.Disk = 1.0
	W.Image = "Stateful:latest"
	W.Network = "host"
}
