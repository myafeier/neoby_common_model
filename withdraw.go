package model

import "fmt"

const (
	WithDrawStatusSubmitted  = 1  //提现已提交
	WithDrawStatusSuccess    = 2  //提现成功
	WithDrawStatusProcessing = 3  //提现处理中
	WithDrawStatusFail       = -1 //提现失败
)

// 提现表
type Withdraw struct {
	Id             int64  `json:"id"`                                                         //int(11) NOT NULL AUTO_INCREMENT,
	MemberId       int64  `xorm:"'mid' default 0 index"`                                      //int(11) NOT NULL DEFAULT '0' COMMENT '用户id',
	MerchantId     int64  `xorm:"'bossid' default 0 index"`                                   //int(11) NOT NULL DEFAULT '0' COMMENT '机构ID',
	StoreId        int64  `xorm:"'storeid' default 0 index"`                                  //int(11) NOT NULL DEFAULT '0' COMMENT '门店ID',
	OutreachId     int64  `xorm:"'outreachid' default 0 index"`                               //int(11) NOT NULL DEFAULT '0' COMMENT '外部推广ID',
	OrderSn        string `xorm:"'order_sn' varchar(100) default 0 index"`                    //varchar(100) NOT NULL DEFAULT '' COMMENT '订单号',
	Money          int64  `xorm:"'money' default 0"`                                          //char(100) NOT NULL DEFAULT '' COMMENT '交易金额',
	AccountMoney   int64  `xorm:"'countmoney' default 0  "`                                   //varchar(50) NOT NULL COMMENT '账户余额',
	BankCardId     string `xorm:"'bankcardid' varchar(255) default ''"`                       //varchar(255) NOT NULL COMMENT '银行卡卡号id',
	BankNo         string `xorm:"'bankno'  varchar(255) default ''"`                          //varchar(255) NOT NULL COMMENT '银行类别',
	TaxAmount      int64  `xorm:"'taxfree' default 0"`                                        //char(100) NOT NULL DEFAULT '' COMMENT '代缴税金',
	AfterTaxAmount int64  `xorm:"'sjprice' default 0"`                                        //char(100) NOT NULL DEFAULT '' COMMENT '实际提现金额',
	WithDrawStatus int    `xorm:"'withdraw_status' tinyint(1)"`                               //tinyint(1) NOT NULL DEFAULT '1' COMMENT '提现状态1提现申请2提现成功-1提现失败',
	ProfitAmountId int64  `json:"profit_amount_id" xorm:"'profit_amount_id' default 0 index"` //归入账户id
	WithDrawType   int    `xorm:"'type' tinyint(1)"`                                          //tinyint(1) NOT NULL DEFAULT '0' COMMENT '0门店1机构',
	CreateTime     int64  `xorm:"'create_time' "`                                             //int(11) NOT NULL DEFAULT '0' COMMENT '订单创建时间',
	UpdateTime     int64  `xorm:"'update_time' "`                                             //int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
}

func (self Withdraw) TableName() string {
	return "jyb_withdraw"
}

func GetProfitAmountOfWithdraw(withdraw *Withdraw) (profitAmountId int64, err error) {
	if withdraw.ProfitAmountId != 0 {
		profitAmountId = withdraw.ProfitAmountId
		return
	}

	session := LocalDB.NewSession()

	if withdraw.MemberId > 0 {
		session.Where("type=?", ProfitAmountAccountTypeOfCustomer).And("member_id=?", withdraw.MemberId)
	} else if withdraw.OutreachId > 0 {
		session.Where("type=?", ProfitAmountAccountTypeOfOutreach).And("outreach_id=?", withdraw.OutreachId)
	} else if withdraw.StoreId > 0 {
		session.Where("type=?", ProfitAmountAccountTypeOfStore).And("store_id=?", withdraw.StoreId)
	}else if withdraw.MerchantId>0{
		session.Where("type=?", ProfitAmountAccountTypeOfBoss).And("boss_id=? and store_id=0", withdraw.MerchantId)
	}
	pa := new(ProfitAmount)
	has, err := session.Get(pa)
	if err != nil {
		return
	}
	if !has {
		err = fmt.Errorf("account not found,withdraw:%+v", *withdraw)
	}
	profitAmountId = pa.Id
	return
}
