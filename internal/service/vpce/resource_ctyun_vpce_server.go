package vpce

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
	"terraform-provider-ctyun/internal/utils"
)

var (
	_ resource.Resource                = &ctyunVpceServer{}
	_ resource.ResourceWithConfigure   = &ctyunVpceServer{}
	_ resource.ResourceWithImportState = &ctyunVpceServer{}
)

type ctyunVpceServer struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpceServer() resource.Resource {
	return &ctyunVpceServer{}
}

func (c *ctyunVpceServer) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpce_server"
}

type CtyunVpceServerConfig struct {
	ID             types.String `tfsdk:"id"`
	RegionID       types.String `tfsdk:"region_id"`
	VpcID          types.String `tfsdk:"vpc_id"`
	Type           types.String `tfsdk:"type"`
	Name           types.String `tfsdk:"name"`
	InstanceType   types.String `tfsdk:"instance_type"`
	InstanceID     types.String `tfsdk:"instance_id"`
	SubnetID       types.String `tfsdk:"subnet_id"`
	AutoConnection types.Bool   `tfsdk:"auto_connection"`
	Rules          types.Set    `tfsdk:"rules"`
	WhitelistEmail types.Set    `tfsdk:"whitelist_email"`
	whitelist      []string
	rules          []CtyunVpceServerRule
}

type CtyunVpceServerRule struct {
	Protocol     types.String `tfsdk:"protocol"`
	ServerPort   types.Int32  `tfsdk:"server_port"`
	EndpointPort types.Int32  `tfsdk:"endpoint_port"`
}

func (c *ctyunVpceServer) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：**`,
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
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "关联的vpcID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "接口还是反向，interface:接口，reverse:反向",
				Validators: []validator.String{
					stringvalidator.OneOf("interface", "reverse"),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "支持拉丁字母、中文、数字，下划线，连字符，中文/英文字母开头，不能以http:/https:开头，长度2-32",
			},
			"instance_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "服务后端实例类型，vm:虚机类型,bm:物理机,vip:vip类型,lb:负载均衡类型,当type为interface时，必填",
				Validators: []validator.String{
					stringvalidator.OneOf("vm", "bm", "vip", "lb"),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("type"),
						types.StringValue("interface"),
					),
				},
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "服务后端实例id,当type为interface时，必填",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("type"),
						types.StringValue("interface"),
					),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "服务后端子网id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auto_connection": schema.BoolAttribute{
				Required:    true,
				Description: "是否自动连接，true表示自动链接，false表示非自动链接",
			},
			"whitelist_email": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "白名单邮箱",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(validator2.Email()),
					setvalidator.SizeAtMost(10),
				},
			},
			"rules": schema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Set{
					validator2.AlsoRequiresEqualSet(
						path.MatchRoot("type"),
						types.StringValue("interface"),
					),
				},
				Description: "节点服务规则,当type为interface时，必填",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"protocol": schema.StringAttribute{
							Required:    true,
							Description: "协议，TCP:TCP协议,UDP:UDP协议",
							Validators: []validator.String{
								stringvalidator.OneOf("TCP", "UDP"),
							},
						},
						"server_port": schema.Int32Attribute{
							Required:    true,
							Description: "服务端口(用于创建backend传入)(1-65535)",
							Validators: []validator.Int32{
								int32validator.Between(1, 65535),
							},
						},
						"endpoint_port": schema.Int32Attribute{
							Required:    true,
							Description: "节点端口(用于创建rule传入)(1-65535)",
							Validators: []validator.Int32{
								int32validator.Between(1, 65535),
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunVpceServer) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunVpceServerConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建
	endpointServerID, err := c.create(ctx, plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(endpointServerID)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)

	err = c.checkAfterCreate(ctx, plan)
	if err != nil {
		return
	}
	err = c.calcWhitelist(ctx, &plan)
	if err != nil {
		return
	}
	err = c.addWhitelist(ctx, plan)
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

func (c *ctyunVpceServer) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceServerConfig
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

func (c *ctyunVpceServer) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunVpceServerConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunVpceServerConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	plan.ID, plan.RegionID = state.ID, state.RegionID
	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	// 更新白名单
	err = c.updateWhitelist(ctx, plan, state)
	if err != nil {
		return
	}

	// 更新规则
	err = c.updateRule(ctx, plan, state)
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

func (c *ctyunVpceServer) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceServerConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
	response.State.RemoveResource(ctx)
}

func (c *ctyunVpceServer) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [endpointServerID],[regionID]
func (c *ctyunVpceServer) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunVpceServerConfig
	var endpointServerID, regionID string
	err = terraform_extend.Split(request.ID, &endpointServerID, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.ID = types.StringValue(endpointServerID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建
func (c *ctyunVpceServer) create(ctx context.Context, plan CtyunVpceServerConfig) (endpointSeverID string, err error) {
	params := &ctvpc.CtvpcCreateEndpointServiceRequest{
		ClientToken:    uuid.NewString(),
		RegionID:       plan.RegionID.ValueString(),
		VpcID:          plan.VpcID.ValueString(),
		Name:           plan.Name.ValueString(),
		RawType:        plan.Type.ValueStringPointer(),
		InstanceType:   plan.InstanceType.ValueStringPointer(),
		InstanceID:     plan.InstanceID.ValueStringPointer(),
		SubnetID:       plan.SubnetID.ValueStringPointer(),
		AutoConnection: plan.AutoConnection.ValueBool(),
	}

	err = c.calcRule(ctx, &plan)
	if err != nil {
		return
	}
	for _, r := range plan.rules {
		params.Rules = append(params.Rules, &ctvpc.CtvpcCreateEndpointServiceRulesRequest{
			Protocol:     r.Protocol.ValueString(),
			ServerPort:   r.ServerPort.ValueInt32(),
			EndpointPort: r.EndpointPort.ValueInt32(),
		})
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateEndpointServiceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil || resp.ReturnObj.EndpointService == nil {
		err = common.InvalidReturnObjError
		return
	}
	endpointSeverID = utils.SecString(resp.ReturnObj.EndpointService.EndpointServiceID)
	return
}

// checkAfterCreate 创建后检查
func (c *ctyunVpceServer) checkAfterCreate(ctx context.Context, plan CtyunVpceServerConfig) (err error) {
	endpointServer, err := c.show(ctx, plan)
	if err != nil {
		return
	}
	// 后端信息不正确
	if len(endpointServer.Backends) == 0 {
		err = fmt.Errorf("终端节点服务创建成功，但后端服务不正确，请使用正确instance_id重新创建")
		return
	}
	return
}

// buildRules 构造终端节点服务规则
func (c *ctyunVpceServer) buildRules(ctx context.Context, plan CtyunVpceServerConfig) (rulesReq []*ctvpc.CtvpcCreateEndpointServiceRulesRequest, err error) {
	if plan.Rules.IsUnknown() || plan.Rules.IsNull() {
		return
	}
	var rules []CtyunVpceServerRule
	diags := plan.Rules.ElementsAs(ctx, &rules, false)
	if diags.HasError() {
		err = fmt.Errorf(diags.Errors()[0].Detail())
		return
	}
	for _, r := range rules {
		item := &ctvpc.CtvpcCreateEndpointServiceRulesRequest{
			Protocol:     r.Protocol.ValueString(),
			ServerPort:   r.ServerPort.ValueInt32(),
			EndpointPort: r.EndpointPort.ValueInt32(),
		}
		rulesReq = append(rulesReq, item)
	}
	return
}

// getAndMerge 从远端查询
func (c *ctyunVpceServer) getAndMerge(ctx context.Context, plan *CtyunVpceServerConfig) (err error) {
	endpointServer, err := c.show(ctx, *plan)
	if err != nil {
		return
	}
	plan.VpcID = utils.SecStringValue(endpointServer.VpcID)
	plan.Name = utils.SecStringValue(endpointServer.Name)
	plan.Type = utils.SecStringValue(endpointServer.RawType)
	plan.AutoConnection = utils.SecBoolValue(endpointServer.AutoConnection)

	if len(endpointServer.Backends) != 0 {
		backend := endpointServer.Backends[0]
		plan.InstanceType = utils.SecStringValue(backend.InstanceType)
		plan.InstanceID = utils.SecStringValue(backend.InstanceID)
	}

	err = c.mergeRules(ctx, plan, endpointServer)
	if err != nil {
		return
	}

	err = c.mergeWhitelist(ctx, plan)
	if err != nil {
		return
	}

	return
}

// update 更新
func (c *ctyunVpceServer) update(ctx context.Context, plan, state CtyunVpceServerConfig) (err error) {
	endpointServerID, regionID := state.ID.ValueString(), state.RegionID.ValueString()
	params := &ctvpc.CtvpcModifyEndpointServiceRequest{
		RegionID:          regionID,
		EndpointServiceID: endpointServerID,
		Name:              plan.Name.ValueStringPointer(),
		AutoConnection:    plan.AutoConnection.ValueBoolPointer(),
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcModifyEndpointServiceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}

// delete 删除
func (c *ctyunVpceServer) delete(ctx context.Context, plan CtyunVpceServerConfig) (err error) {
	endpointServerID, regionID := plan.ID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcDeleteEndpointServiceRequest{
		RegionID:    regionID,
		ID:          endpointServerID,
		ClientToken: uuid.NewString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteEndpointServiceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}

// show 查询VPCEs详情
func (c *ctyunVpceServer) show(ctx context.Context, plan CtyunVpceServerConfig) (endpointServer ctvpc.CtvpcShowEndpointServiceReturnObjResponse, err error) {
	endpointServerID, regionID := plan.ID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcShowEndpointServiceRequest{
		RegionID:          regionID,
		EndpointServiceID: endpointServerID,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowEndpointServiceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	endpointServer = *resp.ReturnObj
	return
}

// calcWhitelist 将types.Set类型的白名单转换为[]string
func (c *ctyunVpceServer) calcWhitelist(ctx context.Context, plan *CtyunVpceServerConfig) (err error) {
	if plan.WhitelistEmail.IsNull() || plan.WhitelistEmail.IsUnknown() {
		return
	}
	plan.whitelist = []string{}
	diags := plan.WhitelistEmail.ElementsAs(ctx, &plan.whitelist, true)
	if diags.HasError() {
		err = fmt.Errorf(diags.Errors()[0].Detail())
	}
	return
}

// addWhitelist 添加白名单
func (c *ctyunVpceServer) addWhitelist(ctx context.Context, plan CtyunVpceServerConfig) (err error) {
	for _, email := range plan.whitelist {
		params := &ctvpc.CtvpcCreateEndpointServiceWhitelistRequest{
			ClientToken:       uuid.NewString(),
			RegionID:          plan.RegionID.ValueString(),
			EndpointServiceID: plan.ID.ValueString(),
			Email:             &email,
		}
		var resp *ctvpc.CtvpcCreateEndpointServiceWhitelistResponse
		resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcCreateEndpointServiceWhitelistApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		}
	}
	return
}

// delWhitelist 删除白名单
func (c *ctyunVpceServer) delWhitelist(ctx context.Context, plan CtyunVpceServerConfig) (err error) {
	for _, email := range plan.whitelist {
		params := &ctvpc.CtvpcDeleteEndpointServiceWhitelistRequest{
			ClientToken:       uuid.NewString(),
			RegionID:          plan.RegionID.ValueString(),
			EndpointServiceID: plan.ID.ValueString(),
			Email:             &email,
		}
		var resp *ctvpc.CtvpcDeleteEndpointServiceWhitelistResponse
		resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcDeleteEndpointServiceWhitelistApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		}
	}
	return
}

// updateWhitelist 更新白名单
func (c *ctyunVpceServer) updateWhitelist(ctx context.Context, plan, state CtyunVpceServerConfig) (err error) {
	err = c.calcWhitelist(ctx, &plan)
	if err != nil {
		return
	}
	err = c.calcWhitelist(ctx, &state)
	if err != nil {
		return
	}

	add, del := utils.DifferenceStrArray(plan.whitelist, state.whitelist)
	plan.whitelist = del
	err = c.delWhitelist(ctx, plan)
	if err != nil {
		return
	}
	plan.whitelist = add
	err = c.addWhitelist(ctx, plan)
	if err != nil {
		return
	}
	return
}

// mergeWhitelist 查询当前白名单
func (c *ctyunVpceServer) mergeWhitelist(ctx context.Context, plan *CtyunVpceServerConfig) (err error) {
	params := ctvpc.CtvpcNewEndpointServiceWhiteListRequest{
		RegionID:          plan.RegionID.ValueString(),
		EndpointServiceID: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewEndpointServiceWhiteListApi.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	whitelist := []string{}
	for _, item := range resp.ReturnObj.Whitelist {
		if item != nil && item.Email != nil {
			whitelist = append(whitelist, *item.Email)
		}
	}
	t, diags := types.SetValueFrom(ctx, types.StringType, whitelist)
	if diags.HasError() {
		err = fmt.Errorf(diags.Errors()[0].Detail())
		return
	}
	plan.WhitelistEmail = t
	plan.whitelist = whitelist
	return
}

// addRule 新增端口映射
func (c *ctyunVpceServer) addRule(ctx context.Context, plan CtyunVpceServerConfig) (err error) {
	for _, rule := range plan.rules {
		params := &ctvpc.CtvpcCreateEndpointServiceRuleRequest{
			ClientToken:       uuid.NewString(),
			RegionID:          plan.RegionID.ValueString(),
			EndpointServiceID: plan.ID.ValueString(),
			Protocol:          rule.Protocol.ValueString(),
			EndpointPort:      rule.EndpointPort.ValueInt32(),
			ServerPort:        rule.ServerPort.ValueInt32(),
		}
		var resp *ctvpc.CtvpcCreateEndpointServiceRuleResponse
		resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcCreateEndpointServiceRuleApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		}
	}
	return
}

// delRule 删除端口映射
func (c *ctyunVpceServer) delRule(ctx context.Context, plan CtyunVpceServerConfig) (err error) {
	for _, rule := range plan.rules {
		params := &ctvpc.CtvpcDeleteEndpointServiceRuleRequest{
			ClientToken:       uuid.NewString(),
			RegionID:          plan.RegionID.ValueString(),
			EndpointServiceID: plan.ID.ValueString(),
			Protocol:          rule.Protocol.ValueString(),
			EndpointPort:      rule.EndpointPort.ValueInt32(),
			ServerPort:        rule.ServerPort.ValueInt32(),
		}
		var resp *ctvpc.CtvpcDeleteEndpointServiceRuleResponse
		resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcDeleteEndpointServiceRuleApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		}
	}
	return
}

// delRule 删除端口映射
func (c *ctyunVpceServer) calcRule(ctx context.Context, plan *CtyunVpceServerConfig) (err error) {
	if plan.Rules.IsUnknown() || plan.Rules.IsNull() {
		return
	}
	plan.rules = []CtyunVpceServerRule{}
	diags := plan.Rules.ElementsAs(ctx, &plan.rules, false)
	if diags.HasError() {
		err = fmt.Errorf(diags.Errors()[0].Detail())
		return
	}
	return
}

// updateRule 更新端口映射
func (c *ctyunVpceServer) updateRule(ctx context.Context, plan, state CtyunVpceServerConfig) (err error) {
	err = c.calcRule(ctx, &plan)
	if err != nil {
		return
	}
	err = c.calcRule(ctx, &state)
	if err != nil {
		return
	}

	add, del := utils.DifferenceStructArray[CtyunVpceServerRule](plan.rules, state.rules)
	plan.rules = del
	err = c.delRule(ctx, plan)
	if err != nil {
		return
	}
	plan.rules = add
	err = c.addRule(ctx, plan)
	if err != nil {
		return
	}
	return
}

// mergeRules 计算当前rules
func (c *ctyunVpceServer) mergeRules(ctx context.Context, plan *CtyunVpceServerConfig, endpointServer ctvpc.CtvpcShowEndpointServiceReturnObjResponse) (err error) {
	rules := []CtyunVpceServerRule{}
	for _, item := range endpointServer.Rules {
		if item != nil {
			rules = append(rules, CtyunVpceServerRule{
				Protocol:     utils.SecStringValue(item.Protocol),
				EndpointPort: types.Int32Value(item.EndpointPort),
				ServerPort:   types.Int32Value(item.ServerPort),
			})
		}
	}
	ruleObj := utils.StructToTFObjectTypes(CtyunVpceServerRule{})
	t, diags := types.SetValueFrom(ctx, ruleObj, rules)
	if diags.HasError() {
		err = fmt.Errorf(diags.Errors()[0].Detail())
		return
	}
	plan.Rules = t
	plan.rules = rules
	return
}
