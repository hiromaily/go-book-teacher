package storage

type Mode string

const (
	TextMode  Mode = "text"
	RedisMode Mode = "redis"
	DummyMode Mode = "dummy"
)
