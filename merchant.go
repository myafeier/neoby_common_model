package model

const (
	MerchantSinopecId =1  //中石化商户id
	MerchantCNPCId=115  //中石油商户id
	MerchantCategoryIdOfDoubleGas=14 //两油的行业分类ID
)

// 对应表jyb_bossinfo
type Merchant struct {
	Id int64 `gorm:"primary_key"`
	CategoryId int64 `xorm:"'categoryid'" gorm:"column:categoryid;type:int(11)"`  //categoryid 行业分类id
	EmployeeUserId string `xorm:"cmpyuserid" gorm:"column:cmpyuserid;type:varchar(100)"`  //cmpyuserid  关联的公司业务员
}

func GetMerchantById(merchantId int64)(merchant *Merchant,err error ){
	merchant=new(Merchant)
	_,err=RemoteDB.Table("jyb_bossinfo").Select("id,categoryid,cmpyuserid").Where("id=?",merchantId).Get(merchant)
	return
}