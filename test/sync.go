package main

import (
	"github.com/myafeier/neoby_common_model"
	"log"
)

func main()  {

//host: 10.8.63.216
//port: 3306
//username: zl
//passwd: cssgdddbt
//db: jyb
	err:=model.InitRemoteDB("10.8.63.216","3306","zl","cssgdddbt","jybtest",true)
	if err!=nil{
		log.Fatal(err)
	}
	err=model.SyncCommissionTable()
	if err!=nil{
		log.Fatal(err)
	}

}

