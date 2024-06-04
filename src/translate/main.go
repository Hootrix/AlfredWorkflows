package main

import (
	"AlfredWorkflows/common"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// 有道翻译
type YoudaoTranslationResult struct {
	ErrorCode   string   `json:"errorCode"`
	Query       string   `json:"query"`
	Translation []string `json:"translation"`
	L           string   `json:"l"`
	Dict        *Dict    `json:"dict,omitempty"`
	Webdict     *Dict    `json:"webdict,omitempty"`
	TSpeakUrl   string   `json:"tSpeakUrl"`
	SpeakUrl    string   `json:"speakUrl"`
}

// Dict 表示词典和web词典的URL
type Dict struct {
	URL string `json:"url"`
}

// {"code":200,"message":"success","data":"hello,buddy"}
type DeeplxTranslationResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type TranslateAlfredWorkflow struct {
	*Config
	*common.AlfredWorkflow
}

type ConfigItem struct {
	Name      string `yaml:"name"`
	URL       string `yaml:"url,omitempty"`
	Token     string `yaml:"token,omitempty"`
	AppKey    string `yaml:"app_key,omitempty"`
	AppSecret string `yaml:"app_secret,omitempty"`
}

// Config 定义整体配置结构体
type Config struct {
	Services []ConfigItem `yaml:"services"`
	Timeout  int          `yaml:"timeout"`
}

func (taw *TranslateAlfredWorkflow) GetInputQuery() string {
	return taw.Args
}

func main() {

	// 设置日志前缀和标志
	log.SetPrefix("translate workflow log: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	mTranslateAlfredWorkflow := &TranslateAlfredWorkflow{
		AlfredWorkflow: common.NewAlfredWorkflow(),
	}
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
		os.Exit(1)
	}
	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config.yaml")
	log.Println(configPath)
	err = mTranslateAlfredWorkflow.LoadConfig(configPath)
	if err != nil {
		log.Println(err)
	}
	resp := mTranslateAlfredWorkflow.Do()
	resp.Print()
}

func (taw *TranslateAlfredWorkflow) LoadConfig(path string) error {
	result := Config{}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &result)
	if err != nil {
		return err
	}
	taw.Config = &result
	return nil
}

func (taw *TranslateAlfredWorkflow) GetConfigItemWithName(name string) *ConfigItem {
	if len(name) == 0 {
		return nil
	}

	result := ConfigItem{}
	for index, i := range taw.Config.Services {
		if i.Name == name {
			result = taw.Config.Services[index]
		}
	}

	return &result
}

func (taw *TranslateAlfredWorkflow) Do() common.AlfredResponse {
	items := []common.AlfredItem{}
	if taw.Config == nil {
		items = append(items, common.AlfredItem{
			Title:    "config.yaml配置错误",
			Subtitle: "请检查workflow目录的配置文件...",
		})
		return common.AlfredResponse{Items: items}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(taw.Config.Timeout)*time.Second) //总超时控制
	defer cancel()

	datachan := make(chan common.AlfredItem, 10)

	var wg sync.WaitGroup

	wg.Add(1)
	go func(ctx context.Context, itemchan chan common.AlfredItem) {
		defer wg.Done()
		yd := taw.QueryYoudao(ctx, taw.GetInputQuery())
		for _, i := range yd {
			select {
			case itemchan <- i:
			case <-ctx.Done():
				return
			}
		}
	}(ctx, datachan)

	wg.Add(1)
	go func(ctx context.Context, itemchan chan common.AlfredItem) {
		defer wg.Done()
		deeplx := taw.QueryDeeplx(ctx, taw.GetInputQuery())
		for _, i := range deeplx {
			select {
			case itemchan <- i:
			case <-ctx.Done():
				return
			}
		}
	}(ctx, datachan)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(datachan)
		done <- struct{}{}
	}()

	timeout := false
	select {
	case <-done: //正常完成
	case <-ctx.Done(): //超时完成
		timeout = true
	}
	for item := range datachan {
		items = append(items, item)
	}
	//TODO
	if timeout && len(items) == 0 && taw.GetInputQuery() != "" {
		items = append(items, common.AlfredItem{
			Title:    fmt.Sprintf("网络请求超时 %ds", taw.Config.Timeout),
			Subtitle: "请检查网络连接或者重试...",
		})
	}

	return common.AlfredResponse{Items: items}
}
func (taw *TranslateAlfredWorkflow) Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func (taw *TranslateAlfredWorkflow) QueryYoudao(ctx context.Context, argsQuery string) []common.AlfredItem {
	result := []common.AlfredItem{}
	baseYoudao := "https://openapi.youdao.com/api"

	config := taw.GetConfigItemWithName("youdao")
	if config == nil {
		return result
	}

	client := http.Client{
		// Timeout: time.Second * 5,
	}

	url, err := url.Parse(baseYoudao)
	if err != nil {
		return result
	}
	salt := strconv.Itoa(int(time.Now().Unix()))
	sign := taw.Md5(fmt.Sprintf(`%s%s%s%s`, config.AppKey, argsQuery, salt, config.AppSecret))
	query := url.Query()
	query.Set("from", "auto")
	query.Set("to", "auto")
	query.Set("q", argsQuery)
	query.Set("appKey", config.AppKey)
	query.Set("salt", salt)
	query.Set("sign", sign)
	url.RawQuery = query.Encode()

	req, _ := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	resp, err := client.Do(req)
	if err != nil {
		return result
	}
	defer resp.Body.Close()
	if body, err := io.ReadAll(resp.Body); err == nil {

		res := YoudaoTranslationResult{}
		json.Unmarshal(body, &res)
		if res.ErrorCode == "0" {
			// hasChineseChar := taw.HasChineseChar(argsQuery)
			for _, item := range res.Translation {
				// arg := res.Query
				// if hasChineseChar {
				// 	arg = item
				// }
				reviewUrl := ""
				if res.Webdict != nil {
					reviewUrl = res.Webdict.URL
				}
				result = append(result, common.AlfredItem{
					Title:        item,
					Subtitle:     "有道：" + res.Query,
					Arg:          item,
					Quicklookurl: reviewUrl,
				})
			}
		}

	}
	return result
}

func (taw *TranslateAlfredWorkflow) HasChineseChar(str string) bool {
	find, _ := regexp.MatchString("[\u4e00-\u9fa5]", str) // 匹配中文
	return find
}

func (taw *TranslateAlfredWorkflow) QueryDeeplx(ctx context.Context, query string) []common.AlfredItem {
	result := []common.AlfredItem{}

	config := taw.GetConfigItemWithName("deeplx")
	if config == nil {
		return result
	}
	url := config.URL + "?token=" + config.Token

	client := http.Client{
		// Timeout: time.Second * 5,
	}

	souce_map := []string{"auto", "zh"}
	find := taw.HasChineseChar(query)
	if find {
		souce_map = []string{"zh", "en"}
	}

	body := fmt.Sprintf(`{
    "text": "%s",
    "source_lang": "%s",
    "target_lang": "%s"
}`, query, souce_map[0], souce_map[1])
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return result
	}
	defer resp.Body.Close()

	if body, err := io.ReadAll(resp.Body); err == nil {
		res := DeeplxTranslationResult{}
		json.Unmarshal(body, &res)
		if res.Code == 200 {

			// arg := query
			// hasChineseChar := taw.HasChineseChar(query)
			// if hasChineseChar {
			// 	arg = res.Data
			// }

			result = append(result, common.AlfredItem{
				Title:    res.Data,
				Subtitle: "DeeplX：" + query,
				Arg:      res.Data,
			})
		}
	}

	return result
}
