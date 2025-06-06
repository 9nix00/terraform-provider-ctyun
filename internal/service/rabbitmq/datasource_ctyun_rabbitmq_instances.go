package rabbitmq

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/amqp"
)

var (
	_ datasource.DataSource              = &ctyunRabbitmqInstances{}
	_ datasource.DataSourceWithConfigure = &ctyunRabbitmqInstances{}
)

type ctyunRabbitmqInstances struct {
	meta *common.CtyunMetadata
}

func NewCtyunRabbitmqInstances() datasource.DataSource {
	return &ctyunRabbitmqInstances{}
}

func (c *ctyunRabbitmqInstances) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rabbitmq_instances"
}

type CtyunRabbitmqInstancesModel struct {
	//Subnet        string      `json:"subnet"`        // 子网名称？
	//          string      `json:"prod"`          // 规格
	//EngineType    string      `json:"engineType"`    // 引擎类型
	//BillMode      string      `json:"billMode"`      // 账单
	//SecurityGroup string      `json:"securityGroup"` // 安全组名称
	//Type      interface{} `json:"prodType"`
	//Network       string      `json:"network"`     // vpc名称?
	//ExpireTime    string      `json:"expireTime"`  // 过期时间
	//CreateTime    string      `json:"createTime"`  // 创建时间
	//ClusterName   string      `json:"clusterName"` // 实例名称
	//InstId    string      `json:"prodInstId"`  // 实例id
	//Status        int32       `json:"status"`      // 状态
}

type CtyunRabbitmqInstancesConfig struct {
	RegionID  types.String                  `tfsdk:"region_id"`
	PageNo    types.Int32                   `tfsdk:"page_no"`
	PageSize  types.Int32                   `tfsdk:"page_size"`
	Instances []CtyunRabbitmqInstancesModel `tfsdk:"instances"`
}

func (c *ctyunRabbitmqInstances) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页数据量大小",
			},
			"instances": schema.ListNestedAttribute{
				Description:  "List of RabbitMQ specifications.",
				Computed:     true,
				NestedObject: schema.NestedAttributeObject{},
			},
		},
	}
}

func (c *ctyunRabbitmqInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRabbitmqInstancesConfig
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
	params := &amqp.AmqpInstanceQueryRequest{
		RegionID: regionId}
	// 调用API
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstanceQueryApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	config.Instances = []CtyunRabbitmqInstancesModel{}
	// 解析返回值

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRabbitmqInstances) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
