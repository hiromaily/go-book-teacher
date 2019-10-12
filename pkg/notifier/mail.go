package notifier

import (
	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/models"
	ml "github.com/hiromaily/golibs/mail"
	"github.com/hiromaily/golibs/tmpl"
)

// Mail is Mail object
type Mail struct {
	info *ml.Info
}

var (
	tmplMails = `
The following tachers are available now!
{{range .Teachers}}
{{$.URL}}teacher/index/{{.ID}}/ [{{.Name}} / {{.Country}}]
{{end}}
Enjoy!`
	subject = "[ENGLISH LESSON] It's Available."
)

// Setup is settings for sending mail
func NewMail(conf *config.MailConfig) *Mail {
	//get environment variable

	smtp := ml.SMTP{Address: conf.SMTP.Address, Pass: conf.SMTP.Pass,
		Server: conf.SMTP.Server, Port: conf.SMTP.Port}

	info := &ml.Info{ToAddress: []string{conf.MailTo}, FromAddress: conf.MailFrom,
		Subject: subject, Body: "", SMTP: smtp}

	return &Mail{
		info: info,
	}
}

// Send is to send mail
func (m *Mail) Send(ths []models.TeacherInfo) error {
	//make body
	//FIXME: handle as interface
	si := &models.SiteInfo{URL: "http://eikaiwa.dmm.com/", Teachers: ths}
	body, err := tmpl.StrTempParser(tmplMails, &si)
	if err != nil {
		return errors.Wrap(err, "fail to parse message for mail")
	}

	m.info.Body = body
	m.info.SendMail("10s")

	return nil
}
