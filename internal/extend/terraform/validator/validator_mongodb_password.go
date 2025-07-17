package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
	"unicode"
)

type validatorMongodbPassword struct{}

func (v validatorMongodbPassword) Description(ctx context.Context) string {
	return "输入的密码不满足要求，实例密码要求：8-32位由大写字母、小写字母、数字、特殊字符中的任意三种组成 特殊字符为!@#$%^&*()_+-="
}

func (v validatorMongodbPassword) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorMongodbPassword) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	password := request.ConfigValue.ValueString()
	// 检查长度
	if len(password) < 8 || len(password) > 32 {
		errMessage := "mongodb实例密码长度需要保持8~32位"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
	// 定义允许的特殊字符
	specialChars := "!@#$%^&*()_+-="
	// 初始化字符类型标记
	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
		validChars = true
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		default:
			validChars = false // 发现非法字符
		}
	}

	// 统计满足的字符类型数量
	typeCount := 0
	if hasUpper {
		typeCount++
	}
	if hasLower {
		typeCount++
	}
	if hasDigit {
		typeCount++
	}
	if hasSpecial {
		typeCount++
	}

	// 验证结果
	if !validChars {
		errMessage := "存在非法字符，密码由大写字母、小写字母、数字、特殊字符中的任意三种组成 特殊字符为!@#$%^&*()_+-="
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
	if typeCount < 3 {
		errMessage := "密码组合必须包括大写字母、小写字母、数字、特殊字符中的任意三种及以上"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}

}

func MongodbPassword() validator.String {
	return &validatorMongodbPassword{}
}
