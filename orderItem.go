package model

type OrderItem struct {
	Id int64
	OrderId int64	`xorm:"'orderid' int(11) default 0 index"`
	ProductId int64 `xorm:"'productid' int(11) default 0 index"`
	ProductName string `xorm:"'productname' varchar(255) default ''"`
}

func (self OrderItem)GetName()string{
	return "jyb_shoporderitem"
}
