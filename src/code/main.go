package main

import (
	"AlfredWorkflows/common"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"

	// _ "net/http/pprof"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
)

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	mCodeAlfredWorkflow := CodeAlfredWorkflow{
		AlfredWorkflow: common.NewAlfredWorkflow(),
	}
	resp := mCodeAlfredWorkflow.Do()
	// fmt.Println(mCodeAlfredWorkflow.ActionItemName)
	resp.Print()

}

type CodeAlfredWorkflow struct {
	// Query      string
	// ActionName []string
	*common.AlfredWorkflow
}

func (caw *CodeAlfredWorkflow) Do() common.AlfredResponse {

	caw.AddItem("Length", caw.Length())
	caw.AddItem("upper", strings.ToUpper(caw.Args))
	caw.AddItem("lower", strings.ToLower(caw.Args))

	caw.AddItem("reverse", caw.Reverse())

	number, lucky := caw.LuckyNumber()
	pre := "❌"
	if lucky {
		pre = "✅"
	}
	caw.AddItem(pre+"LUCKY NUMBER", number)

	caw.AddItem("MD5", caw.Md5())
	caw.AddItem("SHA256", caw.SHA256())
	caw.AddItem("EncodeBase64", caw.EncodeBase64())
	caw.AddItem("DecodeBase64", caw.DecodeBase64())
	caw.AddItem("EncodeStandardURL", caw.EncodeStandardURL())
	caw.AddItem("EncodeAllURL", caw.EncodeAllURL())
	caw.AddItem("DecodeURL", caw.DecodeURL())
	caw.AddItem(`ToHEX`, caw.ToHEX())
	caw.AddItem(`FromHEX`, caw.FromHEX())
	caw.AddItem(`EncodeHTMLEntities`, caw.EncodeHTMLEntities())
	caw.AddItem(`DecodeHTMLEntities`, caw.DecodeHTMLEntities())
	utf16 := caw.UnicodeEscapeUTF16()
	utf32 := caw.UnicodeEscapeUTF32()
	caw.AddItem(`UnicodeUTF16Escape 转义`, utf16)
	if utf32 != utf16 {
		caw.AddItem(`UnicodeUTF32Escape 转义`, utf32)
	}

	// support mix UTF16/UTF32
	caw.AddItem(`UnicodeUnEscape 兼容UTF16/UTF32 反转义`, caw.UnicodeUnEscape())

	//U+XXXX 混合
	// 😄1😄2😄#😄¥ <==> U+1F6041U+1F6042U+1F604#U+1F604U+00A5

	//UTF32混合UTF16代理对
	// 哈😄你好😄1 <==> \u54C8\uD83D\uDE04\u4F60\u597D\U0001F604\u0031
	// 啊哈哈哈哈😄你好😀 <==> \u554a\u54c8\u54c8\u54c8\u54c8\ud83d\ude04\u4f60\u597d\U0001F600

	return common.AlfredResponse{
		Items: caw.Items,
	}
}

func (caw *CodeAlfredWorkflow) Length() string {
	args := caw.Args
	return strconv.Itoa(len([]rune(args)))
}

func (caw *CodeAlfredWorkflow) Reverse() string {
	runes := []rune(caw.Args)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// 幸运数字计算
func (caw *CodeAlfredWorkflow) LuckyNumber() (string, bool) {
	args := caw.Args
	find, _ := regexp.MatchString(`^\d+$`, args)

	if find {
		runes := []rune(args)
		var result int
		for {
			result = 0
			for index := range runes {
				number := runes[index]
				num, _ := strconv.Atoi(string(number))
				result = result + num
			}

			if result < 10 {
				if result == 3 || result == 6 || result == 9 {
					return strconv.Itoa(result), true
				} else {
					return strconv.Itoa(result), false
				}
			}
			runes = []rune(strconv.Itoa(result))
		}

	}

	return "", false
}

func (caw *CodeAlfredWorkflow) Md5() string {
	hash := md5.New()
	hash.Write([]byte(caw.Args))
	return hex.EncodeToString(hash.Sum(nil))
}

func (caw *CodeAlfredWorkflow) SHA256() string {
	hash := sha256.New()
	hash.Write([]byte(caw.Args))
	return hex.EncodeToString(hash.Sum(nil))
}

func (cae *CodeAlfredWorkflow) EncodeBase64() string {
	return base64.StdEncoding.EncodeToString([]byte(cae.Args))
}

func (cae *CodeAlfredWorkflow) DecodeBase64() string {
	if bt, err := base64.StdEncoding.DecodeString(cae.Args); err != nil {
		return ""
	} else {
		return string(bt)
	}
}

// 编码为标准 URL 字符串
func (cae *CodeAlfredWorkflow) EncodeStandardURL() string {
	return url.QueryEscape(cae.Args)
}

// 编码所有为%XX格式
func (cae *CodeAlfredWorkflow) EncodeAllURL() string {
	input := cae.Args
	result := ""
	for _, word := range input {
		// %02X 表示两位十六进制
		result += fmt.Sprintf("%%%02X", word)
	}
	return result
}

// 编码所有为\XX格式
func (cae *CodeAlfredWorkflow) ToHEX() string {
	input := cae.Args
	var escapedHex strings.Builder
	for _, r := range input {
		escapedHex.WriteString(fmt.Sprintf(`\X%02X`, r))
	}
	return escapedHex.String()
}

func (cae *CodeAlfredWorkflow) FromHEX() string {
	input := cae.Args
	// Split the input string on '\X' to get individual hex codes
	re := regexp.MustCompile(`(?i)\\X`)
	parts := re.Split(input, -1)
	var decodedBytess []byte

	for _, part := range parts {
		if part == "" {
			continue
		}

		codePoint, err := strconv.ParseUint(part, 16, 8)
		if err != nil {
			return ""
		}

		decodedBytess = append(decodedBytess, byte(codePoint))
	}

	// Convert runes to a string
	decodedString := string(decodedBytess)
	return decodedString
}

func (cae *CodeAlfredWorkflow) DecodeURL() string {
	result, _ := url.QueryUnescape(cae.Args)
	return result
}

func (cae *CodeAlfredWorkflow) EncodeHTMLEntities() string {
	return html.EscapeString(cae.Args)
}
func (cae *CodeAlfredWorkflow) DecodeHTMLEntities() string {
	return html.UnescapeString(cae.Args)
}

// unicode转义
// 超过标准字符组平面（BMP）的 unicode 有两种表示方法：  比如emoji 😀
// 1. \uD83D\uDE00  UTF-16编码 代理对的变长编码表示。高代理项（high surrogate）：\uD83D，低代理项（low surrogate）：\uDE00
// 2.1 \u0001F600 UTF-32编码，一般是编程语言中表示Unicode字符的表示
// 2.2 U+1F600 Unicode标准对字符码位的通用表示法，常用于文档、规范和描述 Unicode 字符

func (cae *CodeAlfredWorkflow) UnicodeEscape(utfbase int) string {
	result := ""
	input := cae.Args
	for _, r := range input {
		// 判断 r 是否在基本多文种平面（BMP）内。即标准 Unicode 字符集范围
		if r <= 0xFFFF {
			// 使用 \uXXXX 形式表示
			result += fmt.Sprintf("\\u%04X", r)
		} else {
			switch utfbase {
			case 16:
				// 一般是这种
				// UTF-16编码 代理对形式
				r1, r2 := utf16.EncodeRune(r)
				result += fmt.Sprintf("\\u%04X\\u%04X", r1, r2)
			case 32:
				//UTF-32编码 使用 \UXXXXXXXX 形式表示。即扩展字符集范围
				result += fmt.Sprintf("\\U%08X", r)
			}

		}
	}
	return result
}

func (caw *CodeAlfredWorkflow) UnicodeEscapeUTF16() string {
	return caw.UnicodeEscape(16)
}
func (caw *CodeAlfredWorkflow) UnicodeEscapeUTF32() string {
	return caw.UnicodeEscape(32)
}

func (caw *CodeAlfredWorkflow) UnicodeUnEscape() string {
	input := caw.Args
	var result string

	//处理utf32编码的 code  支持 \u0001F600  U+1F600
	// 注意 U+ 格式的两种处理形式
	// 参考: U+hex https://r12a.github.io/app-conversion/
	re := regexp.MustCompile(`(?i)\\U([0-9A-Fa-f]{8})|U\+10([A-Fa-f0-9]{4})|U\+([0-9A-Fa-f]{1,5})`)
	result = re.ReplaceAllStringFunc(input, func(match string) string {
		// 提取 十六进制部分
		code, err := strconv.ParseInt(match[2:], 16, 32)
		if err != nil {
			return match
		}
		// 转换为实际的 Unicode 字符
		return string(rune(code))
	})

	// 处理utf16
	// 优先 json方式解码
	// 其次手动解码
	var str string
	if err := json.Unmarshal([]byte(`"`+result+`"`), &str); err != nil {
		if manualDecode, err := decodeUnicodeForUTF16(result); err != nil {
			return result
		} else {
			return manualDecode
		}
		// return result
	}
	// 处理正常
	if str != "" {
		result = str
	}

	return result
}

// 手动解码 unicode UTF16编码,包含代理对
func decodeUnicodeForUTF16(s string) (string, error) {
	var result strings.Builder

	// Handle surrogate pairs
	for i := 0; i < len(s); {
		if s[i] == '\\' && i+5 < len(s) && s[i+1] == 'u' {
			hex := s[i+2 : i+6]
			codePoint, err := strconv.ParseUint(hex, 16, 16)
			if err != nil {
				return "", err
			}

			if codePoint >= 0xD800 && codePoint <= 0xDBFF && i+11 < len(s) && s[i+6] == '\\' && s[i+7] == 'u' {
				// Handle high surrogate
				lowHex := s[i+8 : i+12]
				lowCodePoint, err := strconv.ParseUint(lowHex, 16, 16)
				if err != nil {
					return "", err
				}

				if lowCodePoint >= 0xDC00 && lowCodePoint <= 0xDFFF {
					// Valid low surrogate
					fullCodePoint := 0x10000 + ((codePoint - 0xD800) << 10) + (lowCodePoint - 0xDC00)
					result.WriteRune(rune(fullCodePoint))
					i += 12
					continue
				}
			}

			// Single Unicode code point
			result.WriteRune(rune(codePoint))
			i += 6
		} else {
			// Regular character
			result.WriteByte(s[i])
			i++
		}
	}

	return result.String(), nil
}
