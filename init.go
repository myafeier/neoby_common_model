package model

import (
	"github.com/go-xorm/xorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var LocalDB *xorm.Engine
var RemoteDB *xorm.Engine

func InitRemoteDB(ip, port, user, password, dbname string, showSQL bool) (err error) {
	RemoteDB, err = xorm.NewEngine("mysql", user+":"+password+"@tcp("+ip+":"+port+")/"+dbname+"?charset=utf8mb4")
	if err != nil {
		return err
	}
	RemoteDB.SetMaxIdleConns(10)
	RemoteDB.SetMaxOpenConns(100)
	RemoteDB.SetConnMaxLifetime(100 * time.Second)
	RemoteDB.ShowSQL(showSQL)
	return
}

func InitLocalDB(ip, port, user, password, dbname string, showSQL bool) (err error) {
	LocalDB, err = xorm.NewEngine("mysql", user+":"+password+"@tcp("+ip+":"+port+")/"+dbname+"?charset=utf8mb4")
	if err != nil {
		return
	}
	LocalDB.SetMaxIdleConns(10)
	LocalDB.SetMaxOpenConns(100)
	LocalDB.SetConnMaxLifetime(100 * time.Second)
	LocalDB.ShowSQL(showSQL)
	return
}


//初始化本元账户of 收益
func InitNeobyAccountOfProfitAmount() (err error) {
	//初始化公司账户
	pa := &ProfitAmount{}
	pa.Id = 1
	pa.IsCompany = 1
	pa.Type = ProfitAmountAccountTypeOfNeoby
	has, err := LocalDB.Get(pa)
	if err != nil {
		return
	}
	if !has {
		_, err = LocalDB.Insert(pa) //初始化
		if err != nil {
			return
		}
	}
	return
}

func UpdateRemoteDbStructOfTable() (err error) {
	_, err = RemoteDB.Exec("alter table jyb_store add column profit_freeze_day int(11) default 0")
	return
}

func SyncLocalDbTable()(err error){
	//var tables = []interface{}{Order{}}
	err=LocalDB.Sync2(Order{})
	return
}

func CreateProfitAmountsTable()(err error){
	var tables = []interface{}{ProfitAmount{}, ProfitAmountLog{}}
	err=createTable(RemoteDB,tables)
	return
}
func SyncCommissionTable()(err error){
	var tables = []interface{}{Commission{}, CommissionRule{},CommissionMapping{}}
	err=createTable(RemoteDB,tables)
	return
}


func CreateLocalDBTable() (err error) {
	var tables = []interface{}{ProfitAmount{}, ProfitAmountLog{}, Reward{}, Withdraw{}, Order{}}
	err=createTable(LocalDB,tables)
	return
}

func createTable(engine *xorm.Engine,tables []interface{})(err error){
	var isExist bool
	for _, v := range tables {
		isExist, err = engine.IsTableExist(v)
		if err != nil {
			return
		}
		if !isExist {
			err = engine.CreateTables(v)
			if err != nil {
				return
			}
			err = engine.CreateIndexes(v)
			if err != nil {
				return
			}
		}
	}
	return

}
