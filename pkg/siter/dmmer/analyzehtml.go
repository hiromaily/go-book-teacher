package dmmer

import (
	"encoding/json"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const timeRangeFrom string = "00:00:00"
const timeRangeTo string = "24:00:00"

// isTimeApplicable is to check within range for applicable time
func isTimeApplicable(strDate string) bool {
	//e.g. 2016-02-27 03:30:00
	strTarget := strings.Split(strDate, " ")[1]

	return strTarget >= timeRangeFrom && strTarget <= timeRangeTo
}

// htmlStringDecode is replace string in HTML into json
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

// isTeacherActive to check HTML (empty or not)
func isTeacherActive(htmldata *goquery.Document) bool {
	ret := htmldata.Find("#fav_count").Text()
	return ret != ""
}

// parseDate is to parse html date
func parseDate(htmldata *goquery.Document) []string {
	var dates []string

	htmldata.Find("a.bt-open").Each(func(_ int, s *goquery.Selection) {
		if jsonData, ok := s.Attr("id"); ok {

			//decode
			htmlStringDecode(&jsonData)

			//analyze json object
			var jsonObject map[string]interface{}
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
