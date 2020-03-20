package service

import (
	"big-text-data/model"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

// 按照省份划分数据
func divideByProvince() (err error){
	provinces := []string{"北京市11", "天津市12", "河北13", "山西14", "内蒙15", "辽宁21"}
	provinceMap := make(map[string]*model.Province)
	filePath := "C:/workspace/go/src/big-text-data/data/"
	// 创建34个省份实例
	for _, value := range provinces {
		name := value[0 : len(value)-2] // 省份名称
		id := value[len(value)-2:]      // 省份编号
		file, err := os.OpenFile(filePath+name+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("os openFile err=", err)
			return err
		}
		dataChan := make(chan string)
		province := model.Province{
			Id:       id,
			Name:     name,
			File:     file,
			DataChan: dataChan,
		}
		provinceMap[id] = &province
	}

	// 起协程
	for _, province := range provinceMap {
		go writeFileToProvince(province)
	}

	// 读入数据，逐行判断身份证的前两位
	filename := filePath + "goodData.txt"
	file, _ := os.Open(filename)
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF { // 读到文件末尾，跳出循环
			break
		}
		lineStr := string(lineBytes)
		fields := strings.Split(lineStr, ",")
		id := fields[1][0:2] // id号是身份证号前两位
		if province, ok := provinceMap[id]; ok {
			// 对号入座，写入相应的管道
			wg.Add(1)
			province.DataChan <- lineStr
		} else {
			fmt.Println("省份id不正确...")
		}
	}
	return
}

// 写数据到按省份分类文件
func writeFileToProvince(province *model.Province) {
	// 死循环从管道中取数据，写入文件
	for lineStr := range province.DataChan {
		_, _ = province.File.WriteString(lineStr + "\n")
		wg.Done()
		fmt.Println(province.Name, "写入:", lineStr)
	}
}

// 按照年龄划分
func divideByAge() (err error){
	filePath := "C:/workspace/go/src/big-text-data/data/"
	// 创建一堆年代对象
	agesMap := make(map[string]*model.Age)
	for i := 195; i < 200; i++ {
		fileName := filePath + strconv.Itoa(i) + "x.txt"
		age := model.Age{
			Decade:   strconv.Itoa(i),
			DataChan: nil,
		}
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		defer file.Close()
		if err != nil {
			fmt.Println("open file err=", err)
			return err
		}
		age.File = file
		age.DataChan = make(chan string)
		agesMap[age.Decade] = &age
	}
	// 为每一个年代开启一个写入协程
	for _, age := range agesMap {
		go writeFileToAge(age)
	}
	// 读入分类数据，断行，判断年代，对号入座
	filename := filePath + "goodData.txt"
	file, _ := os.Open(filename)
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF { // 读到文件末尾，跳出循环
			break
		}
		lineStr := string(lineBytes)
		fields := strings.Split(lineStr, ",")
		decade := fields[1][6:9]
		if age, ok := agesMap[decade]; ok {
			// 对号入座，写入相应的管道
			wg.Add(1)
			age.DataChan <- lineStr
		} else {
			fmt.Println("年代不正确...")
		}
	}
	return
}

// 写数据到按年代分类文件
func writeFileToAge(age *model.Age) {
	// 死循环从管道中取数据
	for contentStr := range age.DataChan {
		_, _ = age.File.WriteString(contentStr + "\n")
		wg.Done()
		fmt.Println(age.Decade+"x", "写入：", contentStr)
	}
}

func DivideByProvince() (err error){
	err = divideByProvince()
	wg.Wait()
	return
}

func DivideByAge() (err error){
	err = divideByAge()
	wg.Wait()
	return
}