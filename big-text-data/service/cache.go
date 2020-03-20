package service

import (
	"big-text-data/model"
	"time"
)

// 声明一个CacheTimedData接口，该接口包含了一个获取缓存时间的方法
type CacheTimedData interface {
	GetCacheTime() int64 // 加入缓存的时间纳秒
}

// 缓存结果
type QueryResult struct {
	Value     *[]model.Guest // 住客数据切片指针
	CacheTime int64          // 加入缓存的时间
	Count     int            // 查询次数
}

// 实现CacheTimedData接口
func (this *QueryResult) GetCacheTime() int64 {
	return this.CacheTime
}

// 缓存框架，传入一个实现了带缓存时间的数据接口
func UpdateCache2(cacheMap *map[string]CacheTimedData) (delKey string) {
	// 预设一个最早时间
	earliestTime := time.Now().UnixNano()
	for k, v := range *cacheMap {
		if v.GetCacheTime() < earliestTime {
			earliestTime = v.GetCacheTime()
			delKey = k
		}
	}
	delete(*cacheMap, delKey)
	return
}

// 3级缓存？？文件/redis
// 整理缓存策略：删除加入最早的缓存
// 考虑QueryResult改为interface
func UpdateCache(cacheMap *map[string]*QueryResult) (delKey string) {
	// 预设一个最早时间
	earliestTime := time.Now().UnixNano()
	for k, v := range *cacheMap {
		if v.CacheTime < earliestTime {
			earliestTime = v.CacheTime
			delKey = k
		}
	}
	delete(*cacheMap, delKey)
	return
}
