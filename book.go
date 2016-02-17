package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	st "github.com/hiromaily/booking-teacher/settings"
	"github.com/hiromaily/golang-libraries/goroutine"
	"github.com/hiromaily/golang-libraries/json"
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

var num_teachers int = len(st.TEACHERS_ID)
var processing_count int = 0
var m *sync.Mutex
var savedTeacherIds []int

// Main Processing
func processing(index int) {
	var flg bool = false
	teacher_list := st.TEACHERS_ID[:]

	//HTTP connection
	doc, err := goquery.NewDocument(st.URL + "teacher/index/" + strconv.Itoa(teacher_list[index].Id) + "/")
	if err != nil {
		log.Fatal(err)
		return
	} else if checkHtml(doc) {
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

// For goroutine
func parentProcessing(s chan<- int, core_num int) {
	var idx int = 0
	for processing_count < num_teachers {
		m.Lock()
		idx = processing_count //set in advance
		processing_count++     //add
		m.Unlock()
		processing(idx)
	}
	goroutine.CallbackGoRoutine(s)
}

// Check html (empty or not)
func checkHtml(htmldata *goquery.Document) bool {
	ret := htmldata.Find("#fav_count").Text()
	return ret != ""
}

// Parse html
func perseHtml(htmldata *goquery.Document) []string {
	var dates []string

	htmldata.Find("a.bt-open").Each(func(_ int, s *goquery.Selection) {
		if json_data, ok := s.Attr("id"); ok {
			//fmt.Println(reflect.TypeOf(json_data))

			//decode
			htmlStringDecode(&json_data)

			//analyze json object
			var json_object map[string]interface{}
			json.JsonAnalyze(json_data, &json_object)

			//extract date from json object
			dates = append(dates, json_object["field19"].(string))
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

func saveTeacerId(id int) {
	savedTeacherIds = append(savedTeacherIds, id)
}

func openBrowser(ids []int) {
	for index := range ids {
		//out, err := exec.Command("open /Applications/Google\\ Chrome.app", fmt.Sprintf("http://eikaiwa.dmm.com/teacher/index/%d/", id)).Output()
		err := exec.Command("open", fmt.Sprintf("http://eikaiwa.dmm.com/teacher/index/%d/", ids[index])).Start()
		if err != nil {
			panic(fmt.Sprintf("open browser error: %v", err))
		}
	}
}

//save teacher status to log
func saveStatus(ids []int) bool {
	openFileName := "./status.log"

	//create string from ids slice
	var sum int = 0
	for index := range ids {
		sum += ids[index]
	}
	newData := strconv.Itoa(sum)

	//open saved log
	fp, err := os.OpenFile(openFileName, os.O_CREATE, 0666)
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
	//ioutil.WriteFile(openFileName, content, os.ModePerm)
	ioutil.WriteFile(openFileName, content, 0666)

	return true
}

// Main
func main() {
	fmt.Println("getting teacher's information")

	m = new(sync.Mutex)

	//Time
	t := time.Now()
	//fmt.Printf("%02d:%02d:%02d\n", t.Hour(), t.Minute(), t.Second())

	//go routine
	c := make(chan int)
	goroutine.RegisterStartRoutine(parentProcessing, c, 20)

	//receiver
	for {
		_, ok := <-c
		if !ok {
			break
		}
	}

	//open browser
	if len(savedTeacherIds) != 0 {
		//save status
		openFlg := saveStatus(savedTeacherIds)
		fmt.Println(openFlg)
		if openFlg {
			openBrowser(savedTeacherIds)
		}
	}

	t2 := time.Now()
	//fmt.Printf("%02d:%02d:%02d\n", t2.Hour(), t2.Minute(), t2.Second())
	fmt.Println(t2.Sub(t))

}
