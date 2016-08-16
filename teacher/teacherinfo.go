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

type TeacherInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

type SiteInfo struct {
	Url      string        `json:"url"`
	Teachers []TeacherInfo `json:"teachers"`
}

var siteInfo SiteInfo

var savedTeachers []TeacherInfo

//init saved teacher
func InitSavedTeachers() {
	savedTeachers = make([]TeacherInfo, 0)
}

//save teacher id to variable
func saveTeacer(t *TeacherInfo) {
	savedTeachers = append(savedTeachers, *t)
}

func GetSiteInfo() *SiteInfo {
	return &siteInfo
}

func GetsavedTeachers() []TeacherInfo {
	return savedTeachers
}

func (t *TeacherInfo) GetHTML(url string) {
	var flg bool = false

	//HTTP connection
	doc, err := goquery.NewDocument(fmt.Sprintf("%steacher/index/%d/", url, t.Id))
	if err != nil {
		lg.Fatal(err)
		return
	} else if isTeacherActive(doc) {
		parsed_html := perseHtml(doc)

		//show teacher's id, name, date
		fmt.Printf("----------- %s / %s / %d ----------- \n", t.Name, t.Country, t.Id)
		for _, dt := range parsed_html {
			fmt.Println(dt)
			flg = true
		}
		//save teacher
		if flg {
			saveTeacer(t)
		}
	} else {
		//no teacher
		fmt.Printf("teacher [%d]%s quit \n", t.Id, t.Name)
	}
}

func LoadJsonFile(filePath string) *SiteInfo {
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
		//filePath = fmt.Sprintf("%s/settings.json", dir)
		lg.Fatal("json filepath have to be set.")
		return nil
	}

	file, err := ioutil.ReadFile(filePath)
	err = json.Unmarshal(file, &siteInfo)

	if err != nil {
		lg.Fatalf("json format is invalid: %v, filepath is %s", err, filePath)
		return nil
	}
	lg.Debugf("SiteInfo.Url: %v", siteInfo.Url)
	lg.Debugf("SiteInfo.Teachers[0].Id: %d, Name: %s, Country: %s", siteInfo.Teachers[0].Id, siteInfo.Teachers[0].Name, siteInfo.Teachers[0].Country)

	return &siteInfo
}

func GetDefinedData() *SiteInfo {
	ti := []TeacherInfo{
		{Id: 6214, Name: "Aleksandra S", Country: "Serbia"},
		{1381, "Anna O", "Rossiya"},
		{1411, "Jekaterina", "Latvia"},
		{2464, "Emilia", "Serbia"},
		//{3141, "Janja", "Serbia"},
		{3293, "Edina", "Serbia"},
		{4107, "Milica Ml", "Serbia"},
		//{4565, "Mana", "Serbia"},
		{4806, "Jekica", "Serbia"},
		{4808, "Joxyly", "Serbia"},
		{5252, "Gagga", "Serbia"},
		{5380, "Olivera V", "Serbia"},
		{5656, "Lavinija", "Serbia"},
		//{5809, "Sandra Z", "Serbia"},
		{6294, "Milica J", "Serbia"},
		{6550, "Yovana", "Serbia"},
		{7002, "Sanja J", "Serbia"},
		//{7888, "Ducica", "Serbia"},
		{8101, "Milica Ja", "Serbia"},
		{8146, "Nejla K", "Serbia"},
		{8160, "Gaja", "Serbia"},
		//{8358, "Sanndra", "Serbia"},
		{9250, "Maria", "Serbia"},
		{11307, "Katherine B", "Serbia"},
		{3486, "Indre", "Lithuania"},
		{3645, "Egliukas", "Lithuania"},
		{7093, "Rita M", "Portugal"},
		//{6466, "Rachel L", "Shingapore"},
		{8519, "Marine", "France"},
		//{11283, "Alice Burthe", "France"},
		{11854, "Daniela D", "Malta"},
		//{8261, "Ela T", "Germany"},
		//{8912, "Pascale", "Netherland"},
	}
	siteInfo = SiteInfo{Url: "http://eikaiwa.dmm.com/", Teachers: ti}
	return &siteInfo
}

func CreateSiteInfo(ti []TeacherInfo) *SiteInfo {
	return &SiteInfo{Url: "http://eikaiwa.dmm.com/", Teachers: ti}
}
