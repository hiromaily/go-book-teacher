package mail

import (
	conf "github.com/hiromaily/go-book-teacher/config"
	th "github.com/hiromaily/go-book-teacher/teacher"
	lg "github.com/hiromaily/golibs/log"
	ml "github.com/hiromaily/golibs/mail"
	"github.com/hiromaily/golibs/tmpl"
)

var (
	mi *ml.Info

	tmplMails = `
The following tachers are available now!
{{range .Teachers}}
{{$.URL}}teacher/index/{{.ID}}/ [{{.Name}} / {{.Country}}]
{{end}}
Enjoy!`
)

// Setup is settings for sending mail
func Setup() {
	//get environment variable
	subject := "[ENGLISH LESSON] It's Available."
	body := ""
	//mails
	smt := conf.GetConf().Mail.SMTP
	m := conf.GetConf().Mail

	smtp := ml.SMTP{Address: smt.Address, Pass: smt.Pass,
		Server: smt.Server, Port: smt.Port}

	mi = &ml.Info{ToAddress: []string{m.MailTo}, FromAddress: m.MailFrom,
		Subject: subject, Body: body, SMTP: smtp}
}

// Send is to send mail
func Send(ths []th.Info) {
	//make body
	si := th.CreateSiteInfo(ths)
	body, err := tmpl.StrTempParser(tmplMails, &si)
	if err != nil {
		lg.Debugf("mail couldn't be send caused by err : %s\n", err)
	} else {
		mi.Body = body
		mi.SendMail("10s")
	}
}
