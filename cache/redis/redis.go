package redis

import (
	"dataProcess/config"
	"dataProcess/constants"
	"dataProcess/logger"
	dbDo "dataProcess/tcp_socket/db_do"
	"github.com/garyburd/redigo/redis"
)

func initConnect() (c redis.Conn, err error) {
	if len(config.CacheConfig.Host) == 0 || len(config.CacheConfig.Port) == 0 {
		c, err = redis.Dial("tcp", dbDo.RedisConfig.Host+constants.Colon+dbDo.RedisConfig.Port)
	} else {
		c, err = redis.Dial("tcp", config.CacheConfig.Host+constants.Colon+config.CacheConfig.Port)
	}
	if err != nil {
		logger.Error("conn redis failed,", err)
		return
	}
	return
}

func SetOperate(key, val string, expireTime int) (err error) {
	c, err := initConnect()
	if err != nil {
		logger.Error("conn redis failed,", err)
		return
	}
	defer c.Close() // 关闭链接

	_, err = c.Do("Set", key, val)
	if expireTime != -1 {
		_, err = c.Do("expire", key, expireTime)
	}
	if err != nil {
		logger.Error(err.Error())
		return
	}
	return
}

func GetOperate(key string) (r string, err error) {
	c, err := initConnect()
	if err != nil {
		logger.Error("conn redis failed,", err)
		return
	}
	defer c.Close() // 关闭链接

	r, err = redis.String(c.Do("Get", key))
	if err != nil {
		return
	}
	return
}

func DelOperate(key string) (err error) {
	c, err := initConnect()
	if err != nil {
		logger.Error("conn redis failed,", err)
		return
	}
	defer c.Close() // 关闭链接

	_, err = c.Do("Del", key)
	if err != nil {
		return
	}
	return
}
