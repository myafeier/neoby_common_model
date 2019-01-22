package model

import "log"

//折扣券及折扣券规则，每天一次性读入内存
// jyb_coupons
type Coupon struct {
	Id int64 `gorm:"primary_key"`
	Rule string `xorm:"'rules' char(100)" gorm:"column:rules;type:char(100)"` //rules 卡券使用规则，可能是1对1关系。存的实际上是int型
	RuleF string `gorm:"-" xorm:"-"`
}

//折扣券规则，如果没有记录，则为"套餐分润"
// jyb_cpsrules
type Rule struct {
	Id int64 `gorm:"primary_key"`
	Name string `xorm:"char(100) default '' " gorm:"column:name;type:char(100)"` // name 规则名称
}

func GetACouponsWithRule()(coupons []*Coupon,err error){

	err=RemoteDB.Table("jyb_coupons").Select("id,rules").Find(&coupons)
	if err!=nil{
		return
	}
	for k,v:=range coupons{
		if v.Rule!=""{
			rule:=new(Rule)
			_,err=RemoteDB.Table("jyb_cpsrules").Select("id,name").Where("id = ?",v.Rule).Get(rule)
			if err!=nil&&err.Error()!="record not found"{
				log.Println("Error :",err)
				return
			}
			err=nil
			coupons[k].RuleF=rule.Name
		}
	}
	return
}

