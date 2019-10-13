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

// MaxGoRoutine is number of goroutine running at the same time
//const MaxGoRoutine uint16 = 20 //FIXME: this should be defined in config

type DMM struct {
	maxGoRoutine int
	url          string
	jsonFile     string
	*models.SiteInfo
	savedTeachers []models.TeacherInfo
}

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

// definedTeachers is defined teachers info
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

func (d *DMM) FetchInitialData() error {
	if d.jsonFile != "" {
		//call json file
		lg.Debugf("Load json file: %s", d.jsonFile)
		siteInfo, err := models.LoadJSON(d.jsonFile)
		if err != nil {
			return err
		}
		d.SiteInfo = siteInfo
	}
	d.SiteInfo = d.definedTeachers()
	return nil
}

func (d *DMM) InitializeSavedTeachers() {
	d.savedTeachers = make([]models.TeacherInfo, 0)
}

func (d *DMM) FindTeachers() []models.TeacherInfo {
	defer tm.Track(time.Now(), "handleTeachers()")

	wg := &sync.WaitGroup{}
	chanSemaphore := make(chan bool, d.maxGoRoutine)

	//d.Teachers
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
			err := d.getHTML(&teacher)
			if err != nil {
				//TODO: this err shouold emit by channel
				lg.Error(err)
			}
		}()
	}
	wg.Wait()

	return d.savedTeachers
}

// GetHTML is to get scraped HTML from web page
func (d *DMM) getHTML(th *models.TeacherInfo) error {
	var flg = false

	//HTTP connection
	url := fmt.Sprintf("%steacher/index/%d/", d.URL, th.ID)
	doc, err := httpdoc.GetHTMLDocs(url)
	if err != nil {
		return errors.Wrapf(err, "fail to call GetHTMLDocs() %s", url)
	} else if isTeacherActive(doc) {
		parsedHTML := perseHTML(doc)

		//show teacher's id, name, date
		fmt.Printf("----------- %s / %s / %d ----------- \n", th.Name, th.Country, th.ID)
		for _, dt := range parsedHTML {
			fmt.Println(dt)
			flg = true
		}
		//save teacher
		if flg {
			//FIXME: mutex
			d.saveTeacer(th)
		}
	} else {
		//no teacher
		fmt.Printf("teacher [%d]%s quit \n", th.ID, th.Name)
	}
	return nil
}

//save teacher id to variable
func (d *DMM) saveTeacer(th *models.TeacherInfo) {
	//FIXME: mutex
	d.savedTeachers = append(d.savedTeachers, *th)
}
