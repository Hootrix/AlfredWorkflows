package utils

import (
	"strings"
	"unicode"
)

// TrimSpace 去除字符串两端的空白字符
func TrimSpace(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return len(TrimSpace(s)) == 0
}
