package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type DistributedCache struct {
	hashLoop *HashLoop
}

func NewDistributedCache(hashLoop *HashLoop) *DistributedCache{
	return &DistributedCache{
		hashLoop: hashLoop,
	}
}

func (d *DistributedCache) Set(key string, value interface{}) error {
	conn := d.getConnOfStoringKey(key)
	defer func() {
		if err := conn.Close(); err != nil {
			logs.Error(err)
		}
	}()
	valueJson, err := json.Marshal(value)
	if err != nil {
		logs.Error(err)
		return err
	}
	if _, err := conn.Do("set", key, valueJson); err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (d *DistributedCache) Get(key string) (interface{}, error) {
	conn := d.getConnOfStoringKey(key)
	defer func() {
		if err := conn.Close(); err != nil {
			logs.Error(err)
		}
	}()
	reply, err := conn.Do("get", key)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return reply, nil
}

func (d *DistributedCache) getConnOfStoringKey(key string) redis.Conn {
	aroundRedisPoolNum := getHashModNum(key, d.hashLoop.getMaxNum()+1)
	redisPool := d.hashLoop.GetRedisPool(aroundRedisPoolNum)
	return redisPool.Get()
}

func getHashModNum(hashedStr string, mod int64) int64 {
	const md5HashStringLength = 32
	md5Hash := md5.New()
	md5Hash.Write([]byte(hashedStr))
	hexHashStr := hex.EncodeToString(md5Hash.Sum(nil))
	var hashNum int64
	var countOfPart = 4
	for i:=0;i<countOfPart;i++{
		hashNumPart, err := strconv.ParseInt(
			hexHashStr[i*(md5HashStringLength/countOfPart):(i+1)*(md5HashStringLength/countOfPart)],
			16,
			64,
		)
		if err != nil {
			logs.Error(err)
			return 0
		}
		hashNum += hashNumPart%mod
	}
	return hashNum %mod
}
