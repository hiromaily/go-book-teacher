package dmmer

import (
	"github.com/hiromaily/go-book-teacher/pkg/models"
	lg "github.com/hiromaily/golibs/log"
)

// Fetcher is interface to fetch initial data
type Fetcher interface {
	FetchInitialData() (*models.SiteInfo, error)
}

func newFetcher(jsonFile, url string) Fetcher {
	if jsonFile != "" {
		return newJSONFetcher(jsonFile)
	}
	return newDummyFetcher(url)
}

// ----------------------------------------------------------------------------
// jsonFetcher
// ----------------------------------------------------------------------------

type jsonFetcher struct {
	jsonFile string
}

func newJSONFetcher(jsonFile string) Fetcher {
	return &jsonFetcher{
		jsonFile: jsonFile,
	}
}

// FetchInitialData is to return siteInfo by loading json
func (f *jsonFetcher) FetchInitialData() (*models.SiteInfo, error) {
	//call json file
	lg.Debugf("Load json file: %s", f.jsonFile)
	siteInfo, err := models.LoadJSON(f.jsonFile)
	if err != nil {
		return nil, err
	}
	return siteInfo, nil
}

// ----------------------------------------------------------------------------
// dummy fetcher
// ----------------------------------------------------------------------------

type dummyFetcher struct {
	url string
}

func newDummyFetcher(url string) Fetcher {
	return &dummyFetcher{
		url: url,
	}
}

// FetchInitialData is to return defined teachers info
func (d *dummyFetcher) FetchInitialData() (*models.SiteInfo, error) {
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
		{ID: 32141, Name: "Colleen Marie", Country: "USA"},
	}

	return &models.SiteInfo{
		URL:      d.url,
		Teachers: ti,
	}, nil
}
