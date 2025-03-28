package translate

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

// YoudaoTranslationResult 有道翻译结果
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

// DeeplxTranslationResult DeepLX翻译结果
type DeeplxTranslationResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// TranslationResult 通用翻译结果
type TranslationResult struct {
	Title    string
	Subtitle string
	Value    string
	Url      *string
}

// Service 翻译服务接口
type Service interface {
	Translate(ctx context.Context, query string) ([]TranslationResult, error)
}

// YoudaoService 有道翻译服务
type YoudaoService struct {
	AppKey    string
	AppSecret string
}

// DeeplxService DeepLX翻译服务
type DeeplxService struct {
	URL   string
	Token string
}

// NewYoudaoService 创建有道翻译服务
func NewYoudaoService(appKey, appSecret string) *YoudaoService {
	return &YoudaoService{
		AppKey:    appKey,
		AppSecret: appSecret,
	}
}

// NewDeeplxService 创建DeepLX翻译服务
func NewDeeplxService(serviceURL, token string) *DeeplxService {
	return &DeeplxService{
		URL:   serviceURL,
		Token: token,
	}
}

// Md5 计算字符串的MD5值
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// HasChineseChar 检查字符串是否包含中文字符
// func HasChineseChar(str string) bool {
// 	for _, r := range str {
// 		if r >= 0x4e00 && r <= 0x9fff {
// 			return true
// 		}
// 	}
// 	return false
// }

func HasChineseChar(str string) bool {
	find, _ := regexp.MatchString("[\u4e00-\u9fa5]", str) // 匹配中文
	return find
}

// Translate 使用有道翻译服务翻译
func (s *YoudaoService) Translate(ctx context.Context, query string) ([]TranslationResult, error) {
	var results []TranslationResult

	salt := strconv.FormatInt(time.Now().Unix(), 10)
	curtime := strconv.FormatInt(time.Now().Unix(), 10)
	sign := Md5(s.AppKey + query + salt + s.AppSecret)

	// to := "en"
	// if HasChineseChar(query) {
	// 	to = "zh-CHS"
	// }

	params := url.Values{}
	params.Add("from", "auto")
	// params.Add("to", to)
	params.Add("to", "auto")
	params.Add("q", query)
	params.Add("appKey", s.AppKey)
	params.Add("salt", salt)
	params.Add("sign", sign)
	// params.Add("signType", "v3")
	params.Add("curtime", curtime)

	apiURL := "https://openapi.youdao.com/api"
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result YoudaoTranslationResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.ErrorCode != "0" {
		return nil, fmt.Errorf("youdao translation error: %s", result.ErrorCode)
	}

	for _, translation := range result.Translation {
		reviewUrl := ""
		if result.Webdict != nil {
			reviewUrl = result.Webdict.URL
		}
		results = append(results, TranslationResult{
			Title:    translation,
			Subtitle: "有道翻译: " + query,
			Value:    translation,
			Url:      &reviewUrl,
		})
	}

	return results, nil
}

// Translate 使用DeepLX翻译服务翻译
func (s *DeeplxService) Translate(ctx context.Context, query string) ([]TranslationResult, error) {
	var results []TranslationResult

	// 确定源语言和目标语言
	sourceLang := "auto"
	targetLang := "en"
	if HasChineseChar(query) {
		targetLang = "en"
	} else {
		targetLang = "zh"
	}

	// 构建请求体
	requestBody := map[string]interface{}{
		"text":        query,
		"source_lang": sourceLang,
		"target_lang": targetLang,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if s.Token != "" {
		req.Header.Set("Authorization", "Bearer "+s.Token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result DeeplxTranslationResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("deeplx translation error: %s", result.Message)
	}

	// 清理结果中的HTML标签
	cleanResult := regexp.MustCompile("<[^>]*>").ReplaceAllString(result.Data, "")

	results = append(results, TranslationResult{
		Title:    cleanResult,
		Subtitle: "DeepLX翻译: " + query,
		Value:    cleanResult,
	})

	return results, nil
}
