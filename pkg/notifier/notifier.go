package notifier

import (
	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/models"
)

// Mailer is to Notice
type Notifier interface {
	Send([]models.TeacherInfo) error
}

func NewNotifier(conf *config.Config) Notifier {
	if conf.ValidateSlack() {
		//slack mode
		return NewSlack(conf.Slack)
	}
	if conf.ValidateMail() {
		//mail mode
		return NewMail(conf.Mail)
	}
	//set dummy
	return NewDummy()
}
