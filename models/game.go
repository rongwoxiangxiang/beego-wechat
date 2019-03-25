package models

import (
	"github.com/astaxie/beego/orm"
	"gowe/common"
	"time"
)

type Game struct {
	Id    int64 `orm:"pk;auto"`
	Name  string `orm:"size(64);default('');unique"`
	Desc  string `orm:"size(255);default('')"`
	Icon  string `orm:"size(255);default('')"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (g *Game) TableName() string {
	return "games"
}

func (g Game) Insert() (Game, error){
	id, err := GetDB().Insert(&g)
	g.Id = id
	if err != nil{
		return Game{},common.ErrDataCreate
	}
	return g,nil
}

func (g Game) Update() (Game, error){
	_, err := GetDB().Update(&g)
	if err != nil{
		return Game{},common.ErrDataCreate
	}
	return g,nil
}

//根据id查询Game{Id:id}
func (g Game) Get() (Game, error){
	err := GetDB().Read(&g)
	if err != nil{
		return Game{},common.ErrDataUnExist
	}
	return g,nil
}

func (g Game) DeleteById() bool{
	_, err := GetDB().Delete(&Game{Id: g.Id})
	if err != nil{
		return false
	}
	return true
}

func (g Game) LimitList(pagesize int,pageno int) (games []Game) {
	qs := GetDB().QueryTable(g.TableName())
	cnt, err :=  qs.Limit(pagesize, (pageno-1)*pagesize).All(&games)
	if err != nil || cnt < 1 {
		games = []Game{}
	}
	return games
}

func init() {
	orm.RegisterModel(new(Game))
}