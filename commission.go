package model

const (
	CommissionTypeFixRateOfProfit     = 1 //从毛利润中分配利润， 其它角色按固定比例分配，剩余的利润归公司单位： 万分之几（整形）
	CommissionTypeOfFixAmountOfProfit = 2 //2: 按固定金额分配，与毛利无关。单位：分
	CommissionTypeOfFixRateOfTrade    = 3 //3:暂时无用

	CommissionStatOk    = 1  //规则正常状态
	CommissionStatPause = -1 //规则停用状态

	CommissionRuleTypeOfWithoutReference                  = 1 //无分享分润规则
	CommissionRuleTypeOfWithCustomerReference             = 2 //个人推荐消费规则
	CommissionRuleTypeOfWithMerchantReference             = 3 //商户推荐消费分润规则
	CommissionRuleTypeOfWithCustomerFromMerchantReference = 4 //来自商户的顾客推荐消费
	CommissionRuleTypeOfWithPartnerReference              = 5 //合作伙伴介绍分润规则

)

//分润规则
type Commission struct {
	Id         int64
	Name       string            `xorm:"varchar(100) default ''"`
	Status     int               `xorm:"tinyint(1) default 1"`
	Desc       string            `xorm:"text"`
	CreateTime int64             `xorm:"default 0"`
	Rules      []*CommissionRule `xorm:"-"`
}

func (self Commission) TableName() string {
	return "jyb_commission"
}

func (self *Commission) GetMap(withRules bool) (result map[int64]*Commission, err error) {
	result = make(map[int64]*Commission)
	rows, err := LocalDB.Table("jyb_commission").Rows(self)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := new(Commission)
		err = rows.Scan(t)
		if err != nil {
			return
		}
		if withRules {
			err = LocalDB.Where("commission_id=?", t.Id).And("status=?", self.Status).Find(&t.Rules)
			if err != nil {
				return
			}
		}
		result[t.Id] = t
	}
	return
}

//分润规则映射
type CommissionMapping struct {
	Id           int64
	BossId       int64 `json:"boss_id" xorm:"default 0 index"`
	StoreId      int64 `json:"store_id" xorm:"default 0 index"`
	GoodsId      int64 `json:"goods_id" xorm:"default 0 index"`
	CouponId     int64 `json:"coupon_id" xorm:"default 0 index"`
	CommissionId int64 `json:"commission_id" xorm:"default 0 index"`
	CreateTime   int64 `json:"create_time" xorm:"default 0"`
}

func (self CommissionMapping) TableName() string {
	return "jyb_commission_mapping"
}

//分润细则
type CommissionRule struct {
	Id              int64
	CommissionId    int64 `xorm:"index"`
	Default         int   `xorm:"tinyint(1) default 0"` //'是否默认，0:非，1:是'
	Status          int   `xorm:"tinyint(1) default 1"` //'细则状态  1:正常，-1:停用'
	Type            int   `xorm:"tinyint(1) default 0"` //'分润细则类型：\n1 //无分享分润规则\n2 //个人推荐消费规则\n3 //商户推荐消费分润规则\n4 //来自商>     户的顾客推荐消费\n5 //合作伙伴介绍分润规则'
	Store           int64 `xorm:"'store' default 0 "`
	StoreType       int   `xorm:"'store_type' default 0 "`
	Outreach        int64 `xorm:"'outreach' default 0 "`
	OutreachType    int   `xorm:"'outreach_type' default 0 "`
	Member          int64 `xorm:"'member' default 0 "`
	MemberType      int   `xorm:"'member_type' default 0 "`
	Vip             int64 `xorm:"'vip' default 0 "`
	VipType         int   `xorm:"'vip_type' default 0 "`
	CompanyUser     int64 `xorm:"'cmpyuser' default 0 "`
	CompanyUserType int   `xorm:"'cmpyuser_type' default 0 "`
	Company         int64 `xorm:"'cmpy' default 0 "`
	CompanyType     int   `xorm:"'cmpy_type' default 0 "`
	Reference       int64 `xorm:" default 0"`
	ReferenceType   int   `xorm:"'reference_type' default 0 "`
	ChargeType      int   `xorm:"'charge_type' default 0 "` // 手续费扣除方式，0:不扣除，1 :从商品毛利中扣除
	CreateTime      int64 `xorm:"'create_time' default 0 "`
	UpdateTime      int64 `xorm:"'update_time' default 0 "`
}

func (self CommissionRule) TableName() string {
	return "jyb_commission_rule"
}
