package utils

import (
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"time"
)

const flagOfRedisDoSuccess int64= 1
const valueOfLock = 1

type DistributedMutex struct {
	connPool *redis.Pool
	lockKey string
	lockValue int
	retryLockInterval time.Duration
}

func NewDistributedMutex(redisHost string,lockKey string,retryLockInterval time.Duration) *DistributedMutex{
	pool := &redis.Pool{
		MaxActive:200,
		Dial:func()(redis.Conn,error) {
			conn,err := redis.Dial("tcp", redisHost)
			if err!=nil{
				logs.Error(err)
				return nil,err
			}
			return conn,nil
		},
		IdleTimeout:time.Second,
		Wait:true,
	}
	return &DistributedMutex{
		connPool: pool,
		lockValue: valueOfLock,
		lockKey:lockKey,
		retryLockInterval:retryLockInterval,
	}
}

func (d *DistributedMutex) Lock() {
	conn := d.connPool.Get()
	defer func() {
		if err := conn.Close();err!=nil{
			logs.Error(err)
		}
	}()
	for {
		//logs.Debug(d.connPool.ActiveCount(),d.connPool.IdleCount())
		reply,err := conn.Do("setnx",d.lockKey,d.lockValue)
		if err != nil{
			logs.Error(err)
			continue
		}
		if reply == flagOfRedisDoSuccess{
			//logs.Info("[%d] 已上锁",i)
			break
		}
		//logs.Info("[%d] 锁住中",i)
		time.Sleep(d.retryLockInterval)
	}
}

func (d *DistributedMutex) Unlock() {
	conn := d.connPool.Get()
	defer func() {
		if err := conn.Close();err!=nil{
			logs.Error(err)
		}
	}()
	reply,err := conn.Do("del",d.lockKey)
	if err != nil{
		logs.Error(err)
		return
	}
	if reply == flagOfRedisDoSuccess{
		//logs.Info("[%d] 已解锁",i)
		return
	}
}

