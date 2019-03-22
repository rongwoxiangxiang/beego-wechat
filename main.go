package main

import (
	"gowe/models"
	_ "gowe/routers"
	"github.com/astaxie/beego"
)

func main() {
	models.Connect()
	beego.Run()
}

