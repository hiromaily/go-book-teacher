package notifier

import (
	"fmt"
	"os/exec"

	"github.com/hiromaily/go-book-teacher/pkg/models"
	lg "github.com/hiromaily/golibs/log"
)

// NewConsole is to return Console object
func NewConsole() *Console {
	return &Console{mode: "console"}
}

// Console is Console object
type Console struct {
	mode string
}

// Send is notification by executing say command and stdout
// TODO: time should be displayed
func (c *Console) Send(ths []models.TeacherInfo) error {
	lg.Debugf("Send by %s", c.mode)

	//emit a sound
	_ = exec.Command("say", "Found").Start()

	for _, th := range ths {
		fmt.Printf("----------- %s / %s / %d ----------- \n", th.Name, th.Country, th.ID)
	}

	return nil
}
