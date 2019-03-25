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

    beego.Router("/wechat/", &controllers.WechatController{}, "get:List")
    beego.Router("/wechat/:id:int", &controllers.WechatController{}, "get:View")
	beego.Router("/wechat/:id:int", &controllers.WechatController{}, "delete:DeleteWechat")
	beego.Router("/wechat/:id:int", &controllers.WechatController{}, "put:UpdateWechat")
	beego.Router("/wechat", &controllers.WechatController{}, "post:InsertWechat")

	beego.Any("/service/:flag", controllers.Service)
}
