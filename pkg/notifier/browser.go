package notifier

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/models"
)

const dmmURL = "http://eikaiwa.dmm.com/" //FIXME: it should be dynamic

// NewBrowser is settings for executing open command
func NewBrowser() *Browser {
	return &Browser{
		url: dmmURL,
	}
}

// Browser is Browser object
type Browser struct {
	url string
}

func (b *Browser) Send(ths []models.TeacherInfo) error {
	for _, t := range ths {
		//during test, it doesn't work.

		//out, err := exec.Command("open /Applications/Google\\ Chrome.app", fmt.Sprintf("http://eikaiwa.dmm.com/teacher/index/%d/", id)).Output()
		err := exec.Command("open", fmt.Sprintf("%steacher/index/%d/", b.url, t.ID)).Start()
		if err != nil {
			return errors.Wrapf(err, "fail to exec.Command(open)")
		}
	}
	return nil
}
