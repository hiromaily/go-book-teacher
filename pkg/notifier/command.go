package notifier

import (
	"os/exec"

	"github.com/hiromaily/go-book-teacher/pkg/models"
)

// NewCommand is settings for executing say command
func NewCommand() *Command {
	return &Command{}
}

// Command is Command object
type Command struct{}

func (c *Command) Send(ths []models.TeacherInfo) error {
	//emit a sound
	_ = exec.Command("say", "Found").Start()

	return nil
}
