package planmodif

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type ImmutableStringModifier struct{}

func (m ImmutableStringModifier) Description(ctx context.Context) string {
	return "字段在资源创建后不可修改"
}

func (m ImmutableStringModifier) MarkdownDescription(ctx context.Context) string {
	return "字段在资源创建后不可修改"
}

func (m ImmutableStringModifier) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	if !req.StateValue.IsNull() && !req.PlanValue.Equal(req.StateValue) {
		resp.Diagnostics.AddAttributeError(
			path.Root(req.Path.String()),
			"字段不可修改",
			fmt.Sprintf("字段 `%s` 在资源创建后禁止修改。如需变更，请删除当前资源并重新创建。", req.Path.String()),
		)
	}
}
