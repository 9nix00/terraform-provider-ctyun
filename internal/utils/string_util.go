package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

// IsLetter 判断是否为英文的字符
func IsLetter(target rune) bool {
	return IsUpper(target) || IsLower(target)
}

// IsUpper 是否为英文大写字母
func IsUpper(target rune) bool {
	return target >= 'A' && target <= 'Z'
}

// IsLower 是否为英文小写字母
func IsLower(target rune) bool {
	return target >= 'a' && target <= 'z'
}

// IsDigit 是否为数字
func IsDigit(target rune) bool {
	return target >= '0' && target <= '9'
}

// SecStringValue 避免nil
func SecStringValue(str *string) types.String {
	if str == nil {
		return types.StringValue("")
	}
	return types.StringValue(*str)
}

// SecLowerStringValue 避免nil, 返回全小写
func SecLowerStringValue(str *string) types.String {
	if str == nil {
		return types.StringValue("")
	}
	return types.StringValue(strings.ToLower(*str))
}

// SecUpperStringValue 避免nil, 返回全大写
func SecUpperStringValue(str *string) types.String {
	if str == nil {
		return types.StringValue("")
	}
	return types.StringValue(strings.ToLower(*str))
}

// StrPointerArrayToStrArray 字符串指针数组转字符串数组
func StrPointerArrayToStrArray(array []*string) []string {
	ret := []string{}
	for _, str := range array {
		if str != nil {
			ret = append(ret, *str)
		} else {
			ret = append(ret, "")
		}
	}
	return ret
}
