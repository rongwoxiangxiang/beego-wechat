package models

import (
	"github.com/astaxie/beego/orm"
	"gowe/common"
	"time"
)

type Activity struct {
	Id     int64 `orm:"pk;auto"`
	Wid int64 `orm:"default(0);index"`
	Name   string `orm:"size(50)"`
	Desc  string `orm:"size(200)"`
	Type int8
	events string `orm:"size(250)"`
	TimeStarted time.Time `orm:"null"`
	TimeEnd time.Time `orm:"null"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func (a *Activity) TableName() string {
	return "activities"
}

func (a Activity) Insert() (Activity, error){
	id, err := GetDB().Insert(&a)
	a.Id = id
	if err != nil{
		return Activity{},common.ErrDataCreate
	}
	return a,nil
}

func (a Activity) Update() (Activity, error){
	_, err := GetDB().Update(&a)
	if err != nil{
		return Activity{},common.ErrDataCreate
	}
	return a,nil
}

func (a Activity) Get() (Activity, error){
	if a.Id <= 0 {
		return Activity{},common.ErrDataGet
	}
	err := GetDB().Read(&a)
	if err != nil{
		return Activity{},common.ErrDataUnExist
	}
	return a,nil
}

func (a Activity) DeleteById() bool{
	_, err := GetDB().Delete(&Activity{Id: a.Id})
	if err != nil{
		return false
	}
	return true
}

func (a Activity) FindByWid() (activities []Activity) {
	cnt, err := GetDB().QueryTable(a.TableName()).Filter("wid", a.Wid).All(&activities)
	if err != nil || cnt < 1 {
		activities = []Activity{}
	}
	return activities
}

func (a Activity) LimitUnderWidList(pagesize int,pageno int) (activities []Activity) {
	qs := GetDB().QueryTable(a.TableName()).Filter("wid", a.Wid)
	cnt, err :=  qs.Limit(pagesize, (pageno-1)*pagesize).All(&activities)
	if err != nil || cnt < 1 {
		activities = []Activity{}
	}
	return activities
}

func init() {
	orm.RegisterModel(new(Activity))
}
