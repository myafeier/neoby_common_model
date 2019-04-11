package model

import (
	"log"
)

type Youpiao struct {
	Id int64
	Name string `xorm:"varchar(100)"`
	Price int64 `xorm:"default 0"`  //单位：元
	Info string `xorm:"varchar(255)" default ''`  //油票说明
	PriceCent int64 `xorm:"-"`  //转换成分
}

func (self *Youpiao)GetOne()(has bool,err error){
	has,err=LocalDB.Table("jyb_yprules").Get(self)
	if err!=nil {
	    log.Println("[Error] ",err)
	    return
	}
	if has{
		self.PriceCent=self.Price*100
	}
	return
}