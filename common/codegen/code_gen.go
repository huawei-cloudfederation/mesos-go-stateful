package main

import (
	"log"
	"os"
	"flag"
	"text/template"
)

func main() {
	var err error


flag.Usage = func() {
        fmt.Printf("Usage: ./code-gen workloadname\n")
    }
    flag.Parse()
	if flag.NArg() == 0 {
        flag.Usage()
        os.Exit(1)
    }
	// Define a template.
	const workload = `
package main 

import (
	"mesos-go-stateful"
	"mesos-go"
	"logs"
)

type DefaultStatefulScheduler struct {
	Name string
	Master string
	CPU, MEM, DISK int64
	JB JobLIST
}


func (S * {{.Name}}StatefulScheduler Start (Instance) {
	logs.Println("start scheduler")

}

func (S *{{.Name}}StatefulScheduler) StartMaster (Instance) {
	logs.Println("start  master")

}

func (S *{{.Name}}StatefulScheduler) StartSlave (Instance) {
	logs.Println("start slaves")

}

func (S *{{.Name}}StatefulScheduler) MasterRunning (Instance) {
	logs.Println(" master is running")

}

func (S *{{.Name}}StatefulScheduler) SlaveRunning (Instance) {
	logs.Println(" slave is running")

}

func (S *{{.Name}}StatefulScheduler) MasterLost (Instance) {
	logs.Println(" master has been loast")

}

func (S *{{.Name}}StatefulScheduler) SlaveLost (Instance) {
	logs.Println("slave  has been loast")


}

func main() {
	Stateful := NewDefaultStatefulScheduler()
	
	Sched := NewStatefulScheduler (DefaultStatefulScheduler)
	
	Sched.Start()
}

`

	// Prepare some data to insert into the template.
	type Workload struct {
		Name string
	}

	var workinfo = []Workload{
		{os.Args[1]},
	}

	// Create a new template and parse the letter into it.
	t := template.Must(template.New("workload").Parse(workload))

	// Execute the template for each workload.
	f , err := os.Create(os.Args[1]+"_workload"+".go")
	if err != nil {
		log.Print(err)
		return
	}
	
	for _, r := range workinfo {
	//	file, err1 := os.Create("file.go") // For read access.
		f, err = os.OpenFile(os.Args[1]+"_workload"+".go", os.O_WRONLY, 0777)
		if err != nil {
			log.Fatal(err)
		}

		err = t.Execute(f, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}

}

