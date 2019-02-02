package model

// 对应表 jyb_store
type Store struct {
	Id int64 `gorm:"primary_key"`
	ProfitFreezeDay int `xorm:"int(11) default 0" gorm:"column:profit_freeze_day;type:int(11)"`  //分润冻结天数
}


func (self Store)TableName() string  {
	return "jyb_store"
}

func GetStoreById(storeId int64)(store *Store,err error){
	store=new(Store)
	_,err=RemoteDB.Table("jyb_store").Where("id=?",storeId).Get(store)
	return
}