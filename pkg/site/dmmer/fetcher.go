package dmmer

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-book-teacher/pkg/models"
)

// Fetcher is interface to fetch initial data
type Fetcher interface {
	FetchInitialData() (*models.SiteInfo, error)
}

func newFetcher(logger *zap.Logger, jsonFile, url string) Fetcher {
	if jsonFile != "" {
		return newJSONFetcher(logger, jsonFile)
	}
	return newDummyFetcher(logger, url)
}

// ----------------------------------------------------------------------------
// jsonFetcher
// ----------------------------------------------------------------------------

type jsonFetcher struct {
	logger   *zap.Logger
	jsonFile string
}

func newJSONFetcher(logger *zap.Logger, jsonFile string) Fetcher {
	return &jsonFetcher{
		logger:   logger,
		jsonFile: jsonFile,
	}
}

// FetchInitialData is to return siteInfo by loading json
func (f *jsonFetcher) FetchInitialData() (*models.SiteInfo, error) {
	// call json file
	f.logger.Debug("load json file", zap.String("file", f.jsonFile))
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
	logger *zap.Logger
	url    string
}

func newDummyFetcher(logger *zap.Logger, url string) Fetcher {
	return &dummyFetcher{
		logger: logger,
		url:    url,
	}
}

// FetchInitialData is to return defined teachers info
func (d *dummyFetcher) FetchInitialData() (*models.SiteInfo, error) {
	d.logger.Debug("use defined data for dmm teacher")
	ti := []models.TeacherInfo{
		{ID: 6214, Name: "Aleksandra S", Country: "Serbia"},
		{ID: 4808, Name: "Joxyly", Country: "Serbia"},
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
