package models

import (
	"github.com/astaxie/beego/orm"
	"gowe/common"
	"time"
)

const (
	REPLY_TYPE_TEXT = "text"
	REPLY_TYPE_CODE = "code"
	REPLY_TYPE_LUCKY = "luck"
	REPLY_TYPE_CHECKIN = "checkin"
)

const PLEASE_TRY_AGAIN = "活动太火爆了，请稍后重试"

type Reply struct {
	Id  int64 `orm:"pk;auto"`
	Wid int64 `orm:"default(0);index"`
	ActivityId int64 `orm:"default(0)"`
	Alias string `orm:"size(200)"`
	ClickKey string `orm:"size(50)"`
	Success string `orm:"size(250)"`
	Fail string `orm:"size(250)"`
	Extra string
	Type string `orm:"size(50)"`
	Disabled int8 `orm:"default(2)"`
	Match int8 `orm:"default(1)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func (r *Reply) TableName() string {
	return "replies"
}

func (r Reply) Insert() (Reply, error){
	id, err := GetDB().Insert(&r)
	r.Id = id
	if err != nil{
		return Reply{},common.ErrDataCreate
	}
	return r,nil
}

func (r Reply) Update() (Reply, error){
	_, err := GetDB().Update(&r)
	if err != nil{
		return Reply{},common.ErrDataCreate
	}
	return r,nil
}

func (r Reply) Get() (Reply, error){
	if r.Id <= 0 {
		return Reply{},common.ErrDataGet
	}
	err := GetDB().Read(&r)
	if err != nil{
		return Reply{},common.ErrDataUnExist
	}
	return r,nil
}

func (r Reply) DeleteById() bool{
	_, err := GetDB().Delete(&Reply{Id: r.Id})
	if err != nil{
		return false
	}
	return true
}

/**
 * @Find
 * @Param Reply.Id int
 * @Param Reply.Alias string
 * @Param Reply.ClickKey string
 * @Success []Reply
 */
func (r Reply) FindOne() (replies Reply) {
	if "" == r.Alias && "" == r.ClickKey {
		return
	}
	qs := GetDB().QueryTable(r.TableName()).Filter("wid", r.Wid)
	if "" != r.Alias {
		qs = qs.Filter("alias", r.Alias)
	}
	if "" != r.ClickKey {
		qs = qs.Filter("click_key", r.ClickKey)
	}
	err := qs.Filter("disabled", common.NO_VALUE).One(&replies)
	if err != nil{
		return Reply{}
	}
	return replies
}

func (r Reply) LimitUnderWidList(pagesize int,pageno int) (replies []Reply) {
	qs := GetDB().QueryTable(r.TableName()).Filter("wid", r.Wid)
	cnt, err :=  qs.Limit(pagesize, (pageno-1)*pagesize).All(&replies)
	if err != nil || cnt < 1 {
		replies = []Reply{}
	}
	return replies
}

func init() {
	orm.RegisterModel(new(Reply))
}
