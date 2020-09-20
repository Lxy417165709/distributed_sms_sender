package middleware

import (
	"distributed/sms/src/utils"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"time"
)

const flagOfRedisDoSuccess int64 = 1
const valueOfLock = 1

type DistributedMutex struct {
	connPool          *redis.Pool
	lockKey           string
	lockValue         int
	retryLockInterval time.Duration
}

func NewDistributedMutex(redisHost string, lockKey string, retryLockInterval time.Duration) *DistributedMutex {
	return &DistributedMutex{
		connPool:          utils.GetRedisConnPool(redisHost),
		lockValue:         valueOfLock,
		lockKey:           lockKey,
		retryLockInterval: retryLockInterval,
	}
}

func (d *DistributedMutex) Lock() {
	conn := d.connPool.Get()
	defer func() {
		if err := conn.Close(); err != nil {
			logs.Error(err)
		}
	}()
	for {
		if d.connPool.ActiveCount() >= d.connPool.MaxActive / 2{
			logs.Warn(
				"active conn over 50% max active conn, current active count of conn is %d",
				d.connPool.ActiveCount(),
			)
		}
		reply, err := conn.Do("setnx", d.lockKey, d.lockValue)
		if err != nil {
			logs.Error(err)
			continue
		}
		if reply == flagOfRedisDoSuccess {
			break
		}
		time.Sleep(d.retryLockInterval)
	}
}

func (d *DistributedMutex) Unlock() {
	conn := d.connPool.Get()
	defer func() {
		if err := conn.Close(); err != nil {
			logs.Error(err)
		}
	}()
	reply, err := conn.Do("del", d.lockKey)
	if err != nil {
		logs.Error(err)
		return
	}
	if reply == flagOfRedisDoSuccess {
		return
	}
}
