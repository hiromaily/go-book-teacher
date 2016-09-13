package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	conf "github.com/hiromaily/go-book-teacher/config"
	rds "github.com/hiromaily/golibs/db/redis"
	hrk "github.com/hiromaily/golibs/heroku"
	lg "github.com/hiromaily/golibs/log"
)

// StoreRedis is Redis object
type StoreRedis struct {
	RD *rds.RD
}

var (
	redisKey = "bookteacher:save"
	rd       StoreRedis
)

// Setup is settings
func Setup() (*StoreRedis, error) {
	redisURL := conf.GetConf().Redis.URL
	host, pass, port, err := hrk.GetRedisInfo(redisURL)
	if err != nil {
		return nil, err
	}
	rd = StoreRedis{}
	rd.RD = rds.New(host, uint16(port), pass, 0)
	//rd.RD.Connection(0)

	return &rd, nil
}

// Get is to get StoreRedis instance
func Get() *StoreRedis {
	if rd.RD.Pool == nil {
		//panic("Before call this, call New in addition to arguments")
		return nil
	}
	return &rd
}

// Save is to save data on Redis
func (rd *StoreRedis) Save(newData string) bool {
	lg.Debug("Using Redis")

	//close
	//defer rd.RD.Close()

	c := rd.RD.Conn
	val, err := redis.String(c.Do("GET", redisKey))

	if err != nil {
		lg.Errorf("redis error is %s\n", err)
	}
	lg.Debugf("new value is %s, old value is %s\n", newData, val)

	if err != nil || newData != val {
		//save
		c.Do("SET", redisKey, newData)
		return true
	}
	return false
}

// Delete is to delete value by key
func (rd *StoreRedis) Delete() error {
	c := rd.RD.Conn
	_, err := c.Do("DEL", redisKey)
	if err != nil {
		return fmt.Errorf("%s", "delete key on redis is failed.")
	}
	return nil
}
