package model

import (
	"github.com/go-xorm/xorm"
	"time"
)

//订单
type Order struct {
	Id              int64   `json:"id"`
	Mid             int64   `json:"mid,omitempty" xorm:"'mid' int(11)" `                        // mit 用户id
	BossId          int64   `json:"bossid,omitempty" xorm:"'bossid' int(11)"`                   //bossid  机构id
	StoreId         int64   `json:"storeid,omitempty" xorm:"'storeid' int(11)"`                 //门店id
	CashierId       int64   `json:"cashierid,omitempty" xorm:"'cashierid' int(11)"`             //收银员id
	OutreachId      int64   `json:"outreachid,omitempty" xorm:"'outreachid' int(11)"`           //外部推广机构id
	CompanyUserIds  string  `json:"cmpyuserid,omitempty" xorm:"'cmpyuserid' char(100)"`         //公司业务人员
	RecommendId     int64   `json:"recommendid,omitempty" xorm:"'recommendid' int(11)"`         //C端推荐人id
	RecommendBossId int64   `json:"recommendbossid,omitempty" xorm:"'recommendbossid' int(11)"` //推荐商家id
	CouponsId       int64   `json:"couponsid,omitempty" xorm:"'couponsid' int(11)"`             //couponsid 机构卡券id
	GoodsId         int64   `json:"goodsid" xorm:"-"`                                           //商品id
	PayMoney        int64   `json:"paymoney,omitempty" xorm:"'paymoney' char(100)"`             //paymoney 实付金额
	Profit          int64   `json:"profit,omitempty" xorm:"'profit' char(100)"`                 //profit 毛利润总金额，单位：分
	CompanyMoney    string  `json:"cmpymoney,omitempty" xorm:"'cmpymoney' char(100)"`           //cmpymoney char(100) 公司订单收入
	PayStatus       int     `json:"paystatus,omitempty" xorm:"'paystatus' tinyint(1)"`          //paystatus 支付状态，1：未支付、2：已支付、-1：支付失败、-2：已退款
	UserCouponsId   string  `json:"usercouponsid,omitempty" xorm:"'usercouponsid' char(100)"`   //usercouponsid 用户卡券id
	AwardStatus     int     `json:"cashbackstatus,omitempty" xorm:"'cashbackstatus' int(1)"`    //cashbackstatus 奖励处理处理状态 0:未返现，1:已返现
	SettlementTime  int64   `json:"settlement_time,omitempty" xorm:"'settlement_time' int(11)"` //settlement_time 佣金结算时间
	ProfitId        int64   `json:"profit_id" xorm:"'profitid' int(11) default  0" `            // 分润id
	Money           int64   `json:"money,omitempty" xorm:"'money' default 0"`                   //`money` char(100) NOT NULL DEFAULT '' COMMENT '交易金额',
	Payfeerate      string  `json:"payfeerate,omitempty" xorm:"'payfeerate' default ''"`        //`payfeerate` char(100) NOT NULL DEFAULT '' COMMENT '支付手续费率',
	Paytype         int     `json:"paytype,omitempty" xorm:"'paytype' default 0"`               //`paytype` tinyint(1) NOT NULL DEFAULT '1' COMMENT '支付类型1微信2支付宝3本元快捷4本元卡',
	Channelid       int     `json:"channelid,omitempty" xorm:"'channelid' default 0"`           //`channelid` tinyint(1) NOT NULL DEFAULT '1' COMMENT '渠道id',
	Services        float64 `json:"services,omitempty" xorm:"'services' default 0"`             //`services` varchar(100) NOT NULL DEFAULT '' COMMENT '服务费',
	Iswh            int     `json:"iswh,omitempty" xorm:"'iswh' default 0"`                     //`iswh` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0-默认(滇惠通)，1-官渡-文惠，2-其他未知',
	Sxf             string  `json:"sxf,omitempty" xorm:"'sxf' default 0"`                       //`sxf` varchar(100) NOT NULL DEFAULT '' COMMENT '交易手续费',
	YPid            int64   `json:"ypid,omitempty" xorm:"'ypid' default 0"`                     // `ypid` int(11)
	UpdateTime      int64   `json:"update_time,omitempty" xorm:"'update_time' int(11)"`         //update_time 更新时间
	CreateTime      int64   `json:"create_time,omitempty" xorm:"'create_time' int(11)"`         //update_time 更新时间
}


func (self *Order) UpdateMatchId() (num int64, err error) {
	return LocalDB.ID(self.Id).Cols("match_cnpc_id").Update(self)
}

func (self *Order) Write() (err error) {
	
	return
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
