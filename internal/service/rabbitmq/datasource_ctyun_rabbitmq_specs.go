package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/amqp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunRabbitmqSpecs{}
	_ datasource.DataSourceWithConfigure = &ctyunRabbitmqSpecs{}
)

type ctyunRabbitmqSpecs struct {
	meta *common.CtyunMetadata
}

func NewCtyunRabbitmqSpecs() datasource.DataSource {
	return &ctyunRabbitmqSpecs{}
}

func (c *ctyunRabbitmqSpecs) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rabbitmq_specs"
}

type CtyunRabbitmqSpecsModel struct {
	FlavorID      types.String  `tfsdk:"flavor_id"`
	SpecName      types.String  `tfsdk:"spec_name"`
	FlavorType    types.String  `tfsdk:"flavor_type"`
	FlavorName    types.String  `tfsdk:"flavor_name"`
	CpuNum        types.Int32   `tfsdk:"cpu_num"`
	MemSize       types.Int32   `tfsdk:"mem_size"`
	MultiQueue    types.Int32   `tfsdk:"multi_queue"`
	Pps           types.Int32   `tfsdk:"pps"`
	BandwidthBase types.Float64 `tfsdk:"bandwidth_base"`
	BandwidthMax  types.Int32   `tfsdk:"bandwidth_max"`
	Series        types.String  `tfsdk:"series"`
	AzList        []string      `tfsdk:"az_list"`
	//CpuArch       interface{}   `tfsdk:"cpuArch"`

}

type CtyunRabbitmqSpecsConfig struct {
	RegionID types.String              `tfsdk:"region_id"`
	Specs    []CtyunRabbitmqSpecsModel `tfsdk:"specs"`
}

func (c *ctyunRabbitmqSpecs) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10029625/10032819**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"specs": schema.ListNestedAttribute{
				Description: "List of RabbitMQ specifications.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"flavor_id": schema.StringAttribute{
							Description: "规格id",
							Computed:    true,
						},
						"spec_name": schema.StringAttribute{
							Description: "套餐名称",
							Computed:    true,
						},
						"flavor_type": schema.StringAttribute{
							Description: "规格类型",
							Computed:    true,
						},
						"flavor_name": schema.StringAttribute{
							Description: "规格名称",
							Computed:    true,
						},
						"cpu_num": schema.Int32Attribute{
							Description: "CPU",
							Computed:    true,
						},
						"mem_size": schema.Int32Attribute{
							Description: "内存",
							Computed:    true,
						},
						"multi_queue": schema.Int32Attribute{
							Description: "并发队列数量",
							Computed:    true,
						},
						"pps": schema.Int32Attribute{
							Description: "每秒包数",
							Computed:    true,
						},
						"bandwidth_base": schema.Float64Attribute{
							Description: "基准带宽",
							Computed:    true,
						},
						"bandwidth_max": schema.Int32Attribute{
							Description: "最大带宽",
							Computed:    true,
						},
						"series": schema.StringAttribute{
							Description: "系列",
							Computed:    true,
						},
						"az_list": schema.ListAttribute{
							Description: "可用区",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRabbitmqSpecs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRabbitmqSpecsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}

	config.RegionID = types.StringValue(regionId)
	// 组装请求体
	params := &amqp.AmqpInstancesQueryProdRequest{regionId}
	// 调用API
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstancesQueryProdApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	config.Specs = []CtyunRabbitmqSpecsModel{}
	// 解析返回值
	for _, spec := range resp.ReturnObj.Data {
		item := CtyunRabbitmqSpecsModel{
			FlavorID:      types.StringValue(spec.FlavorID),
			FlavorType:    types.StringValue(spec.FlavorType),
			FlavorName:    types.StringValue(spec.FlavorName),
			SpecName:      types.StringValue(spec.SpecName),
			CpuNum:        types.Int32Value(spec.CpuNum),
			MemSize:       types.Int32Value(spec.MemSize),
			MultiQueue:    types.Int32Value(spec.MultiQueue),
			Pps:           types.Int32Value(spec.Pps),
			BandwidthBase: types.Float64Value(spec.BandwidthBase),
			BandwidthMax:  types.Int32Value(spec.BandwidthMax),
			Series:        types.StringValue(spec.Series),
			AzList:        spec.AzList,
		}
		config.Specs = append(config.Specs, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRabbitmqSpecs) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
