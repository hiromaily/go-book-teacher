package models

type SiteInfo struct {
	URL      string        `json:"url"`
	Teachers []TeacherInfo `json:"teachers"`
}

// Info is json structure for teacher information
type TeacherInfo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}
