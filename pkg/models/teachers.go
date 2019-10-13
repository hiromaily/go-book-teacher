package models

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

// SiteInfo is site and teacher information
type SiteInfo struct {
	URL      string        `json:"url"`
	Teachers []TeacherInfo `json:"teachers"`
}

// TeacherInfo is json structure of teacher information
type TeacherInfo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

// LoadJSON is to load json file
func LoadJSON(jsonFile string) (*SiteInfo, error) {
	siteInfo := SiteInfo{}
	file, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call ReadFile() %s", jsonFile)
	}
	err = json.Unmarshal(file, &siteInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to Unmarshal json binary: %s", jsonFile)
	}
	//lg.Debugf("SiteInfo.Url: %v", siteInfo.URL)
	//lg.Debugf("SiteInfo.Teachers[0].Id: %d, Name: %s, Country: %s", siteInfo.Teachers[0].ID, siteInfo.Teachers[0].Name, siteInfo.Teachers[0].Country)

	return &siteInfo, nil
}
