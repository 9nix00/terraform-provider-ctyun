package lb

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
)

var (
	_ datasource.DataSource              = &ctyunElbLoadBalancers{}
	_ datasource.DataSourceWithConfigure = &ctyunElbLoadBalancers{}
)

type ctyunElbLoadBalancers struct {
	meta *common.CtyunMetadata
}

func (c *ctyunElbLoadBalancers) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_loadbalancers"
}

func (c *ctyunElbLoadBalancers) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "**详情文档可查看：https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=24&api=5647&data=88&isNormal=1&vid=82",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "区域ID",
			},
			"az_name": schema.StringAttribute{
				Computed:    true,
				Description: "可用区名称",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "负载均衡ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "描述",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "名称",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "描述",
			},
			"vpc_id": schema.StringAttribute{
				Computed:    true,
				Description: "VPC ID",
			},
			"subnet_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "子网ID",
			},
			"port_id": schema.StringAttribute{
				Computed:    true,
				Description: "负载均衡实例默认创建port ID",
			},
			"private_ip_address": schema.StringAttribute{
				Computed:    true,
				Description: "负载均衡实例的内网VIP",
			},
			"ipv6_address": schema.StringAttribute{
				Computed:    true,
				Description: "负载均衡实例的IPv6地址",
			},
			"sla_name": schema.StringAttribute{
				Computed:    true,
				Description: "规格名称",
			},
			"delete_protection": schema.BoolAttribute{
				Computed:    true,
				Description: "删除保护。开启，不开启",
			},
			"admin_status": schema.StringAttribute{
				Computed:    true,
				Description: "管理状态: DOWN / ACTIVE",
				Validators: []validator.String{
					stringvalidator.OneOf(business.AdminStatusName...),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "负载均衡状态: DOWN / ACTIVE",
				Validators: []validator.String{
					stringvalidator.OneOf(business.AdminStatusName...),
				},
			},
			"resource_type": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "负载均衡类型: external / internal",
				Validators: []validator.String{
					stringvalidator.OneOf(business.LbResourceType...),
				},
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "created_time",
			},
			"updated_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间，为UTC格式",
			},
			"ids": schema.ListNestedAttribute{
				Optional:    true,
				Description: "负载均衡ID列表，以,分隔",
			},
			"eip_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"resource_id": schema.StringAttribute{
						Computed:    true,
						Description: "计费类资源ID",
					},
					"eip_id": schema.StringAttribute{
						Computed:    true,
						Description: "弹性公网IP的ID",
					},
					"bandwidth": schema.Int64Attribute{
						Computed:    true,
						Description: "弹性公网IP的带宽",
					},
					"is_talk_order": schema.BoolAttribute{
						Computed:    true,
						Description: "是否按需资源",
					},
				},
			},
		},
	}
}

func (c *ctyunElbLoadBalancers) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunElbLoadBalancersConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	fmt.Println(regionId)
	//params = &
	if !config.IDs.IsNull() {

	}
}

func (c *ctyunElbLoadBalancers) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunElbLoadBalancersConfig struct {
	RegionID         types.String `tfsdk:"region_id"`          //区域ID
	AzName           types.String `tfsdk:"az_name"`            //可用区名称
	ID               types.String `tfsdk:"id"`                 //负载均衡ID
	ProjectID        types.String `tfsdk:"project_id"`         //项目ID
	Name             types.String `tfsdk:"name"`               //名称
	Description      types.String `tfsdk:"description"`        //描述
	VpcID            types.String `tfsdk:"vpc_id"`             //VPC ID
	SubnetID         types.String `tfsdk:"subnet_id"`          //子网ID
	PortID           types.String `tfsdk:"port_id"`            //负载均衡实例默认创建port ID
	privateIpAddress types.String `tfsdk:"private_ip_address"` //负载均衡实例的内网VIP
	Ipv6Address      types.String `tfsdk:"ipv6_address"`       //负载均衡实例的IPv6地址
	EipInfo          EipInfoModel `tfsdk:"eip_info"`           //弹性公网IP信息
	SlaName          types.String `tfsdk:"sla_name"`           //规格名称
	DeleteProtection types.Bool   `tfsdk:"delete_protection"`  //删除保护。开启，不开启
	AdminStatus      types.String `tfsdk:"admin_status"`       //管理状态: DOWN / ACTIVE
	Status           types.String `tfsdk:"status"`             //负载均衡状态: DOWN / ACTIVE
	ResourceType     types.String `tfsdk:"resource_type"`      //负载均衡类型: external / internal
	CreatedTime      types.String `tfsdk:"created_time"`       //创建时间，为UTC格式
	UpdatedTime      types.String `tfsdk:"updated_time"`       //更新时间，为UTC格式
	// 查询的参数
	IDs types.String `tfsdk:"ids"` //负载均衡ID列表，以,分隔

}

type EipInfoModel struct {
	ResourceID  types.String `tfsdk:"resource_id"`
	EipID       types.String `tfsdk:"eip_id"`
	Bandwidth   types.Int64  `tfsdk:"bandwidth"`
	IsTalkOrder types.Bool   `tfsdk:"is_talk_order"`
}
