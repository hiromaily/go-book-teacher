package dmm

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-book-teacher/pkg/httpdoc"
	"github.com/hiromaily/go-book-teacher/pkg/site"
	"github.com/hiromaily/go-book-teacher/pkg/teachers"
)

// MaxGoroutine is number of goroutine to run scraping func
const MaxGoroutine = 20

// siteDMM object
type siteDMM struct {
	teacherRepo  teachers.Teacher
	teachers     []teachers.TeacherRepo
	logger       *zap.Logger
	maxGoRoutine int
	siteURL      string
	day          int
}

// NewDMM returns Siter interface
func NewDMM(
	logger *zap.Logger,
	teacherRepo teachers.Teacher,
	siteURL string,
	day int,
) site.Siter {
	return &siteDMM{
		teacherRepo:  teacherRepo,
		logger:       logger,
		maxGoRoutine: MaxGoroutine,
		siteURL:      siteURL,
		day:          day,
	}
}

// Fetch fetches target teachers to search schedule
func (d *siteDMM) Fetch() error {
	teachers, err := d.teacherRepo.Fetch()
	if err != nil {
		return err
	}

	d.teachers = teachers
	return nil
}

// FindTeachers finds available teachers by scraping web site
func (d *siteDMM) FindTeachers() []teachers.TeacherRepo {
	// defer times.Track(time.Now(), "dmm.FindTeachers()")

	wg := &sync.WaitGroup{}
	chSemaphore := make(chan struct{}, d.maxGoRoutine)
	chTeacher := make(chan teachers.TeacherRepo) // response of found teacher by channel

	for _, teacher := range d.teachers {
		teacher := teacher

		wg.Add(1)
		chSemaphore <- struct{}{}

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
	// wait until chTeacher channel is closed.
	for teacher := range chTeacher {
		savedTeachers = append(savedTeachers, teacher)
	}

	return savedTeachers
}

// getHTML gets teacher information from HTML document and send found teacher by channel
func (d *siteDMM) getHTML(teacher *teachers.TeacherRepo, chTeacher chan teachers.TeacherRepo) error {
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
