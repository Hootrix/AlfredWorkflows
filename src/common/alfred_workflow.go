package common

import (
	"os"
	"strings"
	"unicode"
)

type AlfredWorkflow struct {
	Args string
	// ActionItemName []string
	Items []AlfredItem
}

func NewAlfredWorkflow() *AlfredWorkflow {
	res := AlfredWorkflow{}
	res.Query(nil)
	return &res
}

func (aw *AlfredWorkflow) Query(args []string) {
	if len(args) < 1 {
		args = os.Args[1:]
	}
	input := strings.Join(args, " ")
	input = strings.TrimFunc(input, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	aw.Args = input
}

func (aw *AlfredWorkflow) AddItem(ItemName string, value string, opts ...func(*AlfredItem)) {
	if value != "" {
		item := AlfredItem{
			Title:    value,
			Subtitle: ItemName,
			Arg:      value,
		}
		for _, opt := range opts {
			opt(&item)
		}

		//仅提交有变化的数据
		if value != aw.Args || strings.Contains(ItemName, "LUCKY NUMBER") { //幸运数字item 必须显示
			aw.Items = append(aw.Items, item)
		}
	}
}
