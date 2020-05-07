package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var (
	DB *gorm.DB
)

func InitDB() {
 	var err error
	DB,err = gorm.Open("mysql","root:@tcp(127.0.0.1:3306)/secKill?parseTime=true&charset=utf8&loc=Local")
	if err != nil {
		log.Println("Fail to open mysql :",err)
		return
	}

	if !DB.HasTable(&Order{}){
		DB.CreateTable(&Order{})
	}
}

 

