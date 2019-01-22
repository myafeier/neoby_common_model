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
	ProfitAmountAccountTypeOfMerchant     = 6 //商户及门店

)

//
type ProfitAmount struct {
	Id             int64
	Type         int  `xorm:"tinyint(2) default 0 index"`  //账户类型
	BossId         int64  `xorm:"default 0 index"`            // 商户ID
	StoreId        int64  `xorm:"default 0 index"`            //门店id
	OutreachId     int64  `xorm:"default 0 index"`            //外部推广id
	MemberId       int64  `xorm:"default 0 index"`            //会员id
	VipId          int64  `xorm:"default 0 index"`            //vip id
	SellerIds      string `xorm:"char(100) default '' index"` //销售人员id，多个用,分割
	IsCompany      int    `xorm:"default 0 index"`            //是否是滇汇通公司
	TotalAmount    int64  `xorm:"default 0"`                  //历史合计收益
	FreeAmount     int64  `xorm:"default 0"`                  // 可提现收益
	FreezeAmount   int64  `xorm:"default 0"`                  // 冻结中收益
	WithdrewAmount int64  `xorm:"default 0"`                  //已提现收益
}

type ProfitAmountLog struct {
	Id                   int64
	Action               string    `xorm:"char(20) default ''  index"` //操作名称
	ProfitAmountId int64  `xorm:"default 0 index"` //账户id
	Amount               int64     `xorm:"default 0"`                  //操作前历史合计收益
	UnfreezeTime         int64     `xorm:"default 0 index"`            //预设解冻时间，只有新分配利润才有预设的解冻时间
	BeforeTotalAmount    int64     `xorm:"default 0"`                  //操作前历史合计收益
	BeforeFreeAmount     int64     `xorm:"default 0"`                  //操作前可提现收益
	BeforeFreezeAmount   int64     `xorm:"default 0"`                  //操作前冻结中收益
	BeforeWithdrewAmount int64     `xorm:"default 0"`                  //操作前已提现收益
	//AfterTotalAmount     int64     `xorm:"default 0"`                  //操作后历史合计收益
	//AfterFreeAmount      int64     `xorm:"default 0"`                  //操作后可提现收益
	//AfterFreezeAmount    int64     `xorm:"default 0"`                  //操作后冻结中收益
	//AfterWithdrewAmount  int64     `xorm:"default 0"`                  //操作后已提现收益
	Remark               string    `xorm:"varchar(255) default ''"`    //备注
	RewardId             int64     `xorm:"default 0"`                  //奖励ID，如有
	WithdrawId           int64     `xorm:"default 0"`                  // 提现id，如有
	Created              time.Time `xorm:"created"`
	Updated              time.Time `xorm:"updated"`
}

// bossId 和storeId 必须成对出现,  如果有vipId 依赖于 boosId
// 如果是新增收益需同时提交解冻时间
func TChangeProfitAmount(session *xorm.Session, action string, reward *Reward,withDraw *Withdraw) (accountId int64,err error) {

	//检查账户存在性
	pa := new(ProfitAmount)
	var amount int64
	var withdrawId,rewardId int64

	if reward!=nil{
		//因分润涉及多种id组合形式，故确定账户Id需要根据组合进行判断
		if reward.IsCompany == 1 { //公司内部分润
			if sellerIds := strings.TrimSpace(reward.CompanyUserId); sellerIds == ""||sellerIds=="0" {
				pa.Id = ProfitAmountOfDhtId //公司账户
				pa.Type=ProfitAmountAccountTypeOfNeoby
			} else {
				pa.IsCompany = 1
				pa.SellerIds = sellerIds
				pa.Type=ProfitAmountAccountTypeOfNeobySellers
			}
		} else if reward.OutreachId > 0 { //合作机构 账户
			pa.OutreachId = reward.OutreachId
			pa.Type=ProfitAmountAccountTypeOfOutreach
		} else if reward.BigCustomerId > 0 { //两油VIP账户
			pa.VipId = reward.BigCustomerId
			pa.BossId = reward.BossId
			pa.Type=ProfitAmountAccountTypeOfDoubleOilVIP
		} else if reward.RecommendId > 0 { //顾客用户
			pa.MemberId = reward.RecommendId
			pa.Type=ProfitAmountAccountTypeOfCustomer
		} else if reward.RecommendBossId > 0 {
			if reward.RecommendBossId==1{  //历史bug，推荐id为1 的实际为公司收益
				pa.Id=1
			}else if reward.StoreId != reward.RecommendBossId { //非两油的商户
				pa.StoreId = reward.RecommendBossId
			} else { //两油的油站
				pa.StoreId = reward.StoreId //自己推荐的自己
				pa.BossId = reward.BossId
			}
			pa.Type=ProfitAmountAccountTypeOfMerchant
		} else { //纯碎的商户/门店账户
			pa.StoreId = reward.StoreId
			pa.BossId = reward.BossId
			pa.Type=ProfitAmountAccountTypeOfMerchant
		}
		amount=reward.RewardMoney
		rewardId=reward.Id
	}

	if withDraw!=nil{
		if withDraw.OutreachId!=0{
			pa.OutreachId=withDraw.OutreachId
			pa.Type=ProfitAmountAccountTypeOfOutreach
		}else if withDraw.MemberId!=0{
			pa.MemberId=withDraw.MemberId
			pa.Type=ProfitAmountAccountTypeOfCustomer
		}else if withDraw.MerchantId!=0||withDraw.StoreId!=0{
			pa.Type=ProfitAmountAccountTypeOfMerchant
			pa.BossId=withDraw.MerchantId
			pa.StoreId=withDraw.StoreId
		}
		amount=withDraw.Money
		withdrawId=withDraw.Id
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

	//锁定
	npa := new(ProfitAmount)
	npa.Id = pa.Id
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

	if reward!=nil{
		paLog.UnfreezeTime = reward.UnfreezeTime
	}
	paLog.WithdrawId = withdrawId
	paLog.RewardId = rewardId

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
		affectNum, err = session.ID(npa.Id).Cols("free_amount", "total_amount", ).Decr("free_amount", amount).Decr("total_amount", amount).Where("free_amount-?>=0", amount).And("total_amount-?>=0", amount).Update(npa)
		if err != nil {
			return
		}
		if affectNum != 1 {
			err = fmt.Errorf("ProfitAmount did not changed")
			return
		}
	case ProfitActionReduceFree: //从可用中撤销收入
		affectNum, err = session.ID(npa.Id).Cols("free_amount", "total_amount", ).Decr("free_amount", amount).Decr("total_amount", amount).Where("free_amount-?>=0", amount).And("total_amount-?>=0", amount).Update(npa)
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
	paLog.ProfitAmountId=npa.Id
	affectNum, err = session.Insert(paLog)
	if affectNum != 1 {
		err = fmt.Errorf("insert prefit amount log error")
		return
	}
	accountId=pa.Id
	return
}
