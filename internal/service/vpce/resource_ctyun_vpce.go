package vpce

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunVpce{}
	_ resource.ResourceWithConfigure   = &ctyunVpce{}
	_ resource.ResourceWithImportState = &ctyunVpce{}
)

type ctyunVpce struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpce() resource.Resource {
	return &ctyunVpce{}
}

func (c *ctyunVpce) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpce"
}

type CtyunVpceConfig struct {
	ID                types.String `tfsdk:"id"`
	MasterOrderID     types.String `tfsdk:"master_order_id"`
	EndpointServiceID types.String `tfsdk:"endpoint_service_id"`
	RegionID          types.String `tfsdk:"region_id"`
	VpcID             types.String `tfsdk:"vpc_id"`
	Name              types.String `tfsdk:"name"`
	SubnetID          types.String `tfsdk:"subnet_id"`
	SubnetIP          types.String `tfsdk:"subnet_ip"`
	WhitelistFlag     types.Bool   `tfsdk:"whitelist_flag"`
	WhitelistCidr     types.Set    `tfsdk:"whitelist_cidr"`
	Status            types.Int32  `tfsdk:"status"`
	CreateTime        types.String `tfsdk:"create_time"`
	UpdateTime        types.String `tfsdk:"update_time"`
}

func (c *ctyunVpce) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10042658/10217121`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "主订单号",
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
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"endpoint_service_id": schema.StringAttribute{
				Required:    true,
				Description: "终端节点服务ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "虚拟私有云ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"subnet_ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "子网IP",
				Validators: []validator.String{
					validator2.Ip(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "支持拉丁字母、中文、数字，下划线，连字符，中文/英文字母开头，不能以http:/https:开头，长度2-32，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 32),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\\x{4e00}-\\x{9fa5}][0-9a-zA-Z_\\x{4e00}-\\x{9fa5}fa5}-]+$"), "名称不符合规则"),
				},
			},
			"status": schema.Int32Attribute{
				Computed:    true,
				Description: "endpoint状态, 1 表示已链接，2 表示未链接",
			},
			"whitelist_flag": schema.BoolAttribute{
				Required:    true,
				Description: "是否开启白名单，支持更新",
			},
			"whitelist_cidr": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				Description: "白名单列表，当whitelist_flag=true是必填，最多同时支持20个地址，最少输入一个，支持更新",
				Validators: []validator.Set{
					validator2.AlsoRequiresEqualSet(
						path.MatchRoot("whitelist_flag"),
						types.BoolValue(true),
					),
					validator2.ConflictsWithEqualSet(
						path.MatchRoot("whitelist_flag"),
						types.BoolValue(false),
					),
					setvalidator.SizeAtLeast(1),
					setvalidator.SizeAtMost(20),
					setvalidator.ValueStringsAre(validator2.Cidr()),
				},
			},
			"create_time": schema.StringAttribute{
				Description: "创建时间，为UTC格式",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Description: "更新时间，为UTC格式",
				Computed:    true,
			},
		},
	}
}

func (c *ctyunVpce) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunVpceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建
	masterOrderID, endpointID, err := c.loopCreate(ctx, plan)
	if err != nil {
		return
	}
	plan.MasterOrderID = types.StringValue(masterOrderID)
	plan.ID = types.StringValue(endpointID)
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunVpce) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "resource not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunVpce) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunVpceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunVpceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunVpce) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.loopDelete(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunVpce) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [vpceID],[regionID]
func (c *ctyunVpce) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunVpceConfig
	var endpointID, regionID string
	err = terraform_extend.Split(request.ID, &endpointID, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.ID = types.StringValue(endpointID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// loopCreate 循环执行create
func (c *ctyunVpce) loopCreate(ctx context.Context, plan CtyunVpceConfig) (masterOrderID, id string, err error) {
	clientToken := uuid.NewString()
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			masterOrderID, id, err = c.create(ctx, clientToken, plan)
			if err != nil {
				return false
			}
			if id != "" {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = errors.New("创建时未获取到终端节点ID")
	}
	return
}

// create 创建
func (c *ctyunVpce) create(ctx context.Context, clientToken string, plan CtyunVpceConfig) (masterOrderID, endpointID string, err error) {
	params := &ctvpc.CtvpcCreateEndpointRequest{
		ClientToken:       clientToken,
		RegionID:          plan.RegionID.ValueString(),
		VpcID:             plan.VpcID.ValueString(),
		SubnetID:          plan.SubnetID.ValueString(),
		EndpointName:      plan.Name.ValueString(),
		EndpointServiceID: plan.EndpointServiceID.ValueString(),
		CycleType:         "on_demand",
	}
	WhitelistFlag := plan.WhitelistFlag.ValueBool()
	if WhitelistFlag {
		params.WhitelistFlag = 1
		plan.WhitelistCidr.ElementsAs(ctx, &params.Whitelist, true)
	} else {
		params.WhitelistFlag = 0
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateEndpointApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	masterOrderID = *resp.ReturnObj.MasterOrderID
	endpointID = *resp.ReturnObj.EndpointID
	return
}

// getAndMerge 从远端查询
func (c *ctyunVpce) getAndMerge(ctx context.Context, plan *CtyunVpceConfig) (err error) {
	endpointID, regionID := plan.ID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcShowEndpointRequest{
		RegionID:   regionID,
		EndpointID: endpointID,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowEndpointApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	endpoint := resp.ReturnObj.Endpoint

	plan.VpcID = utils.SecStringValue(endpoint.VpcID)
	plan.Name = utils.SecStringValue(endpoint.Name)
	plan.EndpointServiceID = utils.SecStringValue(endpoint.EndpointServiceID)
	plan.Status = types.Int32Value(endpoint.Status)
	plan.CreateTime = types.StringValue(utils.SecString(endpoint.CreatedTime))
	plan.UpdateTime = types.StringValue(utils.SecString(endpoint.UpdatedTime))

	Whitelist := utils.SecString(endpoint.Whitelist)
	if len(Whitelist) > 0 {
		t := strings.Split(Whitelist, ",")
		plan.WhitelistCidr, _ = types.SetValueFrom(ctx, types.StringType, t)
		plan.WhitelistFlag = types.BoolValue(true)
	} else {
		plan.WhitelistCidr, _ = types.SetValueFrom(ctx, types.StringType, []string{})
		plan.WhitelistFlag = types.BoolValue(false)
	}

	if endpoint.EndpointObj != nil {
		plan.SubnetID = utils.SecStringValue(endpoint.EndpointObj.SubnetID)
		plan.SubnetIP = utils.SecStringValue(endpoint.EndpointObj.Ip)
	}

	return
}

// update 更新
func (c *ctyunVpce) update(ctx context.Context, plan, state CtyunVpceConfig) (err error) {
	//if !plan.SubnetIP.IsUnknown() && !plan.SubnetIP.Equal(state.SubnetIP) {
	//	err = fmt.Errorf("子网ip地址不支持修改")
	//	return
	//}

	endpointID, regionID := state.ID.ValueString(), state.RegionID.ValueString()
	params := &ctvpc.CtvpcUpdateEndpointRequest{
		ClientToken: uuid.NewString(),
		RegionID:    regionID,
		EndpointID:  endpointID,
	}
	if !plan.Name.Equal(state.Name) {
		params.EndpointName = plan.Name.ValueStringPointer()
	}

	flag := plan.WhitelistFlag.ValueBool()
	// 白名单开关状态变化
	if !plan.WhitelistFlag.Equal(state.WhitelistFlag) {
		params.EnableWhitelist = &flag
		if flag {
			plan.WhitelistCidr.ElementsAs(ctx, &params.Whitelist, true)
		} else {
			params.Whitelist = nil
		}
	} else if !plan.WhitelistCidr.Equal(state.WhitelistCidr) { // 开关状态没变，白名单变了
		plan.WhitelistCidr.ElementsAs(ctx, &params.Whitelist, true)
		params.EnableWhitelist = &flag
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdateEndpointApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}

// delete 删除
func (c *ctyunVpce) delete(ctx context.Context, clientToken string, plan CtyunVpceConfig) (status string, err error) {
	endpointID, regionID := plan.ID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcDeleteEndpointRequest{
		RegionID:    regionID,
		EndpointID:  endpointID,
		ClientToken: clientToken,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteEndpointApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	status = utils.SecString(resp.ReturnObj.MasterResourceStatus)
	return
}

// loopDelete 循环执行delete
func (c *ctyunVpce) loopDelete(ctx context.Context, plan CtyunVpceConfig) (err error) {
	clientToken := uuid.NewString()
	var executeSuccessFlag bool
	var status string
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			status, err = c.delete(ctx, clientToken, plan)
			if err != nil {
				return false
			}
			if status == "refunded" {
				time.Sleep(30 * time.Second)
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("删除终端节点 %s 失败", plan.ID.String())
	}
	return
}
