package ebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctebs"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"time"
)

type ctyunEbs struct {
	meta       *common.CtyunMetadata
	ebsService *business.EbsService
}

func NewCtyunEbs() resource.Resource {
	return &ctyunEbs{}
}

func (c *ctyunEbs) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs"
}

func (c *ctyunEbs) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027696**`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "磁盘命名，单账户单资源池下，命名需唯一，长度为2-63个字符，只能由数字、字母、-组成，不能以数字、-开头，且不能以-结尾，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 63),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][0-9a-zA-Z_-]+$"), "磁盘名称不符合规则"),
				},
			},
			"mode": schema.StringAttribute{
				Required:    true,
				Description: "磁盘模式，vbd，iscsi，fcsan",
				Validators: []validator.String{
					stringvalidator.OneOf(business.EbsDiskModes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "磁盘类型，sata：普通IO，sas：高IO，ssd：超高IO，ssd-genric：通用型SSD，fast-ssd：极速型SSD",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.EbsDiskTypes...),
				},
			},
			"size": schema.Int64Attribute{
				Required:    true,
				Description: "磁盘大小，单位GB，取值范围[10, 32768]，支持更新（不支持缩容）",
				Validators: []validator.Int64{
					int64validator.Between(10, 32768),
				},
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month：按月，year：按年、on_demand：按需。当此值为month或者year时，cycle_count为必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypes...),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "订购时长，该参数在cycle_type为month或year时才生效，当cycle_type=month，支持订购1-11个月；当cycle_type=year，支持订购1-5年",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
						types.StringValue(business.OrderCycleTypeYear),
					),
					validator2.ConflictsWithEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					validator2.CycleCount(1, 11, 1, 5),
				},
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订购的受理单id",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "磁盘id",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "云硬盘使用状态，deleting：删除中，creating：资源创建中，detaching：解绑中，detached：未绑定云主机，attaching：绑定中，attached：已绑定，extending：扩容中，error：错误状态，backup：备份中，backupRestoring：从备份恢复中，expired：包周期已结束，freezing：按需计费，处于冻结状态，可能账户受限或余额不足，available：可用，in-use：已挂载云主机，resizing：扩容中",
			},
			"expire_time": schema.StringAttribute{
				Computed:    true,
				Description: "到期时间",
			},
			"multi_attach": schema.BoolAttribute{
				Computed:    true,
				Description: "是否共享云硬盘",
			},
			"encrypted": schema.BoolAttribute{
				Computed:    true,
				Description: "是否加密盘",
			},
			"kms_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "加密盘密钥UUID，是加密盘时才返回",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区id，如果不填则默认使用provider ctyun中的az_name或环境变量中的CTYUN_AZ_NAME",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraAzName, false),
			},
		},
	}
}

func (c *ctyunEbs) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunEbsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	azName := plan.AzName.ValueString()
	onDemand := plan.CycleType.ValueString() == business.OrderCycleTypeOnDemand

	diskMode, err := business.EbsDiskModeMap.FromOriginalScene(plan.Mode.ValueString(), business.EbsDiskModeMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	diskType, err := business.EbsDiskTypeMap.FromOriginalScene(plan.Type.ValueString(), business.EbsDiskTypeMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	resp, err2 := c.meta.Apis.CtEbsApis.EbsCreateApi.Do(ctx, c.meta.Credential, &ctebs.EbsCreateRequest{
		RegionId:    regionId,
		AzName:      azName,
		ProjectId:   projectId,
		ClientToken: uuid.NewString(),
		DiskName:    plan.Name.ValueString(),
		DiskMode:    diskMode.(string),
		DiskType:    diskType.(string),
		DiskSize:    plan.Size.ValueInt64(),
		OnDemand:    onDemand,
		CycleType:   plan.CycleType.ValueString(),
		CycleCount:  plan.CycleCount.ValueInt64(),
	})

	var id, masterOrderId string
	if err2 == nil {
		id = resp.Resources[0].DiskId
		masterOrderId = resp.MasterOrderId
	} else {
		// 判断返回信息是否需要轮询
		if err2.ErrorCode() != common.EbsOrderInProgress {
			response.Diagnostics.AddError(err2.Error(), err2.Error())
			return
		}
		// 获取主订单
		moi, err := c.getMasterOrderIdIfOrderInProgress(err2)
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		response.Diagnostics.Append(response.State.Set(ctx, plan)...)
		// 轮询结果
		helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
		loop, err := helper.OrderLoop(ctx, c.meta.Credential, moi)
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		id = loop.Uuid[0]
		masterOrderId = moi
	}

	plan.Id = types.StringValue(id)
	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	plan.AzName = types.StringValue(azName)
	plan.MasterOrderId = types.StringValue(masterOrderId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, ctyunRequestError := c.getAndMergeEbs(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEbs) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunEbsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	if !c.acquireAndSetIdIfOrderNotFinished(ctx, &state, response) {
		return
	}
	instance, err := c.getAndMergeEbs(ctx, state)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEbs) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan CtyunEbsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var state CtyunEbsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 判断名字是否相同
	if !plan.Name.Equal(state.Name) {
		_, err := c.meta.Apis.CtEbsApis.EbsChangeNameApi.Do(ctx, c.meta.Credential, &ctebs.EbsChangeNameRequest{
			RegionId: state.RegionId.ValueString(),
			DiskId:   state.Id.ValueString(),
			DiskName: plan.Name.ValueString(),
		})
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
	}
	// 判断硬盘大小是否相同，不同要走修改ebs接口
	err := c.ebsService.UpdateSize(ctx, state.Id.ValueString(), state.RegionId.ValueString(), int(state.Size.ValueInt64()), int(plan.Size.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	instance, ctyunRequestError := c.getAndMergeEbs(ctx, state)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEbs) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunEbsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	resp, err := c.meta.Apis.CtEbsApis.EbsDeleteApi.Do(ctx, c.meta.Credential, &ctebs.EbsDeleteRequest{
		RegionId:    state.RegionId.ValueString(),
		DiskId:      state.Id.ValueString(),
		ClientToken: uuid.NewString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	err2 := helper.RefundLoop(ctx, c.meta.Credential, resp.MasterOrderId)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}
}

func (c *ctyunEbs) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ebsService = business.NewEbsService(meta)
}

// getAndMergeEbs 查询ebs
func (c *ctyunEbs) getAndMergeEbs(ctx context.Context, cfg CtyunEbsConfig) (*CtyunEbsConfig, error) {
	resp, err := c.meta.Apis.CtEbsApis.EbsShowApi.Do(ctx, c.meta.Credential, &ctebs.EbsShowRequest{
		RegionId: cfg.RegionId.ValueString(),
		DiskId:   cfg.Id.ValueString(),
	})
	if err != nil {
		if err.ErrorCode() == common.EbsEbsInfoDataDamaged {
			return nil, nil
		}
		return nil, err
	}

	diskMode, err2 := business.EbsDiskModeMap.ToOriginalScene(resp.DiskMode, business.EbsDiskModeMapScene1)
	if err2 != nil {
		return nil, err2
	}
	diskType, err2 := business.EbsDiskTypeMap.ToOriginalScene(resp.DiskType, business.EbsDiskTypeMapScene1)
	if err2 != nil {
		return nil, err2
	}

	cfg.Name = types.StringValue(resp.DiskName)
	cfg.Id = types.StringValue(resp.DiskID)
	cfg.Size = types.Int64Value(resp.DiskSize)
	cfg.Type = types.StringValue(diskType.(string))
	cfg.Mode = types.StringValue(diskMode.(string))
	cfg.Status = types.StringValue(resp.DiskStatus)
	cfg.ExpireTime = types.StringValue(time.UnixMilli(resp.ExpireTime).Format(time.DateTime))
	cfg.MultiAttach = types.BoolValue(resp.MultiAttach)
	cfg.Encrypted = types.BoolValue(resp.IsEncrypt)
	cfg.KmsUuid = types.StringValue(resp.KmsUUID)
	if resp.OnDemand {
		cfg.CycleType = types.StringValue(business.OrderCycleTypeOnDemand)
	} else {
		cfg.CycleType = types.StringValue(*resp.CycleType)
		cfg.CycleCount = types.Int64Value(*resp.CycleCount)
	}
	return &cfg, nil
}

// getMasterOrderIdIfOrderInProgress 获取masterOrderId
func (c *ctyunEbs) getMasterOrderIdIfOrderInProgress(err ctyunsdk.CtyunRequestError) (string, error) {
	resp := struct {
		MasterOrderId string `json:"masterOrderID"`
		MasterOrderNo string `json:"masterOrderNO"`
	}{}
	if err.CtyunResponse() == nil {
		return "", err
	}
	_, err = err.CtyunResponse().ParseByStandardModel(&resp)
	if err != nil {
		return "", err
	}
	return resp.MasterOrderId, err
}

// acquireIdIfOrderNotFinished 重新获取id，如果前订单状态有问题需要重新轮询
// 返回值：数据是否有效
func (c *ctyunEbs) acquireAndSetIdIfOrderNotFinished(ctx context.Context, state *CtyunEbsConfig, response *resource.ReadResponse) bool {
	id := state.Id.ValueString()
	masterOrderId := state.MasterOrderId.ValueString()
	if id != "" {
		// 数据是完整的，无需处理
		return true
	}
	if state.MasterOrderId.ValueString() == "" {
		// 没有受理的订购单id，数据是不可恢复的，直接把当前状态移除并且返回
		response.State.RemoveResource(ctx)
		return false
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	resp, err := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
	if err != nil || len(resp.Uuid) == 0 {
		// 报错了，或者受理没有返回数据的情况，那么意思是这个单子并没有开通出来，此时数据无法恢复
		response.State.RemoveResource(ctx)
		return false
	}

	// 成功把id恢复出来
	state.Id = types.StringValue(resp.Uuid[0])
	response.State.Set(ctx, state)
	return true
}

type CtyunEbsConfig struct {
	Name          types.String `tfsdk:"name"`
	Mode          types.String `tfsdk:"mode"`
	Type          types.String `tfsdk:"type"`
	Size          types.Int64  `tfsdk:"size"`
	CycleType     types.String `tfsdk:"cycle_type"`
	CycleCount    types.Int64  `tfsdk:"cycle_count"`
	MasterOrderId types.String `tfsdk:"master_order_id"`
	Id            types.String `tfsdk:"id"`           // 磁盘ID
	Status        types.String `tfsdk:"status"`       // 云硬盘使用状态 deleting/creating/detaching，具体请参考云硬盘使用状态
	ExpireTime    types.String `tfsdk:"expire_time"`  // 过期时刻，epoch时戳，精度毫秒
	MultiAttach   types.Bool   `tfsdk:"multi_attach"` // 是否共享云硬盘
	Encrypted     types.Bool   `tfsdk:"encrypted"`    // 是否加密盘
	KmsUuid       types.String `tfsdk:"kms_uuid"`     // 加密盘密钥UUID，是加密盘时才返回
	ProjectId     types.String `tfsdk:"project_id"`
	RegionId      types.String `tfsdk:"region_id"`
	AzName        types.String `tfsdk:"az_name"`
}
