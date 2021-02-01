package dmm

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-book-teacher/pkg/httpdoc"
	"github.com/hiromaily/go-book-teacher/pkg/teachers"
	"github.com/hiromaily/go-book-teacher/pkg/times"
)

const MaxGoroutine = 20

// DMM object
type DMM struct {
	teacherFetcher teachers.Teacher
	teachers       []teachers.TeacherRepo
	logger         *zap.Logger
	maxGoRoutine   int
	siteURL        string
	day            int
}

// NewDMM is to return DMM object
func NewDMM(
	logger *zap.Logger,
	teacherFetcher teachers.Teacher,
	siteURL string,
	day int,
) *DMM {
	return &DMM{
		teacherFetcher: teacherFetcher,
		logger:         logger,
		maxGoRoutine:   MaxGoroutine,
		siteURL:        siteURL,
		day:            day,
	}
}

// Fetch fetches target teachers to search schedule
func (d *DMM) Fetch() error {
	teachers, err := d.teacherFetcher.Fetch()
	if err != nil {
		return nil
	}

	d.teachers = teachers
	return nil
}

// FindTeachers finds available teachers by scraping web site
func (d *DMM) FindTeachers() []teachers.TeacherRepo {
	defer times.Track(time.Now(), "dmm.FindTeachers()")

	wg := &sync.WaitGroup{}
	chSemaphore := make(chan bool, d.maxGoRoutine)
	chTeacher := make(chan teachers.TeacherRepo) // response of found teacher by channel

	for _, teacher := range d.teachers {
		teacher := teacher

		wg.Add(1)
		chSemaphore <- true

		go func() {
			defer func() {
				<-chSemaphore
				wg.Done()
			}()
			// concurrent func
			err := d.getHTML(&teacher, chTeacher)
			if err != nil {
				d.logger.Error("fail to call getHTML()",
					zap.Any("target_teacher", teacher),
					zap.Error(err),
				)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(chTeacher)
	}()

	savedTeachers := make([]teachers.TeacherRepo, 0)
	// wait until results channel is closed.
	for teacher := range chTeacher {
		savedTeachers = append(savedTeachers, teacher)
	}

	return savedTeachers
}

// getHTML is to get teacher information from HTML document
func (d *DMM) getHTML(teacher *teachers.TeacherRepo, chTeacher chan teachers.TeacherRepo) error {

	// HTTP connection
	targetURL := fmt.Sprintf("%steacher/index/%d/", d.siteURL, teacher.ID)
	doc, err := httpdoc.GetHTMLDocs(targetURL)
	if err != nil {
		return errors.Wrapf(err, "fail to call GetHTMLDocs() %s", targetURL)
	}
	if !isTeacherActive(doc) {
		// no teacher
		fmt.Printf("teacher [%d]%s quit \n", teacher.ID, teacher.Name)
		return nil
	}

	parsedHTML := parseDate(doc, d.day)

	// show teacher's id, name, date
	fmt.Printf("----------- %s / %s / %d ----------- \n", teacher.Name, teacher.Country, teacher.ID)
	var isFound bool
	for _, dt := range parsedHTML {
		fmt.Println(dt)
		isFound = true
	}
	// send teacher by channel
	if isFound {
		chTeacher <- *teacher
	}

	return nil
}
