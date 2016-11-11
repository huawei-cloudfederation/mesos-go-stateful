package main

import (
	"flag"
	"log"
	"os"
	"text/template"
)

const SrcTemplate = `
//This is an autoGenerated code by the template

package main

import (
	"flag"

	"github.com/huawei-cloudfederation/mesos-go-stateful/common/logs"
	typ "github.com/huawei-cloudfederation/mesos-go-stateful/common/types"
	"github.com/huawei-cloudfederation/mesos-go-stateful/sched"
)

type {{.Name}}Scheduler struct {
	Name           string
	Master         string
	CPU, MEM, DISK int64
}

func (S *{{.Name}}Scheduler) Start(I *typ.Instance) {
	logs.Println("start scheduler")

}

func (S *{{.Name}}Scheduler) StartMaster(I *typ.Instance) {
	logs.Println("start  master")

}

func (S *{{.Name}}Scheduler) StartSlave(I *typ.Instance) {
	logs.Println("start slaves")

}

func (S *{{.Name}}Scheduler) MasterRunning(I *typ.Instance) {
	logs.Println(" master is running")

}

func (S *{{.Name}}Scheduler) SlaveRunning(I *typ.Instance) {
	logs.Println(" slave is running")

}

func (S *{{.Name}}Scheduler) MasterLost(I *typ.Instance) {
	logs.Println(" master has been loast")

}

func (S *{{.Name}}Scheduler) SlaveLost(I *typ.Instance) {
	logs.Println("slave  has been loast")

}

func New{{.Name}}Scheduler() *{{.Name}}Scheduler {
	return &{{.Name}}Scheduler{}
}

func main() {

	//Parse Config and command arguments
	ConfigFileName := flag.String("config", "../Config/config.json", "Location of the config file")
	flag.Parse()

	logs.Printf("Scheduler terminated")

	sched.Init(*ConfigFileName)

	logs.Printf("Scheduler terminated")

}
`

// Define a template.
// Prepare some data to insert into the template.
type Project struct {
	Name string
}

func main() {
	var err error

	Name := flag.String("name", "Example", "Name of the scheduler eg: RedisScheduler or MySQLScheduler")
	Path := flag.String("path", "./", "Path where the project needs to be created")
	flag.Parse()

	/* Defile Project Paths */
	ProjDir := *Path + "/" + *Name
	SchedDir := ProjDir + "/Scheduler"
	ExecDir := ProjDir + "/Executor"
	ConfDir := ProjDir + "/Config"

	// Create all the nessary directories for the project
	err = os.Mkdir(ProjDir, os.ModePerm|os.ModeDir)
	if err != nil {
		log.Printf("Unable to create dir %s%s", *Path, *Name)
		return
	}

	err = os.Mkdir(SchedDir, os.ModePerm|os.ModeDir)
	if err != nil {
		log.Printf("Unable to create dir %s%s", *Path, *Name)
		return
	}
	err = os.Mkdir(ExecDir, os.ModePerm|os.ModeDir)
	if err != nil {
		log.Printf("Unable to create dir %s%s", *Path, *Name)
		return
	}
	err = os.Mkdir(ConfDir, os.ModePerm|os.ModeDir)
	if err != nil {
		log.Printf("Unable to create dir %s%s", *Path, *Name)
		return
	}

	//Process the template
	workinfo := Project{Name: *Name}
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("SrcTemplate").Parse(SrcTemplate))

	// Execute the template for the scheduler
	f, err := os.Create(SchedDir + "/Scheduler.go")
	if err != nil {
		log.Printf("Unable to create scheduler file %v", err)
		return
	}
	err = t.Execute(f, workinfo)
	if err != nil {
		log.Println("executing template:", err)
	}
}
