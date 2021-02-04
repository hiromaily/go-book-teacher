package save

import (
	"net"
	"net/url"
	"os"
	"strconv"
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

// NewRedisSaver returns Saver
func NewRedisSaver(logger *zap.Logger, redisURL, env string) (Saver, error) {
	var redisConn redis.Conn
	if redisURL == "" {
		u, err := getEnvURL(logger, env)
		if err != nil {
			return nil, err
		}
		redisConn, err = redis.DialURL(u.String(), redis.DialTLSSkipVerify(true))
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		redisConn, err = redis.DialURL(redisURL)
		if err != nil {
			return nil, err
		}
	}

	return &redisSaver{
		mode:     RedisMode,
		logger:   logger,
		conn:     redisConn,
		redisKey: "bookteacher:save",
	}, nil
}

// `REDIS_URL` is expected in env
// https://devcenter.heroku.com/articles/securing-heroku-redis#using-go
func getEnvURL(logger *zap.Logger, env string) (*url.URL, error) {
	if env == "" {
		return nil, errors.New("env is empty")
	}
	redisURL := os.Getenv(env)
	logger.Debug("getEnvURL", zap.String(env, redisURL))
	u, err := url.Parse(redisURL)
	if err != nil {
		return nil, err
	}

	host, strPort, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, err
	}

	port, _ := strconv.Atoi(strPort)
	port++
	u.Scheme = "rediss"
	u.Host = net.JoinHostPort(host, strconv.Itoa(port))
	logger.Debug("getEnvURL", zap.String("url", u.String()))

	return u, nil
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
