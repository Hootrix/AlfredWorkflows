package core

// Item 表示一个通用的结果项
type Item interface {
	GetTitle() string
	GetSubtitle() string
	GetValue() string
}

// Response 表示一个通用的响应
type Response interface {
	AddItem(item Item)
	Print()
}

// Command 表示一个通用的命令
type Command interface {
	Execute(args []string) Response
}
