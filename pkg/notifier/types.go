package notifier

// Mode is notifier mode
type Mode string

const (
	// ConsoleMode is console mode
	ConsoleMode Mode = "console"
	// SlackMode is slack mode
	SlackMode Mode = "slack"
	// DummyMode is dummy mode
	DummyMode Mode = "dummy"
)

// String converts Mode to string
func (m Mode) String() string {
	return string(m)
}
