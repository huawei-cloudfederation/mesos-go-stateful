package httplib

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.BConfig.CopyRequestBody = true
	beego.Router("/v1/CREATE/:INSTANCENAME", &MainController{}, "post:CreateInstance")
	beego.Router("/v1/DELETE/:INSTANCENAME", &MainController{}, "delete:DeleteInstance")
	beego.Router("/v1/STATUS/:INSTANCENAME", &MainController{}, "get:StatusOfInstance")
	beego.Router("/v1/STATUS/", &MainController{}, "get:ListAllInstances")
	beego.Router("/v1/UPDATE/:INSTANCENAME/SLAVES/:SLAVES", &MainController{}, "put:AddSlaves")
}
