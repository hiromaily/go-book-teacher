package notifier

import (
	"fmt"
	"os/exec"

	"go.uber.org/zap"

	"github.com/hiromaily/go-book-teacher/pkg/teachers"
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
func (c *console) Notify(teachers []teachers.TeacherRepo) error {
	c.logger.Debug("notify", zap.String("mode", c.mode.String()), zap.Any("teachers", teachers))

	// emit a sound (maybe macOS only)
	_ = exec.Command("say", "Found").Start()

	fmt.Println("teachers are found !!")
	for _, teacher := range teachers {
		fmt.Printf("----------- %s / %s / %d ----------- \n", teacher.Name, teacher.Country, teacher.ID)
	}

	return nil
}
