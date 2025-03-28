package alfred

import (
	"encoding/json"
	"fmt"
)

// AlfredResponse 表示 Alfred Workflow 的响应
type AlfredResponse struct {
	Items []AlfredItem `json:"items"`
}

// NewResponse 创建一个新的 AlfredResponse
func NewResponse() *AlfredResponse {
	return &AlfredResponse{
		Items: []AlfredItem{},
	}
}

// AddItem 添加一个项目到响应中
func (resp *AlfredResponse) AddItem(item AlfredItem) {
	resp.Items = append(resp.Items, item)
}

// Print 将响应打印为 JSON 格式
func (resp *AlfredResponse) Print() {
	result, _ := json.Marshal(resp)
	fmt.Println(string(result))
}
