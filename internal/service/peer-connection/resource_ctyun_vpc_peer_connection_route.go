package peer_connection

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type CtyunVpcPeerConnectionRoute struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunVpcPeerConnectionRoute) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_peer_connection_route"
}

func (c *CtyunVpcPeerConnectionRoute) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunVpcPeerConnectionRoute() resource.Resource {
	return &CtyunVpcPeerConnectionRoute{}
}

func (c *CtyunVpcPeerConnectionRoute) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunVpcPeerConnectionRouteConfig
	var ID, regionId, projectId, vpcId, subnetId string
	err = terraform_extend.Split(request.ID, &ID, &regionId, &projectId, &vpcId, &subnetId)
	if err != nil {
		return
	}
	config.ID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionId)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunVpcPeerConnectionRoute) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id,如果不填这默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"ip_version": schema.StringAttribute{
				Required:    true,
				Description: "ip版本，ipv4或ipv6",
				Validators: []validator.String{
					stringvalidator.OneOf("ipv4", "ipv6"),
				},
			},
			"next_hop_id": schema.StringAttribute{
				Required:    true,
				Description: "下一跳设备id",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "路由表所在 vpc id",
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"destination": schema.StringAttribute{
				Required:    true,
				Description: "目的 cidr",
				Validators: []validator.String{
					validator2.Cidr(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "对等连接路由规则id",
			},
		},
	}
}

func (c *CtyunVpcPeerConnectionRoute) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunVpcPeerConnectionRouteConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunVpcPeerConnectionRoute) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpcPeerConnectionRouteConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunVpcPeerConnectionRoute) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunVpcPeerConnectionRoute) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunVpcPeerConnectionRouteConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunVpcPeerConnectionRoute) create(ctx context.Context, config *CtyunVpcPeerConnectionRouteConfig) error {
	params := &ctvpc.CtvpcCreateVpcPeerRouteRequest{
		IpVersion:   business.IPVersionDict[config.IPVersion.ValueString()],
		NextHopID:   config.NextHopID.ValueString(),
		VpcID:       config.VpcID.ValueString(),
		Destination: config.Destination.ValueString(),
		RegionID:    config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateVpcPeerRouteApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建vpc对待连接路由失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	config.ID = types.StringValue(*resp.ReturnObj.RouteRule)
	return nil
}

func (c *CtyunVpcPeerConnectionRoute) getAndMerge(ctx context.Context, config *CtyunVpcPeerConnectionRouteConfig) error {
	detail, err := c.getPeerConnectRoute(ctx, config)
	if err != nil {
		return err
	}
	config.NextHopID = types.StringValue(*detail.NextHopID)
	config.VpcID = types.StringValue(*detail.VpcID)
	config.Destination = types.StringValue(*detail.Destination)
	return nil
}

func (c *CtyunVpcPeerConnectionRoute) getPeerConnectRoute(ctx context.Context, config *CtyunVpcPeerConnectionRouteConfig) (*ctvpc.CtvpcGetVpcPeerRouteReturnObjResultsResponse, error) {
	params := &ctvpc.CtvpcGetVpcPeerRouteRequest{
		VpcID:      config.VpcID.ValueString(),
		RegionID:   config.RegionID.ValueString(),
		PageSize:   10,
		PageNumber: 1,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcGetVpcPeerRouteApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取vpc对等连接路由详情失败(id=%s)，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	for _, rule := range resp.ReturnObj.Results {
		if *rule.RouteRuleID == config.ID.ValueString() {
			return rule, nil
		}
	}
	//if len(resp.ReturnObj.Results) < 1 {
	//	err = fmt.Errorf("获取vpc对=等连接路由详情失败(id=%s)，未查询到该待连接路由信息，请检查资源是否存在！", config.ID.ValueString())
	//	return nil, err
	//} else if len(resp.ReturnObj.Results) > 1 {
	//	err = fmt.Errorf("通过id=%s查询vpc对等连接路由详情失败， 接口返回多个对等路由信息！", config.ID.ValueString())
	//}
	return resp.ReturnObj.Results[0], err
}

func (c *CtyunVpcPeerConnectionRoute) delete(ctx context.Context, config CtyunVpcPeerConnectionRouteConfig) error {
	params := &ctvpc.CtvpcDeleteVpcPeerRouteRequest{
		RouteRuleID: config.ID.ValueString(),
		RegionID:    config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteVpcPeerRouteApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除vpc对等连接路由失败(id=%s)，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return err
	}
	return nil
}

type CtyunVpcPeerConnectionRouteConfig struct {
	RegionID    types.String `tfsdk:"region_id"`
	IPVersion   types.String `tfsdk:"ip_version"`
	NextHopID   types.String `tfsdk:"next_hop_id"`
	VpcID       types.String `tfsdk:"vpc_id"`
	Destination types.String `tfsdk:"destination"`
	ID          types.String `tfsdk:"id"`
}
