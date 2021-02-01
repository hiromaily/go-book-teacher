package save

// Saver interface
type Saver interface {
	Save(string) (bool, error)
	Delete() error
	Close()
}

type Mode string

const (
	TextMode  Mode = "text"
	RedisMode Mode = "redis"
	DummyMode Mode = "dummy"
)

func (m Mode) String() string {
	return string(m)
}
