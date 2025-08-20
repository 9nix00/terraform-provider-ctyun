package scaling

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/scaling"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type ctyunScalingEcs struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *ctyunScalingEcs) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_scaling_ecs"
}

func (c *ctyunScalingEcs) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)
}

func NewCtyunScalingEcs() resource.Resource {
	return &ctyunScalingEcs{}
}

func (c *ctyunScalingEcs) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "区域ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "伸缩组ID",
			},
			"instance_uuid_list": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "云主机ID列表。update阶段会和state阶段做对比，与state一致，不变；state中有，update阶段没有触发移除；state中没有，update阶段有触发新增。支持更新",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"protect_status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "云主机保护状态，设置了保护状态的云主机实例，在伸缩组进行缩容活动时将不会被移出。disable-关闭云主机保护，enable-开启云主机保护。支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ScalingPolicyStatuses...),
				},
			},
			"is_destroy": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "移除时是否销毁，仅当移除云主机时生效，true-ecs移出伸缩组时销毁， false-ecs移出伸缩组时不销毁",
			},
		},
	}
}

func (c *ctyunScalingEcs) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunScalingEcsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	//创建前检查
	isValid, err := c.checkBeforeScalingEcs(ctx, plan)
	if !isValid || err != nil {
		return
	}
	err = c.addScalingEcs(ctx, &plan)
	if err != nil {
		return
	}
	// 创建后，通过创建的请求轮询，确认创建成功
	//err = c.createLoop(ctx, &plan, createParams, 60)
	if err != nil {
		return
	}
	// 创建后反查创建后的证书信息
	err = c.getAndMergeScalingEcs(ctx, &plan)
	if err != nil {
		return
	}
	// 若创建时候，就需要开启/关闭云主机保护
	err = c.updateProtectStatus(ctx, &plan, &plan)

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunScalingEcs) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunScalingEcsConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeScalingEcs(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "NotExists") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunScalingEcs) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取 plan -tf文件中配置
	var plan CtyunScalingEcsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunScalingEcsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
	}
	// 控制云主机保护开关
	err = c.updateProtectStatus(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 增删云主机
	err = c.updateInstanceUUIDList(ctx, &plan, &plan)
	if err != nil {
		return
	}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeScalingEcs(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunScalingEcs) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 获取state
	var state CtyunScalingEcsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.removeEcs(ctx, &state)
	if err != nil {
		return
	}
}

func (c *ctyunScalingEcs) checkBeforeScalingEcs(ctx context.Context, config CtyunScalingEcsConfig) (bool, error) {
	// 校验最大实例数量
	// 获取scaling group 详情
	params := &scaling.ScalingGroupListRequest{
		RegionID: config.RegionID.ValueString(),
		GroupID:  config.GroupID.ValueInt64(),
		PageNo:   1,
		PageSize: 10,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return false, err
	} else if resp == nil {
		err = errors.New("获取弹性伸缩列表信息返回nil，请稍后重试或联系研发人员！")
		return false, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return false, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return false, err
	}
	if len(resp.ReturnObj.ScalingGroups) > 1 {
		err = fmt.Errorf("根据groupid: %d 获取的弹性伸缩详情返回多个实例。具体如下:%#v\n", config.GroupID.ValueInt64(), resp.ReturnObj.ScalingGroups)
		return false, err
	}
	// 确认maxCount
	maxCount := resp.ReturnObj.ScalingGroups[0].MaxCount
	// 获取原本弹性组下有多少台机器
	instanceList, err := c.getInstanceListByGroupID(ctx, &config)
	if err != nil {
		return false, err
	}
	var instanceUUIDList []string
	diags := config.InstanceUUIDList.ElementsAs(ctx, &instanceUUIDList, true)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return false, err
	}
	if len(instanceList)+len(instanceUUIDList) > int(maxCount) {
		err = fmt.Errorf("弹性伸缩组（id：%d）的max_count=%d。当前伸缩组下有云主机%d台，若手动移入%d台，将移入失败！", config.GroupID.ValueInt64(), maxCount, len(instanceList), len(instanceUUIDList))
		return false, err
	}
	return true, nil
}

func (c *ctyunScalingEcs) addScalingEcs(ctx context.Context, config *CtyunScalingEcsConfig) error {
	var instanceUUIDList []string
	diags := config.InstanceUUIDList.ElementsAs(ctx, &instanceUUIDList, true)
	if diags.HasError() {
		err := errors.New(diags[0].Detail())
		return err
	}
	err := c.addScalingEcsList(ctx, instanceUUIDList, config)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunScalingEcs) getAndMergeScalingEcs(ctx context.Context, config *CtyunScalingEcsConfig) error {
	instanceList, err := c.getInstanceListByGroupID(ctx, config)
	if err != nil {
		return err
	}
	var diags diag.Diagnostics
	config.InstanceUUIDList, diags = types.SetValueFrom(ctx, types.StringType, instanceList)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return err
	}

	return nil
}

func (c *ctyunScalingEcs) requestsScalingEcsList(ctx context.Context, config *CtyunScalingEcsConfig, pageNo int32, pageSize int32) (*scaling.ScalingGroupQueryInstanceListResponse, error) {
	params := &scaling.ScalingGroupQueryInstanceListRequest{
		RegionID: config.RegionID.ValueString(),
		GroupID:  config.GroupID.ValueInt64(),
		PageNo:   pageNo,
		PageSize: pageSize,
	}

	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupQueryInstanceListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询弹性伸缩组的列表失败，接口返回nil，伸缩组id：%d", config.GroupID.ValueInt64())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return nil, err
	}
	return resp, nil
}

func (c *ctyunScalingEcs) updateProtectStatus(ctx context.Context, state *CtyunScalingEcsConfig, plan *CtyunScalingEcsConfig) error {
	if !plan.ProtectStatus.IsNull() && !plan.ProtectStatus.IsUnknown() && !plan.ProtectStatus.Equal(state.ProtectStatus) {
		var instanceUUIDs []string
		diags := plan.InstanceUUIDList.ElementsAs(ctx, &instanceUUIDs, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}

		instanceIds, err := c.getInstanceAssocIdByUUID(ctx, state, instanceUUIDs)
		if err != nil {
			return err
		}
		// 关闭云主机保护
		if plan.ProtectStatus.ValueString() == business.StatusDisabledStr {
			err = c.disableProtectEcs(ctx, state, instanceIds, instanceUUIDs)
			if err != nil {
				return err
			}
		} else if plan.ProtectStatus.ValueString() == business.StatusEnabledStr {
			err = c.enableProtectEcs(ctx, state, instanceIds, instanceUUIDs)
			if err != nil {
				return err
			}
		}
	}
	state.ProtectStatus = plan.ProtectStatus
	return nil
}

func (c *ctyunScalingEcs) getInstanceAssocIdByUUID(ctx context.Context, state *CtyunScalingEcsConfig, UUIDs []string) ([]int32, error) {
	instanceMap, err := c.getInstanceMap(ctx, state)
	if err != nil {
		return nil, err
	}
	var instanceIdList []int32
	for _, uuid := range UUIDs {
		ecsInfo := instanceMap[uuid]
		instanceIdList = append(instanceIdList, ecsInfo.Id)
	}
	return instanceIdList, nil
}

func (c *ctyunScalingEcs) disableProtectEcs(ctx context.Context, state *CtyunScalingEcsConfig, instanceIDs []int32, instanceUUIDs []string) error {

	params := &scaling.ScalingGroupProtectDisableRequest{
		RegionID:       state.RegionID.ValueString(),
		GroupID:        state.GroupID.ValueInt64(),
		InstanceIDList: instanceIDs,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupProtectDisableApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("关闭云主机保护失败，接口返回nil，ecs列表：%s", strings.Join(instanceUUIDs, ", "))
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

func (c *ctyunScalingEcs) getInstanceMap(ctx context.Context, config *CtyunScalingEcsConfig) (map[string]*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse, error) {
	var pageNo, pageSize int32
	pageNo = 1
	pageSize = 100
	pageEndNo := pageNo
	ecsListResp, err := c.requestsScalingEcsList(ctx, config, pageNo, pageSize)
	if err != nil {
		return nil, err
	}

	totalCount := ecsListResp.ReturnObj.TotalCount

	if totalCount > pageSize {
		pageEndNo = totalCount / pageSize
	}
	// 先获取所有ecs列表，并设置成map
	var instanceMap map[string]*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse
	for pageNo <= pageEndNo {
		ecsList := ecsListResp.ReturnObj.InstanceList
		for _, ecs := range ecsList {
			instanceId := ecs.InstanceID
			instanceMap[instanceId] = ecs
		}

		pageNo++
		if pageNo > pageEndNo {
			break
		}
		ecsListResp, err = c.requestsScalingEcsList(ctx, config, pageNo, pageSize)
		if err != nil {
			return nil, err
		}
	}
	return instanceMap, nil
}

func (c *ctyunScalingEcs) enableProtectEcs(ctx context.Context, state *CtyunScalingEcsConfig, instanceIds []int32, instanceUUIDs []string) error {
	params := scaling.ScalingGroupProtectEnableRequest{
		RegionID:       state.RegionID.ValueString(),
		GroupID:        state.GroupID.ValueInt64(),
		InstanceIDList: instanceIds,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupProtectEnableApi.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("开启云主机保护失败，接口返回nil。ecs列表：%s", strings.Join(instanceUUIDs, ", "))
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

func (c *ctyunScalingEcs) removeEcs(ctx context.Context, config *CtyunScalingEcsConfig) error {
	var instanceUUIDs []string
	diags := config.InstanceUUIDList.ElementsAs(ctx, &instanceUUIDs, true)
	if diags.HasError() {
		err := errors.New(diags[0].Detail())
		return err
	}
	err := c.removeEcsRequest(ctx, instanceUUIDs, config, config)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunScalingEcs) updateInstanceUUIDList(ctx context.Context, state *CtyunScalingEcsConfig, plan *CtyunScalingEcsConfig) error {

	if plan.InstanceUUIDList.IsNull() || plan.InstanceUUIDList.IsUnknown() {
		return nil
	}
	if plan.InstanceUUIDList.Equal(state.InstanceUUIDList) {
		return nil
	}
	// 比对list， 挑出需要删除/新增的机器
	var stateInstanceList []string
	var planInstanceList []string
	diags := plan.InstanceUUIDList.ElementsAs(ctx, &stateInstanceList, true)
	if diags.HasError() {
		err := errors.New(diags[0].Detail())
		return err
	}
	diags = plan.InstanceUUIDList.ElementsAs(ctx, &planInstanceList, true)
	if diags.HasError() {
		err := errors.New(diags[0].Detail())
		return err
	}
	add, remove := c.getDiffInstanceList(stateInstanceList, planInstanceList)
	// 处理新增
	err := c.addScalingEcsList(ctx, add, state)
	if err != nil {
		return err
	}
	// 处理删除
	err = c.removeEcsRequest(ctx, remove, state, plan)
	if err != nil {
		return err
	}
	return nil

}

func (c *ctyunScalingEcs) getDiffInstanceList(state, plan []string) (toAdd, toRemove []string) {
	// 使用 map 快速查找差异
	planSet := make(map[string]bool)
	stateSet := make(map[string]bool)

	// 填充计划集
	for _, item := range plan {
		planSet[item] = true
	}

	// 填充状态集
	for _, item := range state {
		stateSet[item] = true
	}

	// 找出需要新增的项目（在 plan 中但不在 state 中）
	for _, item := range plan {
		if !stateSet[item] {
			toAdd = append(toAdd, item)
		}
	}

	// 找出需要删除的项目（在 state 中但不在 plan 中）
	for _, item := range state {
		if !planSet[item] {
			toRemove = append(toRemove, item)
		}
	}

	return toAdd, toRemove
}

func (c *ctyunScalingEcs) addScalingEcsList(ctx context.Context, instanceUUIDList []string, config *CtyunScalingEcsConfig) error {
	params := &scaling.ScalingGroupInstanceMoveInRequest{
		RegionID: config.RegionID.ValueString(),
		GroupID:  config.GroupID.ValueInt64(),
	}

	params.InstanceUUIDList = instanceUUIDList
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupInstanceMoveInApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("手动添加云主机失败，接口返回nil，伸缩组id：%d。需添加的云主机id列表为：%s", config.GroupID.ValueInt64(), strings.Join(instanceUUIDList, ", "))
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

func (c *ctyunScalingEcs) removeEcsRequest(ctx context.Context, instanceUUIDs []string, config *CtyunScalingEcsConfig, plan *CtyunScalingEcsConfig) error {
	if plan.IsDestroy.ValueBool() {
		// 移除并释放
		instanceIds, err := c.getInstanceAssocIdByUUID(ctx, config, instanceUUIDs)
		if err != nil {
			return err
		}
		params := &scaling.ScalingGroupInstanceMoveOutReleaseRequest{
			RegionID:       config.RegionID.ValueString(),
			GroupID:        config.GroupID.ValueInt64(),
			InstanceIDList: instanceIds,
		}
		resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupInstanceMoveOutReleaseApi.Do(ctx, c.meta.SdkCredential, params)

		if err != nil {
			return err
		} else if resp == nil {
			err = fmt.Errorf("ecs移除并释放失败，接口返回nil。ecs uuid 列表：%s", strings.Join(instanceUUIDs, ", "))
		} else if resp.StatusCode != common.NormalStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
			return err
		}
	} else {
		// 移除不释放
		params := &scaling.ScalingGroupInstanceMoveOutRequest{
			RegionID:       config.RegionID.ValueString(),
			GroupID:        config.GroupID.ValueInt64(),
			InstanceIDList: instanceUUIDs,
		}
		resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupInstanceMoveOutApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return err
		} else if resp == nil {
			err = fmt.Errorf("ecs移除失败，接口返回nil。ecs uuid 列表：%s", strings.Join(instanceUUIDs, ", "))
		} else if resp.StatusCode != common.NormalStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
			return err
		}
	}
	return nil
}

func (c *ctyunScalingEcs) getInstanceListByGroupID(ctx context.Context, config *CtyunScalingEcsConfig) ([]string, error) {
	var pageSize, pageNo int32
	pageSize = 100
	pageNo = 1
	resp, err := c.requestInstanceListByGroup(ctx, config, pageSize, pageNo)
	if err != nil {
		return nil, err
	}

	totalCount := resp.ReturnObj.TotalCount
	totalPageNo := pageNo
	if totalCount > pageSize {
		totalPageNo = totalCount/pageSize + 1
	}
	var instanceUUIDs []string
	for pageNo <= totalPageNo {
		instanceList := resp.ReturnObj.InstanceList
		for _, instance := range instanceList {
			instanceUUIDs = append(instanceUUIDs, instance.InstanceID)
		}
		pageNo++
		if pageNo > totalPageNo {
			break
		}
		resp, err = c.requestInstanceListByGroup(ctx, config, pageSize, pageNo)
		if err != nil {
			return nil, err
		}
	}
	return instanceUUIDs, nil
}

func (c *ctyunScalingEcs) requestInstanceListByGroup(ctx context.Context, config *CtyunScalingEcsConfig, pageSize, pageNo int32) (*scaling.ScalingGroupQueryInstanceListResponse, error) {
	params := &scaling.ScalingGroupQueryInstanceListRequest{
		RegionID: config.RegionID.ValueString(),
		GroupID:  config.GroupID.ValueInt64(),
		PageNo:   pageNo,
		PageSize: pageSize,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupQueryInstanceListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询group id为%d下的云主机列表失败，接口范围nil。请联系研发，或稍后重试！", config.GroupID.ValueInt64())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

type CtyunScalingEcsConfig struct {
	RegionID         types.String `tfsdk:"region_id"`          // 资源池ID
	GroupID          types.Int64  `tfsdk:"group_id"`           // 伸缩组ID
	InstanceUUIDList types.Set    `tfsdk:"instance_uuid_list"` // 云主机ID列表
	ProtectStatus    types.String `tfsdk:"protect_status"`     // 保护状态。1：已保护。2：未保护。
	IsDestroy        types.Bool   `tfsdk:"is_destroy"`         // 移除时是否销毁
}
