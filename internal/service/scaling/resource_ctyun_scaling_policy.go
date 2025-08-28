package scaling

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/scaling"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type ctyunScalingPolicy struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
	imageService  *business.ImageService
}

func (c *ctyunScalingPolicy) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_scaling_policy"
}

func (c *ctyunScalingPolicy) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)
	c.imageService = business.NewImageService(c.meta)
}

func NewCtyunScalingPolicy() resource.Resource {
	return &ctyunScalingPolicy{}
}

func (c *ctyunScalingPolicy) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "弹性伸缩策略管理，支持伸缩策略的创建、修改和删除。具体细节可参考文档：https://www.ctyun.cn/document/10027725/10241454",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "区域ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "伸缩组ID",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "伸缩策略名称",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"policy_type": schema.StringAttribute{
				Required:    true,
				Description: "策略类型: alert-告警策略, regular-定时策略, period-周期策略, target-目标追踪策略",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ScalingPolicyTypes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"operate_unit": schema.StringAttribute{
				Optional:    true,
				Description: "操作单位: count-个数, percent-百分比",
				Validators: []validator.String{
					stringvalidator.OneOf(business.OperateUnits...),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("policy_type"),
						types.StringValue("alert"),
						types.StringValue("regular"),
						types.StringValue("period"),
					),
				},
			},
			"operate_count": schema.Int32Attribute{
				Optional:    true,
				Description: "调整值",
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("policy_type"),
						types.StringValue("alert"),
						types.StringValue("regular"),
						types.StringValue("period"),
					),
				},
			},
			"action": schema.StringAttribute{
				Optional:    true,
				Description: "执行动作: increase-增加, decrease-减少, set-设置为",
				Validators: []validator.String{
					stringvalidator.OneOf(business.Actions...),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("policy_type"),
						types.StringValue("alert"),
						types.StringValue("regular"),
						types.StringValue("period"),
					),
				},
			},
			// todo 当policy_type = 3必填
			"cycle": schema.StringAttribute{
				Optional:    true,
				Description: "循环方式: monthly-按月循环, weekly-按周循环, daily-按天循环",
				Validators: []validator.String{
					stringvalidator.OneOf(business.Cycles...),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("policy_type"),
						types.StringValue("period"),
					),
				},
			},
			// todo ,做校验，当cycle=monthly，取值范围【1，31】， 当cycle=weekly，取值范围【1，7】
			"day": schema.SetAttribute{
				ElementType: types.Int32Type,
				Optional:    true,
				Description: "执行日期",
				Validators: []validator.Set{
					validator2.AlsoRequiresEqualSet(
						path.MatchRoot("cycle"),
						types.StringValue("monthly"),
						types.StringValue("weekly"),
					),
				},
			},
			// todo 需要做校验, type=3必填
			"effective_from": schema.StringAttribute{
				Optional:    true,
				Description: "周期策略生效开始时间 (格式: 2006-01-02 15:04:05)",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("policy_type"),
						types.StringValue("period"),
					),
				},
			},
			// todo 需要做校验, type=3必填
			"effective_till": schema.StringAttribute{
				Optional:    true,
				Description: "周期策略生效截止时间 (格式: 2006-01-02 15:04:05)",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("policy_type"),
						types.StringValue("period"),
					),
				},
			},
			// todo 需要做校验, type=2， 3必填
			"execution_time": schema.StringAttribute{
				Optional:    true,
				Description: "执行时间 (格式: 2006-01-02 15:04:05)",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("policy_type"),
						types.StringValue("regular"),
						types.StringValue("period"),
					),
				},
			},
			"cooldown": schema.Int32Attribute{
				Optional:    true,
				Description: "冷却/预热时间 (单位: 秒)",
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("policy_type"),
						types.StringValue("alert"),
						types.StringValue("target"),
					),
				},
			},
			"trigger_name": schema.StringAttribute{
				Optional:    true,
				Description: "告警规则名称",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("policy_type"),
						types.StringValue("alert"),
					),
				},
			},
			"trigger_metric_name": schema.StringAttribute{
				Optional: true,
				Description: "监控指标名称，取值范围：cpu_util-CPU使用率，mem_util-内存使用率，network_incoming_bytes_rate_inband-网络流入速率，network_outing_bytes_rate_inband-网络流出速率，" +
					"disk_read_bytes_rate-磁盘读速率，disk_write_bytes_rate-磁盘写速率，disk_read_requests_rate-磁盘读请求速率，disk_write_requests_rate-磁盘写请求速率",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"cpu_util",
						"mem_util",
						"network_incoming_bytes_rate_inband",
						"network_outing_bytes_rate_inband",
						"disk_read_bytes_rate",
						"disk_write_bytes_rate",
						"disk_read_requests_rate",
						"disk_write_requests_rate",
					),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("policy_type"),
						types.StringValue("alert"),
					),
				},
			},
			"trigger_statistics": schema.StringAttribute{
				Optional:    true,
				Description: "聚合方法: avg-平均值, max-最大值, min-最小值",
				Validators: []validator.String{
					stringvalidator.OneOf("avg", "max", "min"),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("policy_type"),
						types.StringValue("alert"),
					),
				},
			},
			"trigger_comparison_operator": schema.StringAttribute{
				Optional:    true,
				Description: "比较符: ge-大于等于, le-小于等于, gt-大于, lt-小于",
				Validators: []validator.String{
					stringvalidator.OneOf("ge", "le", "gt", "lt"),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("policy_type"),
						types.StringValue("alert"),
					),
				},
			},
			"trigger_threshold": schema.Int32Attribute{
				Optional:    true,
				Description: "阈值",
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("policy_type"),
						types.StringValue("alert"),
					),
				},
			},
			"trigger_period": schema.StringAttribute{
				Optional:    true,
				Description: "监控周期 (如: 5m)",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("policy_type"),
						types.StringValue("alert"),
					),
				},
			},
			"trigger_evaluation_count": schema.Int32Attribute{
				Optional:    true,
				Description: "连续出现次数",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("policy_type"),
						types.StringValue("alert"),
					),
				},
			},
			"target_metric_name": schema.StringAttribute{
				Optional:    true,
				Description: "监控指标名称，取值范围：cpu_util-CPU使用率，network_incoming_bytes_rate_inband-网络流入速率，network_outing_bytes_rate_inband-网络流出速率",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"cpu_util",
						"network_incoming_bytes_rate_inband",
						"network_outing_bytes_rate_inband",
					),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("policy_type"),
						types.StringValue("target"),
					),
				},
			},
			"target_value": schema.Int32Attribute{
				Optional:    true,
				Description: "追踪目标值",
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("policy_type"),
						types.StringValue("target"),
					),
				},
			},
			"target_scale_out_evaluation_count": schema.Int32Attribute{
				Optional:    true,
				Description: "扩容连续告警次数 (范围: 1-100)",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("policy_type"),
						types.StringValue("target"),
					),
				},
			},
			"target_scale_in_evaluation_count": schema.Int32Attribute{
				Optional:    true,
				Description: "缩容连续告警次数 (范围: 1-100)",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("policy_type"),
						types.StringValue("target"),
					),
				},
			},
			"target_operate_range": schema.Int32Attribute{
				Optional:    true,
				Description: "缩容波动范围 (范围: 10-20)，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(10, 20),
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("policy_type"),
						types.StringValue("target"),
					),
				},
			},
			"target_disable_scale_in": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "是否禁用缩容，支持更新。默认为false。",
				Validators: []validator.Bool{
					validator2.AlsoRequiresEqualBool(
						path.MatchRoot("policy_type"),
						types.StringValue("target"),
					),
				},
			},
			"is_execute": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "控制是否需要执行弹性伸缩策略，true表示执行，false表示不执行。默认为false",
				Default:     booldefault.StaticBool(false),
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "伸缩策略id",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "告警规则状态：enable：启用。disable：停用",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ScalingPolicyStatuses...),
				},
			},
		},
	}
}

func (c *ctyunScalingPolicy) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunScalingPolicyConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	//创建前检查,参数有效性
	isValid, err := c.checkBeforeScalingPolicy(ctx, plan)
	if !isValid || err != nil {
		return
	}
	err = c.createScalingPolicy(ctx, &plan)
	if err != nil {
		return
	}
	// 创建后，通过创建的请求轮询，确认创建成功
	//err = c.createLoop(ctx, &plan, createParams, 60)
	if err != nil {
		return
	}
	// 创建后反查创建后的证书信息
	err = c.getAndMergeScalingPolicy(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunScalingPolicy) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunScalingPolicyConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeScalingPolicy(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "未找到") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunScalingPolicy) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取 plan -tf文件中配置
	var plan CtyunScalingPolicyConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunScalingPolicyConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 停用/启用伸缩策略
	err = c.updatePolicyStatus(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新基本信息
	err = c.updateScalingPolicy(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeScalingPolicy(ctx, &state)
	if err != nil {
		return
	}

	// 执行该条伸缩策略
	err = c.executePolicy(ctx, &state, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunScalingPolicy) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunScalingPolicyConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &scaling.ScalingRuleDeleteRequest{
		RegionID: state.RegionID.ValueString(),
		GroupID:  state.GroupID.ValueInt64(),
		RuleID:   state.ID.ValueInt64(),
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingRuleDeleteApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("删除弹性伸缩组:%d 中的伸缩策略：%d的失败，接口返回nil。请联系研发确认。或稍后重试！", state.GroupID.ValueInt64(), state.ID)
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *ctyunScalingPolicy) checkBeforeScalingPolicy(ctx context.Context, plan CtyunScalingPolicyConfig) (bool, error) {
	return true, nil

}

func (c *ctyunScalingPolicy) createScalingPolicy(ctx context.Context, config *CtyunScalingPolicyConfig) error {
	params := &scaling.ScalingRuleCreateRequest{
		RegionID: config.RegionID.ValueString(),
		GroupID:  config.GroupID.ValueInt64(),
		Name:     config.Name.ValueString(),
		RawType:  business.ScalingPolicyTypeDict[config.PolicyType.ValueString()],
	}
	// 当伸缩策略为告警策略
	if config.PolicyType.ValueString() == business.ScalingPolicyAlertStr {
		params.OperateUnit = business.OperateUnitDict[config.OperateUnit.ValueString()]
		params.OperateCount = config.OperateCount.ValueInt32()
		params.Action = business.ActionDict[config.Action.ValueString()]
		params.Cooldown = config.Cooldown.ValueInt32()
		triggerObj := scaling.ScalingRuleCreateTriggerObjRequest{
			Name:               config.TriggerName.ValueString(),
			MetricName:         config.TriggerMetricName.ValueString(),
			Statistics:         config.TriggerStatistics.ValueString(),
			ComparisonOperator: config.TriggerComparisonOperator.ValueString(),
			Threshold:          config.TriggerThreshold.ValueInt32(),
			Period:             config.TriggerPeriod.ValueString(),
			EvaluationCount:    config.TriggerEvaluationCount.ValueInt32(),
		}
		params.TriggerObj = &triggerObj

	} else if config.PolicyType.ValueString() == business.ScalingPolicyRegularStr {
		// 当伸缩策略为定时策略
		params.OperateUnit = business.OperateUnitDict[config.OperateUnit.ValueString()]
		params.OperateCount = config.OperateCount.ValueInt32()
		params.Action = business.ActionDict[config.Action.ValueString()]
		params.ExecutionTime = config.ExecutionTime.ValueString()
	} else if config.PolicyType.ValueString() == business.ScalingPolicyPeriodStr {
		// 当伸缩策略为周期策略
		params.OperateUnit = business.OperateUnitDict[config.OperateUnit.ValueString()]
		params.OperateCount = config.OperateCount.ValueInt32()
		params.Action = business.ActionDict[config.Action.ValueString()]
		params.Cycle = business.CycleDict[config.Cycle.ValueString()]
		var day []int32
		diags := config.Day.ElementsAs(ctx, &day, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		params.Day = day
		params.EffectiveFrom = config.EffectiveFrom.ValueString()
		params.EffectiveTill = config.EffectiveTill.ValueString()
		params.ExecutionTime = config.ExecutionTime.ValueString()

	} else if config.PolicyType.ValueString() == business.ScalingPolicyTargetStr {
		// 当伸缩策略为目标追踪策略
		params.Cooldown = config.Cooldown.ValueInt32() // 预热时间
		targetObj := scaling.ScalingRuleCreateTargetObjRequest{
			MetricName:              config.TargetMetricName.ValueString(),
			TargetValue:             config.TargetValue.ValueInt32(),
			ScaleOutEvaluationCount: config.TargetScaleOutEvaluationCount.ValueInt32(),
			ScaleInEvaluationCount:  config.TargetScaleInEvaluationCount.ValueInt32(),
			OperateRange:            config.TargetOperateRange.ValueInt32(),
			DisableScaleIn:          config.TargetDisableScaleIn.ValueBool(),
		}
		params.TargetObj = &targetObj
	} else {
		err := fmt.Errorf("创建参数policy type 输入有误，输入值为：%s", config.PolicyType.ValueString())
		return err
	}

	resp, err := c.meta.Apis.SdkScalingApis.ScalingRuleCreateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = errors.New("创建伸缩策略失败，接口返回nil。请联系研发确认！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	config.ID = types.Int64Value(resp.ReturnObj.RuleID)
	return nil
}

func (c *ctyunScalingPolicy) getAndMergeScalingPolicy(ctx context.Context, config *CtyunScalingPolicyConfig) error {
	rule, err := c.getScalingPolicyDetail(ctx, config)
	if err != nil {
		return err
	}
	if business.ScalingPolicyTypeDictRev[rule.RuleType] != config.PolicyType.ValueString() {
		err = fmt.Errorf("伸缩策略详情有误，id为：%d，本地和控制台上策略类型不一致。本地策略类型为：%s，但是控制台上策略类型为：%s", config.ID.ValueInt64(), config.PolicyType.ValueString(), business.ScalingPolicyTypeDictRev[rule.RuleType])
		return err
	}
	config.Name = types.StringValue(rule.Name)
	config.Status = types.StringValue(business.ScalingPolicyStatusDictRev[rule.Status])

	if config.PolicyType.ValueString() == business.ScalingPolicyAlertStr {
		// 告警策略
		// 触发字段
		config.TriggerName = types.StringValue(rule.TriggerObj.Name)
		config.TriggerMetricName = types.StringValue(rule.TriggerObj.MetricName)
		config.TriggerStatistics = types.StringValue(rule.TriggerObj.Statistics)
		config.TriggerComparisonOperator = types.StringValue(rule.TriggerObj.ComparisonOperator)
		config.TriggerThreshold = types.Int32Value(rule.TriggerObj.Threshold)
		config.TriggerPeriod = types.StringValue(rule.TriggerObj.Period)
		config.TriggerEvaluationCount = types.Int32Value(rule.TriggerObj.EvaluationCount)
		// 冷却时间 (秒)
		config.Cooldown = types.Int32Value(rule.Cooldown)
		// 执行动作
		config.Action = types.StringValue(business.ActionDictRev[rule.Action])
		config.OperateCount = types.Int32Value(rule.OperateCount)
		config.OperateUnit = types.StringValue(business.OperateUnitDictRev[rule.OperateUnit])
	} else if config.PolicyType.ValueString() == business.ScalingPolicyRegularStr {
		// 定时策略
		// 触发时间
		config.ExecutionTime = types.StringValue(rule.ExecutionTime)
		// 执行动作
		config.Action = types.StringValue(business.ActionDictRev[rule.Action])
		config.OperateCount = types.Int32Value(rule.OperateCount)
		config.OperateUnit = types.StringValue(business.OperateUnitDictRev[rule.OperateUnit])
	} else if config.PolicyType.ValueString() == business.ScalingPolicyPeriodStr {
		// 周期策略
		config.Cycle = types.StringValue(business.CycleDictRev[rule.Cycle])
		// 触发时间
		config.ExecutionTime = types.StringValue(rule.ExecutionTime)
		// 生效时间
		config.EffectiveFrom = types.StringValue(rule.EffectiveFrom)
		config.EffectiveTill = types.StringValue(rule.EffectiveTill)
		// 执行动作
		config.Action = types.StringValue(business.ActionDictRev[rule.Action])
		config.OperateCount = types.Int32Value(rule.OperateCount)
		config.OperateUnit = types.StringValue(business.OperateUnitDictRev[rule.OperateUnit])
		var diags diag.Diagnostics
		config.Day, diags = types.SetValueFrom(ctx, types.Int32Type, rule.Day)
		if diags.HasError() {
			err = fmt.Errorf(diags[0].Detail())
			return err
		}
	} else if config.PolicyType.ValueString() == business.ScalingPolicyTargetStr {
		// 目标追踪策略
		// 目标追踪字段
		// 目标值
		config.TargetMetricName = types.StringValue(rule.TargetObj.MetricName)
		// 目标值
		config.TargetValue = types.Int32Value(rule.TargetObj.TargetValue)
		// 缩容波动范围
		config.TargetOperateRange = types.Int32Value(rule.TargetObj.OperateRange)
		// 预热时间
		config.Cooldown = types.Int32Value(rule.Cooldown)
		// 扩容连续告警次数
		config.TargetScaleOutEvaluationCount = types.Int32Value(rule.TargetObj.ScaleOutEvaluationCount)
		// 缩容连续告警次数
		config.TargetScaleInEvaluationCount = types.Int32Value(rule.TargetObj.ScaleInEvaluationCount)
		//缩容波动范围
		config.TargetDisableScaleIn = types.BoolValue(*rule.TargetObj.DisableScaleIn)
	}
	return nil
}

func (c *ctyunScalingPolicy) getScalingPolicyDetail(ctx context.Context, config *CtyunScalingPolicyConfig) (*scaling.ScalingRuleListReturnObjRuleListResponse, error) {
	var pageNo, pageSize int32
	pageNo = 1
	pageSize = 100
	pageEndNo := pageNo

	resp, err := c.requestRuleList(ctx, config, pageNo, pageSize)
	if err != nil {
		return nil, err
	}
	// 伸缩策略只有列表查询，需要先查询到再通过遍历定位到具体策略
	policyNum := resp.ReturnObj.NumberOfAll
	if policyNum <= 0 {
		return nil, fmt.Errorf("未查询到弹性伸缩组：%d下有任何伸缩策略", config.GroupID)
	}

	// 如果策略数量大于页面大小，则需要翻页获取。
	if policyNum > pageSize {
		pageEndNo = policyNum/pageSize + 1
	}

	for pageNo <= pageEndNo {
		ruleList := resp.ReturnObj.RuleList
		for _, rule := range ruleList {
			if rule.RuleID != config.ID.ValueInt64() {
				continue
			}
			return rule, nil
		}

		pageNo++
		if pageNo > pageEndNo {
			break
		}
		resp, err = c.requestRuleList(ctx, config, pageNo, pageSize)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (c *ctyunScalingPolicy) requestRuleList(ctx context.Context, config *CtyunScalingPolicyConfig, pageNo int32, pageSize int32) (*scaling.ScalingRuleListResponse, error) {
	params := &scaling.ScalingRuleListRequest{
		RegionID: config.RegionID.ValueString(),
		GroupID:  config.GroupID.ValueInt64(),
		PageNo:   pageNo,
		PageSize: pageSize,
	}

	resp, err := c.meta.Apis.SdkScalingApis.ScalingRuleListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询伸缩策略列表失败，接口返回nil。对应的策略组id为：{%d}, 请联系研发确认！", config.GroupID.ValueInt64())
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

func (c *ctyunScalingPolicy) updateScalingPolicy(ctx context.Context, state *CtyunScalingPolicyConfig, plan *CtyunScalingPolicyConfig) error {

	params := &scaling.ScalingRuleUpdateRequest{
		RegionID: state.RegionID.ValueString(),
		GroupID:  state.GroupID.ValueInt64(),
		RuleID:   state.ID.ValueInt64(),
		Name:     plan.Name.ValueString(),
	}
	// policy_type 不可变，policy_type
	if state.PolicyType.ValueString() == business.ScalingPolicyAlertStr {
		params.OperateUnit = business.OperateUnitDict[plan.OperateUnit.ValueString()]
		params.OperateCount = plan.OperateCount.ValueInt32()
		params.Action = business.ActionDict[plan.Action.ValueString()]

		params.Cooldown = plan.Cooldown.ValueInt32()

		triggerObj := scaling.ScalingRuleUpdateTriggerObjRequest{
			Name:               plan.TriggerName.ValueString(),
			MetricName:         plan.TriggerMetricName.ValueString(),
			Statistics:         plan.TriggerStatistics.ValueString(),
			ComparisonOperator: plan.TriggerComparisonOperator.ValueString(),
			Threshold:          plan.TriggerThreshold.ValueInt32(),
			Period:             plan.TriggerPeriod.ValueString(),
			EvaluationCount:    plan.TriggerEvaluationCount.ValueInt32(),
		}
		params.TriggerObj = &triggerObj

	} else if state.PolicyType.ValueString() == business.ScalingPolicyRegularStr {
		// 当伸缩策略为定时策略
		params.OperateUnit = business.OperateUnitDict[plan.OperateUnit.ValueString()]
		params.OperateCount = plan.OperateCount.ValueInt32()
		params.Action = business.ActionDict[plan.Action.ValueString()]
		params.ExecutionTime = plan.ExecutionTime.ValueString()
	} else if state.PolicyType.ValueString() == business.ScalingPolicyPeriodStr {
		// 当伸缩策略为周期策略
		params.OperateUnit = business.OperateUnitDict[plan.OperateUnit.ValueString()]
		params.OperateCount = plan.OperateCount.ValueInt32()
		params.Action = business.ActionDict[plan.Action.ValueString()]
		params.Cycle = business.CycleDict[plan.Cycle.ValueString()]
		var day []int32
		diags := plan.Day.ElementsAs(ctx, &day, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		params.Day = day
		params.EffectiveFrom = plan.EffectiveFrom.ValueString()
		params.EffectiveTill = plan.EffectiveTill.ValueString()
		params.ExecutionTime = plan.ExecutionTime.ValueString()

	} else if state.PolicyType.ValueString() == business.ScalingPolicyTargetStr {
		// 当伸缩策略为目标追踪策略
		params.Cooldown = plan.Cooldown.ValueInt32() // 预热时间
		targetObj := scaling.ScalingRuleUpdateTargetObjRequest{
			MetricName:              plan.TargetMetricName.ValueString(),
			TargetValue:             plan.TargetValue.ValueInt32(),
			ScaleOutEvaluationCount: plan.TargetScaleOutEvaluationCount.ValueInt32(),
			ScaleInEvaluationCount:  plan.TargetScaleInEvaluationCount.ValueInt32(),
			OperateRange:            plan.TargetOperateRange.ValueInt32(),
			DisableScaleIn:          plan.TargetDisableScaleIn.ValueBoolPointer(),
		}
		params.TargetObj = &targetObj
	} else {
		err := fmt.Errorf("创建参数policy type 输入有误，输入值为：%s", state.PolicyType.ValueString())
		return err
	}

	resp, err := c.meta.Apis.SdkScalingApis.ScalingRuleUpdateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("伸缩策略更新失败，id:%d", state.ID.ValueInt64())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}

	return nil
}

func (c *ctyunScalingPolicy) updatePolicyStatus(ctx context.Context, state *CtyunScalingPolicyConfig, plan *CtyunScalingPolicyConfig) error {
	// 若plan阶段status发生变化，触发停用/启用
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() && !plan.Status.Equal(state.Status) {
		// 启用
		if plan.Status.ValueString() == business.StatusEnabledStr {
			params := &scaling.ScalingRuleStartRequest{
				RegionID: state.RegionID.ValueString(),
				GroupID:  state.GroupID.ValueInt64(),
				RuleID:   state.ID.ValueInt64(),
			}
			resp, err := c.meta.Apis.SdkScalingApis.ScalingRuleStartApi.Do(ctx, c.meta.SdkCredential, params)
			if err != nil {
				return err
			} else if resp == nil {
				err = fmt.Errorf("启用弹性伸缩策略失败，其id:%d", state.ID.ValueInt64())
				return err
			} else if resp.StatusCode != common.NormalStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
				return err
			}
		} else if plan.Status.ValueString() == business.StatusDisabledStr {
			params := &scaling.ScalingRuleStopRequest{
				RegionID: state.RegionID.ValueString(),
				GroupID:  state.GroupID.ValueInt64(),
				RuleID:   state.ID.ValueInt64(),
			}
			resp, err := c.meta.Apis.SdkScalingApis.ScalingRuleStopApi.Do(ctx, c.meta.SdkCredential, params)
			if err != nil {
				return err
			} else if resp == nil {
				err = fmt.Errorf("停用弹性伸缩策略失败，其id:%d", state.ID.ValueInt64())
				return err
			} else if resp.StatusCode != common.NormalStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
				return err
			}
		}
	}
	return nil
}

func (c *ctyunScalingPolicy) executePolicy(ctx context.Context, state *CtyunScalingPolicyConfig, plan *CtyunScalingPolicyConfig) error {
	if plan.IsExecute.ValueBool() {
		params := &scaling.ScalingRuleExecuteRequest{
			RegionID: state.RegionID.ValueString(),
			RuleID:   state.ID.ValueInt64(),
			GroupID:  state.GroupID.ValueInt64(),
		}
		resp, err := c.meta.Apis.SdkScalingApis.ScalingRuleExecuteApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return err
		} else if resp == nil {
			err = fmt.Errorf("伸缩策略执行失败，id：:%d。接口返回nil，具体原因可联系研发咨询。", state.ID.ValueInt64())
			return err
		} else if resp.StatusCode != common.NormalStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
			return err
		}
	}
	state.IsExecute = plan.IsExecute
	return nil
}

type CtyunScalingPolicyConfig struct {
	RegionID                      types.String `tfsdk:"region_id"`                         // 资源池id
	GroupID                       types.Int64  `tfsdk:"group_id"`                          // 伸缩组ID
	Name                          types.String `tfsdk:"name"`                              // 伸缩策略名称
	PolicyType                    types.String `tfsdk:"policy_type"`                       // 策略类型
	OperateUnit                   types.String `tfsdk:"operate_unit"`                      // 操作单位
	OperateCount                  types.Int32  `tfsdk:"operate_count"`                     // 调整值
	Action                        types.String `tfsdk:"action"`                            // 执行动作
	Cycle                         types.String `tfsdk:"cycle"`                             // 循环方式
	Day                           types.Set    `tfsdk:"day"`                               // 执行日期
	EffectiveFrom                 types.String `tfsdk:"effective_from"`                    // 生效开始时间
	EffectiveTill                 types.String `tfsdk:"effective_till"`                    // 生效截止时间
	ExecutionTime                 types.String `tfsdk:"execution_time"`                    // 执行时间
	Cooldown                      types.Int32  `tfsdk:"cooldown"`                          // 冷却/预热时间
	TriggerName                   types.String `tfsdk:"trigger_name"`                      // 告警策略-告警规则名称
	TriggerMetricName             types.String `tfsdk:"trigger_metric_name"`               // 告警策略-监控指标名称
	TriggerStatistics             types.String `tfsdk:"trigger_statistics"`                // 告警策略-聚合方法
	TriggerComparisonOperator     types.String `tfsdk:"trigger_comparison_operator"`       // 告警策略-比较符
	TriggerThreshold              types.Int32  `tfsdk:"trigger_threshold"`                 // 告警策略-阈值
	TriggerPeriod                 types.String `tfsdk:"trigger_period"`                    // 告警策略-监控周期
	TriggerEvaluationCount        types.Int32  `tfsdk:"trigger_evaluation_count"`          // 告警策略-连续出现次数
	TargetMetricName              types.String `tfsdk:"target_metric_name"`                // 目标追踪策略-监控指标名称
	TargetValue                   types.Int32  `tfsdk:"target_value"`                      // 目标追踪策略-追踪目标值
	TargetScaleOutEvaluationCount types.Int32  `tfsdk:"target_scale_out_evaluation_count"` // 目标追踪策略-扩容连续告警次数
	TargetScaleInEvaluationCount  types.Int32  `tfsdk:"target_scale_in_evaluation_count"`  // 目标追踪策略-缩容连续告警次数
	TargetOperateRange            types.Int32  `tfsdk:"target_operate_range"`              // 目标追踪策略-缩容波动范围
	TargetDisableScaleIn          types.Bool   `tfsdk:"target_disable_scale_in"`           // 目标追踪策略-是否禁用缩容
	ID                            types.Int64  `tfsdk:"id"`                                // 弹性伸缩策略id
	Status                        types.String `tfsdk:"status"`                            // 策略状态
	IsExecute                     types.Bool   `tfsdk:"is_execute"`                        // 是否执行当前的伸缩策略
}
