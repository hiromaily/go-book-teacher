package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/models"
	lg "github.com/hiromaily/golibs/log"
	"github.com/hiromaily/golibs/tmpl"
)

type Slack struct {
	mode string
	url  string
}

type Message struct {
	Text string `json:"text"`
}

//{"text": "New comic book alert! _The Further Adventures of Slackbot_, Volume 1, Issue 3."}

var (
	tmplSlackMsg = `
🤓😎😴 The following tachers are available now! 🤓😎😴
{{range .Teachers}}
*[{{.Name}} / {{.Country}}]*
{{$.URL}}teacher/index/{{.ID}}/
{{end}}
Enjoy!😄

`
)

func NewSlack(conf *config.SlackConfig) *Slack {
	return &Slack{
		mode: "slack",
		url:  getURL(conf.Key),
	}
}

// getURL is to get URL
func getURL(key string) string {
	return fmt.Sprintf("https://hooks.slack.com/services/%s", key)
}

// Send is to send mail
func (s *Slack) Send(ths []models.TeacherInfo) error {
	lg.Debugf("Send by %s", s.mode)

	//make body
	//FIXME: handle as interface
	si := &models.SiteInfo{URL: "http://eikaiwa.dmm.com/", Teachers: ths}
	msg, err := tmpl.StrTempParser(tmplSlackMsg, &si)
	if err != nil {
		return errors.Wrap(err, "fail to parse message for slack")
	}

	//crate json
	sm := Message{Text: msg}
	data, err := json.Marshal(&sm)
	if err != nil {
		return errors.Wrap(err, "fail to call json.Marshal")
	}
	//send
	body, err := s.sendPost(data)
	if err != nil {
		return err
	}
	lg.Debugf("body: %s", string(body))

	return nil
}

func (s *Slack) sendPost(data []byte) ([]byte, error) {

	//1. prepare NewRequest data
	req, err := http.NewRequest(
		"POST",
		s.url,
		//bytes.NewBuffer(jsonStr),
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call http.NewRequest()")
	}

	//2. set http header
	// Content-Type:application/json; charset=utf-8
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	//3. send
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call client.Do()")
	}
	defer resp.Body.Close()

	//5. read response
	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}
