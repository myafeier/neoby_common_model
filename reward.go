package model

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"log"
	"time"
)

const (
	RewardStatNotYet   = 0  //未分润
	RewardStatOk       = 1  //已分润，已解冻
	RewardStatFreezing = 2  //已分润，冻结中
	RewardStatRefund   = -1 //退单后退分润
)

//奖励
type Reward struct {
	Id              int64  `json:"id"`
	OrderId         int64  `json:"orderid" xorm:"'orderid' default 0 index"`                   //orderid 订单ID
	Stat            int    `json:"type" xorm:"'type' default 0 index "`                        //奖励状态，0 未分润，1已分润、可提现，2，已分润，冻结中，-1：已退还
	BossId          int64  `json:"bossid" xorm:"'bossid' default 0 index"`                     //bossid
	RecommendBossId int64  `json:"recommendbossid" xorm:"'recommendbossid' default 0 index"`   //推荐商家id`
	RecommendId     int64  `json:"recommendid" xorm:"'recommendid' default 0 index"`           //recommendid C端推荐人id
	StoreId         int64  `json:"storeid" xorm:"'storeid' default 0 index"`                   //storeid 门店id
	OutreachId      int64  `json:"outreachid" xorm:"'outreachid' default 0 index"`             //outreachid 外部推广机构id
	CouponsId       int64  `json:"couponsid" xorm:"'couponsid' default 0 index"`               //couponsid 卡id
	GoodsId         int64  `json:"goodsid" xorm:"'goodsid' int(11) default 0 index"`           //商品id
	UserCouponsId   string `json:"usercouponsid" xorm:"'usercouponsid' char(100) default 0 "`  //usercouponsid 用户卡券id
	OrderMoney      int64  `json:"ordermoney" xorm:"'ordermoney' default 0"`                   //ordermoney 订单交易金额
	ProfitMoney     int64  `json:"profitmoney" xorm:"'profitmoney' default 0"`                 //profitmoney 订单毛利
	MzPrice         string `json:"mzprice" xorm:"'mzprice' char(100)"`                         //mzprice 满足条件
	Reward          string `json:"reward" xorm:"'reward' char(100)"`                           //reward 佣金比例
	RewardMoney     int64  `json:"rewardmoney" xorm:"'rewardmoney' default 0"`                 //rewardmoney 奖励佣金
	IsCompany       int    `json:"iscmpy" xorm:"'iscmpy' tinyint(1)"`                          //iscmpy 0 外部收入，1 公司收入
	CompanyUserId   string `json:"cmpyuserid" xorm:"'cmpyuserid' varchar(100)"`                //cmpyuserid 公司业务员id
	BigCustomerId   int64  `json:"bigcustid" xorm:"'bigcustid' default 0 "`                    // bigcustid 大客户ID
	SettlementTime  int64  `json:"settlement_time" xorm:"'settlement_time' int(11)"`           //settlement_time  结算时间
	UnfreezeTime    int64  `json:"unfreeze_time" xorm:"'unfreeze_time' int(11)"`               //冻结到的时间
	ProfitAmountId  int64  `json:"profit_amount_id" xorm:"'profit_amount_id' default 0 index"` //归入账户id
	CreateTime      int64  `json:"create_time" xorm:"int(11) 'create_time'"`                   //create_time 插入时间
	UpdateTime      int64  `json:"update_time" xorm:"int(11) 'update_time'"`                   //update_time 更新时间
}

func (self Reward) TableName() string {
	return "jyb_reward"
}

//退单退分润
func TRewardChargeBack(session *xorm.Session, reward *Reward) (err error) {

	err = tUpdateRewardStat(session, reward.Id, reward.Stat, RewardStatRefund)
	if err != nil {
		return
	}
	if reward.Stat == RewardStatFreezing { //从冻结状态扣除
		err = updateAccountAmount(session, reward.ProfitAmountId, ProfitActionReduceFreeze, reward.RewardMoney, 0, 0)
		if err != nil {
			return
		}

	} else if reward.Stat == RewardStatOk { //从解冻状态扣除
		err = updateAccountAmount(session, reward.ProfitAmountId, ProfitActionReduceFree, reward.RewardMoney, 0, 0)
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("wrong reward stat:%s", reward.Stat)
		return
	}
	return
}

//解冻分润
func TUnfreezeReward(session *xorm.Session, reward *Reward) (err error) {
	err = tUpdateRewardStat(session, reward.Id, reward.Stat, RewardStatOk)
	if err != nil {
		return
	}

	err = updateAccountAmount(session, reward.ProfitAmountId, ProfitActionUnFreeze, reward.RewardMoney, 0, 0)
	if err != nil {
		return
	}
	return
}

func tUpdateRewardStat(session *xorm.Session, id int64, fromStat, toStat int) (err error) {
	reward := new(Reward)
	has, err := session.ID(id).Where("`type`=?", fromStat).Select("id").ForUpdate().Get(&Reward{})
	if err != nil {
		log.Println(err)
		return
	}
	if !has {
		err = fmt.Errorf("lock reward:%d error! \n", id)
		session.Rollback()
		return
	}

	reward.Stat = toStat
	reward.UpdateTime = time.Now().Unix()
	aftNum, err := session.ID(id).Cols("type", "update_time").Update(reward)
	if err != nil {
		log.Println(err)
		return
	}
	if aftNum != 1 {
		err = fmt.Errorf("update reward:%d fail!\n", reward.Id)
		return
	}
	return
}
