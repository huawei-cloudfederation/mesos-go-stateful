package httplib

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"../../common/wlogs"
)

//MainController of the HTTP server
type MainController struct {
	beego.Controller
}

//CreateInstance Handles a Create Instance
func (this *MainController) CreateInstance() {

	//Parse the input URL

	var name string

	var data map[string]interface{}
	name = this.Ctx.Input.Param(":INSTANCENAME")

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	if err != nil {
		wlogs.Info("Cannot Unmarshal\n", err)
		return
	}

	this.Ctx.ResponseWriter.WriteHeader(201)
	this.Ctx.WriteString(wlogs.Infoln("Request Accepted, %s Instance will be created", name))

}

//DeleteInstance handles a delete instance REST call
func (this *MainController) DeleteInstance() {

	var name string
	//Parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instance

	this.Ctx.ResponseWriter.WriteHeader(200)
	this.Ctx.WriteString(wlogs.Infoln("Request Placed for destroying %s instance", name))

}

//Status handles a STATUS REST call
func (this *MainController) StatusOfInstance() {

	//Parse the input URL
	//var name string
	var name string

	//Parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instance

	this.Ctx.WriteString(wlogs.Infoln("jsoninfo is empty for the instance %s", name))
}

//StatusAll handles StatusAll REST call
func (this *MainController) ListAllInstances() {

	this.Ctx.WriteString("jsoninfo is empty for all the instances")

}

//UpdateSlaves Not yet implemented
func (this *MainController) AddSlaves() {

	//var name string
	var name string

	//parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instance

	this.Ctx.WriteString(wlogs.Infoln("Adding the instance slaves for the instance %s", name))

}

//Run main function that starts the HTTP server
func Run(config string) {

	wlogs.Info("Starting the HTTP server at port %s", config)
	beego.Run(":" + config)
}
