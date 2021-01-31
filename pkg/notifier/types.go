package notifier

type Mode string

const (
	ConsoleMode Mode = "console"
	SlackMode   Mode = "slack"
	DummyMode   Mode = "dummy"
)

func (m Mode) String() string {
	return string(m)
}
