package routers

import (
	"github.com/astaxie/beego"
	"github.com/royburns/goTestLinkReport/controllers"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/plan", &controllers.PlanController{})
	beego.Router("/report", &controllers.ReportController{})
	beego.Router("/test", &controllers.TestController{})
}
