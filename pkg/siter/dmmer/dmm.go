package dmmer

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/httpdoc"
	"github.com/hiromaily/go-book-teacher/pkg/models"
	lg "github.com/hiromaily/golibs/log"
	tm "github.com/hiromaily/golibs/time"
)

// DMM is DMM object
type DMM struct {
	maxGoRoutine int
	url          string
	jsonFile     string
	*models.SiteInfo
	fetcher Fetcher
}

// NewDMM is to return DMM object
func NewDMM(jsonFile, url string, concurrency int) *DMM {
	if concurrency < 2 {
		lg.Warnf("concurrency in config is invalid: %d", concurrency)
		concurrency = 2
	}

	return &DMM{
		maxGoRoutine: concurrency,
		url:          url,
		jsonFile:     jsonFile,
		fetcher:      newFetcher(jsonFile, url),
	}
}

// FetchInitialData is to fetch target teacher data
func (d *DMM) FetchInitialData() error {
	siteInfo, err := d.fetcher.FetchInitialData()
	if err != nil {
		return nil
	}

	d.SiteInfo = siteInfo
	return nil
}

// FindTeachers is to find available teachers by scraping web site
func (d *DMM) FindTeachers(day int) []models.TeacherInfo {
	defer tm.Track(time.Now(), "dmm.FindTeachers()")

	wg := &sync.WaitGroup{}
	chanSemaphore := make(chan bool, d.maxGoRoutine)
	chanTh := make(chan *models.TeacherInfo) // response of found teacher by channel

	for _, teacher := range d.Teachers {
		teacher := teacher

		wg.Add(1)
		chanSemaphore <- true

		go func() {
			defer func() {
				<-chanSemaphore
				wg.Done()
			}()
			// concurrent func
			err := d.getHTML(&teacher, chanTh, day)
			if err != nil {
				// TODO: this err shouold emit by channel
				lg.Error(err)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(chanTh)
	}()

	savedTeachers := make([]models.TeacherInfo, 0)
	// wait until results channel is closed.
	for th := range chanTh {
		savedTeachers = append(savedTeachers, *th)
	}

	return savedTeachers
}

// getHTML is to get teacher information from HTML document
func (d *DMM) getHTML(th *models.TeacherInfo, chTh chan *models.TeacherInfo, day int) error {
	flg := false

	// HTTP connection
	url := fmt.Sprintf("%steacher/index/%d/", d.url, th.ID)
	doc, err := httpdoc.GetHTMLDocs(url)
	if err != nil {
		return errors.Wrapf(err, "fail to call GetHTMLDocs() %s", url)
	} else if isTeacherActive(doc) {
		parsedHTML := parseDate(doc, day)

		// show teacher's id, name, date
		fmt.Printf("----------- %s / %s / %d ----------- \n", th.Name, th.Country, th.ID)
		for _, dt := range parsedHTML {
			fmt.Println(dt)
			flg = true
		}
		// send teacher by channel
		if flg {
			chTh <- th
		}
	} else {
		// no teacher
		fmt.Printf("teacher [%d]%s quit \n", th.ID, th.Name)
	}
	return nil
}
