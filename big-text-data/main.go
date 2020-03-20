package main

import (
	"big-text-data/model"
	"big-text-data/service"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {

	//err := service.DataWash()
	//if err != nil {
	//	fmt.Println("DataWash err=",err)
	//}
	//err = service.DivideByProvince()
	//if err != nil {
	//	fmt.Println("DivideByProvince err=",err)
	//}
	//err = service.DivideByAge()
	//if err != nil {
	//	fmt.Println("DivideByAge err=",err)
	//}

	dbName := "mysql"
	dsn := "root:wxm19990516@tcp(127.0.0.1:3306)/bigdata?charset=utf8"
	maxOpen := 200
	maxIdle := 100
	maxLifeTime := time.Second * 1000
	err := model.InitMysql(dbName, dsn, maxOpen, maxIdle, maxLifeTime)
	if err != nil {
		fmt.Println("db init err=",err)
	}
	err = model.InitData()
	if err != nil {
		fmt.Println("InitData err=",err)
	}
	err = service.Search()

}
