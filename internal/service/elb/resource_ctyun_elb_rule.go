package elb

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	strings "strings"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	ctelb "terraform-provider-ctyun/internal/core/ctelb"
	"terraform-provider-ctyun/internal/utils"
)

var (
	_ resource.Resource                = &CtyunElbRule{}
	_ resource.ResourceWithConfigure   = &CtyunElbRule{}
	_ resource.ResourceWithImportState = &CtyunElbRule{}
)

type CtyunElbRule struct {
	meta *common.CtyunMetadata
}

func NewCtyunSnatResource() resource.Resource {
	return &CtyunElbRule{}
}

func (c *CtyunElbRule) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunElbRule) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunElbRule) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_rule"
}

func (c *CtyunElbRule) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "区域ID",
			},
			"listener_id": schema.StringAttribute{
				Required:    true,
				Description: "监听器ID",
			},
			"priority": schema.Int32Attribute{
				Optional:    true,
				Description: "优先级，数字越小优先级越高，取值范围为：1-100(目前不支持配置此参数,只取默认值100)",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:'{},./;'[,]·~！@#￥%……&*（） —— -+={}",
			},
			"conditions": schema.ListNestedAttribute{
				Required:    true,
				Description: "匹配规则数据",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "类型。取值范围：server_name（服务名称）、url_path（匹配路径）",
							Validators: []validator.String{
								stringvalidator.OneOf(business.ElbRuleConditionTypes...),
							},
						},
						"condition_server_name": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "服务名称",
						},
						"condition_url_paths": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "匹配路径",
						},
						"condition_match_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "匹配类型。取值范围：ABSOLUTE，PREFIX，REG",
							Validators: []validator.String{
								stringvalidator.OneOf(business.ElbRuleMatchTypes...),
							},
						},
					},
				},
			},
			"action_type": schema.StringAttribute{
				Required:    true,
				Description: "默认规则动作类型。取值范围：forward、redirect、deny(目前暂不支持配置为deny)",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ElbRuleActionType...),
				},
			},
			"action_target_groups": schema.ListNestedAttribute{
				Optional:    true,
				Description: "后端服务组",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"target_group_id": schema.StringAttribute{
							Required:    true,
							Description: "后端服务组ID",
						},
						"weight": schema.Int32Attribute{
							Optional:    true,
							Description: "权重，取值范围：1-256。默认为100",
							Validators: []validator.Int32{
								int32validator.Between(1, 256),
							},
						},
					},
				},
			},
			"action_redirect_listener_id": schema.StringAttribute{
				Optional:    true,
				Description: "重定向监听器ID，当type为redirect时，此字段必填",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "转发规则 ID",
			},
			"az_name": schema.StringAttribute{
				Computed:    true,
				Description: "可用区名称",
			},
			"project_id": schema.StringAttribute{
				Computed:    true,
				Description: "项目ID",
			},
			"load_balancer_id": schema.StringAttribute{
				Computed:    true,
				Description: "负载均衡ID",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "状态: ACTIVE / DOWN",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ElbRuleStatus...),
				},
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
			},
			"updated_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间，为UTC格式",
			},
		},
	}
}

func (c *CtyunElbRule) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunElbRuleConfig

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 开始创建
	err = c.createElbRule(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后反查创建后的Rule信息
	err = c.getAndMergeRule(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbRule) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunElbRuleConfig
	if response.Diagnostics.HasError() {
		return
	}
	// 确认该rule是否异常
	if !c.acquireAndSetIdIfOrderNotFinished(ctx, &state, response) {
		return
	}

	//查询远端并同步state
	err = c.getAndMergeRule(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunElbRule) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置
	var plan CtyunElbRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunElbRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
	}

	// 更新rule信息
	err = c.updateElbRule(ctx, state, plan)
	if err != nil {
		return
	}
	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeRule(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (c *CtyunElbRule) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunElbRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	params := &ctelb.CtelbDeleteRuleRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionID.ValueString(),
		PolicyID:    state.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbDeleteRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
}

func (c *CtyunElbRule) createElbRule(ctx context.Context, plan *CtyunElbRuleConfig) (err error) {
	//配置创建接口所需请求参数
	regionId := c.meta.GetExtraIfEmpty(plan.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("创建转发规则时，regionID不能为空")
		return
	}
	if plan.ListenerID.IsNull() {
		err = errors.New("创建转发规则时，ListenerID不能为空")
		return
	}

	params := &ctelb.CtelbCreateRuleRequest{
		ClientToken: uuid.NewString(),
		RegionID:    regionId,
		ListenerID:  plan.ListenerID.ValueString(),
		Priority:    100, //目前不支持配置此参数,只取默认值100
	}
	// 构建condition参数体
	var conditionList []ConditionsModel
	var conditions []*ctelb.CtelbCreateRuleConditionsRequest
	if plan.Conditions.IsNull() {
		err = errors.New("创建转发规则时，conditions不能为空")
		return
	}
	// 将types.list->[]ConditionsModel
	diags := plan.Conditions.ElementsAs(ctx, &conditionList, false)
	if diags.HasError() {
		return
	}
	for _, conditionItem := range conditionList {
		var condition ctelb.CtelbCreateRuleConditionsRequest
		if conditionItem.Type.IsNull() {
			err = errors.New("创建转发规则时，condition_type不能为空")
		}

		condition.RawType = conditionItem.Type.ValueString()
		if !conditionItem.ConditionServerName.IsNull() {
			condition.ServerNameConfig.ServerName = conditionItem.ConditionServerName.ValueString()
		}
		if !conditionItem.ConditionMatchType.IsNull() {
			condition.UrlPathConfig.MatchType = conditionItem.ConditionMatchType.ValueString()
		}
		if !conditionItem.ConditionUrlPaths.IsNull() {
			condition.UrlPathConfig.UrlPaths = conditionItem.ConditionUrlPaths.ValueString()
		}
		conditions = append(conditions, &condition)
	}
	params.Conditions = conditions

	// 构建Action请求体
	var action *ctelb.CtelbCreateRuleActionRequest
	if plan.ActionType.IsNull() {
		err = errors.New("创建转发规则时，action type不能为空")
	}

	action.RawType = plan.ActionType.ValueString()
	if plan.ActionType.ValueString() == business.ElbRuleActionTypeRedirect && plan.ActionRedirectListenerID.IsNull() {
		err = errors.New("创建转发规则时，若action type = redirect, redirectListenerID不能为空")
		return
	}
	if !plan.ActionRedirectListenerID.IsNull() {
		action.RedirectListenerID = plan.ActionRedirectListenerID.ValueString()
	}

	// 构建action.forwardConfig请求体
	var targetGroupList []TargetGroupModel
	var targetGroups []*ctelb.CtelbCreateRuleActionForwardConfigTargetGroupsRequest
	diags = plan.ActionTargetGroups.ElementsAs(ctx, &targetGroupList, false)
	if diags.HasError() {
		return
	}
	for _, targetGroupItem := range targetGroupList {
		var targetGroup ctelb.CtelbCreateRuleActionForwardConfigTargetGroupsRequest
		if targetGroupItem.TargetGroupID.IsNull() {
			err = errors.New("创建转发规则时，targetGroupID不能为空")
			return
		}
		targetGroup.TargetGroupID = targetGroupItem.TargetGroupID.ValueString()
		if !targetGroupItem.Weight.IsNull() {
			targetGroup.Weight = targetGroupItem.Weight.ValueInt32()
		}
		targetGroups = append(targetGroups, &targetGroup)
	}
	action.ForwardConfig.TargetGroups = targetGroups
	params.Action = action

	// 调用创建接口
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbCreateRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
	var ids []string
	// 获取规则id
	for _, returnOjbItem := range resp.ReturnObj {
		ids = append(ids, returnOjbItem.ID)
	}

	idsTmp := strings.Join(ids, ",")
	plan.ID = types.StringValue(idsTmp)
	return
}

func (c *CtyunElbRule) getAndMergeRule(ctx context.Context, plan *CtyunElbRuleConfig) (err error) {
	//查看rule详情
	params := &ctelb.CtelbShowRuleRequest{
		RegionID: plan.RegionID.ValueString(),
		PolicyID: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbShowRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
	returnObj := resp.ReturnObj
	//解析rule明细，将远端明细合并至本地
	plan.AzName = types.StringValue(returnObj.AzName)
	plan.ProjectID = types.StringValue(returnObj.ProjectID)
	plan.LoadBalancerID = types.StringValue(returnObj.LoadBalancerID)
	plan.Status = types.StringValue(returnObj.Status)
	plan.CreatedTime = types.StringValue(returnObj.CreatedTime)
	plan.UpdatedTime = types.StringValue(returnObj.UpdatedTime)

	// 合并conditions
	conditionList := returnObj.Conditions
	var conditions []ConditionsModel
	for _, conditionItem := range conditionList {
		var condition ConditionsModel
		condition.ConditionServerName = types.StringValue(conditionItem.ServerNameConfig.ServerName)
		condition.ConditionUrlPaths = types.StringValue(conditionItem.UrlPathConfig.UrlPaths)
		condition.ConditionMatchType = types.StringValue(conditionItem.UrlPathConfig.MatchType)
		conditions = append(conditions, condition)
	}
	plan.Conditions, _ = types.ListValueFrom(ctx, utils.StructToTFObjectTypes(ConditionsModel{}), conditions)

	// 处理action
	plan.ActionType = types.StringValue(returnObj.Action.RawType)
	plan.ActionRedirectListenerID = types.StringValue(returnObj.Action.RedirectListenerID)

	targetGroupList := returnObj.Action.ForwardConfig.TargetGroups
	var targetGroups []TargetGroupModel
	for _, targetGroupItem := range targetGroupList {
		var targetGroup TargetGroupModel
		targetGroup.TargetGroupID = types.StringValue(targetGroupItem.TargetGroupID)
		targetGroup.Weight = types.Int32Value(targetGroupItem.Weight)
		targetGroups = append(targetGroups, targetGroup)
	}
	plan.ActionTargetGroups, _ = types.ListValueFrom(ctx, utils.StructToTFObjectTypes(TargetGroupModel{}), targetGroups)
	return
}

func (c *CtyunElbRule) updateElbRule(ctx context.Context, state CtyunElbRuleConfig, plan CtyunElbRuleConfig) (err error) {
	params := &ctelb.CtelbUpdateRuleRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionID.ValueString(),
		PolicyID:    state.ID.ValueString(),
		Conditions:  nil,
		Action:      nil,
	}
	// 处理condition更新值
	var conditionList []ConditionsModel
	var conditions []*ctelb.CtelbUpdateRuleConditionsRequest

	// 将types.list->[]ConditionsModel
	diags := plan.Conditions.ElementsAs(ctx, &conditionList, false)
	if diags.HasError() {
		return
	}
	if len(conditionList) > 0 {
		for _, conditionItem := range conditionList {
			var condition ctelb.CtelbUpdateRuleConditionsRequest
			if conditionItem.Type.IsNull() {
				err = errors.New("更新转发规则时，condition_type不能为空")
			}
			condition.RawType = conditionItem.Type.ValueString()
			if !conditionItem.ConditionServerName.IsNull() {
				condition.ServerNameConfig.ServerName = conditionItem.ConditionServerName.ValueString()
			}
			if !conditionItem.ConditionMatchType.IsNull() {
				condition.UrlPathConfig.MatchType = conditionItem.ConditionMatchType.ValueString()
			}
			if !conditionItem.ConditionUrlPaths.IsNull() {
				condition.UrlPathConfig.UrlPaths = conditionItem.ConditionUrlPaths.ValueString()
			}
			conditions = append(conditions, &condition)
		}
		params.Conditions = conditions
	}

	// 处理action更新值

	var action *ctelb.CtelbUpdateRuleActionRequest
	if plan.ActionType.IsNull() {
		err = errors.New("修改转发规则时，action type不能为空")
	}

	action.RawType = plan.ActionType.ValueString()
	if plan.ActionType.ValueString() == business.ElbRuleActionTypeRedirect && plan.ActionRedirectListenerID.IsNull() {
		err = errors.New("修改转发规则时，若action type = redirect, redirectListenerID不能为空")
		return
	}
	if !plan.ActionRedirectListenerID.IsNull() {
		action.RedirectListenerID = plan.ActionRedirectListenerID.ValueString()
	}

	// 构建action.forwardConfig请求体
	var targetGroupList []TargetGroupModel
	var targetGroups []*ctelb.CtelbUpdateRuleActionForwardConfigTargetGroupsRequest
	diags = plan.ActionTargetGroups.ElementsAs(ctx, &targetGroupList, false)
	if diags.HasError() {
		return
	}
	for _, targetGroupItem := range targetGroupList {
		var targetGroup ctelb.CtelbUpdateRuleActionForwardConfigTargetGroupsRequest
		if targetGroupItem.TargetGroupID.IsNull() {
			err = errors.New("修改转发规则时，targetGroupID不能为空")
			return
		}
		targetGroup.TargetGroupID = targetGroupItem.TargetGroupID.ValueString()
		if !targetGroupItem.Weight.IsNull() {
			targetGroup.Weight = targetGroupItem.Weight.ValueInt32()
		}
		targetGroups = append(targetGroups, &targetGroup)
	}
	action.ForwardConfig.TargetGroups = targetGroups
	params.Action = action

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
	return
}

func (c *CtyunElbRule) acquireAndSetIdIfOrderNotFinished(ctx context.Context, state *CtyunElbRuleConfig, response *resource.ReadResponse) bool {
	if state.ID.IsNull() {
		// 该rule没有id，为非法id。移除当前状态并返回
		response.State.RemoveResource(ctx)
		return false
	}
	return true
}

type CtyunElbRuleConfig struct {
	RegionID                 types.String `tfsdk:"region_id"`                   //区域ID
	ListenerID               types.String `tfsdk:"listener_id"`                 //监听器ID
	Priority                 types.Int32  `tfsdk:"priority"`                    //优先级，数字越小优先级越高，取值范围为：1-100(目前不支持配置此参数,只取默认值100)
	Description              types.String `tfsdk:"description"`                 //支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:'{},./;'[,]·~！@#￥%……&*（） —— -+={},
	Conditions               types.List   `tfsdk:"conditions"`                  //匹配规则数据
	ActionType               types.String `tfsdk:"action_type"`                 //默认规则动作类型。取值范围：forward、redirect、deny(目前暂不支持配置为deny)
	ActionTargetGroups       types.List   `tfsdk:"action_target_groups"`        //后端服务组
	ActionRedirectListenerID types.String `tfsdk:"action_redirect_listener_id"` //重定向监听器ID，当type为redirect时，此字段必填
	ID                       types.String `tfsdk:"id"`                          //转发规则 ID
	AzName                   types.String `tfsdk:"az_name"`                     //可用区名称
	ProjectID                types.String `tfsdk:"project_id"`                  //	项目ID
	LoadBalancerID           types.String `tfsdk:"load_balancer_id"`            //负载均衡ID
	Status                   types.String `tfsdk:"status"`                      //状态: ACTIVE / DOWN
	CreatedTime              types.String `tfsdk:"created_time"`                //创建时间，为UTC格式
	UpdatedTime              types.String `tfsdk:"updated_time"`                //更新时间，为UTC格式
}

type ConditionsModel struct {
	Type                types.String `tfsdk:"type"`                  //类型。取值范围：server_name（服务名称）、url_path（匹配路径）
	ConditionServerName types.String `tfsdk:"condition_server_name"` //服务名称
	ConditionUrlPaths   types.String `tfsdk:"condition_url_paths"`   //匹配路径
	ConditionMatchType  types.String `tfsdk:"condition_match_type"`  //匹配类型。取值范围：ABSOLUTE，PREFIX，REG
}

type TargetGroupModel struct {
	TargetGroupID types.String `tfsdk:"target_group_id"` //后端服务组ID
	Weight        types.Int32  `tfsdk:"weight"`          //权重，取值范围：1-256。默认为100
}
