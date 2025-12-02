package sdwan

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sdwan"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var (
	_ resource.Resource                = &CtyunSdwanAcl{}
	_ resource.ResourceWithConfigure   = &CtyunSdwanAcl{}
	_ resource.ResourceWithImportState = &CtyunSdwanAcl{}
)

func NewCtyunSdwanAcl() resource.Resource {
	return &CtyunSdwanAcl{}
}

type CtyunSdwanAcl struct {
	meta *common.CtyunMetadata
}

type CtyunSdwanAclConfig struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
	Rules     types.List   `tfsdk:"rules"`
}

type SdwanAclRule struct {
	Direction    types.String `tfsdk:"direction"`
	Protocol     types.Int32  `tfsdk:"protocol"`
	IpVersion    types.String `tfsdk:"ip_version"`
	DstCidr      types.String `tfsdk:"dst_cidr"`
	DstPortRange types.String `tfsdk:"dst_port_range"`
	Priority     types.String `tfsdk:"priority"`
	Action       types.String `tfsdk:"action"`
	SrcCidr      types.String `tfsdk:"src_cidr"`
	SrcPortRange types.String `tfsdk:"src_port_range"`
}

func (c *CtyunSdwanAcl) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdwan_acl"
}

func (c *CtyunSdwanAcl) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `**SD-WAN访问控制资源,详细说明请见文档**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, true),
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "访问控制名称",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"rules": schema.ListNestedAttribute{
				Required:    true,
				Description: "ACL规则列表",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"direction": schema.StringAttribute{
							Required:    true,
							Description: "控制方向，取值范围: in(入方向), out(出方向)",
							Validators: []validator.String{
								stringvalidator.OneOf("in", "out"),
							},
						},
						"protocol": schema.StringAttribute{
							Required:    true,
							Description: "协议类型，取值范围: udp(UDP), icmp(ICMP), all(ALL), tcp(TCP)",
							Validators: []validator.String{
								stringvalidator.OneOf("udp", "icmp", "all", "tcp"),
							},
						},
						"ip_version": schema.StringAttribute{
							Required:    true,
							Description: "IP协议版本，取值范围: IPv4(IPv4), IPv6(IPv6)",
							Validators: []validator.String{
								stringvalidator.OneOf("IPv4", "IPv6"),
							},
						},
						"dst_cidr": schema.StringAttribute{
							Required:    true,
							Description: "目的网段",
						},
						"dst_port_range": schema.StringAttribute{
							Required:    true,
							Description: "目的端口范围（例如1-200， -1/-1为默认值，表示1-65535）",
						},
						"priority": schema.Int32Attribute{
							Required:    true,
							Description: "优先级",
						},
						"action": schema.StringAttribute{
							Required:    true,
							Description: "策略类型，取值范围: allow(允许), deny(拒绝)",
							Validators: []validator.String{
								stringvalidator.OneOf("allow", "deny"),
							},
						},
						"src_cidr": schema.StringAttribute{
							Required:    true,
							Description: "源网段",
						},
						"src_port_range": schema.StringAttribute{
							Required:    true,
							Description: "源端口范围（例如1-200， -1/-1为默认值，表示1-65535）",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunSdwanAcl) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunSdwanAcl) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunSdwanAclConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunSdwanAcl) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSdwanAclConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 查询远端确认资源是否存在
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunSdwanAcl) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// SD-WAN ACL 不支持更新操作，直接返回
	var plan CtyunSdwanAclConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *CtyunSdwanAcl) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSdwanAclConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunSdwanAcl) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunSdwanAclConfig
	cfg.ID = types.StringValue(req.ID)

	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}

	// 导入时不设置 rules 字段，保持其为未知状态
	cfg.Rules = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"direction":      types.StringType,
			"protocol":       types.StringType,
			"ip_version":     types.StringType,
			"dst_cidr":       types.StringType,
			"dst_port_range": types.StringType,
			"priority":       types.Int32Type,
			"action":         types.StringType,
			"src_cidr":       types.StringType,
			"src_port_range": types.StringType,
		},
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, cfg)...)
}

func (c *CtyunSdwanAcl) create(ctx context.Context, plan *CtyunSdwanAclConfig) (err error) {
	var rules []*sdwan.SdwanCreateSdwanAclRulesRequest

	// 遍历规则列表
	for _, ruleValue := range plan.Rules.Elements() {
		// 将规则值转换为对象
		ruleObj := ruleValue.(types.Object)
		attributes := ruleObj.Attributes()

		rules = append(rules, &sdwan.SdwanCreateSdwanAclRulesRequest{
			Direction:    attributes["direction"].(types.String).ValueString(),
			Protocol:     attributes["protocol"].(types.String).ValueString(),
			IpVersion:    attributes["ip_version"].(types.String).ValueString(),
			DstCidr:      attributes["dst_cidr"].(types.String).ValueString(),
			DstPortRange: attributes["dst_port_range"].(types.String).ValueString(),
			Priority:     attributes["priority"].(types.Int32).ValueInt32(),
			Action:       attributes["action"].(types.String).ValueString(),
			SrcCidr:      attributes["src_cidr"].(types.String).ValueString(),
			SrcPortRange: attributes["src_port_range"].(types.String).ValueString(),
		})
	}

	createReq := &sdwan.SdwanCreateSdwanAclRequest{
		AclName:   plan.Name.ValueString(),
		ProjectID: plan.ProjectID.ValueString(),
		Rules:     rules,
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanCreateSdwanAclApi.Do(ctx, c.meta.SdkCredential, createReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}

	// 等待创建完成
	time.Sleep(3 * time.Second)
	return
}

func (c *CtyunSdwanAcl) getAndMerge(ctx context.Context, plan *CtyunSdwanAclConfig) (err error) {
	listReq := &sdwan.SdwanGetSdwanAclRequest{
		PageNo:   1,
		PageSize: 1000,
	}
	if plan.ID.ValueStringPointer() != nil {
		listReq.AclID = plan.ID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkSdwanApis.SdwanGetSdwanAclApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}

	// 查找对应的ACL
	var found bool
	if plan.ID.ValueString() != "" {
		for _, aclItem := range resp.ReturnObj.Result {
			if aclItem.AclID != nil && *aclItem.AclID == plan.ID.ValueString() {
				plan.Name = types.StringValue(*aclItem.Name)
				found = true
				break
			}
		}

		if !found {
			return common.ResourceNotExistError
		}
		return
	} else if plan.Name.ValueString() != "" {
		for _, aclItem := range resp.ReturnObj.Result {
			if aclItem.Name != nil && *aclItem.Name == plan.Name.ValueString() {
				plan.ID = types.StringValue(*aclItem.AclID)
				found = true
				break
			}
		}
		if !found {
			return common.ResourceNotExistError
		}
		return
	}

	return
}

func (c *CtyunSdwanAcl) delete(ctx context.Context, state *CtyunSdwanAclConfig) (err error) {
	deleteReq := &sdwan.SdwanDeleteSdwanAclRequest{
		AclID: state.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanDeleteSdwanAclApi.Do(ctx, c.meta.SdkCredential, deleteReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}

	// 等待删除完成
	time.Sleep(3 * time.Second)

	return
}

func (c *CtyunSdwanAcl) update(ctx context.Context, plan *CtyunSdwanAclConfig) (err error) {
	updateReq := &sdwan.SdwanUpdateSdwanAclRequest{
		AclID:   plan.ID.ValueString(),
		AclName: plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanUpdateSdwanAclApi.Do(ctx, c.meta.SdkCredential, updateReq)
	if err != nil {
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}

	// 等待更新完成
	time.Sleep(3 * time.Second)
	return nil
}
