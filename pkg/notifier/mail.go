package notifier

import (
	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/models"
	lg "github.com/hiromaily/golibs/log"
	ml "github.com/hiromaily/golibs/mail"
	"github.com/hiromaily/golibs/tmpl"
)

// Mail is Mail object
type Mail struct {
	mode string
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

// NewMail is to return Mail object
func NewMail(conf *config.MailConfig) *Mail {
	// lg.Debug(conf.SMTP.Address)
	// lg.Debug(conf.SMTP.Pass)
	// lg.Debug(conf.SMTP.Server)
	// lg.Debug(conf.MailTo)
	// lg.Debug(conf.MailFrom)

	smtp := ml.SMTP{
		Address: conf.SMTP.Address, Pass: conf.SMTP.Pass,
		Server: conf.SMTP.Server, Port: conf.SMTP.Port,
	}

	info := &ml.Info{
		ToAddress: []string{conf.MailTo}, FromAddress: conf.MailFrom,
		Subject: subject, Body: "", SMTP: smtp,
	}

	return &Mail{
		mode: "mail",
		info: info,
	}
}

// Send is notification by sending mail
func (m *Mail) Send(ths []models.TeacherInfo) error {
	lg.Debugf("Send by %s", m.mode)

	// make body
	// FIXME: handle as interface
	si := &models.SiteInfo{URL: "http://eikaiwa.dmm.com/", Teachers: ths}
	body, err := tmpl.StrTempParser(tmplMails, &si)
	if err != nil {
		return errors.Wrap(err, "fail to parse message for mail")
	}

	m.info.Body = body
	m.info.SendMail("10s")

	return nil
}
