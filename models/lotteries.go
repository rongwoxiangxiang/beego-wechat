package models

import (
	"github.com/astaxie/beego/orm"
	"gowe/common"
	"math/rand"
	"time"
)

const MAX_LUCKY_NUM = 10000

type Lottery struct {
	Id  int64 `orm:"pk;auto"`
	Wid int64 `orm:"default(0)"`
	ActivityId int64 `orm:"default(0)"`
	Name string `orm:"size(50)"`
	Desc string `orm:"size(200)"`
	TotalNum int64
	ClaimedNum int64
	Probability int
	FirstCodeId int64
	Level int8
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func (l *Lottery) TableName() string {
	return "lotteries"
}

func (l Lottery) Update() (Lottery, error){
	_, err := GetDB().Update(&l)
	if err != nil{
		return Lottery{},common.ErrDataCreate
	}
	return l,nil
}

func (l Lottery) Insert() (Lottery, error){
	id, err := GetDB().Insert(&l)
	l.Id = id
	if err != nil{
		return Lottery{},common.ErrDataCreate
	}
	return l,nil
}

func (l Lottery) List() (lotteries []Lottery) {
	if l.ActivityId == 0 || l.Wid == 0 {
		return
	}
	cnt, err := GetDB().QueryTable(l.TableName()).
		Filter("wid", l.Wid).
		Filter("activity_id", l.ActivityId).
		OrderBy("probability").
		All(&lotteries)
	if err != nil || cnt < 1 {
		lotteries = []Lottery{}
	}
	return lotteries
}

//抽奖
func (l Lottery) Luck() (lottery Lottery, err error) {
	lotteries := l.List()
	if len(lotteries) < 1 {
		err = common.ErrDataUnExist
		return
	}
	max := MAX_LUCKY_NUM
	actvityFinished := true
	for _, lot := range lotteries {
		if lot.ClaimedNum >= lot.TotalNum {//当前奖品发放完毕
			max -= lot.Probability
			continue
		}
		actvityFinished = false
		random := rand.Intn(max)
		if random <= lot.Probability {
			lottery = lot
			break
		}
		max -= lot.Probability
	}
	if actvityFinished {//全部奖品发放完毕，自动结束活动
		GetDB().QueryTable("replies").
			Filter("wid", l.Wid).
			Filter("activity_id", l.ActivityId).
			Update(orm.Params{"disabled": common.YES_VALUE})
		err = common.ErrLuckFinal
		return
	}

	if lottery.Id == 0 {
		err = common.ErrLuckFail
		return
	}
	_, err = GetDB().QueryTable(l.TableName()).
		Filter("id", lottery.Id).
		Filter("claimed_num", lottery.ClaimedNum).
		Update(orm.Params{"claimed_num": lottery.ClaimedNum+1})//乐观锁
	if err != nil {
		err = common.ErrDataUpdate
		return
	}
	return
}

func init() {
	orm.RegisterModel(new(Lottery))
}
