package utils

import (
	"time"
)

// 从RFC3339转换到本地时间格式
func FromRFC3339ToLocal(timeStr string) string {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t.In(time.Local).Format(time.DateTime)
}

// FromUnixToUTC 将Unix时间戳转换为UTC时间字符串
// 输入的时间值：1717574263,
// 返回值: 转换后的UTC时间字符串，如"2025-09-09T03:42:52Z"
func FromUnixToUTC(timestamp int64) string {
	// 检查时间戳是否合理（避免过大的值）
	if timestamp < 0 || timestamp > 9999999999 {
		return ""
	}

	// 将Unix时间戳转换为时间对象
	t := time.Unix(timestamp, 0).UTC()

	// 格式化为RFC3339格式的UTC时间字符串
	utcStr := t.Format(time.RFC3339)

	return utcStr
}

// ConvertToUTCZ 将带时区的时间字符串转换为UTC时间，并格式化为2006-01-02T15:04:05Z格式
// input: 输入时间字符串，如"2025-09-09T11:42:52+08:00"
// 返回值: 转换后的UTC时间字符串，如"2025-09-09T03:42:52Z"，及可能的错误
func ConvertToUTCZ(input string) string {
	// 解析输入时间
	t, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return ""
	}
	// 转换为UTC时区并格式化
	return t.UTC().Format(time.RFC3339)
}

// FromLocalToUTCZ 将本地时间格式(yyyy-MM-dd HH:mm:ss)转换为UTC时间格式(yyyy-MM-ddTHH:mm:ssZ)
// input: 输入时间字符串，如"2022-08-20 11:53:51"
// 返回值: 转换后的UTC时间字符串，如"2022-08-20T11:53:51Z"
func FromLocalToUTCZ(input string) string {
	// 解析输入时间
	t, err := time.Parse(time.DateOnly+" "+time.TimeOnly, input)
	if err != nil {
		return ""
	}

	// 转换为UTC时区并格式化为RFC3339格式
	return t.UTC().Format(time.RFC3339)
}
