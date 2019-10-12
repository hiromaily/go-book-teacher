package notifier

import (
	"fmt"
	lg "github.com/hiromaily/golibs/log"
	"os/exec"

	"github.com/hiromaily/go-book-teacher/pkg/models"
)

// NewCommand is settings for executing say command
func NewCommand() *Command {
	return &Command{mode: "command"}
}

// Command is Command object
type Command struct {
	mode string
}

func (c *Command) Send(ths []models.TeacherInfo) error {
	lg.Debugf("Send by %s", c.mode)

	//emit a sound
	_ = exec.Command("say", "Found").Start()

	for _, th := range ths {
		fmt.Printf("----------- %s / %s / %d ----------- \n", th.Name, th.Country, th.ID)
	}

	return nil
}
