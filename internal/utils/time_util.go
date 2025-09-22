package utils

import (
	"time"
)

// 从RFC3339转换到本地时间格式
func FromRFC3339ToLocal(timeStr string) string {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t.In(time.Local).Format(time.DateTime)
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

// BeijingToUTCZ 东八区转UTC
func BeijingToUTCZ(input string) string {
	// 解析输入时间
	t, err := time.Parse(time.DateOnly+" "+time.TimeOnly, input)
	if err != nil {
		return ""
	}
	// 转换为UTC时区并格式化
	return t.UTC().Format(time.RFC3339)
}
