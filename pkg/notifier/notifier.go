package notifier

import (
	"github.com/hiromaily/go-book-teacher/pkg/models"
)

// Notifier interface
type Notifier interface {
	Notify([]models.TeacherInfo) error
}
