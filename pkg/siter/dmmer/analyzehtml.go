package dmmer

import (
	"encoding/json"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const timeRangeFrom string = "00:00:00"
const timeRangeTo string = "24:00:00"

// check within range for applicable time
func isTimeApplicable(strDate string) bool {
	//e.g. 2016-02-27 03:30:00
	strTarget := strings.Split(strDate, " ")[1]

	return strTarget >= timeRangeFrom && strTarget <= timeRangeTo
}

// Html text replace into json
func htmlStringDecode(jsondata *string) {
	lst := [10][2]string{
		{"&amp;", "&"},
		{"&lt;", "<"},
		{"&gt;", ">"},
		{"&quot;", "\""},
		{"a:3:", ""},
		{"a:4:", ""},
		{"s:", "\"field"},
		{";", ","},
		{":\"", "\":\""},
		{",N,", ""},
	}

	for _, data := range lst {
		*jsondata = strings.Replace(*jsondata, data[0], data[1], -1)
	}
}

// Check html (empty or not)
func isTeacherActive(htmldata *goquery.Document) bool {
	ret := htmldata.Find("#fav_count").Text()
	return ret != ""
}

// Parse html
func perseHTML(htmldata *goquery.Document) []string {
	var dates []string

	htmldata.Find("a.bt-open").Each(func(_ int, s *goquery.Selection) {
		if jsonData, ok := s.Attr("id"); ok {

			//decode
			htmlStringDecode(&jsonData)

			//analyze json object
			var jsonObject map[string]interface{}
			//json.JsonAnalyze(jsonData, &jsonObject)
			json.Unmarshal([]byte(jsonData), &jsonObject)

			//extract date from json object
			//e.g. 2016-02-27 03:30:00

			strDate := jsonObject["field19"].(string)
			if isTimeApplicable(strDate) {
				dates = append(dates, strDate)
			}
		}
	})

	return dates
}