package main

import (
	"flag"
	"log"
	"os"
	"text/template"
)

const SrcTemplate = `
package main 

import (

	"fmt"
	"flag"

	"github.com/huawei-cloudfederation/mesos-go-stateful/sched"
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/logs"
)

type {{.Name}}Scheduler struct {
	Name string
	Master string
	CPU, MEM, DISK int64
	JB JobLIST
}


func (S * {{.Name}}Scheduler) Start (Instance) {
	logs.Println("start scheduler")

}

func (S *{{.Name}}Scheduler) StartMaster (Instance) {
	logs.Println("start  master")

}

func (S *{{.Name}}Scheduler) StartSlave (Instance) {
	logs.Println("start slaves")

}

func (S *{{.Name}}Scheduler) MasterRunning (Instance) {
	logs.Println(" master is running")

}

func (S *{{.Name}}Scheduler) SlaveRunning (Instance) {
	logs.Println(" slave is running")

}

func (S *{{.Name}}Scheduler) MasterLost (Instance) {
	logs.Println(" master has been loast")

}

func (S *{{.Name}}ulScheduler) SlaveLost (Instance) {
	logs.Println("slave  has been loast")


}

func main() {
	Stateful := NewDefaultStatefulScheduler()
	
	Sched := NewStatefulScheduler (DefaultStatefulScheduler)
	
	Sched.Start()
}
`

func main() {
	var err error

	Name := flag.String("Name", "example", "Name of the scheduler eg: RedisScheduler or MySQLScheduler")

	flag.Parse()
	// Define a template.
	// Prepare some data to insert into the template.
	type Workload struct {
		Name string
	}

	workinfo := Workload{Name: *Name}

	// Create a new template and parse the letter into it.
	t := template.Must(template.New("SrcTemplate").Parse(SrcTemplate))

	// Execute the template for the scheduler
	f, err := os.Create(*Name + "_scheduler" + ".go")
	if err != nil {
		log.Print(err)
		return
	}
	err = t.Execute(f, workinfo)
	if err != nil {
		log.Println("executing template:", err)
	}
}
