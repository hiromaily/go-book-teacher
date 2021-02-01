package save

import (
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// redisSaver object
type redisSaver struct {
	mode     Mode
	logger   *zap.Logger
	conn     redis.Conn
	redisKey string
}

// NewRedis returns Saver
func NewRedisSaver(logger *zap.Logger, redisURL string) (Saver, error) {
	redisConn, err := redis.DialURL(redisURL)
	if err != nil {
		return nil, err
	}

	return &redisSaver{
		mode:     RedisMode,
		logger:   logger,
		conn:     redisConn,
		redisKey: "bookteacher:save",
	}, nil
}

// Save saves data on Redis
func (r *redisSaver) Save(newData string) (bool, error) {
	r.logger.Debug("save", zap.String("mode", r.mode.String()))

	currentData, err := redis.String(r.conn.Do("GET", r.redisKey))
	if err != nil && !strings.Contains(err.Error(), "nil returned") {
		return false, errors.Wrapf(err, "fail to call r.conn.Do(`GET`) key: %s", r.redisKey)
	}
	r.logger.Debug("value",
		zap.String("new", newData),
		zap.String("old", currentData),
	)
	if newData != currentData {
		// save
		r.conn.Do("SET", r.redisKey, newData)
		return true, nil
	}
	return false, nil
}

// Delete deletes data from redis
func (r *redisSaver) Delete() error {
	_, err := r.conn.Do("DEL", r.redisKey)
	if err != nil {
		return errors.New("fail to call r.conn.Do(`DEL`)")
	}
	return nil
}

// Close closes redis connection
func (r *redisSaver) Close() {
	if r != nil && r.conn != nil {
		r.conn.Close()
	}
}
