package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type data struct {
	id int
	start, end int
}

type record struct {
	cat byte
	id int
	t time.Time
}

var records = make([]record, 0)
var timing = make(map[int]*[60]int)
var sums = make(map[int]int)

func main() {
	file, err := os.Open("data/test.data")
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := Solve(file, true)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func Solve(rd io.Reader, first bool) (res int, err error) {
	err = read(rd)
	if err != nil {
		return
	}

	sort.Sort(byTime(records))
	process()

	if first {
		res = result()
	} else {
		res = result2()
	}
	return
}

func read(rd io.Reader) (err error) {
	reader := bufio.NewReader(rd)
	var line string
	var rec record

	for {
		line, err = reader.ReadString('\n')

		if strings.Contains(line, "begins shift") {
			rec.id, err = parseId(line)
			rec.t, err = parseTime(line)
			rec.cat = 0
		} else if strings.Contains(line, "falls asleep") {
			rec.t, err = parseTime(line)
			rec.cat = 1
		} else if strings.Contains(line, "wakes up") {
			rec.t, err = parseTime(line)
			rec.cat = 2
		}
		records = append(records, rec)

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return
		}
	}
}

func parseTime(line string) (t time.Time, err error) {
	r := regexp.MustCompile(`\[(.*)\]`)
	matches := r.FindStringSubmatch(line)
	if len(matches) != 2 {
		return t, errors.New("error for parsing date")
	}

	layout := "2006-01-02 15:04"
	t, err = time.Parse(layout, matches[1])
	return
}

func parseId(line string) (num int, err error) {
	r := regexp.MustCompile(`#(\d+)`)
	matches := r.FindStringSubmatch(line)
	if len(matches) != 2 {
		return num, errors.New("error for parsing guard id")
	}

	num, err = strconv.Atoi(matches[1])
	return
}

func process() {
	var d data
	for i := 0; i < len(records); i++ {
		switch records[i].cat {
		case 0:
			d.id = records[i].id
		case 1:
			d.start = records[i].t.Minute()
		case 2:
			d.end = records[i].t.Minute()
			update(d)
		}
	}
}

func update(d data) {
	sums[d.id] += d.end - d.start
	if timing[d.id] == nil {
		timing[d.id] = new([60]int)
	}
	for i := d.start; i < d.end; i++ {
		(*timing[d.id])[i]++
	}
}

func result() int {
	max := -1
	id := -1
	for key, val := range sums {
		if val > max {
			max = val
			id = key
		}
	}

	max = -1
	minute := -1
	for i := 0; i < len(timing[id]); i++ {
		if timing[id][i] > max {
			max = timing[id][i]
			minute = i
		}
	}

	return id * minute
}

func result2() int {
	max := -1
	minute := -1
	id := -1
	for key, minutes := range timing {
		for i := 0; i < len(minutes); i++ {
			if minutes[i] > max {
				max = minutes[i]
				minute = i
				id = key
			}
		}
	}

	return id * minute
}

type byTime []record

func (r byTime) Len() int {
	return len(r)
}

func (r byTime) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r byTime) Less(i, j int) bool {
	return r[i].t.Before(r[j].t)
}