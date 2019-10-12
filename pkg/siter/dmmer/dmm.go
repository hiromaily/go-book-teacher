package dmmer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/models"
	lg "github.com/hiromaily/golibs/log"
	tm "github.com/hiromaily/golibs/time"
)

// MaxGoRoutine is number of goroutine running at the same time
const MaxGoRoutine uint16 = 20 //FIXME: this should be defined in config
const dmmURL = "http://eikaiwa.dmm.com/"

type DMM struct {
	url      string
	jsonFile string
	*models.SiteInfo
	savedTeachers []models.TeacherInfo
}

func NewDMM(jsonFile string) *DMM {
	return &DMM{
		url:      dmmURL,
		jsonFile: jsonFile,
	}
}

// LoadJSONFile is to load json file
func loadJSON(jsonFile string) (*models.SiteInfo, error) {
	lg.Debugf("load json file: %s", jsonFile)
	siteInfo := models.SiteInfo{}
	file, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call ReadFile() %s", jsonFile)
	}
	err = json.Unmarshal(file, &siteInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to Unmarshal json binary: %s", jsonFile)
	}
	lg.Debugf("SiteInfo.Url: %v", siteInfo.URL)
	lg.Debugf("SiteInfo.Teachers[0].Id: %d, Name: %s, Country: %s", siteInfo.Teachers[0].ID, siteInfo.Teachers[0].Name, siteInfo.Teachers[0].Country)

	return &siteInfo, nil
}

// definedTeachers is defined teachers info
func definedTeachers() *models.SiteInfo {
	lg.Debug("use defined data for dmm teacher")
	ti := []models.TeacherInfo{
		{ID: 6214, Name: "Aleksandra S", Country: "Serbia"},
		{ID: 4808, Name: "Joxyly", Country: "Serbia"},
		{ID: 10157, Name: "Patrick B", Country: "New Zealand"},
		{ID: 25473, Name: "Oliver B", Country: "Ireland"},
		{ID: 25622, Name: "Shannon J", Country: "UK"},
		{ID: 24397, Name: "Elisabeth L", Country: "USA"},
		{ID: 25475, Name: "Dan Cr", Country: "USA"},
	}

	return &models.SiteInfo{
		URL:      dmmURL,
		Teachers: ti,
	}
}

func (d *DMM) FetchInitialData() error {
	if d.jsonFile != "" {
		//call json file
		siteInfo, err := loadJSON(d.jsonFile)
		if err != nil {
			return err
		}
		d.SiteInfo = siteInfo
	}
	d.SiteInfo = definedTeachers()
	return nil
}

func (d *DMM) InitializeSavedTeachers() {
	d.savedTeachers = make([]models.TeacherInfo, 0)
}

func (d *DMM) HandleTeachers() {
	defer tm.Track(time.Now(), "handleTeachers()")

	wg := &sync.WaitGroup{}
	chanSemaphore := make(chan bool, MaxGoRoutine)

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
			d.getHTML(&teacher)
		}()
	}
	wg.Wait()
}

// GetsavedTeachers is to get savedTeachers
func (d *DMM) GetSavedTeachers() []models.TeacherInfo {
	return d.savedTeachers
}

// GetHTML is to get scraped HTML from web page
func (d *DMM) getHTML(th *models.TeacherInfo) {
	var flg = false

	//HTTP connection
	doc, err := goquery.NewDocument(fmt.Sprintf("%steacher/index/%d/", d.URL, th.ID))
	if err != nil {
		lg.Fatal(err) //FIXME: change
		return
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
}

//save teacher id to variable
func (d *DMM) saveTeacer(th *models.TeacherInfo) {
	//FIXME: mutex
	d.savedTeachers = append(d.savedTeachers, *th)
}
