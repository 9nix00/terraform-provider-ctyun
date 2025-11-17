package ec

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &CtyunEcCdaInstance{}
	_ resource.ResourceWithConfigure   = &CtyunEcCdaInstance{}
	_ resource.ResourceWithImportState = &CtyunEcCdaInstance{}
)

func NewCtyunEcCdaInstance() resource.Resource {
	return &CtyunEcCdaInstance{}
}

type CtyunEcCdaInstance struct {
	meta *common.CtyunMetadata
}

type CtyunEcCdaInstanceConfig struct {
	ID            types.String   `tfsdk:"id"`
	EcID          types.String   `tfsdk:"ec_id"`
	CgwID         types.String   `tfsdk:"cgw_id"`
	CdaID         types.String   `tfsdk:"cda_id"`
	CdaName       types.String   `tfsdk:"cda_name"`
	CdaCidrV4List []types.String `tfsdk:"cda_cidr_v4_list"`
	CdaCidrV6List []types.String `tfsdk:"cda_cidr_v6_list"`
	RtbID         types.String   `tfsdk:"rtb_id"`
	CdaInfo       types.String   `tfsdk:"cda_info"`
	Account       types.String   `tfsdk:"account"`
	Weights       types.Int64    `tfsdk:"weights"`
	RouteLearn    types.Int64    `tfsdk:"route_learn"`
	RouteSync     types.Int64    `tfsdk:"route_sync"`
	RegionId      types.String   `tfsdk:"region_id"`
}

func (c *CtyunEcCdaInstance) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ec_cda_instance"
}

func (c *CtyunEcCdaInstance) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `**CDA网络实例资源**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "网络实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cgw_id": schema.StringAttribute{
				Required:    true,
				Description: "云网关ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cda_id": schema.StringAttribute{
				Required:    true,
				Description: "云专线ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cda_name": schema.StringAttribute{
				Required:    true,
				Description: "云间高速侧显示的云专线名称（建议保持和cda创建时的名称一致）",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cda_cidr_v4_list": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "已选择的V4子网列表",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"cda_cidr_v6_list": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "已选择的V6子网列表",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"rtb_id": schema.StringAttribute{
				Required:    true,
				Description: "路由表ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cda_info": schema.StringAttribute{
				Required:    true,
				Description: "云专线信息，json格式的字符串",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"account": schema.StringAttribute{
				Required:    true,
				Description: "天翼云客户邮箱",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"weights": schema.Int64Attribute{
				Optional:    true,
				Description: "权重，专线默认50，无冗余实例则不传",
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
			},
			"route_learn": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "路由学习开关，开启后云网关自动学习网络实例路由，取值范围: 1:学习 0:不学习，默认学习",
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
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
			"route_sync": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "路由同步开关，开启后云网关路由自动同步到网络实例，取值范围: 1:同步 0:不同步，默认同步",
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *CtyunEcCdaInstance) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunEcCdaInstance) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunEcCdaInstanceConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}

	// 设置资源ID
	plan.ID = types.StringValue(fmt.Sprintf("%s-%s", plan.EcID.ValueString(), plan.CdaID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunEcCdaInstance) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunEcCdaInstanceConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunEcCdaInstance) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunEcCdaInstanceConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &plan)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunEcCdaInstance) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunEcCdaInstanceConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *CtyunEcCdaInstance) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// 暂不支持导入
	resp.Diagnostics.AddError(
		"Import not supported",
		"This resource does not support import.",
	)
}

func (c *CtyunEcCdaInstance) create(ctx context.Context, plan *CtyunEcCdaInstanceConfig) (err error) {
	// 转换 V4 CIDR 列表
	cidrV4List := make([]string, len(plan.CdaCidrV4List))
	for i, cidr := range plan.CdaCidrV4List {
		cidrV4List[i] = cidr.ValueString()
	}

	// 转换 V6 CIDR 列表
	var cidrV6List []*string
	if len(plan.CdaCidrV6List) > 0 {
		cidrV6List = make([]*string, len(plan.CdaCidrV6List))
		for i, cidr := range plan.CdaCidrV6List {
			cidrStr := cidr.ValueString()
			cidrV6List[i] = &cidrStr
		}
	}

	req := &ec.EcEcAddCDANetworkRequest{
		EcID:          plan.EcID.ValueString(),
		CgwID:         plan.CgwID.ValueString(),
		CdaID:         plan.CdaID.ValueString(),
		CdaName:       plan.CdaName.ValueString(),
		CdaCidrV4List: cidrV4List,
		CdaCidrV6List: cidrV6List,
		RtbID:         plan.RtbID.ValueString(),
		CdaInfo:       plan.CdaInfo.ValueString(),
		Account:       plan.Account.ValueStringPointer(),
		DcID:          plan.RegionId.ValueStringPointer(),
	}

	// 添加 account 参数

	if !plan.Weights.IsNull() {
		weights := int32(plan.Weights.ValueInt64())
		req.Weights = &weights
	}

	if !plan.RouteLearn.IsNull() {
		routeLearn := int32(plan.RouteLearn.ValueInt64())
		req.RouteLearn = &routeLearn
	} else {
		// 默认值为1（学习）
		defaultRouteLearn := int32(1)
		req.RouteLearn = &defaultRouteLearn
	}

	if !plan.RouteSync.IsNull() {
		routeSync := int32(plan.RouteSync.ValueInt64())
		req.RouteSync = &routeSync
	} else {
		// 默认值为1（同步）
		defaultRouteSync := int32(1)
		req.RouteSync = &defaultRouteSync
	}

	tflog.Info(ctx, "创建CDA网络实例", map[string]interface{}{
		"ec_id":  plan.EcID.ValueString(),
		"cgw_id": plan.CgwID.ValueString(),
		"cda_id": plan.CdaID.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcAddCDANetworkApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp == nil {
		return fmt.Errorf("API return error. StatusCode is nil")
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return nil
}

func (c *CtyunEcCdaInstance) getAndMerge(ctx context.Context, state *CtyunEcCdaInstanceConfig) (err error) {
	// 由于列表接口返回的是字符串而不是结构化数据，我们无法解析具体的实例信息
	// 在实际应用中，可能需要解析返回的JSON字符串
	return nil
}

func (c *CtyunEcCdaInstance) update(ctx context.Context, plan *CtyunEcCdaInstanceConfig) (err error) {
	// 构造更新请求
	// 根据API文档，更新操作只支持更新子网信息，通过cdaInfo字段传递
	cdaInfo := make(map[string]interface{})

	// 转换 V4 CIDR 列表
	cidrV4List := make([]string, len(plan.CdaCidrV4List))
	for i, cidr := range plan.CdaCidrV4List {
		cidrV4List[i] = cidr.ValueString()
	}
	cdaInfo["cdaCidrV4List"] = cidrV4List

	// 转换 V6 CIDR 列表
	if len(plan.CdaCidrV6List) > 0 {
		cidrV6List := make([]string, len(plan.CdaCidrV6List))
		for i, cidr := range plan.CdaCidrV6List {
			cidrV6List[i] = cidr.ValueString()
		}
		cdaInfo["cdaCidrV6List"] = cidrV6List
	}

	// 将cdaInfo转换为JSON字符串
	cdaInfoJSON, err := json.Marshal(cdaInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal cdaInfo to JSON: %v", err)
	}

	req := &ec.EcEcUpdateCDANetworkRequest{
		InstanceID: plan.ID.ValueString(),
		CdaInfo:    string(cdaInfoJSON),
	}

	tflog.Info(ctx, "更新CDA网络实例", map[string]interface{}{
		"instance_id": plan.ID.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcUpdateCDANetworkApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp == nil {
		return fmt.Errorf("API return error. StatusCode is nil")
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return
	}
	return nil
}

func (c *CtyunEcCdaInstance) delete(ctx context.Context, state CtyunEcCdaInstanceConfig) (err error) {
	req := &ec.EcEcDeleteCDANetworkRequest{
		InstanceID: state.ID.ValueString(),
	}

	// 如果有CdaID信息，也一并传入
	if !state.CdaID.IsNull() {
		cdaID := state.CdaID.ValueString()
		req.CdaID = &cdaID
	}

	tflog.Info(ctx, "删除CDA网络实例", map[string]interface{}{
		"instance_id": state.ID.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcDeleteCDANetworkApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp == nil {
		return fmt.Errorf("API return error. StatusCode is nil")
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return nil
}
