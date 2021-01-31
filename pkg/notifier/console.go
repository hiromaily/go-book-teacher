package notifier

import (
	"fmt"
	"os/exec"

	"go.uber.org/zap"

	"github.com/hiromaily/go-book-teacher/pkg/models"
)

// console object
type console struct {
	mode   Mode
	logger *zap.Logger
}

// NewConsole returns Notifier interface
func NewConsole(logger *zap.Logger) Notifier {
	return &console{
		mode:   ConsoleMode,
		logger: logger,
	}
}

// Notify notifies on console
// TODO: time should be displayed
func (c *console) Notify(ths []models.TeacherInfo) error {
	c.logger.Debug("notify", zap.String("mode", c.mode.String()))

	// emit a sound (maybe macOS only)
	_ = exec.Command("say", "Found").Start()

	for _, th := range ths {
		fmt.Printf("----------- %s / %s / %d ----------- \n", th.Name, th.Country, th.ID)
	}

	return nil
}
