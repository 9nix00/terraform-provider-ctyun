package redis

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/common"
)

var (
	_ datasource.DataSource              = &ctyunRedisSpecs{}
	_ datasource.DataSourceWithConfigure = &ctyunRedisSpecs{}
)

type ctyunRedisSpecs struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisSpecs() datasource.DataSource {
	return &ctyunRedisSpecs{}
}

func (c *ctyunRedisSpecs) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_specs"
}

type CtyunRedisSpecsModel struct {
}

type CtyunRedisSpecsConfig struct {
	RegionID types.String `tfsdk:"region_id"`
}

func (c *ctyunRedisSpecs) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：**`,
		Attributes:          map[string]schema.Attribute{},
	}
}

func (c *ctyunRedisSpecs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRedisSpecsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}

	// 组装请求体

	// 调用API

	// 解析返回值

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRedisSpecs) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
