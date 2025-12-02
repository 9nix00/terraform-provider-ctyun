package peer_connection

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CtyunVpcPeerConnectionRoutes struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpcPeerConnectionRoutes() datasource.DataSource {
	return &CtyunVpcPeerConnectionRoutes{}
}

func (c *CtyunVpcPeerConnectionRoutes) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunVpcPeerConnectionRoutes) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_peer_connection_routes"
}

func (c *CtyunVpcPeerConnectionRoutes) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"route_table_id": schema.StringAttribute{
				Description: "路由表 id",
				Optional:    true,
			},
			"vpc_id": schema.StringAttribute{
				Description: "VPC id",
				Required:    true,
			},
			"region_id": schema.StringAttribute{
				Description: "资源池 id",
				Optional:    true,
			},
			"page_no": schema.Int32Attribute{
				Description: "页码，默认为1",
				Optional:    true,
			},
			"page_size": schema.Int32Attribute{
				Description: "当前页数据条数，默认为10，最大值为50",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"routes": schema.ListNestedAttribute{
				Description: "路由规则列表",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "路由规则id",
							Computed:    true,
						},
						"destination": schema.StringAttribute{
							Description: "目的 cidr",
							Computed:    true,
						},
						"route_table_id": schema.StringAttribute{
							Description: "路由表 id",
							Computed:    true,
						},
						"vpc_id": schema.StringAttribute{
							Description: "vpc id",
							Computed:    true,
						},
						"next_hop_id": schema.StringAttribute{
							Description: "下一跳设备 id",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (c *CtyunVpcPeerConnectionRoutes) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunVpcPeerConnectionRoutesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("region id 为空")
		return
	}
	config.RegionID = types.StringValue(regionId)

	params := &ctvpc.CtvpcGetVpcPeerRouteRequest{
		VpcID:      config.VpcID.ValueString(),
		RegionID:   config.RegionID.ValueString(),
		PageSize:   10,
		PageNumber: 1,
	}
	if !config.RouteTableID.IsNull() {
		params.RouteTableID = config.RouteTableID.ValueStringPointer()
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}
	if !config.PageNo.IsNull() {
		params.PageNumber = config.PageNo.ValueInt32()
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcGetVpcPeerRouteApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("获取vpc对等连接路由详情失败，接口返回nil，请联系研发确认问题原因！")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 处理路由规则列表
	var routes []CtyunVpcPeerConnectionRouteInfoModel
	for _, routeItem := range resp.ReturnObj.Results {
		var route CtyunVpcPeerConnectionRouteInfoModel
		route.ID = types.StringValue(*routeItem.RouteRuleID)
		route.RouteTableID = types.StringValue(*routeItem.RouteTableID)
		route.Destination = types.StringValue(*routeItem.Destination)
		route.VpcID = types.StringValue(*routeItem.VpcID)
		route.NextHopID = types.StringValue(*routeItem.NextHopID)
		routes = append(routes, route)
	}
	config.Routes = routes
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunVpcPeerConnectionRouteInfoModel struct {
	ID           types.String `tfsdk:"id"`
	Destination  types.String `tfsdk:"destination"`
	RouteTableID types.String `tfsdk:"route_table_id"`
	VpcID        types.String `tfsdk:"vpc_id"`
	NextHopID    types.String `tfsdk:"next_hop_id"`
}

type CtyunVpcPeerConnectionRoutesConfig struct {
	RouteTableID types.String                           `tfsdk:"route_table_id"`
	VpcID        types.String                           `tfsdk:"vpc_id"`
	RegionID     types.String                           `tfsdk:"region_id"`
	PageNo       types.Int32                            `tfsdk:"page_no"`
	PageSize     types.Int32                            `tfsdk:"page_size"`
	Routes       []CtyunVpcPeerConnectionRouteInfoModel `tfsdk:"routes"`
}
