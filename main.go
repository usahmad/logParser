package main

import (
	"bufio"
	"fmt"
	db2 "logParser/db"
	"logParser/db/models"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const directory = "./timesheetLogs"

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
	time.Sleep(time.Second * 3)
}

func main() {
	fmt.Println("TimeSheet parser by Us.@hmad started blэт")
	defer timeTrack(time.Now(), "Execution")
	fmt.Println("Connecting to db")
	db := db2.InitDb()
	if db == nil {
		panic("DB NOT CONNECTED")
	}
	fmt.Println("Scanning directory")
	files := scanDir()
	for _, file := range files {
		fmt.Println("Parsing file: " + file)
		items := parseFile(file)
		if len(items) > 0 {
			fmt.Println("Inserting Logs")
			for _, log := range items {
				err := log.Create(db, log)
				if err != nil {
					panic(err)
				}
			}
			fmt.Println("Calculating results of: " + file)
			results := calculateHours(items)
			fmt.Println("Inserting Data: " + file)
			for _, result := range results {
				err := result.Create(db, result)
				if err != nil {
					fmt.Println(err)
					panic(err)
				}
			}
			err := os.Remove(file)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println(fmt.Sprintf("Empty File %v", file))
			err := os.Remove(file)
			if err != nil {
				panic(err)
			}
		}
	}
}

func scanDir() []string {
	files, err := filepath.Glob(filepath.Join(directory, "*"))
	if err != nil {
		panic(err)
	}

	return files
}

func calculateHours(items []models.TimeSheetLog) []models.TimeSheet {
	groups := make(map[int][]models.TimeSheetLog)

	for _, l := range items {
		id := l.SIP
		if _, ok := groups[id]; !ok {
			groups[id] = make([]models.TimeSheetLog, 0)
		}
		groups[id] = append(groups[id], l)
	}
	var results []models.TimeSheet
	for sip, data := range groups {
		timeSpent := 0
		in := 0
		date := ""
		for _, item := range data {
			t, err := time.Parse("2006-01-02 15:04:05", item.Date)
			if err != nil {
				panic(err)
			}
			unixTime := int(t.Unix())
			temp := t.Truncate(24 * time.Hour)
			date = temp.Format("02-01-2006")
			if item.Type == "in" {
				in = unixTime
			} else if item.Type == "out" {
				if in != 0 {
					timeSpent += unixTime - in
					in = 0
				}
			}
		}
		if timeSpent != 0 {
			results = append(results, models.TimeSheet{
				SIP:        sip,
				TimeWorked: timeSpent,
				Date:       date,
			})
		}
	}

	return results
}

func parseFile(file string) []models.TimeSheetLog {
	readFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)

	var data []models.TimeSheetLog

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "specific") || strings.Contains(line, "Registered SIP") {
			re := regexp.MustCompile(`\[(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})\].*Registered\s+SIP\s+'(\d{3})'`)
			match := re.FindStringSubmatch(line)
			if match != nil {
				i, err := strconv.Atoi(match[2])
				if err != nil {
					panic(err)
				}
				item := models.TimeSheetLog{
					SIP:  i,
					Date: match[1],
					Type: "in",
				}
				data = append(data, item)
			}
		}
		if strings.Contains(line, "specific") || strings.Contains(line, "Unregistered") {
			re := regexp.MustCompile(`\[(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})\].*Unregistered\s+SIP\s+'(\d{3})'`)
			match := re.FindStringSubmatch(line)
			if match != nil {
				i, err := strconv.Atoi(match[2])
				if err != nil {
					panic(err)
				}
				item := models.TimeSheetLog{
					SIP:  i,
					Date: match[1],
					Type: "out",
				}
				data = append(data, item)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return data
}
