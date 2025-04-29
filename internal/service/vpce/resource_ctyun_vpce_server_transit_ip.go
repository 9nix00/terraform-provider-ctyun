package vpce

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
	"terraform-provider-ctyun/internal/utils"
)

var (
	_ resource.Resource                = &ctyunVpceServerTransitIP{}
	_ resource.ResourceWithConfigure   = &ctyunVpceServerTransitIP{}
	_ resource.ResourceWithImportState = &ctyunVpceServerTransitIP{}
)

type ctyunVpceServerTransitIP struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpceServerTransitIP() resource.Resource {
	return &ctyunVpceServerTransitIP{}
}

func (c *ctyunVpceServerTransitIP) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpce_server_transit_ip"
}

type CtyunVpceServerTransitIPConfig struct {
	ID               types.String `tfsdk:"id"`
	EndpointServerID types.String `tfsdk:"endpoint_server_id"`
	RegionID         types.String `tfsdk:"region_id"`
	SubnetID         types.String `tfsdk:"subnet_id"`
	TransitIP        types.String `tfsdk:"transit_ip"`
}

func (c *ctyunVpceServerTransitIP) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID，使用的ip地址",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"endpoint_server_id": schema.StringAttribute{
				Required:    true,
				Description: "终端节点服务id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"transit_ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "中转地址",
				Validators: []validator.String{
					validator2.Ip(),
				},
			},
		},
	}
}

func (c *ctyunVpceServerTransitIP) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunVpceServerTransitIPConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建
	ip, err := c.create(ctx, plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(ip)
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunVpceServerTransitIP) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceServerTransitIPConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunVpceServerTransitIP) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (c *ctyunVpceServerTransitIP) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceServerTransitIPConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunVpceServerTransitIP) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [ip],[endpointServerID],[regionID]
func (c *ctyunVpceServerTransitIP) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunVpceServerTransitIPConfig
	var ip, endpointServerID, regionID string
	err = terraform_extend.Split(request.ID, &ip, &endpointServerID, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.EndpointServerID = types.StringValue(endpointServerID)
	cfg.ID = types.StringValue(ip)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建
func (c *ctyunVpceServerTransitIP) create(ctx context.Context, plan CtyunVpceServerTransitIPConfig) (ip string, err error) {
	params := &ctvpc.CtvpcCreateEndpointServiceTransitIPRequest{
		ClientToken:       uuid.NewString(),
		RegionID:          plan.RegionID.ValueString(),
		SubnetID:          plan.SubnetID.ValueString(),
		EndpointServiceID: plan.EndpointServerID.ValueString(),
	}
	transitIP := plan.TransitIP.ValueString()
	if transitIP != "" {
		params.TransitIP = &transitIP
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateEndpointServiceTransitIPApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	ip = resp.ReturnObj.TransitIP
	return
}

// getAndMerge 从远端查询
func (c *ctyunVpceServerTransitIP) getAndMerge(ctx context.Context, plan *CtyunVpceServerTransitIPConfig) (err error) {
	params := &ctvpc.CtvpcListEndpointServiceTransitIPRequest{
		RegionID:          plan.RegionID.ValueString(),
		EndpointServiceID: plan.EndpointServerID.ValueString(),
		PageSize:          50,
		PageNo:            1,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListEndpointServiceTransitIPApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var exist bool
	for _, ip := range resp.ReturnObj.TransitIPs {
		if utils.SecStringValue(ip.TransitIP) == plan.ID {
			plan.SubnetID = utils.SecStringValue(ip.SubnetID)
			plan.TransitIP = utils.SecStringValue(ip.TransitIP)
			exist = true
		}
	}
	if !exist {
		err = common.InvalidReturnObjResultsError
		return
	}

	return
}

// delete 删除
func (c *ctyunVpceServerTransitIP) delete(ctx context.Context, plan CtyunVpceServerTransitIPConfig) (err error) {
	params := &ctvpc.CtvpcDeleteEndpointServiceTransitIPRequest{
		RegionID:          plan.RegionID.ValueString(),
		EndpointServiceID: plan.EndpointServerID.ValueString(),
		TransitIP:         plan.TransitIP.ValueString(),
		ClientToken:       uuid.NewString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteEndpointServiceTransitIPApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}
