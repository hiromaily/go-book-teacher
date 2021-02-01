package save

// Saver interface
type Saver interface {
	Save(string) (bool, error)
	Delete() error
	Close()
}

// Mode is save mode
type Mode string

const (
	// TextMode is text mode
	TextMode Mode = "text"
	// RedisMode is redis mode
	RedisMode Mode = "redis"
	// DummyMode is dummy mode
	DummyMode Mode = "dummy"
)

// String converts Mode to string
func (m Mode) String() string {
	return string(m)
}
