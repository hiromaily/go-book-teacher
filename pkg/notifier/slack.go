package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-book-teacher/pkg/models"
	"github.com/hiromaily/go-book-teacher/pkg/tmpl"
)

// slack object
type slack struct {
	mode          Mode
	logger        *zap.Logger
	slackURL      string
	targetSiteURL string
}

// Message object
type Message struct {
	Text string `json:"text"`
}

var tmplSlackMsg = `
🤓😎😴 The following tachers are available now! 🤓😎😴
{{range .Teachers}}
*[{{.Name}} / {{.Country}}]*
{{$.URL}}teacher/index/{{.ID}}/
{{end}}
Enjoy!😄

`

// NewSlack is to return Slack object
func NewSlack(logger *zap.Logger, key string, targetSiteURL string) Notifier {
	return &slack{
		mode:          SlackMode,
		logger:        logger,
		slackURL:      fmt.Sprintf("https://hooks.slack.com/services/%s", key),
		targetSiteURL: targetSiteURL,
	}
}

// Send is notification by Slack
func (s *slack) Notify(ths []models.TeacherInfo) error {
	s.logger.Debug("notify", zap.String("mode", s.mode.String()))

	// make body
	// FIXME: handle as interface
	si := &models.SiteInfo{URL: s.targetSiteURL, Teachers: ths}
	msg, err := tmpl.StrTempParser(tmplSlackMsg, &si)
	if err != nil {
		return errors.Wrap(err, "fail to parse message for slack")
	}

	// crate json
	sm := Message{Text: msg}
	data, err := json.Marshal(&sm)
	if err != nil {
		return errors.Wrap(err, "fail to call json.Marshal")
	}
	// send
	body, err := s.sendPost(data)
	if err != nil {
		return err
	}
	s.logger.Debug("body", zap.String("body", string(body)))

	return nil
}

func (s *slack) sendPost(data []byte) ([]byte, error) {
	// 1. prepare NewRequest data
	req, err := http.NewRequest(
		"POST",
		s.slackURL,
		// bytes.NewBuffer(jsonStr),
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call http.NewRequest()")
	}

	// 2. set http header
	// Content-Type:application/json; charset=utf-8
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// 3. send
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call client.Do()")
	}
	defer resp.Body.Close()

	// 4. read response
	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}
