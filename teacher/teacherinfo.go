package teacher

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	lg "github.com/hiromaily/golibs/log"
	"io/ioutil"
	//"os"
	//"path"
)

// Info is json structure for teacher information
type Info struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

// SiteInfo is json structure for teacher information
type SiteInfo struct {
	URL      string `json:"url"`
	Teachers []Info `json:"teachers"`
}

var (
	siteInfo      SiteInfo
	savedTeachers []Info
	printOn       = true
)

// SetPrintOn is to show teacher info using fmt.Printf()
func SetPrintOn(b bool) {
	printOn = b
}

//InitSavedTeachers is to initialize variable for saved teacher
func InitSavedTeachers() {
	savedTeachers = make([]Info, 0)
}

//save teacher id to variable
func saveTeacer(t *Info) {
	savedTeachers = append(savedTeachers, *t)
}

// GetSiteInfo is to get siteInfo
func GetSiteInfo() *SiteInfo {
	return &siteInfo
}

// GetsavedTeachers is to get savedTeachers
func GetsavedTeachers() []Info {
	return savedTeachers
}

// GetHTML is to get scraped HTML from web page
func (t *Info) GetHTML(url string) {
	var flg = false

	//HTTP connection
	doc, err := goquery.NewDocument(fmt.Sprintf("%steacher/index/%d/", url, t.ID))
	if err != nil {
		lg.Fatal(err)
		return
	} else if isTeacherActive(doc) {
		parsedHTML := perseHTML(doc)

		//show teacher's id, name, date
		if printOn {
			fmt.Printf("----------- %s / %s / %d ----------- \n", t.Name, t.Country, t.ID)
		}
		for _, dt := range parsedHTML {
			if printOn {
				fmt.Println(dt)
			}
			flg = true
		}
		//save teacher
		if flg {
			saveTeacer(t)
		}
	} else {
		//no teacher
		if printOn {
			fmt.Printf("teacher [%d]%s quit \n", t.ID, t.Name)
		}
	}
}

// LoadJSONFile is to load json file
func LoadJSONFile(filePath string) *SiteInfo {
	lg.Debug("load json file")
	//initialize
	siteInfo = SiteInfo{}

	//test(get current dir)
	//dir := path.Dir(os.Args[0])
	//lg.Debugf("path.Dir(os.Args[0]): %s", dir)

	// Loading jsonfile
	if filePath == "" {
		//dir := path.Dir(os.Args[0])
		//lg.Debugf("path.Dir(os.Args[0]): %s", dir)
		//filePath = fmt.Sprintf("%s/json/teachers.json", dir)
		lg.Fatal("json filepath have to be set.")
		return nil
	}

	file, _ := ioutil.ReadFile(filePath)
	err := json.Unmarshal(file, &siteInfo)

	if err != nil {
		lg.Fatalf("json format is invalid: %v, filepath is %s", err, filePath)
		return nil
	}
	lg.Debugf("SiteInfo.Url: %v", siteInfo.URL)
	lg.Debugf("SiteInfo.Teachers[0].Id: %d, Name: %s, Country: %s", siteInfo.Teachers[0].ID, siteInfo.Teachers[0].Name, siteInfo.Teachers[0].Country)

	return &siteInfo
}

// GetDefinedData is registered teacher info
func GetDefinedData() *SiteInfo {
	ti := []Info{
		{ID: 6214, Name: "Aleksandra S", Country: "Serbia"},
		{1381, "Anna O", "Rossiya"},
		{2464, "Emilia", "Serbia"},
		{4107, "Milica Ml", "Serbia"},
		{4806, "Jekica", "Serbia"},
		{4808, "Joxyly", "Serbia"},
		{5252, "Gagga", "Serbia"},
		{5380, "Olivera V", "Serbia"},
		{5656, "Lavinija", "Serbia"},
		{6294, "Milica J", "Serbia"},
		{6550, "Yovana", "Serbia"},
		{7646, "Kaytee", "Serbia"},
		{8160, "Gaja", "Serbia"},
		{3486, "Indre", "Lithuania"},
		{7093, "Rita M", "Portugal"},
		{8519, "Marine", "France"},
	}
	siteInfo = SiteInfo{URL: "http://eikaiwa.dmm.com/", Teachers: ti}
	return &siteInfo
}

// CreateSiteInfo is to get SiteInfo
func CreateSiteInfo(ti []Info) *SiteInfo {
	return &SiteInfo{URL: "http://eikaiwa.dmm.com/", Teachers: ti}
}
