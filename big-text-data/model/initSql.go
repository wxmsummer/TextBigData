package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db *sql.DB

// 初始化mysql、创建连接池
func InitMysql(name, dsn string, maxOpen, maxIdle int, maxLifeTime time.Duration) (err error){
	db1, err := sql.Open(name, dsn)
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		db = db1
		db.SetMaxOpenConns(maxOpen)
		db.SetMaxIdleConns(maxIdle)
		db.SetConnMaxLifetime(maxLifeTime)
		CreateTableWithGuests()
	}
	return
}

// 操作数据库
func ModifyDB(sql string, args ...interface{}) (int64, error) {
	result, err := db.Exec(sql, args...)
	if err != nil {
		fmt.Println("db.Exec err =", err)
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		fmt.Println("result.RowsAffected err=", err)
		return 0, err
	}
	return count, err
}

// 查询单行数据
func QueryRowDB(sql string) *sql.Row {
	return db.QueryRow(sql)
}

// 查询多行数据
func QueryRowsDB(sql string) (*sql.Rows, error) {
	return db.Query(sql)
}

// 创建住房用户表
func CreateTableWithGuests() {
	sqlStr := `create table if not exists guests(
	id int primary key auto_increment not null,
	name varchar(50),
	idNumber char(18)
	);`
	_, _ = ModifyDB(sqlStr)
}



