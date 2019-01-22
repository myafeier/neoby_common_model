package model

const (
	CommissionTypeFixRateOfProfit=1 //按毛利润的固定比例分配，单位： 万分之几（整形）
	CommissionTypeOfFixAmountOfProfit=2 //2: 按固定金额分配，单位：分
	CommissionTypeOfFixRateOfTrade=3 //3: 按交易金额的服务费率分配，单位：万分之几（整形）

	CommissionFreightDeductionFalse=-1 //分润前是否扣除运费，1:扣除，-1 不扣除
	CommissionFreightDeductionTrue=1 //分润前是否扣除运费，1:扣除，-1 不扣除
)
//分润规则
type Commission struct {
	Id               int64
	Name             string `xorm:"'name' default 0 "`
	Store            int64  `xorm:"'store' default 0 "`
	StoreType        int    `xorm:"'store_type' default 0 "`
	Outreach         int64  `xorm:"'outreach' default 0 "`
	OutreachType     int    `xorm:"'outreach_type' default 0 "`
	Member           int64  `xorm:"'member' default 0 "`
	MemberType       int    `xorm:"'member_type' default 0 "`
	Vip              int64  `xorm:"'vip' default 0 "`
	VipType          int    `xorm:"'vip_type' default 0 "`
	CompanyUser      int64  `xorm:"'cmpyuser' default 0 "`
	CompanyUserType  int    `xorm:"'cmpyuser_type' default 0 "`
	Company          int64  `xorm:"'cmpy' default 0 "`
	CompanyType      int    `xorm:"'cmpy_type' default 0 "`
	FreightDeduction int    `xorm:"'freight_deduction' default 0 "`
	CreateTime       int64  `xorm:"'create_time' default 0 "`
	UpdateTIme       int64  `xorm:"'update_t_ime' default 0 "`
}

func (self Commission)TableName() string  {
	return "jyb_commission"
}
