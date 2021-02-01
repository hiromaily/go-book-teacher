package notifier

import (
	"github.com/hiromaily/go-book-teacher/pkg/teachers"
)

// Notifier interface
type Notifier interface {
	Notify([]teachers.TeacherRepo) error
}
