package models

import (
	"github.com/astaxie/beego/orm"
	"gowe/common"
	"time"
)

const (
	PRIZE_LEVEL_DEFAULT = 0
)

type Prize struct {
	Id  int64 `orm:"pk;auto"`
	Wid int64 `orm:"default(0)"`
	ActivityId int64 `orm:"default(0)"`
	Code string `orm:"size(64)"`
	Level int8
	Used int8
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func (wu *Prize) TableName() string {
	return "prizes"
}

func (p Prize) FindOneUsedCode() (prize Prize) {
	GetDB().QueryTable(p.TableName()).
		Filter("activity_id", p.ActivityId).
		Filter("level", p.Level).
		Filter("used", common.NO_VALUE).
		One(&prize)
	//TODO 额外判断
	_, err := GetDB().QueryTable(p.TableName()).
		Filter("id", prize.Id).
		Filter("used", common.NO_VALUE).
		Update(orm.Params{"used": common.YES_VALUE})//乐观锁
	if err != nil {

	}
	return
}

func (p Prize) Update() (Prize, error){
	_, err := GetDB().Update(&p)
	if err != nil{
		return Prize{},common.ErrDataCreate
	}
	return p,nil
}

func (p Prize) Insert() (Prize, error){
	id, err := GetDB().Insert(&p)
	p.Id = id
	if err != nil{
		return Prize{},common.ErrDataCreate
	}
	return p,nil
}

func (p Prize) InsertBatch(prizes []Prize) (error){
	_, err := GetDB().InsertMulti(len(prizes), prizes)
	if err != nil{
		return common.ErrDataCreate
	}
	return nil
}

func (r Prize) LimitUnderActivityList(pagesize int,pageno int) (prizes []Prize) {
	qs := GetDB().QueryTable(r.TableName()).Filter("activity_id", r.ActivityId)
	cnt, err :=  qs.Limit(pagesize, (pageno-1)*pagesize).All(&prizes)
	if err != nil || cnt < 1 {
		prizes = []Prize{}
	}
	return prizes
}

//func (r Prize) count() int{
//
//}

func init() {
	orm.RegisterModel(new(Prize))
}
