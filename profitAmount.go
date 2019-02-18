package model

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"strings"
	"time"
)

const (
	ProfitActionNewAward     = "newToFreeze"  //新收入冻结
	ProfitActionUnFreeze     = "unfreeze"     //解冻
	ProfitActionWithdraw     = "widthdraw"    //提现
	ProfitActionReduceFreeze = "reduceFreeze" //扣除冻结
	ProfitActionReduceFree   = "reduceFree"   //扣除可用
)

const (
	ProfitAmountOfDhtId = 1 //滇汇通公司账户ID
)

//账户类型
const (
	ProfitAmountAccountTypeOfNeoby        = 1 //滇汇通账户
	ProfitAmountAccountTypeOfNeobySellers = 2 //滇汇通销售人员账户
	ProfitAmountAccountTypeOfOutreach     = 3 //合作推广机构
	ProfitAmountAccountTypeOfDoubleOilVIP = 4 //两油VIP
	ProfitAmountAccountTypeOfCustomer     = 5 //顾客账户
	ProfitAmountAccountTypeOfBoss         = 6 //机构
	ProfitAmountAccountTypeOfStore        = 7 //门店

)

// 用户账户数据结构
// swagger:model
type ProfitAmount struct {
	Id int64 `json:"id"`
	//账户类型
	Type int `json:"type" xorm:"tinyint(2) default 0 index"`
	// 商户ID
	BossId int64 `json:"boss_id,omitempty" xorm:"default 0 index"`
	//门店id
	StoreId int64 `json:"store_id,omitempty" xorm:"default 0 index"`
	//外部推广id
	OutreachId int64 `json:"outreach_id,omitempty" xorm:"default 0 index"`
	//会员id
	MemberId int64 `json:"member_id,omitempty" xorm:"default 0 index"`
	//vip id
	VipId int64 `json:"vip_id,omitempty" xorm:"default 0 index"`
	//销售人员id，多个用,分割
	SellerIds string `json:"-" xorm:"char(100) default '' index"`
	//是否是滇汇通公司
	IsCompany int `json:"-" xorm:"default 0 index"`
	//历史合计收益
	TotalAmount int64 `json:"total_amount" xorm:"default 0"`
	// 可提现收益
	FreeAmount int64 `json:"free_amount" xorm:"default 0"`
	// 冻结中收益
	FreezeAmount int64 `json:"freeze_amount" xorm:"default 0"`
	//已提现收益
	WithdrewAmount int64  `json:"withdrew_amount" xorm:"default 0"`
	StoreName      string `json:"store_name,omitempty" xorm:"-"`
}

func (self ProfitAmount) TableName() string {
	return "jyb_profit_amount"
}

type ProfitAmountLog struct {
	Id                   int64     `json:"id"`
	Action               string    `json:"action" xorm:"char(20) default ''  index"` //操作名称
	ProfitAmountId       int64     `json:"profit_amount_id" xorm:"default 0 index"`  //账户id
	Amount               int64     `json:"amount" xorm:"default 0"`                  //操作前历史合计收益
	UnfreezeTime         int64     `json:"unfreeze_time" xorm:"default 0 index"`     //预设解冻时间，只有新分配利润才有预设的解冻时间
	BeforeTotalAmount    int64     `json:"before_total_amount" xorm:"default 0"`     //操作前历史合计收益
	BeforeFreeAmount     int64     `json:"before_free_amount" xorm:"default 0"`      //操作前可提现收益
	BeforeFreezeAmount   int64     `json:"before_freeze_amount" xorm:"default 0"`    //操作前冻结中收益
	BeforeWithdrewAmount int64     `json:"before_withdrew_amount" xorm:"default 0"`  //操作前已提现收益
	Remark               string    `json:"remark" xorm:"varchar(255) default ''"`    //备注
	RewardId             int64     `json:"reward_id" xorm:"default 0"`               //奖励ID，如有
	WithdrawId           int64     `json:"withdraw_id" xorm:"default 0"`             // 提现id，如有
	Created              time.Time `json:"created" xorm:"created"`
	Updated              time.Time `json:"updated" xorm:"updated"`
}

func (self ProfitAmountLog) TableName() string {
	return "jyb_profit_amount_log"
}

// bossId 和storeId 必须成对出现,  如果有vipId 依赖于 boosId
// 如果是新增收益需同时提交解冻时间
func TChangeProfitAmount(session *xorm.Session, action string, reward *Reward, withDraw *Withdraw) (accountId int64, err error) {

	//检查账户存在性
	pa := new(ProfitAmount)
	var amount int64

	if reward != nil {
		//因分润涉及多种id组合形式，故确定账户Id需要根据组合进行判断
		if reward.IsCompany == 1 { //公司内部分润
			if sellerIds := strings.TrimSpace(reward.CompanyUserId); sellerIds == "" || sellerIds == "0" {
				pa.Id = ProfitAmountOfDhtId //公司账户
				pa.Type = ProfitAmountAccountTypeOfNeoby
				pa.IsCompany = 1
			} else {
				pa.IsCompany = 1
				pa.SellerIds = sellerIds
				pa.Type = ProfitAmountAccountTypeOfNeobySellers
			}
		} else if reward.OutreachId > 0 { //合作机构 账户
			pa.OutreachId = reward.OutreachId
			pa.Type = ProfitAmountAccountTypeOfOutreach
		} else if reward.BigCustomerId > 0 { //两油VIP账户
			pa.VipId = reward.BigCustomerId
			pa.BossId = reward.BossId
			pa.Type = ProfitAmountAccountTypeOfDoubleOilVIP
		} else if reward.RecommendId > 0 { //顾客用户
			pa.MemberId = reward.RecommendId
			pa.Type = ProfitAmountAccountTypeOfCustomer
		} else if reward.RecommendBossId > 0 {
			if reward.RecommendBossId == 1 { //历史bug，推荐id为1 的实际为公司收益
				pa.Id = ProfitAmountOfDhtId
			} else if reward.StoreId != reward.RecommendBossId { //非两油的商户
				pa.StoreId = reward.RecommendBossId
				pa.Type = ProfitAmountAccountTypeOfStore
			} else { //两油的油站
				pa.StoreId = reward.StoreId //自己推荐的自己
				pa.BossId = reward.BossId
				pa.Type = ProfitAmountAccountTypeOfStore
			}

		} else { //纯碎的商户/门店账户
			pa.StoreId = reward.StoreId
			pa.BossId = reward.BossId
			pa.Type = ProfitAmountAccountTypeOfStore
		}
		amount = reward.RewardMoney
	}

	if withDraw != nil {
		if withDraw.OutreachId != 0 {
			pa.OutreachId = withDraw.OutreachId
			pa.Type = ProfitAmountAccountTypeOfOutreach
		} else if withDraw.MemberId != 0 {
			pa.MemberId = withDraw.MemberId
			pa.Type = ProfitAmountAccountTypeOfCustomer
		} else if withDraw.MerchantId != 0 || withDraw.StoreId != 0 {
			if withDraw.WithDrawType == 0 { //门店
				pa.Type = ProfitAmountAccountTypeOfStore
			} else if withDraw.WithDrawType == 1 { //机构
				pa.Type = ProfitAmountAccountTypeOfBoss
			}

			pa.BossId = withDraw.MerchantId
			pa.StoreId = withDraw.StoreId
		}
		amount = withDraw.Money
	}

	has, err := session.Get(pa)
	if err != nil {
		return
	}
	if !has {
		//新建账户
		_, err = session.Insert(pa)
		if err != nil {
			return
		}
	}
	var unfreezeTime, linkId int64
	if reward != nil {
		unfreezeTime = reward.UnfreezeTime
		linkId = reward.Id
	}

	if withDraw != nil {
		linkId = withDraw.Id
	}

	accountId = pa.Id
	err = updateAccountAmount(session, accountId, action, amount, unfreezeTime, linkId)

	return
}

// 提现的简化实现
func TWithDrawWithFixAccountId(session *xorm.Session, accountId, amount, linkId int64) (err error) {
	return updateAccountAmount(session, accountId, ProfitActionWithdraw, amount, 0, linkId)
}

// 账户金额变动
func updateAccountAmount(session *xorm.Session, accountId int64, action string, amount int64, unfreezeTime int64, linkId int64) (err error) {

	//锁定
	npa := new(ProfitAmount)
	npa.Id = accountId
	_, err = session.ForUpdate().Get(npa)
	if err != nil {
		return
	}
	var affectNum int64

	paLog := new(ProfitAmountLog)
	paLog.Action = action
	paLog.Amount = amount
	paLog.BeforeFreeAmount = npa.FreeAmount
	paLog.BeforeFreezeAmount = npa.FreezeAmount
	paLog.BeforeTotalAmount = npa.TotalAmount
	paLog.BeforeWithdrewAmount = npa.WithdrewAmount
	paLog.UnfreezeTime = unfreezeTime
	if action == ProfitActionNewAward {
		paLog.RewardId = linkId
	}
	if action == ProfitActionWithdraw {
		paLog.WithdrawId = linkId
	}

	switch action {
	case ProfitActionNewAward:
		affectNum, err = session.ID(npa.Id).Cols("total_amount", "freeze_amount", ).Incr("total_amount", amount).Incr("freeze_amount", amount).Update(npa)
		if err != nil {
			return
		}
		if affectNum != 1 {
			err = fmt.Errorf("ProfitAmount did not changed")
			return
		}
	case ProfitActionUnFreeze: //解冻
		affectNum, err = session.ID(npa.Id).Cols("free_amount", "freeze_amount", ).Incr("free_amount", amount).Decr("freeze_amount", amount).Where("freeze_amount-?>=0", amount).Update(npa)
		if err != nil {
			return
		}
		if affectNum != 1 {
			err = fmt.Errorf("ProfitAmount did not changed")
			return
		}
	case ProfitActionWithdraw: //提现
		affectNum, err = session.ID(npa.Id).Cols("free_amount").Decr("free_amount", amount).Incr("withdrew_amount", amount).Where("free_amount-?>=0", amount).Update(npa)
		if err != nil {
			return
		}
		if affectNum != 1 {
			err = fmt.Errorf("ProfitAmount did not changed")
			return
		}
	case ProfitActionReduceFree: //从可用中撤销收入,可能导致账户出现负值
		affectNum, err = session.ID(npa.Id).Cols("free_amount", "total_amount", ).Decr("free_amount", amount).Decr("total_amount", amount).Update(npa)
		if err != nil {
			return
		}
		if affectNum != 1 {
			err = fmt.Errorf("ProfitAmount did not changed")
			return
		}
	case ProfitActionReduceFreeze: //从冻结中撤销收入
		affectNum, err = session.ID(npa.Id).Cols("total_amount", "freeze_amount", ).Decr("freeze_amount", amount).Decr("total_amount", amount).Where("freeze_amount-?>=0", amount).And("total_amount-?>=0", amount).Update(npa)
		if err != nil {
			return
		}
		if affectNum != 1 {
			err = fmt.Errorf("ProfitAmount did not changed")
			return
		}
	}
	paLog.ProfitAmountId = npa.Id
	affectNum, err = session.Insert(paLog)
	if affectNum != 1 {
		err = fmt.Errorf("insert prefit amount log error")
		return
	}
	return
}

// 原有的计算账户余额的算法
func CalAmountOfOrigin(accountType int, cmpusersIds string, storeId, vipId, CustomerId, bossId, outreachId, startId, endId int64) (account map[string]int64, err error) {
	account = make(map[string]int64)
	tableName := "jyb_reward"
	withdraw := new(Withdraw)
	reward := new(Reward)
	var total float64
	if accountType == ProfitAmountAccountTypeOfNeoby {
		se := LocalDB.Table(tableName).Where("iscmpy=?", 1).And("bigcustid=?", 0).And("recommendbossid=?", 0).And("recommendid=?", 0).And("outreachid=?", 0).And("cmpyuserid=?", "0")
		if startId > 0 {
			se.Where("id>=?", startId)
		}
		if endId > 0 {
			se.Where("id<=?", endId)
		}
		total, err = se.Sum(reward, "convert(`rewardmoney`,SIGNED)")
		if err != nil {
			return
		}
		account["award"] = int64(total) //获取分润总额

		//获取已提现总额
		account["withdraw"] = 0

	} else if accountType == ProfitAmountAccountTypeOfNeobySellers {
		se := LocalDB.Table(tableName).Where("iscmpy=?", 1).And("cmpyuserid=?", cmpusersIds).And("bigcustid=?", 0).And("recommendbossid=?", 0).And("recommendid=?", 0).And("outreachid=?", 0)
		if startId > 0 {
			se.Where("id>=?", startId)
		}
		if endId > 0 {
			se.Where("id<=?", endId)
		}
		total, err = se.Sum(reward, "convert(`rewardmoney`,SIGNED)")
		if err != nil {
			return
		}
		account["award"] = int64(total) //获取分润总额

		//获取已提现总额
		account["withdraw"] = 0
	} else if accountType == ProfitAmountAccountTypeOfOutreach {
		se := LocalDB.Table(tableName).Where("iscmpy=?", 0).And("outreachid=?", outreachId).And("cmpyuserid=?", 0).And("bigcustid=?", 0).And("recommendbossid=?", 0).And("recommendid=?", 0)
		if startId > 0 {
			se.Where("id>=?", startId)
		}
		if endId > 0 {
			se.Where("id<=?", endId)
		}
		total, err = se.Sum(reward, "convert(`rewardmoney`,SIGNED)")
		if err != nil {
			return
		}
		account["award"] = int64(total) //获取分润总额
		total, err = LocalDB.Where("storeid=?", 0).And("mid=?", 0).And("outreachid=?", outreachId).And("withdraw_status=?", 2).Sum(withdraw, "money")
		if err != nil {
			return
		}
		//获取已提现总额
		account["withdraw"] = int64(total)

	} else if accountType == ProfitAmountAccountTypeOfDoubleOilVIP {

		se := LocalDB.Table(tableName).Where("iscmpy=?", 0).And("bossid=?", bossId).And("bigcustid=?", 1).And("outreachid=?", 0).And("cmpyuserid=?", 0).And("recommendbossid=?", 0).And("recommendid=?", 0)
		if startId > 0 {
			se.Where("id>=?", startId)
		}
		if endId > 0 {
			se.Where("id<=?", endId)
		}
		total, err = se.Sum(reward, "convert(`rewardmoney`,SIGNED)")
		if err != nil {
			return
		}
		account["award"] = int64(total) //获取分润总额

		//获取已提现总额
		account["withdraw"] = 0
	} else if accountType == ProfitAmountAccountTypeOfCustomer {

		se := LocalDB.Table(tableName).Where("iscmpy=?", 0).And("recommendid=?", CustomerId).And("bigcustid=?", 0).And("outreachid=?", 0).And("cmpyuserid=?", 0).And("recommendbossid=?", 0)
		if startId > 0 {
			se.Where("id>=?", startId)
		}
		if endId > 0 {
			se.Where("id<=?", endId)
		}
		total, err = se.Sum(reward, "convert(`rewardmoney`,SIGNED)")
		if err != nil {
			return
		}
		account["award"] = int64(total) //获取分润总额
		total, err = LocalDB.Where("storeid=?", 0).And("mid=?", CustomerId).And("outreachid=?", 0).And("withdraw_status=?", 2).Sum(withdraw, "money")
		if err != nil {
			return
		}
		//获取已提现总额
		account["withdraw"] = int64(total)

	} else if accountType == ProfitAmountAccountTypeOfStore {

		se := LocalDB.Table(tableName).Where("recommendbossid=?", storeId)
		if startId > 0 {
			se.Where("id>=?", startId)
		}
		if endId > 0 {
			se.Where("id<=?", endId)
		}
		total, err = se.Sum(reward, "convert(`rewardmoney`,SIGNED)")
		if err != nil {
			return
		}
		account["award"] = int64(total) //获取分润总额
		total, err = LocalDB.Where("storeid=?", storeId).And("withdraw_status=?", 2).Sum(withdraw, "money")
		if err != nil {
			return
		}
		//获取已提现总额
		account["withdraw"] = int64(total)
	}

	return
}
