package alfred

// AlfredItem 表示 Alfred Workflow 中的一个结果项
type AlfredItem struct {
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	Arg          string `json:"arg,omitempty"`
	Icon         string `json:"icon,omitempty"`         // 每行显示的 icon
	Quicklookurl string `json:"quicklookurl,omitempty"` // 快速预览的URL
}

// GetTitle 返回项目标题
func (item *AlfredItem) GetTitle() string {
	return item.Title
}

// GetSubtitle 返回项目副标题
func (item *AlfredItem) GetSubtitle() string {
	return item.Subtitle
}

// GetValue 返回项目的值
func (item *AlfredItem) GetValue() string {
	return item.Arg
}
