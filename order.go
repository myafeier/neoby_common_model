package model

import (
	"github.com/go-xorm/xorm"
	"time"
)

//订单
type Order struct {
	Id              int64  `gorm:"primary_key"`
	Mid             int64  `json:"mid,omitempty" xorm:"'mid' int(11)" gorm:"column:mid;type:int(11)"`                                     // mit 用户id
	BossId          int64  `json:"bossid,omitempty" xorm:"'bossid' int(11)" gorm:"column:bossid;type:int(11)"`                            //bossid  机构id
	StoreId         int64  `json:"storeid,omitempty" xorm:"'storeid' int(11)" gorm:"column:storeid;type:int(11)"`                         //门店id
	CashierId       int64  `json:"cashierid,omitempty" xorm:"'cashierid' int(11)" gorm:"column:cashierid;type:int(11)"`                   //收银员id
	OutreachId      int64  `json:"outreachid,omitempty" xorm:"'outreachid' int(11)" gorm:"column:outreachid;type:int(11)"`                //外部推广机构id
	CompanyUserIds  string `json:"cmpyuserid,omitempty" xorm:"'cmpyuserid' char(100)" gorm:"column:cmpyuserid;type:char(100)"`            //公司业务人员
	RecommendId     int64  `json:"recommendid,omitempty" xorm:"'recommendid' int(11)" gorm:"column:recommendid;type:int(11)"`             //C端推荐人id
	RecommendBossId int64  `json:"recommendbossid,omitempty" xorm:"'recommendbossid' int(11)" gorm:"column:recommendbossid;type:int(11)"` //推荐商家id
	CouponsId       int64  `json:"couponsid,omitempty" xorm:"'couponsid' int(11)" gorm:"column:couponsid;type:int(11)"`                   //couponsid 机构卡券id
	PayStatus       int    `json:"paystatus,omitempty" xorm:"'paystatus' tinyint(1)" gorm:"column:paystatus;type:tinyint(1)"`             //paystatus 支付状态，1：未支付、2：已支付、-1：支付失败、-2：已退款
	AwardStatus     int    `json:"cashbackstatus,omitempty" xorm:"'cashbackstatus' int(1)" gorm:"column:cashbackstatus;type:int(1)"`      //cashbackstatus 奖励处理处理状态 0:未返现，1:已返现
	CompanyMoney    string `json:"cmpymoney,omitempty" xorm:"'cmpymoney' char(100)" gorm:"column:cmpymoney;type:char(100)"`               //cmpymoney char(100) 公司订单收入
	Profit          int64  `json:"profit,omitempty" xorm:"'profit' char(100)" gorm:"column:profit;type:char(100)"`                        //profit 毛利润总金额，单位：分
	UserCouponsId   string `json:"usercouponsid,omitempty" xorm:"'usercouponsid' char(100)" gorm:"column:usercouponsid;type:char(100)"`   //usercouponsid 用户卡券id
	PayMoney        int64  `json:"paymoney,omitempty" xorm:"'paymoney' char(100)" gorm:"column:paymoney;type:char(100)"`                  //paymoney 实付金额
	SettlementTime  int64  `json:"settlement_time,omitempty" xorm:"'settlement_time' int(11)" gorm:"column:settlement_time;type:int(11)"` //settlement_time 佣金结算时间
	ProfitId        int64  `json:"profit_id" xorm:"'profitid' int(11)" gorm:"column:profitid;type:int(11)"`
	UpdateTime      int64  `json:"update_time,omitempty" xorm:"'update_time' int(11)" gorm:"column:update_time;type:int(11)"` //update_time 更新时间
	CreateTime      int64  `json:"create_time,omitempty" xorm:"'create_time' int(11)" gorm:"column:create_time;type:int(11)"` //update_time 更新时间
}

func (self Order) TableName() string {
	return "jyb_order"
}

//获取一个月内退款的订单
func GetProfitedOrderWithinMonth(db *xorm.Engine) (orders []*Order, err error) {
	startTime := time.Now().AddDate(0, -1, 0)
	err = db.Table("jyb_order").Where("update_time>?", startTime.Unix()).Where("cashbackstatus=1").Where("paystatus=-2").Select("id").Find(&orders)
	return
}
