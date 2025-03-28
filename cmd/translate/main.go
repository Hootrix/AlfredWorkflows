package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"AlfredWorkflows/internal/core/translate"
	"AlfredWorkflows/internal/platform/alfred"
	"AlfredWorkflows/pkg/utils"

	"gopkg.in/yaml.v3"
)

// TranslateWorkflow 翻译工作流结构体
type TranslateWorkflow struct {
	Config   *translate.Config
	Workflow *alfred.AlfredWorkflow
}

// NewTranslateWorkflow 创建新的翻译工作流
func NewTranslateWorkflow() *TranslateWorkflow {
	return &TranslateWorkflow{
		Config:   &translate.Config{},
		Workflow: alfred.NewWorkflow(),
	}
}

// GetInputQuery 获取输入查询
func (tw *TranslateWorkflow) GetInputQuery() string {
	return tw.Workflow.Args
}

// LoadConfig 加载配置文件
func (tw *TranslateWorkflow) LoadConfig(path string) error {
	// 如果没有指定配置文件路径，则使用默认路径
	if path == "" {
		execPath, err := os.Executable()
		if err != nil {
			return err
		}

		path = filepath.Join(filepath.Dir(execPath), "config.yaml")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, tw.Config)
}

// Execute 执行翻译
func (tw *TranslateWorkflow) Execute() *alfred.AlfredResponse {
	query := tw.GetInputQuery()

	if utils.IsEmpty(query) {
		item := alfred.AlfredItem{
			Title:    "请输入要翻译的文本",
			Subtitle: "支持中英文互译",
			Arg:      "",
		}
		tw.Workflow.Items = append(tw.Workflow.Items, item)
		return tw.Workflow.GetResponse()
	}

	// 创建上下文，设置超时
	timeout := time.Duration(tw.Config.Timeout) * time.Second
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 创建结果通道
	resultChan := make(chan alfred.AlfredItem, 10)
	var wg sync.WaitGroup

	// 查询有道翻译
	wg.Add(1)
	go func(ctx context.Context, itemChan chan<- alfred.AlfredItem) {
		defer wg.Done()
		youdaoConfig := tw.Config.GetConfigItemWithName("youdao")
		if youdaoConfig != nil && youdaoConfig.AppKey != "" && youdaoConfig.AppSecret != "" {
			service := translate.NewYoudaoService(youdaoConfig.AppKey, youdaoConfig.AppSecret)
			results, err := service.Translate(ctx, query)
			if err == nil {
				for _, result := range results {
					u := ""
					if result.Url != nil {
						u = *result.Url
					}
					item := alfred.AlfredItem{
						Title:        result.Title,
						Subtitle:     result.Subtitle,
						Arg:          result.Value,
						Quicklookurl: u,
					}
					// 发送结果到通道，同时检查上下文是否已取消
					select {
					case itemChan <- item:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}(ctx, resultChan)

	// 查询DeepLX翻译
	wg.Add(1)
	go func(ctx context.Context, itemChan chan<- alfred.AlfredItem) {
		defer wg.Done()
		deeplxConfig := tw.Config.GetConfigItemWithName("deeplx")
		if deeplxConfig != nil && deeplxConfig.URL != "" {
			service := translate.NewDeeplxService(deeplxConfig.URL, deeplxConfig.Token)
			results, err := service.Translate(ctx, query)
			if err == nil {
				for _, result := range results {
					item := alfred.AlfredItem{
						Title:    result.Title,
						Subtitle: result.Subtitle,
						Arg:      result.Value,
					}
					// 发送结果到通道，同时检查上下文是否已取消
					select {
					case itemChan <- item:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}(ctx, resultChan)

	// 创建一个通道用于通知所有goroutine已完成
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(resultChan)
		done <- struct{}{}
	}()

	// 收集结果
	var allItems []alfred.AlfredItem
	timeoutOccurred := false

	// 等待所有翻译完成或超时
	select {
	case <-done: // 所有翻译正常完成
	case <-ctx.Done(): // 超时
		timeoutOccurred = true
	}

	// 从通道中获取所有结果
	for item := range resultChan {
		allItems = append(allItems, item)
	}

	// 如果没有结果，显示错误信息
	if len(allItems) == 0 {
		if timeoutOccurred {
			allItems = append(allItems, alfred.AlfredItem{
				Title:    fmt.Sprintf("翻译超时 %d秒", int(timeout.Seconds())),
				Subtitle: "请检查网络连接或稍后重试",
				Arg:      "",
			})
		} else {
			allItems = append(allItems, alfred.AlfredItem{
				Title:    "翻译失败",
				Subtitle: "请检查网络连接和配置",
				Arg:      "",
			})
		}
	}

	tw.Workflow.Items = allItems
	return tw.Workflow.GetResponse()
}

func main() {
	tw := NewTranslateWorkflow()

	// 加载配置文件
	configPath := filepath.Join(filepath.Dir(os.Args[0]), "config.yaml")
	if err := tw.LoadConfig(configPath); err != nil {
		log.Printf("加载配置文件失败: %v", err)
	}

	// 处理命令行参数
	tw.Workflow.Query(os.Args[1:])

	// 执行翻译并输出结果
	response := tw.Execute()
	response.Print()
}
