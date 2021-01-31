package notifier

type Mode string

const (
	ConsoleMode Mode = "console"
	SlackMode   Mode = "slack"
	DummyMode   Mode = "dummy"
)
