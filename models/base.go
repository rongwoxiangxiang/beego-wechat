package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type BaseModel struct {
}

var db orm.Ormer
//数据库连接
func Connect() {
	var dsn string
	db_type := beego.AppConfig.String("db_type")
	db_host := beego.AppConfig.String("db_host")
	db_port := beego.AppConfig.String("db_port")
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_name := beego.AppConfig.String("db_name")

	orm.RegisterDriver("mysql", orm.DRMySQL)
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)

	orm.RegisterDataBase("default", db_type, dsn)

	db = orm.NewOrm()
}

func GetDB() orm.Ormer {
	orm.Debug = true
	return orm.NewOrm()
	//return db
}

//调用时传值
func Get(model interface{}) (interface{}, error){
	err := GetDB().Read(model)
	if err != nil{
		return model, nil
	}
	return model, nil
}
