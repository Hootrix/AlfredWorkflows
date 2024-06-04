package common

import (
	"encoding/json"
	"fmt"
)

type AlfredResponse struct {
	Items []AlfredItem `json:"items"`
}

func (resp *AlfredResponse) Print() {
	result, _ := json.Marshal(resp)
	fmt.Println(string(result))
}

type AlfredItem struct {
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	Arg          string `json:"arg,omitempty"`
	Icon         string `json:"icon,omitempty"`         // 每行显示的 icon
	Quicklookurl string `json:"quicklookurl,omitempty"` // 快速预览的URL
}
