package ecs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ctyunEcsPortAssociation struct {
	meta *common.CtyunMetadata
}

func NewCtyunEcsPortAssociation() resource.Resource {
	return &ctyunEcsPortAssociation{}
}

func (c *ctyunEcsPortAssociation) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_port_association"
}

type CtyunEcsPortAssociationConfig struct {
	ID                 types.String `tfsdk:"id"`
	RegionID           types.String `tfsdk:"region_id"`
	InstanceID         types.String `tfsdk:"instance_id"`
	NetworkInterfaceID types.String `tfsdk:"network_interface_id"`
	AzName             types.String `tfsdk:"az_name"`
	ProjectID          types.String `tfsdk:"project_id"`
}

func (c *ctyunEcsPortAssociation) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**云主机绑定弹性网卡，详细说明请见文档：https://www.ctyun.cn/document/10026730/10597686**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识符",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_interface_id": schema.StringAttribute{
				Required:    true,
				Description: "弹性网卡ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称，如果不填则默认使用provider ctyun中的az_name或环境变量中的CTYUN_AZ_NAME",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraAzName, true),
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
			},
		},
	}
}

func (c *ctyunEcsPortAssociation) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunEcsPortAssociation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CtyunEcsPortAssociationConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := c.createEcsPortAssociation(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"绑定弹性网卡失败",
			fmt.Sprintf("绑定弹性网卡到云主机时发生错误: %s", err.Error()),
		)
		return
	}

	// 设置ID
	plan.ID = types.StringValue(generateEcsPortAssociationId(plan.RegionID.ValueString(), plan.InstanceID.ValueString(), plan.NetworkInterfaceID.ValueString()))
	_, err = c.getAndMergeEcsPortAssociation(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"获取云主机弹性网卡绑定详情失败",
			fmt.Sprintf("获取云主机弹性网卡绑定详情时发生错误: %s", err.Error()),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *ctyunEcsPortAssociation) createEcsPortAssociation(ctx context.Context, plan CtyunEcsPortAssociationConfig) (*ctecs.CtecsPortsAttachInstanceV41Response, error) {
	regionId := c.meta.GetExtraIfEmpty(plan.RegionID.ValueString(), common.ExtraRegionId)
	projectId := c.meta.GetExtraIfEmpty(plan.ProjectID.ValueString(), common.ExtraProjectId)
	azName := c.meta.GetExtraIfEmpty(plan.AzName.ValueString(), common.ExtraAzName)

	// 绑定弹性网卡到云主机
	attachRequest := &ctecs.CtecsPortsAttachInstanceV41Request{
		ClientToken:        uuid.NewString(),
		RegionID:           regionId,
		ProjectID:          projectId,
		AzName:             azName,
		NetworkInterfaceID: plan.NetworkInterfaceID.ValueString(),
		InstanceID:         plan.InstanceID.ValueString(),
		InstanceType:       3, // 3-虚拟机
	}

	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsPortsAttachInstanceV41Api.Do(ctx, c.meta.SdkCredential, attachRequest)
	return resp, err
}

func (c *ctyunEcsPortAssociation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CtyunEcsPortAssociationConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeResponse, err := c.getAndMergeEcsPortAssociation(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"获取云主机详情失败",
			fmt.Sprintf("获取云主机详情时发生错误: %s", err.Error()),
		)
		return
	}

	// 检查弹性网卡是否仍然绑定到该云主机
	portAttached := false
	if describeResponse.ReturnObj != nil {
		for _, networkCard := range describeResponse.ReturnObj.NetworkCardList {
			if networkCard != nil && networkCard.NetworkCardID == state.NetworkInterfaceID.ValueString() {
				portAttached = true
				break
			}
		}
	}

	if !portAttached {
		// 如果弹性网卡未绑定，则从状态中移除该资源
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *ctyunEcsPortAssociation) getAndMergeEcsPortAssociation(ctx context.Context, state CtyunEcsPortAssociationConfig) (*ctecs.CtecsDetailsInstanceV41Response, error) {
	regionId := c.meta.GetExtraIfEmpty(state.RegionID.ValueString(), common.ExtraRegionId)

	// 查询云主机详情以确认弹性网卡是否仍然绑定
	describeRequest := &ctecs.CtecsDetailsInstanceV41Request{
		RegionID:   regionId,
		InstanceID: state.InstanceID.ValueString(),
	}

	describeResponse, err := c.meta.Apis.SdkCtEcsApis.CtecsDetailsInstanceV41Api.Do(ctx, c.meta.SdkCredential, describeRequest)
	return describeResponse, err
}

func (c *ctyunEcsPortAssociation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// 由于所有属性都需要替换资源，实际上不会执行更新操作
	resp.Diagnostics.AddError(
		"不支持更新操作",
		"云主机绑定弹性网卡资源不支持更新操作，如需修改请先删除再重新创建",
	)
}

func (c *ctyunEcsPortAssociation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CtyunEcsPortAssociationConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(state.RegionID.ValueString(), common.ExtraRegionId)

	// 解绑弹性网卡
	detachRequest := &ctecs.CtecsPortsDetachInstanceV41Request{
		ClientToken:        uuid.NewString(),
		RegionID:           regionId,
		NetworkInterfaceID: state.NetworkInterfaceID.ValueString(),
		InstanceID:         state.InstanceID.ValueString(),
	}

	_, err := c.meta.Apis.SdkCtEcsApis.CtecsPortsDetachInstanceV41Api.Do(ctx, c.meta.SdkCredential, detachRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"解绑弹性网卡失败",
			fmt.Sprintf("解绑弹性网卡时发生错误: %s", err.Error()),
		)
		return
	}
}

func (c *ctyunEcsPortAssociation) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var regionId, instanceId, networkInterfaceId string
	err := terraform_extend.Split(req.ID, &regionId, &instanceId, &networkInterfaceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"导入参数不正确",
			fmt.Sprintf("导入参数应为 region_id,instance_id,network_interface_id 格式: %s", err.Error()),
		)
		return
	}

	// 设置导入的属性
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), generateEcsPortAssociationId(regionId, instanceId, networkInterfaceId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region_id"), regionId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), instanceId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_interface_id"), networkInterfaceId)...)
}

func generateEcsPortAssociationId(regionId, instanceId, networkInterfaceId string) string {
	return fmt.Sprintf("%s/%s/%s", regionId, instanceId, networkInterfaceId)
}
