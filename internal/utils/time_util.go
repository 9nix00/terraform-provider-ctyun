package utils

import (
	"time"
)

// 从RFC3339转换到本地时间格式
func FromRFC3339ToLocal(timeStr string) string {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t.In(time.Local).Format(time.DateTime)
}

// FromUnixToUTC 将Unix时间戳（10位秒级/13位毫秒级）转换为UTC时间字符串
// 输入的时间值：1717574263（10位）、1717574263000（13位）
// 返回值: 转换后的UTC时间字符串，如"2025-09-09T03:42:52Z"，无效输入返回空
func FromUnixToUTC(timestamp int64) string {
	// 调整范围校验：包含13位合理最大值（约2286年）
	if timestamp < 0 || timestamp > 3000000000000 {
		return ""
	}

	var t time.Time
	if timestamp > 9999999999 {
		sec := timestamp / 1000
		nsec := (timestamp % 1000) * 1e6
		t = time.Unix(sec, nsec).UTC()
	} else {
		t = time.Unix(timestamp, 0).UTC()
	}

	return t.Format(time.RFC3339)
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
