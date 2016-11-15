package httplib

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"

	"github.com/huawei-cloudfederation/mesos-go-stateful/common/logs"
	typ "github.com/huawei-cloudfederation/mesos-go-stateful/common/types"
)

//MainController of the HTTP server
type MainController struct {
	beego.Controller
}

//CreateInstance Handles a Create Instance
func (this *MainController) CreateInstance() {

	//Parse the input URL
	var name string

	name = this.Ctx.Input.Param(":INSTANCENAME")
	slaves, _ := strconv.Atoi(this.Ctx.Input.Param(":SLAVES"))
	logs.Printf("HTTP: CREATE request for instance %v", name)

	//Check the Cache if we have the Instance already
	//If intance is already available return with Error
	if typ.MemDb.IsValid(name) {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusConflict)
		this.Ctx.WriteString(fmt.Sprintf("Instance %v already exists, Cannot the created again", name))
		logs.Printf("Instance %v already there, Cannot the created again", name)
		return
	}

	dataBytes := this.Ctx.Input.RequestBody
	ToCreator := typ.NewHttpToCR(name, slaves, string(dataBytes))
	typ.Cchan <- ToCreator
	logs.Printf("HTTP-To-CREATOR %v  Sent", ToCreator)

	//Instance is Unavailable Supply the information to CREATE go-routine to be converted as OFFERS

	//Return with CREATED HTTP code
	this.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	this.Ctx.WriteString(fmt.Sprintf("Request Accepted, %s Instance will be created", name))
	logs.Printf("Request Accepted, %s Instance will be created", name)
}

//DeleteInstance handles a delete instance REST call
func (this *MainController) DeleteInstance() {

	var name string
	//Parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instance
	this.Ctx.ResponseWriter.WriteHeader(200)
	this.Ctx.WriteString(fmt.Sprintf("Request Placed for destroying %s instance", name))

}

//Status handles a STATUS REST call
func (this *MainController) StatusOfInstance() {

	//Parse the input URL
	//var name string
	var name string

	//Parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instance

	this.Ctx.WriteString(fmt.Sprintf("jsoninfo is empty for the instance %s", name))
}

//StatusAll handles StatusAll REST call
func (this *MainController) ListAllInstances() {

	this.Ctx.WriteString("jsoninfo is empty for all the instances")

}

//UpdateSlaves handles AddSlaves REST call
func (this *MainController) AddSlaves() {

	//var name string
	var name string

	//parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instance

	this.Ctx.WriteString(fmt.Sprintf("Adding the instance slaves for the instance %s", name))

}

//Run main function that starts the HTTP server
func Run(config string) {

	logs.Printf("Starting the HTTP server at port %s", config)
	beego.Run(":" + config)
}
