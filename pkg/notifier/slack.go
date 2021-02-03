package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-book-teacher/pkg/teachers"
)

// slack object
type slack struct {
	mode          Mode
	logger        *zap.Logger
	slackURL      string
	targetSiteURL string
}

// NewSlack returns Notifier
func NewSlack(logger *zap.Logger, key string, targetSiteURL string) Notifier {
	return &slack{
		mode:          SlackMode,
		logger:        logger,
		slackURL:      fmt.Sprintf("https://hooks.slack.com/services/%s", key),
		targetSiteURL: targetSiteURL,
	}
}

// TeacherInfo is used for template parameter
type TeacherInfo struct {
	URL      string
	Teachers []teachers.TeacherRepo
}

// Message object of slack
type Message struct {
	Text string `json:"text"`
}

// Send is notification by Slack
func (s *slack) Notify(teachers []teachers.TeacherRepo) error {
	s.logger.Debug("notify", zap.String("mode", s.mode.String()))

	// make body
	msgBody, err := templateParser(&TeacherInfo{
		URL:      s.targetSiteURL,
		Teachers: teachers,
	})
	if err != nil {
		return errors.Wrap(err, "fail to parse message body of slack")
	}

	// crate json
	jsonMsg, err := json.Marshal(&Message{Text: msgBody})
	if err != nil {
		return errors.Wrap(err, "fail to call json.Marshal")
	}
	// send
	body, err := s.post(jsonMsg)
	if err != nil {
		return err
	}
	s.logger.Debug("slack_body", zap.String("body", string(body)))

	return nil
}

func (s *slack) post(jsonMsg []byte) ([]byte, error) {
	req, err := http.NewRequest(
		"POST",
		s.slackURL,
		bytes.NewReader(jsonMsg),
	)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call http.NewRequest()")
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call client.Do()")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}

// template
func templateParser(params interface{}) (string, error) {
	slackTemplate := `
🤓😎😴 The following tachers are available now! 🤓😎😴
{{range .Teachers}}
*[{{.Name}} / {{.Country}}]*
{{$.URL}}teacher/index/{{.ID}}/
{{end}}
Enjoy!😄

`

	var writer bytes.Buffer
	tpl := template.Must(template.New("tpl").Parse(slackTemplate))
	if err := tpl.Execute(&writer, params); err != nil {
		return "", err
	}

	return writer.String(), nil
}
