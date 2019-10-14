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
	}
}

// definedTeachers is to return defined teachers info
func (d *DMM) definedTeachers() *models.SiteInfo {
	lg.Debug("use defined data for dmm teacher")
	ti := []models.TeacherInfo{
		{ID: 6214, Name: "Aleksandra S", Country: "Serbia"},
		{ID: 4808, Name: "Joxyly", Country: "Serbia"},
		{ID: 10157, Name: "Patrick B", Country: "New Zealand"},
		{ID: 25473, Name: "Oliver B", Country: "Ireland"}, //quit
		{ID: 25622, Name: "Shannon J", Country: "UK"},     //quit
		{ID: 24397, Name: "Elisabeth L", Country: "USA"},
		{ID: 23979, Name: "Lina Bianca", Country: "USA"},
		{ID: 25070, Name: "Celene", Country: "Australia"},
		{ID: 24721, Name: "Kenzie", Country: "USA"},
		{ID: 27828, Name: "Sanndy", Country: "UK"},
		{ID: 28302, Name: "Danni", Country: "South Africa"},
		{ID: 30216, Name: "Tamm", Country: "UK"},
		{ID: 25302, Name: "Nami", Country: "USA"},
	}

	return &models.SiteInfo{
		URL:      d.url,
		Teachers: ti,
	}
}

// FetchInitialData is to fetch target teacher data
func (d *DMM) FetchInitialData() error {
	if d.jsonFile != "" {
		//call json file
		lg.Debugf("Load json file: %s", d.jsonFile)
		siteInfo, err := models.LoadJSON(d.jsonFile)
		if err != nil {
			return err
		}
		d.SiteInfo = siteInfo
		return nil
	}
	d.SiteInfo = d.definedTeachers()
	return nil
}

// FindTeachers is to find available teachers by scraping web site
func (d *DMM) FindTeachers() []models.TeacherInfo {
	defer tm.Track(time.Now(), "dmm.FindTeachers()")

	wg := &sync.WaitGroup{}
	chanSemaphore := make(chan bool, d.maxGoRoutine)
	chanTh := make(chan *models.TeacherInfo) //response of found teacher by channel

	for _, teacher := range d.Teachers {
		teacher := teacher

		wg.Add(1)
		chanSemaphore <- true

		go func() {
			defer func() {
				<-chanSemaphore
				wg.Done()
			}()
			//concurrent func
			err := d.getHTML(&teacher, chanTh)
			if err != nil {
				//TODO: this err shouold emit by channel
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
func (d *DMM) getHTML(th *models.TeacherInfo, chTh chan *models.TeacherInfo) error {
	var flg = false

	//HTTP connection
	url := fmt.Sprintf("%steacher/index/%d/", d.URL, th.ID)
	doc, err := httpdoc.GetHTMLDocs(url)
	if err != nil {
		return errors.Wrapf(err, "fail to call GetHTMLDocs() %s", url)
	} else if isTeacherActive(doc) {
		parsedHTML := parseDate(doc)

		//show teacher's id, name, date
		fmt.Printf("----------- %s / %s / %d ----------- \n", th.Name, th.Country, th.ID)
		for _, dt := range parsedHTML {
			fmt.Println(dt)
			flg = true
		}
		//send teacher by channel
		if flg {
			chTh <- th
		}
	} else {
		//no teacher
		fmt.Printf("teacher [%d]%s quit \n", th.ID, th.Name)
	}
	return nil
}
