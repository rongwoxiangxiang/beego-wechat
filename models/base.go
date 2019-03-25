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

////更新用户
//func UpdateById(id int,table string,filed map[string] interface{})bool{
//	o := GetDB()
//	_, err := o.QueryTable(
//		table).Filter(
//		"Id", id).Update(
//		filed)
//	if err == nil{
//		return true
//	}
//	return false
//}
//
//
////根据用户ID查询用户
//func QueryById(uid int64) (User, bool){
//
//	o := GetDB()
//	u := User{ID: uid}
//
//	err := o.Read(&u)
//
//	if err == orm.ErrNoRows {
//		fmt.Println("查询不到")
//		return u,false
//	} else if err == orm.ErrMissPK {
//		fmt.Println("找不到主键")
//		return u,false
//	} else {
//		fmt.Println(u.Id, u.Name)
//		return u,true
//	}
//}
//
////根据用户名称查询用户
//func QueryByName(name string) (User, error) {
//	var user User
//
//	o := orm.NewOrm()
//	qs := o.QueryTable("user")
//
//	err := qs.Filter("Name", name).One(&user)
//	fmt.Println(err)
//	if err == nil {
//		fmt.Println(user.Name)
//		return user,nil
//	}
//	return user, err
//}
//
////根据用户数据列表
//func DataList() (users []User) {
//
//	o := orm.NewOrm()
//	qs := o.QueryTable("users")
//
//	var us []User
//	cnt, err :=  qs.Filter("id__gt", 0).OrderBy("-id").Limit(10, 0).All(&us)
//	if err == nil {
//		fmt.Printf("count", cnt)
//	}
//	return us
//}
//
////查询语句，sql语句的执行
////格式类似于:o.Raw("UPDATE user SET name = ? WHERE name = ?", "testing", "slene")
////
//func QueryBySql(sql string, qarms[] string) bool{
//
//	o := orm.NewOrm()
//
//	//执行sql语句
//	o.Raw(sql, qarms)
//
//	return true
//}
////根据用户分页数据列表
//func LimitList(pagesize int,pageno int) (users []User) {
//
//	o := orm.NewOrm()
//	qs := o.QueryTable("user")
//
//	var us []User
//	cnt, err :=  qs.Limit(pagesize, (pageno-1)*pagesize).All(&us)
//	if err == nil {
//		fmt.Printf("count", cnt)
//	}
//	return us
//}
////根据用户数据总个数
//func GetDataNum(tableName string) int64 {
//
//	o := orm.NewOrm()
//	qs := o.QueryTable(tableName)
//
//	var us []User
//	num, err :=  qs.Filter("id__gt", 0).All(&us)
//	if err == nil {
//		return num
//	}else{
//		return 0
//	}
//}