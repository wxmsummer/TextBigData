package model

import (
	"fmt"
	"os"
)

type Guest struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	IdNumber string `db:"idNumber"`
}

type Province struct {
	Id       string      // 省份编号，即身份证号前两位
	Name     string      // 省份名称
	File     *os.File    // 绑定一个文件指针，便于按省份进行文件读写
	DataChan chan string // 本省文件的数据管道
}

type Age struct {
	Decade   string // 年代，取前三位：196x，197x，198x，199x，200x，201x
	File     *os.File
	DataChan chan string
}

// 通过名字查询住户
func QueryGuestsByName(name string) (guests []Guest, err error){
	sqlStr := fmt.Sprintf("select * from guests where name like '%s';", name)
	rows, err := QueryRowsDB(sqlStr)
	if err != nil {
		fmt.Println("db query err=", err)
		return
	}
	for rows.Next() {
		guest := Guest{}
		rows.Scan(&guest.Id, &guest.Name, &guest.IdNumber)
		guests = append(guests, guest)
	}
	return
}


func InsertGuest(guest *Guest) (int64, error) {
	sql := "insert into guests (name, idNumber) values (?, ?);"
	return ModifyDB(sql, guest.Name, guest.IdNumber)
}