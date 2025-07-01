package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
	"strings"
)

type validatorNetworkName struct {
}

func NetworkName() validator.String {
	return &validatorNetworkName{}
}

func (v validatorNetworkName) Description(_ context.Context) string {
	return "name不满足要求"
}

func (v validatorNetworkName) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// 支持拉丁字母、中文、数字，下划线，连字符，中文/英文字母开头，不能以http:/https:开头，长度2-32
func (v validatorNetworkName) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	name := request.ConfigValue.ValueString()
	length := len(name)
	if length < 2 || length > 32 {
		errMessage := "name长度必须在2-32"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
	if strings.HasPrefix(name, "http:") || strings.HasPrefix(name, "https:") {
		errMessage := "name不能以http:/https:开头"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}

	pattern := `^[a-zA-Z\x{4e00}-\x{9fa5}][a-zA-Z\x{4e00}-\x{9fa5}0-9_-]*$`
	rex, err := regexp.Compile(pattern)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), "")
		return
	}
	f := rex.MatchString(name)
	if !f {
		response.Diagnostics.AddError("支持拉丁字母、中文、数字，下划线，连字符，中文/英文字母开头", "")
		return
	}

}
