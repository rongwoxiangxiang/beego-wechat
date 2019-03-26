package models

import (
	"github.com/astaxie/beego/orm"
	"gowe/common"
	"time"
)

type PrizeHistory struct {
	Id  int64 `orm:"pk;auto"`
	ActivityId int64 `orm:"default(0)"`
	Wuid int64 `orm:"default(0);index"`
	Prize string `orm:"size(64)"`
	Level int8
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func (wu *PrizeHistory) TableName() string {
	return "prize_history"
}

func (ph PrizeHistory) Insert() (PrizeHistory, error){
	id, err := GetDB().Insert(&ph)
	ph.Id = id
	if err != nil{
		return PrizeHistory{},common.ErrDataCreate
	}
	return ph,nil
}

func (wu PrizeHistory) GetByActivityWuId() (prizeHistory []PrizeHistory){
	cnt, err := GetDB().QueryTable(wu.TableName()).
		Filter("activity_id", wu.ActivityId).
		Filter("wuid", wu.Wuid).
		All(&prizeHistory)
	if err != nil || cnt < 1 {
		prizeHistory = []PrizeHistory{}
	}
	return prizeHistory
}

func (wu PrizeHistory) PrizeHistory() bool{
	_, err := GetDB().Delete(&PrizeHistory{Id: wu.Id})
	if err != nil{
		return false
	}
	return true
}


func (r PrizeHistory) LimitUnderList(pagesize int,pageno int) (users []PrizeHistory) {
	qs := GetDB().QueryTable(r.TableName())
	cnt, err :=  qs.Limit(pagesize, (pageno-1)*pagesize).All(&users)
	if err != nil || cnt < 1 {
		users = []PrizeHistory{}
	}
	return users
}

func init() {
	orm.RegisterModel(new(PrizeHistory))
}
