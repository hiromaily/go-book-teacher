package storages

import (
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"

	lg "github.com/hiromaily/golibs/log"
)

// RedisRepo is Redis object
type RedisRepo struct {
	mode string
	RD   redis.Conn
}

var (
	redisKey = "bookteacher:save"
	rd       RedisRepo
)

// NewRedis is to return RedisRepo object
func NewRedis(redisURL string) (*RedisRepo, error) {
	c, err := redis.DialURL(redisURL)
	if err != nil {
		return nil, err
	}
	rd = RedisRepo{
		mode: "redis",
		RD:   c,
	}

	//host, pass, port, err := hrk.GetRedisInfo(redisURL)
	//if err != nil {
	//	return nil, err
	//}
	//rd = RedisRepo{
	//	mode: "redis",
	//	RD:   rds.New(host, uint16(port), pass, 0),
	//}
	//rd.RD.Connection(0)
	//_, err = rd.RD.Conn.Do("SELECT", 0)

	return &rd, nil
}

// Save is to save data on Redis
func (rd *RedisRepo) Save(newData string) (bool, error) {
	lg.Debugf("Save by %s", rd.mode)

	// close
	// defer rd.RD.Close()

	c := rd.RD
	val, err := redis.String(c.Do("GET", redisKey))
	if err != nil && !strings.Contains(err.Error(), "nil returned") {
		return false, errors.Wrapf(err, "fail to call redis.GET by %s", redisKey)
	}
	lg.Debugf("new value is %s, old value is %s", newData, val)

	if newData != val {
		// save
		c.Do("SET", redisKey, newData)
		return true, nil
	}
	return false, nil
}

// Delete is to delete value by key
func (rd *RedisRepo) Delete() error {
	c := rd.RD
	_, err := c.Do("DEL", redisKey)
	if err != nil {
		return fmt.Errorf("%s", "delete key on redis is failed.")
	}
	return nil
}

// Close is to close connection from Redis
func (rd *RedisRepo) Close() {
	if rd != nil && rd.RD != nil {
		rd.RD.Close()
	}
}
