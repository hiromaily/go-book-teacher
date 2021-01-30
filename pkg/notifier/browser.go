package notifier

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/models"
	lg "github.com/hiromaily/golibs/log"
)

// NewBrowser is to return Browser object
func NewBrowser(url string) *Browser {
	return &Browser{
		mode: "browser",
		url:  url,
	}
}

// Browser is Browser object
type Browser struct {
	mode string
	url  string
}

// Send is notification by executing open command
// Note:during test, it should not use
func (b *Browser) Send(ths []models.TeacherInfo) error {
	lg.Debugf("Send by %s", b.mode)

	for _, t := range ths {
		// out, err := exec.Command("open /Applications/Google\\ Chrome.app", fmt.Sprintf("http://eikaiwa.dmm.com/teacher/index/%d/", id)).Output()
		err := exec.Command("open", fmt.Sprintf("%steacher/index/%d/", b.url, t.ID)).Start()
		if err != nil {
			return errors.Wrapf(err, "fail to exec.Command(open)")
		}
	}
	return nil
}
