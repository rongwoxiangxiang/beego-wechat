package models

import (
	"github.com/astaxie/beego/orm"
	"gowe/common"
	"time"
)

type Wechat struct {
	Id     int64 `orm:"pk;auto"`
	Gid int64 `orm:"default(0)"`
	Name   string
	Appid  string
	Appsecret string
	EncodingAesKey string
	Token string
	Flag  string
	Type  int8
	Pass  int8
	SaveInput int8
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func (w *Wechat) TableName() string {
	return "wechats"
}

func (u *Wechat) TableUnique() [][]string {
	return [][]string{
		[]string{"GameId", "Name"},
		[]string{"Appid"},
		[]string{"Flag"},
	}
}

func (w Wechat) Insert() (Wechat, error){
	id, err := GetDB().Insert(&w)
	w.Id = id
	if err != nil{
		return Wechat{},common.ErrDataCreate
	}
	return w,nil
}

func (w Wechat) Update() (Wechat, error){
	_, err := GetDB().Update(&w)
	if err != nil{
		return Wechat{},common.ErrDataCreate
	}
	return w,nil
}

func (w Wechat) Get() (Wechat, error){
	err := GetDB().Read(&w)
	if err != nil{
		return Wechat{},common.ErrDataUnExist
	}
	return w,nil
}

func (w Wechat) DeleteById() bool{
	_, err := GetDB().Delete(&Wechat{Id: w.Id})
	if err != nil{
		return false
	}
	return true
}

func (w Wechat) Find() (wechats []Wechat) {
	o := GetDB()
	qs := o.QueryTable(w.TableName())
	if w.Gid > 0 {
		qs = qs.Filter("gid", w.Gid)//orm????!!!!
	}
	if "" != w.Flag {
		qs = qs.Filter("flag", w.Flag)
	}
	cnt, err := qs.All(&wechats)
	if err != nil || cnt < 1 {
		wechats = []Wechat{}
	}
	return wechats
}

func (w Wechat) LimitList(pagesize int,pageno int) (wechats []Wechat) {
	qs := GetDB().QueryTable(w.TableName())
	cnt, err :=  qs.Limit(pagesize, (pageno-1)*pagesize).All(&wechats)
	if err != nil || cnt < 1 {
		wechats = []Wechat{}
	}
	return wechats
}

func init() {
	orm.RegisterModel(new(Wechat))
}