package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"gowe/common"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Response(code int, msg interface{}) {
	data, err := json.Marshal(&msg)
	if err != nil {
		this.Data["json"] = &map[string]interface{}{"code": common.ErrUnKnow.Code, "msg": common.ErrUnKnow.Msg}
	} else {
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": string(data)}
	}
	this.ServeJSON()
}

func (this *BaseController) ResponseStr(err error) {
	this.Data["json"] = err.Error()
	this.ServeJSON()
}

func init() {

}