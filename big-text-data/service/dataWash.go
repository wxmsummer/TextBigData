package service

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// 数据清洗，除去不符合格式的数据
func DataWash() (err error){
	filePath := "C:/workspace/go/src/big-text-data/data/"
	// 原始数据
	srcFile, _ := os.Open(filePath + "srcData.txt")
	defer srcFile.Close()
	// 准备一个优质数据文件
	goodFile, _ := os.OpenFile(filePath+"goodData.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer goodFile.Close()
	// 准备一个劣质数据文件
	badFile, _ := os.OpenFile(filePath+"badData.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer badFile.Close()
	reader := bufio.NewReader(srcFile)
	for {
		// 逐行读取
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		//如果原来是gbk，先转成utf8
		//gbkStr := string(lineBytes)
		//lineStr, _ := utils.ConvertEncoding(gbkStr, "GBK", "UTF-8")
		lineStr := string(lineBytes)
		// 这里可以考虑做一个清洗进度条？？？
		fields := strings.Split(lineStr, ",")
		// 进行身份证数据清洗
		if len(fields) > 1 && len(fields[1]) == 18 {
			// 写入另一个文件中
			_, _ = goodFile.WriteString(lineStr + "\n")
			fmt.Println("Good:", lineStr)
		} else { // 优化细节：开协程，将数据写入不同文件
			// 写入劣质文件
			_, _ = badFile.WriteString(lineStr + "\n")
			fmt.Println("Bad:", lineStr)
		}
	}
	fmt.Println("数据清洗完毕...")
	return
}
