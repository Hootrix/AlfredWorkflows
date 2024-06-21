package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"AlfredWorkflows/common"

	"github.com/araddon/dateparse"
)

func main() {

	items := []common.AlfredItem{}
	// result := make([]AlfredItem, 0)

	args := os.Args[1:]
	switch len(args) {
	case 0:
		now := time.Now()
		time := now.Format(time.DateTime)
		ts := now.Unix()
		items = append(items, common.AlfredItem{Title: strconv.FormatInt(ts, 10), Subtitle: strconv.FormatInt(ts, 10), Arg: strconv.FormatInt(ts, 10)})
		items = append(items, common.AlfredItem{Title: time, Subtitle: time, Arg: time})
	default:
		input := strings.Join(args, " ")
		input = strings.TrimFunc(input, func(r rune) bool {
			return unicode.IsSpace(r)
		})

		//时间戳
		if regex := regexp.MustCompile(`^\s*(\d+)\s*$`); regex != nil {
			matchs := regex.FindStringSubmatch(input)
			if len(matchs) > 1 {
				ts, _ := strconv.ParseInt(matchs[1], 10, 64)
				tm := time.Unix(ts, 0)
				items = append(items, common.AlfredItem{Title: tm.Format(time.DateTime), Subtitle: tm.Format(time.DateTime), Arg: tm.Format(time.DateTime)})
			}
		}

		//其他格式
		if len(items) < 1 {
			if tm, err := dateparse.ParseLocal(input); err == nil {
				items = append(items, common.AlfredItem{Title: tm.Format(time.DateTime), Subtitle: tm.Format(time.DateTime), Arg: tm.Format(time.DateTime)})
				items = append(items, common.AlfredItem{Title: strconv.FormatInt(tm.Unix(), 10), Subtitle: strconv.FormatInt(tm.Unix(), 10), Arg: strconv.FormatInt(tm.Unix(), 10)})
			}
		}

	}

	if len(items) < 1 {
		items = append(items, common.AlfredItem{Title: "ERROR", Subtitle: "nil", Arg: ""})
	}
	// fmt.Printf("%+v\n", result)

	result, _ := json.Marshal(common.AlfredResponse{
		Items: items,
	})
	fmt.Println(string(result))
}
