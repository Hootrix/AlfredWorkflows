package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"AlfredWorkflows/internal/core/timestamp"
	"AlfredWorkflows/internal/platform/alfred"
	"AlfredWorkflows/pkg/utils"
)

func main() {
	// 创建 Alfred 工作流
	workflow := alfred.NewWorkflow()

	args := os.Args[1:]
	if len(args) == 0 {
		// 没有参数，显示当前时间戳和格式化时间
		ts, timeStr := timestamp.GetCurrentTimestamp()
		workflow.AddItem("当前时间戳", strconv.FormatInt(ts, 10))
		workflow.AddItem("当前时间", timeStr)
	} else {
		// 处理输入参数
		input := strings.Join(args, " ")
		input = utils.TrimSpace(input)

		// 尝试解析时间戳
		if regex := regexp.MustCompile(`^\s*(\d+)\s*$`); regex != nil {
			matches := regex.FindStringSubmatch(input)
			if len(matches) > 1 {
				ts, _ := strconv.ParseInt(matches[1], 10, 64)
				timeStr := timestamp.TimestampToTime(ts)
				workflow.AddItem("转换后的时间", timeStr)
			}
		}

		// 如果没有匹配到时间戳，尝试解析其他格式的时间
		if len(workflow.Items) < 1 {
			if tm, err := timestamp.ParseTimeString(input); err == nil {
				workflow.AddItem("格式化时间", tm.Format("2006-01-02 15:04:05"))
				workflow.AddItem("Unix时间戳", timestamp.FormatUnixTimestamp(tm.Unix()))
			}
		}
	}

	// 如果没有任何结果，显示错误信息
	if len(workflow.Items) < 1 {
		workflow.AddItem("错误", "无法解析输入")
	}

	// 输出结果
	workflow.GetResponse().Print()
}
