package utils

import (
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"sort"
)

type HashLoop struct{
	maxNum int64
	cacheRedisPool map[int64]*redis.Pool

	redisPoolSortedNums []int64	// 需要保证非空

}

func NewHashLoop(maxNum int64,cacheRedisPool map[int64]*redis.Pool) *HashLoop{
	var redisPoolSortedNums []int64
	for num,_ := range cacheRedisPool{
		redisPoolSortedNums = append(redisPoolSortedNums,num)
		if num > maxNum{
			logs.Warn("hash loop is invalid as its redis pool num is over max num")
		}
	}
	sort.Slice(redisPoolSortedNums,func(i ,j int) bool{
		return redisPoolSortedNums[i]<redisPoolSortedNums[j]
	})
	if len(redisPoolSortedNums)==0{
		logs.Warn("hash loop is invalid as its redis pool num is empty")
	}
	return &HashLoop{
		maxNum: maxNum,
		cacheRedisPool: cacheRedisPool,
		redisPoolSortedNums: redisPoolSortedNums,
	}
}

func (h *HashLoop)GetRedisPool(aroundRedisPoolNum int64) *redis.Pool{
	redisPoolNumIndex := getIndexOfFirstGreaterOrEqual(h.redisPoolSortedNums,aroundRedisPoolNum)
	redisPoolNum := h.redisPoolSortedNums[redisPoolNumIndex%len(h.redisPoolSortedNums)]
	return h.cacheRedisPool[redisPoolNum]
}

func (h *HashLoop)getMaxNum() int64{
	return h.maxNum
}

func getIndexOfFirstGreaterOrEqual(sortedNums []int64,targetNum int64) int{
	left,right:=0,len(sortedNums)-1
	for left<=right{
		mid := left + (right-left)/2
		if sortedNums[mid]>=targetNum{
			right = mid - 1
		}else{
			left = mid + 1
		}
	}
	return left
}



