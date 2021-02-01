package teachers

// Teacher interface
type Teacher interface {
	Fetch() ([]TeacherRepo, error)
}

// TeacherRepo is teacher json structure
type TeacherRepo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}
