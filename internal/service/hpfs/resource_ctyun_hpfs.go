package hpfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/hpfs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

type ctyunHpfs struct {
	meta *common.CtyunMetadata
}

func (c *ctyunHpfs) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_hpfs"
}

func (c *ctyunHpfs) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func NewCtyunMongodbInstance() resource.Resource {
	return &ctyunHpfs{}
}

func (c *ctyunHpfs) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
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

			"sfs_type": schema.StringAttribute{
				Required:    true,
				Description: "并行文件类型 (hpc/hpc_cache)",
				Validators: []validator.String{
					stringvalidator.OneOf("hpc", "hpc_cache"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sfs_protocol": schema.StringAttribute{
				Required:    true,
				Description: "协议类型",
				Validators: []validator.String{
					stringvalidator.OneOf("NFS", "CIFS"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_type": schema.StringAttribute{
				Optional:    true,
				Description: "包周期类型",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "包周期数",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"sfs_name": schema.StringAttribute{
				Required:    true,
				Description: "并行文件名",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"sfs_size": schema.Int64Attribute{
				Required:    true,
				Description: "文件大小（GB），范围: 500-32768",
				Validators: []validator.Int64{
					int64validator.Between(500, 32768),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"az_name": schema.StringAttribute{
				Required:    true,
				Description: "可用区名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_name": schema.StringAttribute{
				Optional:    true,
				Description: "集群名称",
			},
			"baseline": schema.StringAttribute{
				Optional:    true,
				Description: "性能基线",
				Validators: []validator.String{
					stringvalidator.OneOf("low", "medium", "high"),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "虚拟网 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订单ID",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sfs_status": schema.StringAttribute{
				Computed:    true,
				Description: "并行文件状态",
			},
			"used_size": schema.Int64Attribute{
				Computed:    true,
				Description: "已用大小（MB）",
			},
			"dataflow_list": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "HPFS文件系统下的数据流动策略ID列表",
			},
		},
	}
}

func (c *ctyunHpfs) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunHpfsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	//创建前检查,检查证书有效性
	isValid, err := c.checkBeforeHpfs(ctx, plan)
	if !isValid || err != nil {
		return
	}
	err = c.createHpfs(ctx, &plan)
	if err != nil {
		return
	}
	// 创建后反查创建后的证书信息
	err = c.getAndMergeHpfs(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunHpfs) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunHpfsConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeHpfs(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "NotExists") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunHpfs) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取 plan -tf文件中配置
	var plan CtyunHpfsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunHpfsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
	}

	// 更新基本信息
	err = c.updateHfps(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeHpfs(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunHpfs) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunHpfsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &hpfs.HpfsRenameSfsRequest{
		RegionID: state.RegionID.ValueString(),
		SfsUID:   state.ID.ValueString(),
		SfsName:  state.SfsName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkHpfsApis.HpfsRenameSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = errors.New("hpfs退订失败，返回值为nil")
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	// 异步接口，需要轮询查看是否退订成功
	err = c.deleteLoop(ctx, &state, 600)
	if err != nil {
		return
	}
}

func (c *ctyunHpfs) checkBeforeHpfs(ctx context.Context, plan CtyunHpfsConfig) (inValid bool, err error) {
	// 判断sfs_type，sfs_protocol是否合理，/v4/hpfs/list-cluster
	return true, nil
}

func (c *ctyunHpfs) createHpfs(ctx context.Context, config *CtyunHpfsConfig) error {
	params := &hpfs.HpfsNewSfsRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
		SfsType:     config.SfsType.ValueString(),
		SfsProtocol: config.SfsProtocol.ValueString(),
		CycleType:   config.CycleType.ValueString(),
		SfsName:     config.SfsName.ValueString(),
		SfsSize:     config.SfsSize.ValueInt32(),
		Vpc:         config.VpcID.ValueString(),
		Subnet:      config.SubnetID.ValueString(),
	}
	if config.CycleType.ValueString() == business.HpfsCycleTypeOnDemand {
		onDemand := true
		params.OnDemand = &onDemand
	} else {
		params.CycleCount = config.CycleCount.ValueInt32()
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		params.ProjectID = config.ProjectID.ValueString()
	}
	if !config.AzName.IsNull() && !config.AzName.IsUnknown() {
		params.AzName = config.AzName.ValueString()
	}
	if !config.ClusterName.IsNull() && !config.ClusterName.IsUnknown() {
		params.ClusterName = config.ClusterName.ValueString()
	}
	if !config.Baseline.IsNull() && !config.Baseline.IsUnknown() {
		params.Baseline = config.Baseline.ValueString()
	}
	resp, err := c.meta.Apis.SdkHpfsApis.HpfsNewSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = errors.New("开通hpfs失败，返回nil")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	config.MasterOrderID = types.StringValue(resp.ReturnObj.MasterOrderID)
	config.ID = types.StringValue(resp.ReturnObj.Resources[0].SfsUID)
	return nil
}

func (c *ctyunHpfs) getAndMergeHpfs(ctx context.Context, config *CtyunHpfsConfig) error {
	// 获取hpfs详情
	hpfsResp, err := c.getHpfsDetail(ctx, config)
	if err != nil {
		return err
	}
	hpfsDetail := hpfsResp.ReturnObj
	config.SfsName = types.StringValue(hpfsDetail.SfsName)
	config.SfsSize = types.Int32Value(hpfsDetail.SfsSize)
	config.SfsType = types.StringValue(hpfsDetail.SfsType)
	config.SfsStatus = types.StringValue(hpfsDetail.SfsStatus)
	config.ClusterName = types.StringValue(hpfsDetail.ClusterName)
	dataFlowList, diags := types.SetValueFrom(ctx, types.StringType, hpfsDetail.DataflowList)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return err
	}
	config.DataflowList = dataFlowList
	return nil
}

func (c *ctyunHpfs) getHpfsDetail(ctx context.Context, config *CtyunHpfsConfig) (*hpfs.HpfsInfoSfsResponse, error) {
	params := &hpfs.HpfsInfoSfsRequest{
		SfsUID:   config.ID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkHpfsApis.HpfsInfoSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = errors.New("获取hpfs详情失败，返回为nil")
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

func (c *ctyunHpfs) updateHfps(ctx context.Context, state *CtyunHpfsConfig, plan *CtyunHpfsConfig) error {
	// 并行文件重命名
	err := c.hfpsRename(ctx, state, plan)
	if err != nil {
		return err
	}
	// 并行文件修改规格
	err = c.updateHpfsSize(ctx, state, plan)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunHpfs) hfpsRename(ctx context.Context, state *CtyunHpfsConfig, plan *CtyunHpfsConfig) error {
	if plan.SfsName.IsNull() || state.SfsName == plan.SfsName {
		return nil
	}
	params := &hpfs.HpfsRenameSfsRequest{
		RegionID: state.RegionID.ValueString(),
		SfsUID:   state.ID.ValueString(),
		SfsName:  plan.SfsName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkHpfsApis.HpfsRenameSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = errors.New("hpfs 更名失败，返回nil")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

func (c *ctyunHpfs) updateHpfsSize(ctx context.Context, state *CtyunHpfsConfig, plan *CtyunHpfsConfig) error {
	// 判断是否需要进行修改
	if plan.SfsSize.IsNull() || state.SfsSize == plan.SfsSize {
		return nil
	}
	// 配置修改参数
	params := &hpfs.HpfsResizeSfsRequest{
		SfsSize:     plan.SfsSize.ValueInt32(),
		SfsUID:      state.ID.ValueString(),
		RegionID:    state.RegionID.ValueString(),
		ClientToken: uuid.NewString(),
	}
	resp, err := c.meta.Apis.SdkHpfsApis.HpfsResizeSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = errors.New("hpfs sfs_size修改失败，返回值为Nil。")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	return nil
}

func (c *ctyunHpfs) deleteLoop(ctx context.Context, config *CtyunHpfsConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*5, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, _ := c.getHpfsDetail(ctx, config)
			if resp.StatusCode == common.ErrorStatusCode && strings.Contains(resp.Message, "不存在") {
				return false
			}
			return true

		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未退订成功！")
	}

	return
}

type CtyunHpfsConfig struct {
	RegionID      types.String `tfsdk:"region_id"`       // 资源池 ID
	ClientToken   types.String `tfsdk:"client_token"`    // 客户端存根，用于保证订单幂等性
	ProjectID     types.String `tfsdk:"project_id"`      // 资源所属企业项目 ID
	SfsType       types.String `tfsdk:"sfs_type"`        // 并行文件类型
	SfsProtocol   types.String `tfsdk:"sfs_protocol"`    // 协议类型
	CycleType     types.String `tfsdk:"cycle_type"`      // 包周期类型
	CycleCount    types.Int32  `tfsdk:"cycle_count"`     // 包周期数
	SfsName       types.String `tfsdk:"sfs_name"`        // 并行文件名
	SfsSize       types.Int32  `tfsdk:"sfs_size"`        // 文件大小（GB）
	AzName        types.String `tfsdk:"az_name"`         // 可用区名称
	ClusterName   types.String `tfsdk:"cluster_name"`    // 集群名称
	Baseline      types.String `tfsdk:"baseline"`        // 性能基线
	VpcID         types.String `tfsdk:"vpc_id"`          // 虚拟网 ID
	SubnetID      types.String `tfsdk:"subnet_id"`       // 子网 ID
	MasterOrderID types.String `tfsdk:"master_order_id"` // 订单id
	ID            types.String `tfsdk:"id"`              // 资源 ID
	SfsStatus     types.String `tfsdk:"sfs_status"`      // 并行文件状态
	UsedSize      types.Int32  `tfsdk:"used_size"`       // 已用大小（MB）
	DataflowList  types.Set    `tfsdk:"dataflow_list"`   // HPFS文件系统下的数据流动策略ID列表
}
