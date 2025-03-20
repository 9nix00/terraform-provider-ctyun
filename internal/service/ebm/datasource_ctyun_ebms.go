package ebm

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"terraform-provider-ctyun/internal/common"
)

var (
	_ datasource.DataSource              = &ctyunEbmDeviceTypes{}
	_ datasource.DataSourceWithConfigure = &ctyunEbmDeviceTypes{}
)

type ctyunEbms struct {
	meta *common.CtyunMetadata
}

func CtyunEbms() datasource.DataSource {
	return &ctyunEbms{}
}

func (c *ctyunEbms) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ebms"
}

type CtyunEbmsModel struct {
}

type CtyunEbmsConfig struct {
}

func (c *ctyunEbms) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：**`,
		Attributes:          map[string]schema.Attribute{},
	}
}

func (c *ctyunEbms) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var config CtyunEbmsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 组装请求体

	// 调用API

	// 解析返回值

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEbms) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
