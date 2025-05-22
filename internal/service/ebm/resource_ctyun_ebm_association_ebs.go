package ebm

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctebm"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
	"terraform-provider-ctyun/internal/utils"
	"time"
)

var (
	_ resource.Resource                = &ctyunEbmAssociationEbs{}
	_ resource.ResourceWithConfigure   = &ctyunEbmAssociationEbs{}
	_ resource.ResourceWithImportState = &ctyunEbmAssociationEbs{}
)

type ctyunEbmAssociationEbs struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbmAssociationEbs() resource.Resource {
	return &ctyunEbmAssociationEbs{}
}

func (c *ctyunEbmAssociationEbs) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebm_association_ebs"
}

type CtyunEbmAssociationEbsConfig struct {
	ID         types.String `tfsdk:"id"`
	RegionID   types.String `tfsdk:"region_id"`
	AzName     types.String `tfsdk:"az_name"`
	InstanceID types.String `tfsdk:"instance_id"`
	EbsID      types.String `tfsdk:"ebs_id"`
}

func (c *ctyunEbmAssociationEbs) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027724/10173867**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID",
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
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称",
				Default:     defaults.AcquireFromGlobalString(common.ExtraAzName, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "物理机ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ebs_id": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (c *ctyunEbmAssociationEbs) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbmAssociationEbsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeAssociation(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	err = c.association(ctx, plan)
	if err != nil {
		return
	}
	// 创建后等待绑定成功
	err = c.checkAfterAssociation(ctx, plan)
	if err != nil {
		return
	}
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunEbmAssociationEbs) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbmAssociationEbsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "未关联") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEbmAssociationEbs) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (c *ctyunEbmAssociationEbs) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbmAssociationEbsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 解绑
	err = c.checkBeforeDissociation(ctx, state)
	if err != nil {
		return
	}
	err = c.dissociation(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDissociation(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunEbmAssociationEbs) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [instanceID],[ebsID],[regionID],[azName]
func (c *ctyunEbmAssociationEbs) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEbmAssociationEbsConfig
	var instanceID, ebsID, regionID, azName string
	err = terraform_extend.Split(request.ID, &instanceID, &ebsID, &regionID, &azName)
	if err != nil {
		return
	}

	cfg.InstanceID = types.StringValue(instanceID)
	cfg.EbsID = types.StringValue(ebsID)
	cfg.AzName = types.StringValue(azName)
	cfg.RegionID = types.StringValue(regionID)

	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeAssociation 绑定前检查
func (c *ctyunEbmAssociationEbs) checkBeforeAssociation(ctx context.Context, plan CtyunEbmAssociationEbsConfig) (err error) {
	// 校验物理机
	instance, err := business.NewEbmService(c.meta).GetEbmInfo(
		ctx,
		plan.InstanceID.ValueString(),
		plan.RegionID.ValueString(),
		plan.AzName.ValueString(),
	)
	if err != nil {
		return
	}
	id := utils.SecString(instance.InstanceUUID)
	support := utils.SecBool(instance.DeviceDetail.SupportCloud)
	if !support {
		err = fmt.Errorf("物理机 %s 不支持挂载云硬盘", id)
		return
	}
	status := utils.SecLowerStringValue(instance.EbmState).ValueString()
	if status != business.EbmStatusRunning && status != business.EbmStatusStopping {
		err = fmt.Errorf("物理机 %s 状态必须是运行或开机状态，当前状态 %s", id, status)
		return
	}
	if len(instance.AttachedVolumes) > 9 {
		err = fmt.Errorf("物理机 %s 不能挂载更多云硬盘", id)
		return
	}
	for _, ebsID := range instance.AttachedVolumes {
		if plan.EbsID.ValueString() == utils.SecString(ebsID) {
			err = fmt.Errorf("物理机 %s 和云硬盘 %s 已关联", id, plan.EbsID.ValueString())
			return
		}
	}

	// 校验云硬盘
	err = business.NewEbsService(c.meta).MustExist(ctx, plan.EbsID.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		return
	}
	return
}

// checkAfterAssociation 绑定后检查
func (c *ctyunEbmAssociationEbs) checkAfterAssociation(ctx context.Context, plan CtyunEbmAssociationEbsConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			instance, err := business.NewEbmService(c.meta).GetEbmInfo(
				ctx,
				plan.InstanceID.ValueString(),
				plan.RegionID.ValueString(),
				plan.AzName.ValueString(),
			)
			if err != nil {
				return false
			}
			for _, ebsID := range instance.AttachedVolumes {
				if plan.EbsID.ValueString() == utils.SecString(ebsID) {
					executeSuccessFlag = true
					return false
				}
			}
			return true
		})
	if !executeSuccessFlag {
		return fmt.Errorf("物理机 %s 和云硬盘 %s 未关联", plan.InstanceID.ValueString(), plan.EbsID.ValueString())
	}
	return nil
}

// association 将物理机和云硬盘绑定
func (c *ctyunEbmAssociationEbs) association(ctx context.Context, plan CtyunEbmAssociationEbsConfig) (err error) {
	params := &ctebm.EbmAttachVolumeRequest{
		RegionID:     plan.RegionID.ValueString(),
		AzName:       plan.AzName.ValueString(),
		InstanceUUID: plan.InstanceID.ValueString(),
		VolumeID:     plan.EbsID.ValueString(),
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmAttachVolumeApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}

	return
}

// checkBeforeAssociation 解绑前检查
func (c *ctyunEbmAssociationEbs) checkBeforeDissociation(ctx context.Context, plan CtyunEbmAssociationEbsConfig) (err error) {
	// 校验物理机
	instance, err := business.NewEbmService(c.meta).GetEbmInfo(
		ctx,
		plan.InstanceID.ValueString(),
		plan.RegionID.ValueString(),
		plan.AzName.ValueString(),
	)
	if err != nil {
		return
	}
	id := utils.SecString(instance.InstanceUUID)
	status := utils.SecLowerStringValue(instance.EbmState).ValueString()
	if status != business.EbmStatusRunning && status != business.EbmStatusStopping {
		err = fmt.Errorf("物理机 %s 状态必须是运行或开机状态，当前状态 %s", id, status)
		return
	}
	return
}

// dissociation 将物理机和云硬盘解绑
func (c *ctyunEbmAssociationEbs) dissociation(ctx context.Context, plan CtyunEbmAssociationEbsConfig) (err error) {
	params := &ctebm.EbmDetachVolumeRequest{
		RegionID:     plan.RegionID.ValueString(),
		AzName:       plan.AzName.ValueString(),
		InstanceUUID: plan.InstanceID.ValueString(),
		VolumeID:     plan.EbsID.ValueString(),
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmDetachVolumeApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}

// association 绑定后检查
func (c *ctyunEbmAssociationEbs) checkAfterDissociation(ctx context.Context, plan CtyunEbmAssociationEbsConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			instance, err := business.NewEbmService(c.meta).GetEbmInfo(
				ctx,
				plan.InstanceID.ValueString(),
				plan.RegionID.ValueString(),
				plan.AzName.ValueString(),
			)
			if err != nil {
				return false
			}
			for _, ebsID := range instance.AttachedVolumes {
				if plan.EbsID.ValueString() == utils.SecString(ebsID) {
					return true
				}
			}
			executeSuccessFlag = true
			return false
		})
	if !executeSuccessFlag {
		return fmt.Errorf("物理机 %s 和云硬盘 %s 未解绑", plan.InstanceID.ValueString(), plan.EbsID.ValueString())
	}
	return nil
}

// getAndMerge 查询绑定关系
func (c *ctyunEbmAssociationEbs) getAndMerge(ctx context.Context, plan *CtyunEbmAssociationEbsConfig) (err error) {
	instance, err := business.NewEbmService(c.meta).GetEbmInfo(
		ctx,
		plan.InstanceID.ValueString(),
		plan.RegionID.ValueString(),
		plan.AzName.ValueString(),
	)
	if err != nil {
		return
	}
	ebsID := plan.EbsID.ValueString()
	instanceID := utils.SecString(instance.InstanceUUID)

	for _, attachID := range instance.AttachedVolumes {
		if ebsID == utils.SecString(attachID) {
			plan.ID = types.StringValue(fmt.Sprintf(
				"%s,%s,%s,%s",
				instanceID, ebsID, plan.RegionID.ValueString(), plan.AzName.ValueString()))
			return
		}
	}
	err = fmt.Errorf("物理机 %s 和云硬盘 %s 未关联", instanceID, ebsID)
	return
}
