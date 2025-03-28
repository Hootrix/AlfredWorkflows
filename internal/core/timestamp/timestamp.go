package timestamp

import (
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

// GetCurrentTimestamp 获取当前时间戳和格式化时间
func GetCurrentTimestamp() (int64, string) {
	now := time.Now()
	ts := now.Unix()
	timeStr := now.Format(time.DateTime)
	return ts, timeStr
}

// TimestampToTime 将时间戳转换为格式化时间
func TimestampToTime(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format(time.DateTime)
}

// ParseTimeString 解析时间字符串
func ParseTimeString(timeStr string) (time.Time, error) {
	return dateparse.ParseLocal(timeStr)
}

// FormatUnixTimestamp 格式化Unix时间戳为字符串
func FormatUnixTimestamp(ts int64) string {
	return strconv.FormatInt(ts, 10)
}
