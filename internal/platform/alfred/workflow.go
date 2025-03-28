package alfred

import (
	"os"
	"strings"
	"unicode"
)

// AlfredWorkflow 表示 Alfred Workflow 的主要结构
type AlfredWorkflow struct {
	Args  string
	Items []AlfredItem
}

// NewWorkflow 创建一个新的 AlfredWorkflow
func NewWorkflow() *AlfredWorkflow {
	res := AlfredWorkflow{}
	res.Query(nil)
	return &res
}

// Query 处理查询参数
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

// AddItem 添加一个项目到工作流中
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

		// 仅提交有变化的数据
		if value != aw.Args || strings.Contains(ItemName, "LUCKY NUMBER") { // 幸运数字item 必须显示
			aw.Items = append(aw.Items, item)
		}
	}
}

// GetResponse 获取工作流的响应
func (aw *AlfredWorkflow) GetResponse() *AlfredResponse {
	resp := &AlfredResponse{
		Items: aw.Items,
	}
	return resp
}
