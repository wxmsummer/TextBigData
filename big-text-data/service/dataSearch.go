package service

import (
	"big-text-data/model"
	"fmt"
	"time"
)

// 接受用户想要查询的姓名
// 先查看内存中是否存在该记录
// 如果内存中没有，再去查数据库，查到的结果存入内存：局部效应
// 如果内存增长过大，则按照一定策略清除内存：缓存优化
func Search() (err error){
	var name string
	var guestsMap = make(map[string]*QueryResult, 0)
	for {
		fmt.Println("请输入要查询的住客姓名....")
		fmt.Println("输入exit退出查询，输入cache显示缓存结果")
		fmt.Scanf("%s", &name)
		if name == "exit" {
			break
		}else if name == "cache" {
			// 查看所有缓存
			fmt.Printf("共缓存了%d条结果\n", len(guestsMap))
			for k, _ := range guestsMap {
				fmt.Println(k)
			}
			continue
		}
		// 先查看内存中是否存在该记录
		if queryResult, ok := guestsMap[name]; ok {
			queryResult.Count++
			fmt.Println(queryResult.Value)
			fmt.Printf("查询完毕，共查询到%d条结果。\n", len(*queryResult.Value))
			continue
		}
		guests, err := model.QueryGuestsByName(name)
		if err != nil {
			fmt.Println("QueryGuestsByName err=", err)
		}
		queryResult := QueryResult{
			Value:     &guests,
			CacheTime: time.Now().UnixNano(),
			Count:     1,
		}
		fmt.Println(guests)
		fmt.Printf("查询完毕，共查询到%d条结果。\n", len(guests))
		guestsMap[name] = &queryResult
		// 有必要时，清除一些缓存
		if len(guestsMap) > 2 { // 2可以配置
			delKey := UpdateCache(&guestsMap)
			fmt.Printf("触发缓存清除，%s已经被淘汰出缓存\n", delKey)
		}
	}
	return
}
