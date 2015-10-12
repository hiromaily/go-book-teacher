package main

import (
	"fmt"
	"github.com/hiromaily/golang-libraries/goroutine"
	"github.com/hiromaily/golang-libraries/json"
	"github.com/PuerkitoBio/goquery"
	"log"
	//"reflect"
	st "settings"
	"strconv"
	"strings"
	"time"
)

var num_teachers int = len(st.TEACHERS_ID)
var processing_count int = 0

// Main Processing
func processing(index int) {
	teacher_list := st.TEACHERS_ID[:]

	//HTTP connection
	doc, err := goquery.NewDocument(st.URL + "teacher/index/" + strconv.Itoa(teacher_list[index].Id) + "/")
	if err != nil {
		log.Fatal(err)
		return
	} else {
		parsed_html := perseHtml(doc)

		//show teacher's id, name, date
		fmt.Println("-----------" + teacher_list[index].Name + "/" + teacher_list[index].Country + "/" + strconv.Itoa(teacher_list[index].Id) + "-----------")
		for _, dt := range parsed_html {
			fmt.Println(dt)
		}
	}
}

// For goroutine
func parentProcessing(s chan<- int, core_num int) {
	var idx int = 0
	for processing_count < num_teachers {
		idx = processing_count //set in advance
		processing_count++     //add
		processing(idx)
	}
	goroutine.CallbackGoRoutine(s)
}

// Parse html
func perseHtml(htmldata *goquery.Document) []string {
	var dates []string

	htmldata.Find("a.bt-open").Each(func(_ int, s *goquery.Selection) {
		if json_data, ok := s.Attr("id"); ok {
			//fmt.Println(reflect.TypeOf(json_data))

			//decode
			htmlStringDecode(&json_data)

			//encode txt into json object
			//encoded_json := jsonpkg.JsonEncode(json_data)

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

// Main
func main() {
	fmt.Println("getting teacher's information")

	//Time
	t := time.Now()
	//fmt.Printf("%02d:%02d:%02d\n", t.Hour(), t.Minute(), t.Second())

	//go routine
	c := make(chan int)
	goroutine.RegisterStartRoutine(parentProcessing, c)

	//receiver
	for {
		_, ok := <-c
		if !ok {
			break
		}
	}

	t2 := time.Now()
	//fmt.Printf("%02d:%02d:%02d\n", t2.Hour(), t2.Minute(), t2.Second())

	fmt.Println(t2.Sub(t))

}
