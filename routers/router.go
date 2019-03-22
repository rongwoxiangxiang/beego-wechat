package routers

import (
	"gowe/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{}, )
    beego.Router("/game/", &controllers.GameController{}, "get:List")
    beego.Router("/game/:id:int", &controllers.GameController{}, "get:View")
	beego.Router("/game/:id:int", &controllers.GameController{}, "delete:DeleteGame")
	beego.Router("/game/:id:int", &controllers.GameController{}, "put:UpdateGame")
	beego.Router("/game", &controllers.GameController{}, "post:InsertGame")
}
