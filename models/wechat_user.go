package models

import (
	"github.com/astaxie/beego/orm"
	"gowe/common"
	"time"
)

type WechatUser struct {
	Id  int64 `orm:"pk;auto"`
	Wid int64 `orm:"default(0);index"`
	UserId   int64 `orm:"default(0)"`
	Openid  string `orm:"size(64)"`
	Nickname string `orm:"size(64)"`
	Sex int8
	Province string `orm:"size(20)"`
	City string `orm:"size(20)"`
	Country string `orm:"size(20)"`
	Language string `orm:"default(20)"`
	Headimgurl string `orm:"default(200)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func (wu *WechatUser) TableName() string {
	return "wechat_users"
}

func (wu WechatUser) Insert() (WechatUser, error){
	id, err := GetDB().Insert(&wu)
	wu.Id = id
	if err != nil{
		return WechatUser{},common.ErrDataCreate
	}
	return wu,nil
}

func (wu WechatUser) Update() (WechatUser, error){
	_, err := GetDB().Update(&wu)
	if err != nil{
		return WechatUser{},common.ErrDataCreate
	}
	return wu,nil
}

func (wu WechatUser) GetByOpenid() (WechatUser){
	if wu.Openid == "" || wu.Wid == 0{
		return WechatUser{}
	}
	if _, id, err := GetDB().ReadOrCreate(&wu, "openid","wid"); err == nil {
		wu.Id = id
		return wu
	}
	return WechatUser{}
}

func (wu WechatUser) DeleteById() bool{
	_, err := GetDB().Delete(&WechatUser{Id: wu.Id})
	if err != nil{
		return false
	}
	return true
}


func (r WechatUser) LimitUnderWidList(pagesize int,pageno int) (users []WechatUser) {
	qs := GetDB().QueryTable(r.TableName()).Filter("wid", r.Wid)
	cnt, err :=  qs.Limit(pagesize, (pageno-1)*pagesize).All(&users)
	if err != nil || cnt < 1 {
		users = []WechatUser{}
	}
	return users
}

func init() {
	orm.RegisterModel(new(WechatUser))
}
