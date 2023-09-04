package main

import (
	"encoding/json"
	"fmt"
	"github.com/apognu/gocal"
	"io"
	"net/http"
	"strings"
	"time"
)

type Lesson struct {
	Start time.Time
	End   time.Time
	Where string
	Name  string
}

func constructMessage(lessons []Lesson) string {
	outString := fmt.Sprintf("Сегодня (%v) %v пар:\n", time.Now().Format("02-01"), len(lessons))
	for _, v := range lessons {
		outString += fmt.Sprintf("%v-%v %v %v\n", v.Start.Format("15:04"), v.End.Format("15:04"), v.Name, v.Where)
	}
	return outString

}

func getSchedule() []Lesson {
	client := &http.Client{}
	// TODO: Change hardcoded value of a group to a variable
	req, err := http.NewRequest("GET", "https://schedule-of.mirea.ru/_next/data/eUrUSSNTpalv5LqTcfVyv/index.json?s=1_820", nil)
	if err != nil {
		panic(err)
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	schedule := data["pageProps"].(map[string]interface{})["scheduleLoadInfo"].([]interface{})[0].(map[string]interface{})["iCalContent"].(string)
	return parseLessonsFromICal(strings.NewReader(schedule))
}

func parseLessonsFromICal(r io.Reader) []Lesson {

	start, end := time.Now(), time.Now().Add(24*time.Hour)

	c := gocal.NewParser(r)
	c.Start, c.End = &start, &end
	c.Parse()
	lessons := make([]Lesson, len(c.Events))
	for i, e := range c.Events {
		lessons[i] = Lesson{Start: *e.Start, End: *e.End, Where: e.Location, Name: e.Summary}
	}
	return lessons
}

func main() {
	fmt.Println(constructMessage(getSchedule()))

}
