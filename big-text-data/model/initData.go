package model

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// 数据入库：将文本大数据以适当的结构存入mysql数据库
// 终端循环输入 search -name xxx 对记录进行查询：
// 实现精确查询和模糊查询
// 实现内存-数据库的二级缓存策略，并显示每一次查询的时间消耗
// 入库成功后，做一个文件标记，下一次见到标记就不再执行入库操作
var chanData chan *Guest
var wg sync.WaitGroup
func InitData() (err error) {
	// 如果已经初始化过了，就不执行本函数
	_, err = os.Stat("C:/workspace/go/src/big-text-data/data/dbInit_ok.mark")
	if err == nil {
		fmt.Println("标记文件已存在")
	} else {
		fmt.Println("标记文件不存在...")
	}
	// 分批次读入大数据文本
	filename := "C:/workspace/go/src/big-text-data/data/goodData.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("os.Open err=", err)
	}
	reader := bufio.NewReader(file)
	// 数据管道
	chanData = make(chan *Guest, 100000)
	// 开启协程，从数据管道获取信息，插入数据库，一读多写
	for i := 0; i < 10; i++ {
		go AddGuest()
	}

	for {
		lineBytes, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("reader err=", err)
		}
		if err == io.EOF {
			break
		}
		// 逐条入库（如何并发入库？）
		lineStr := string(lineBytes)
		fields := strings.Split(lineStr, ",")
		name, idNumber := fields[0], fields[1]
		guest := Guest{
			Name:     name,
			IdNumber: idNumber,
		}
		// 方案一：每次插入都开协程??? // 不能这样 行不通,开2000万条协程，耗尽资源，程序崩溃
		// go InsertGuest(db, &guest)

		// 方案二：开有限条协程，从管道中读取数据
		// 主协程负责读，小协程负责写
		wg.Add(1)
		chanData <- &guest
	}
	fmt.Println("数据初始化成功...")
	// 创建一个文件，标记一下初始化成功，以后就不再进行初始化
	//_, err = os.Create("C:/workspace/go/src/big-text-data/data/dbInit_ok.mark")
	//if err == nil {
	//	fmt.Println("初始化标记文件已创建...")
	//}
	wg.Wait()
	return
}

// 如何进一步优化？？？
// 多读多写？？？分布式读写，把文件拆分成多个；；分数据库
// 分多个库写
func AddGuest() {
	for guest := range chanData {
		// 循环插入，至成功为止
		for {
			result, err := InsertGuest(guest)
			if err != nil {
				fmt.Println("InsertGuest err=", err)
				// db 的性能被耗尽了，先让协程歇会
				<-time.After(5 * time.Second)
			} else if result > 0 {
				fmt.Printf("数据库写入%s成功!\n", guest.Name)
				wg.Done()
				break
			}
		}
	}
}
