package model

import "log"

type NewReward struct {
	Id                   int64  `json:"id"`                                                                      //`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	Orderid              int64  `json:"orderid" xorm:"'orderid' default 0 index"`                                //`orderid` INT(11) NOT NULL DEFAULT '0',
	Stat                 int    `json:"type" xorm:"'type' default 0 index "`                                     //`type` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '* 0  //未分润 1  //已分润，已解冻 2 //已分润，冻结中 -1  //退单后退分润',
	Bossid               int64  `json:"bossid" xorm:"'bossid' default 0 index"`                                  //`bossid` INT(11) NOT NULL DEFAULT '0' COMMENT '发生交易的机构ID',
	Storeid              int64  `json:"storeid" xorm:"'storeid' default 0 index"`                                //`storeid` INT(11) NOT NULL DEFAULT '0' COMMENT '发生交易的门店ID',
	Outreachid           int64  `json:"outreachid" xorm:"'outreachid' default 0 index"`                          //outreachid 外部推广机构id
	Outreachreward       string `json:"outreachreward" xorm:"outreachreward char(100) default ''"`               //`outreachreward` CHAR(100) NOT NULL DEFAULT '' COMMENT '外部机构分润比例（例：40%）或固定金额(例：50)',
	Outreachmoney        int64  `json:"outreachmoney" xorm:"'outreachmoney' default 0"`                          //INT(11) NOT NULL DEFAULT '0' COMMENT '外部推广机构佣金',
	Cmpyuserid           string `json:"cmpyuserid" xorm:"'cmpyuserid' varchar(100)"`                             //cmpyuserid 公司业务员id VARCHAR(100) NOT NULL DEFAULT '0' COMMENT '公司业务员ID',
	Cmpyuserreward       string `json:"cmpyuserreward" xorm:"cmpyuserreward char(100) default ''"`               //CHAR(100) NOT NULL DEFAULT '' COMMENT '分润佣金比例（例：40%）或固定金额(例：50)',
	Cmpyusermoney        int64  `json:"cmpyusermoney" xorm:"'cmpyusermoney' default 0"`                          //INT(11) NOT NULL DEFAULT '0' COMMENT '公司业务员佣金',
	Recommendid          int64  `json:"recommendid" xorm:"'recommendid' default 0 index"`                        // INT(11) NOT NULL DEFAULT '0' COMMENT 'C端推荐人id',
	Recommendreward      string `json:"recommendreward" xorm:"'recommendreward' char(100) default ''"`           // CHAR(100) NOT NULL DEFAULT '' COMMENT '分润佣金比例（例：40%）或固定金额(例：50)',
	Recommendmoney       int64  `json:"recommendmoney" xorm:"'recommendmoney' default 0"`                        // INT(11) NOT NULL DEFAULT '0' COMMENT 'C端推荐人佣金',
	Recommendstoreid     int64  `json:"recommendstoreid" xorm:"'recommendstoreid' default 0 index"`              //  INT(11) NOT NULL DEFAULT '0' COMMENT '门店推荐ID',
	Recommendstorereward string `json:"recommendstorereward" xorm:"'recommendstorereward' char(100) default ''"` //  CHAR(100) NOT NULL DEFAULT '' COMMENT '分润佣金比例（例：40%）或固定金额(例：50)',
	Recommendstoremoney  int64  `json:"recommendstoremoney" xorm:"'recommendstoremoney' default 0"`              // INT(11) NOT NULL DEFAULT '0' COMMENT '门店推荐佣金',
	Bigcustid            int64  `json:"bigcustid" xorm:"'bigcustid' default 0 index"`                            //  INT(11) NOT NULL DEFAULT '0' COMMENT '大客户ID',
	Bigcustreward        string `json:"bigcustreward" xorm:"'bigcustreward' char(100) default ''"`               //  CHAR(100) NOT NULL DEFAULT '' COMMENT '分润佣金比例（例：40%）或固定金额(例：50)',
	Bigcustmoney         int64  `json:"bigcustmoney" xorm:"'bigcustmoney' default 0"`                            // INT(11) NOT NULL DEFAULT '0' COMMENT '门店推荐佣金',
	Oilstoreid           int64  `json:"oilstoreid" xorm:"'oilstoreid' default 0 index"`                          //  INT(11) NOT NULL DEFAULT '0' COMMENT '油站ID',
	Oilstorereward       string `json:"oilstorereward" xorm:"'oilstorereward' char(100) default ''"`             //  CHAR(100) NOT NULL DEFAULT '' COMMENT '油站佣金比例（例：40%）或固定金额(例：50)',
	Oilstoremoney        int64  `json:"oilstoremoney" xorm:"'oilstoremoney' default 0"`                          // INT(11) NOT NULL DEFAULT '0' COMMENT '油站佣金',
	Cmpyreward           string `json:"cmpyreward" xorm:"'cmpyreward' char(100) default ''"`                     //CHAR(100) NOT NULL DEFAULT '' COMMENT '分润佣金比例（例：40%）或固定金额(例：50)',
	Cmpymoney            int64  `json:"cmpymoney" xorm:"'cmpymoney' default 0"`                                  //INT(11) NOT NULL DEFAULT '0' COMMENT '公司佣金',
	CouponsId            int64  `json:"couponsid" xorm:"'couponsid' default 0 index"`                            //INT(11) NOT NULL DEFAULT '0' COMMENT '卡id',
	GoodsId              int64  `json:"goodsid" xorm:"'goodsid' int(11) default 0 index"`                        //INT(11) NOT NULL DEFAULT '0' COMMENT '商品id',
	UserCouponsId        string `json:"usercouponsid" xorm:"'usercouponsid' char(100) default 0 "`               //CHAR(100) NOT NULL DEFAULT '' COMMENT '用户卡券id',
	OrderMoney           int64  `json:"ordermoney" xorm:"'ordermoney' default 0"`                                //INT(11) NOT NULL DEFAULT '0' COMMENT '交易金额',
	ProfitMoney          int64  `json:"profitmoney" xorm:"'profitmoney' default 0"`                              //INT(11) NOT NULL DEFAULT '0' COMMENT '毛利',
	SettlementTime       int64  `json:"settlement_time" xorm:"'settlement_time' int(11)"`                        //settlement_time  结算时间
	UnfreezeTime         int64  `json:"unfreeze_time" xorm:"'unfreeze_time' int(11)"`                            //冻结到的时间
	ProfitAmountId       int64  `json:"profit_amount_id" xorm:"'profit_amount_id' default 0 index"`              //归入账户id
	CreateTime           int64  `json:"create_time" xorm:"int(11) 'create_time'"`                                //create_time 插入时间
	UpdateTime           int64  `json:"update_time" xorm:"int(11) 'update_time'"`                                //update_time 更新时间
}

func (self *NewReward)InsertOrUpdate()(err error){
	has,err:=LocalDB.Table("jyb_newreward").Where("orderid=?",self.Orderid).Get(&NewReward{})
	if err!=nil {
	    log.Println("[Error] ",err)
	    return
	}
	if !has{
		_,err=LocalDB.Table("jyb_newreward").Insert(self)
	}else{
		_,err=LocalDB.Table("jyb_newreward").Where("orderid=?",self.Orderid).Update(self)
	}
	return
}

func (self *NewReward)Update()(err error){
	_,err=LocalDB.Table("jyb_newreward").ID(self.Id).Update(self)
	return
}
