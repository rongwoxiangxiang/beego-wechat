package models

import (
	"github.com/astaxie/beego/orm"
	"gowe/common"
	"time"
)

const CHECK_FAIL = "签到失败，请重试！"

type Checkin struct {
	Id     int64 `orm:"pk;auto"`
	Wid int64 `orm:"default(0);index"`
	ActivityId int64 `orm:"default(0)"`
	Wuid int64 `orm:"default(0)"`
	Liner int64 `orm:"default(0)"`
	Total int64 `orm:"default(0)"`
	Lastcheckin time.Time `orm:"null"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func (c *Checkin) TableName() string {
	return "checkins"
}

func (c Checkin) GetCheckinByActivityWuid() (checkin Checkin){
	if c.ActivityId == 0 || c.Wuid == 0{
		return
	}
	GetDB().QueryTable(c.TableName()).
		Filter("activity_id", c.ActivityId).
		Filter("wuid", c.Wuid).
		One(&checkin)
	if checkin.Id == 0 {
		checkin.ActivityId = c.ActivityId
		checkin.Wuid = c.Wuid
		checkin.Lastcheckin = time.Now()
		checkin.Liner = 1
		checkin.Total = 1
		id, _ := GetDB().Insert(&checkin)
		if id > 0 {
			checkin.Id = id
		}
	}
	return checkin
}


func (a Checkin) Update() (Checkin, error){
	_, err := GetDB().Update(&a)
	if err != nil{
		return Checkin{},common.ErrDataCreate
	}
	return a,nil
}

func (a Checkin) Get() (Checkin, error){
	if a.Id <= 0 {
		return Checkin{},common.ErrDataGet
	}
	err := GetDB().Read(&a)
	if err != nil{
		return Checkin{},common.ErrDataUnExist
	}
	return a,nil
}


func init() {
	orm.RegisterModel(new(Checkin))
}
