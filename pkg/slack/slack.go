package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	conf "github.com/hiromaily/go-book-teacher/pkg/config"
	th "github.com/hiromaily/go-book-teacher/pkg/teacher"
	lg "github.com/hiromaily/golibs/log"
	"github.com/hiromaily/golibs/tmpl"
	"io/ioutil"
	"net/http"
)

type SlackMsg struct {
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

// Send is to send mail
func Send(ths []th.Info) error {
	//make body
	si := th.CreateSiteInfo(ths)
	msg, err := tmpl.StrTempParser(tmplSlackMsg, &si)
	if err != nil {
		lg.Debugf("slack couldn't be send caused by err : %s\n", err)
	} else {
		//crate json
		sm := SlackMsg{Text: msg}
		data, err := json.Marshal(&sm)
		if err != nil {
			return fmt.Errorf("[ERROR] When calling `json.Marshal`: %v\n", err)
		}
		//send
		body, err := sendPost(data, getURL())
		if err != nil {
			return err
		}
		lg.Debugf("body: %s", string(body))
	}
	return nil
}

// getURL is to get URL
func getURL() string {
	return fmt.Sprintf("https://hooks.slack.com/services/%s", conf.GetConf().Slack.Key)
}

func sendPost(data []byte, url string) ([]byte, error) {

	//1. prepare NewRequest data
	req, err := http.NewRequest(
		"POST",
		url,
		//bytes.NewBuffer(jsonStr),
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] When calling `http.NewRequest`: %v", err)
	}

	//2. set http header
	// Content-Type:application/json; charset=utf-8
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	//3. send
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] When calling `client.Do`: %v", err)
	}
	defer resp.Body.Close()

	//5. read response
	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}
