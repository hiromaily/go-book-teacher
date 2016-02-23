package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	st "github.com/hiromaily/booking-teacher/settings"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	//"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const MaxGoRoutine uint16 = 20
const OpenFileName string = "/tmp/status.log"

const timeRangeFrom string = "07:00:00"
const timeRangeTo string = "24:00:00"

var teachersNum int = len(st.TEACHERS_ID)
var savedTeacherIds []int

var m *sync.Mutex

// managing semaphore
func manageSemaphore() {
	var processingCount int = 0
	var finishedCount int = 0

	chanSemaphore := make(chan bool, MaxGoRoutine)

	for {
		chanSemaphore <- true

		//TODO:condition when to break is necessary
		if processingCount != teachersNum {
			m.Lock()
			idx := processingCount //set in advance
			processingCount++      //add
			m.Unlock()

			//chanSemaphore <- true
			go func(index int) {
				defer func() {
					<-chanSemaphore
				}()
				handleHtmlProcessing(index)
				finishedCount++
			}(idx)
		} else if finishedCount == teachersNum {
			close(chanSemaphore)
			break
		} else {
			continue
		}
	}
}

// Handling Html Processing
func handleHtmlProcessing(index int) {
	var flg bool = false
	teacher_list := st.TEACHERS_ID[:]

	//HTTP connection
	doc, err := goquery.NewDocument(st.URL + "teacher/index/" + strconv.Itoa(teacher_list[index].Id) + "/")
	if err != nil {
		log.Fatal(err)
		return
	} else if isTeacherActive(doc) {
		parsed_html := perseHtml(doc)

		//show teacher's id, name, date
		fmt.Printf("----------- %s / %s / %d ----------- \n", teacher_list[index].Name, teacher_list[index].Country, teacher_list[index].Id)
		for _, dt := range parsed_html {
			fmt.Println(dt)
			flg = true
		}
		//save teacher
		if flg {
			saveTeacerId(teacher_list[index].Id)
		}
	} else {
		//no teacher
		fmt.Printf("teacher %s quit \n", teacher_list[index].Name)
	}
}

// Check html (empty or not)
func isTeacherActive(htmldata *goquery.Document) bool {
	ret := htmldata.Find("#fav_count").Text()
	return ret != ""
}

// check within range for applicable time
func isTimeApplicable(strDate string) bool {
	//e.g. 2016-02-27 03:30:00
	strTarget := strings.Split(strDate, " ")[1]

	return strTarget >= timeRangeFrom && strTarget <= timeRangeTo
}

// Parse html
func perseHtml(htmldata *goquery.Document) []string {
	var dates []string

	htmldata.Find("a.bt-open").Each(func(_ int, s *goquery.Selection) {
		if jsonData, ok := s.Attr("id"); ok {
			//fmt.Println(reflect.TypeOf(jsonData))

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

// Html text replace into json
func htmlStringDecode(jsondata *string) {
	lst := [9][2]string{
		{"&amp;", "&"},
		{"&lt;", "<"},
		{"&gt;", ">"},
		{"&quot;", "\""},
		{"a:3:", ""},
		{"s:", "\"field"},
		{";", ","},
		{":\"", "\":\""},
		{",}", "}"},
	}

	for _, data := range lst {
		*jsondata = strings.Replace(*jsondata, data[0], data[1], -1)
	}
}

//save teacher id to variable
func saveTeacerId(id int) {
	savedTeacherIds = append(savedTeacherIds, id)
}

//save teacher status to log
func saveStatus(ids []int) bool {

	//create string from ids slice
	var sum int = 0
	for index := range ids {
		sum += ids[index]
	}
	newData := strconv.Itoa(sum)

	//open saved log
	fp, err := os.OpenFile(OpenFileName, os.O_CREATE, 0664)
	if err == nil {
		scanner := bufio.NewScanner(fp)
		scanner.Scan()
		filedata := scanner.Text()

		if newData == filedata {
			return false
		}
	} else {
		panic(err.Error())
	}

	//save latest info
	content := []byte(newData)
	ioutil.WriteFile(OpenFileName, content, 0664)

	return true
}

//open browser on PC
func openBrowser(ids []int) {
	for index := range ids {
		//out, err := exec.Command("open /Applications/Google\\ Chrome.app", fmt.Sprintf("http://eikaiwa.dmm.com/teacher/index/%d/", id)).Output()
		err := exec.Command("open", fmt.Sprintf("http://eikaiwa.dmm.com/teacher/index/%d/", ids[index])).Start()
		if err != nil {
			panic(fmt.Sprintf("open browser error: %v", err))
		}
	}
}

//serial processing
func serialProcessing() {
	//Time
	t := time.Now()
	//fmt.Printf("%02d:%02d:%02d\n", t.Hour(), t.Minute(), t.Second())

	//
	manageSemaphore()

	//open browser
	//fmt.Println("open browser")
	if len(savedTeacherIds) != 0 {
		//save status
		openFlg := saveStatus(savedTeacherIds)
		fmt.Println(openFlg)
		if openFlg {
			openBrowser(savedTeacherIds)
		}
	}
	//reset
	savedTeacherIds = nil

	t2 := time.Now()
	//fmt.Printf("%02d:%02d:%02d\n", t2.Hour(), t2.Minute(), t2.Second())
	fmt.Println(t2.Sub(t))
}

// Main
func main() {
	fmt.Println("getting teacher's information")

	m = new(sync.Mutex)

	for {
		serialProcessing()
		time.Sleep(60 * time.Second)
	}
}
