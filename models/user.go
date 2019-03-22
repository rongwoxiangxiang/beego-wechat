package models

import (
	"github.com/astaxie/beego/orm"
	"gowe/common"
	"time"
)

type User struct {
	Id    int `orm:"pk;auto"`
	User_name  string `orm:"size(64);default('')" json:"user_name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (u *User) TableName() string {
	return "users"
}

// 多字段索引
//func (u *User) TableIndex() [][]string {
//	return [][]string{
//		[]string{"Id", "Name"},
//	}
//}

// 设置引擎为 INNODB
func (u *User) TableEngine() string {
	return "INNODB"
}

// 多字段唯一键
func (u *User) TableUnique() [][]string {
	return [][]string{
		[]string{"Name", "Email"},
	}
}

func (us User) Find() (User, error){
	err := GetDB().Read(&us)
	if err != nil{
		return us,common.ErrUserNoExist
	}
	return us,nil
}

func (us User) DeleteById() bool{
	_, err := GetDB().Delete(&User{Id: us.Id})
	if err != nil{
		return false
	}
	return true
}

func init() {
	orm.RegisterModel(new(User))
}