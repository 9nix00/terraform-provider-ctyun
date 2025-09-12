package mongodb

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"regexp"
)

var (
	_ resource.Resource                = &CtyunMongodbAccount{}
	_ resource.ResourceWithConfigure   = &CtyunMongodbAccount{}
	_ resource.ResourceWithImportState = &CtyunMongodbAccount{}
)

func NewCtyunMongodbAccount() resource.Resource {
	return &CtyunMongodbAccount{}
}

type CtyunMongodbAccount struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMongodbAccount) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_account"
}

func (c *CtyunMongodbAccount) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `**MongoDB数据库账号管理资源,详细说明请见文档 https://www.ctyun.cn/document/10034467/10089535**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识，格式为 {instance_id}:{account_name}",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "MongoDB实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "账号名称，以字母开头，由字母、数字和下划线组成，长度2-16个字符",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 16),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]*$"), "实例名称不符合规范"),
				},
			},
			// 实现一个validator方法
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "实例密码（8-26位由大写字母、小写字母、数字、特殊字符中的任意三种组成 特殊字符为~!@#%^*_=+），RSA公钥加密存储,支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 26),
					validator2.MongodbPassword(),
				},
			},
			"database": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "数据库名称，默认为admin",
				Default:     stringdefault.StaticString("admin"),
			},
			"page_now": schema.Int32Attribute{
				Optional:    true,
				Description: "当前页",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "单页记录条数",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"roles": schema.ListNestedAttribute{
				Required:    true,
				Description: "角色列表  ,支持更新",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"db": schema.StringAttribute{
							Required:    true,
							Description: "数据库名称",
						},
						"role": schema.StringAttribute{
							Required:    true,
							Description: "角色，可选值：read、readWrite、readWriteAnyDatabase等，默认为readWrite",
							Validators: []validator.String{
								stringvalidator.OneOf("read", "readWrite", "readAnyDatabase", "readWriteAnyDatabase", "dbAdmin", "dbAdminAnyDatabase", "userAdmin", "userAdminAnyDatabase", "clusterAdmin"),
							},
						},
					},
				},
			},
		},
	}
}

func (c *CtyunMongodbAccount) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMongodbAccount) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MongodbAccountConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 创建账号权限列表
	var roles []mongodb.MongodbAccountRole

	// 如果定义了roles块，则使用roles块中的权限

	for _, role := range plan.Roles {
		roles = append(roles, mongodb.MongodbAccountRole{
			DB:   role.Database.ValueString(),
			Role: role.Role.ValueString(),
		})
	}

	// 对密码进行Base64编码
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(plan.Password.ValueString()))

	createReq := &mongodb.MongodbCreateAccountRequest{
		AccountName:     plan.Name.ValueString(),
		AccountPassword: encodedPassword,
		Roles:           &roles,
	}

	// 只有当database字段被设置时才添加到请求中
	if !plan.Database.IsNull() && !plan.Database.IsUnknown() {
		createReq.DatabaseName = plan.Database.ValueStringPointer()
	}

	headers := &mongodb.MongodbCreateAccountRequestHeaders{
		RegionID:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		headers.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "开始创建MongoDB账号", map[string]interface{}{
		"instance_id":  plan.InstanceID.ValueString(),
		"account_name": plan.Name.ValueString(),
	})

	createResp, err := c.meta.Apis.SdkMongodbApis.MongodbCreateAccountApi.Do(ctx, c.meta.Credential, createReq, headers)
	if err != nil {
		resp.Diagnostics.AddError("创建MongoDB账号失败", err.Error())
		return
	}

	if createResp.StatusCode != 800 {
		resp.Diagnostics.AddError("创建MongoDB账号失败", fmt.Sprintf("API返回错误: %s", *createResp.Message))
		return
	}

	// 设置ID
	id := fmt.Sprintf("%s:%s", plan.InstanceID.ValueString(), plan.Name.ValueString())
	plan.ID = types.StringValue(id)

	// 查询账号信息
	if err := c.readAccountInfo(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("获取账号信息失败", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunMongodbAccount) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MongodbAccountConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := c.readAccountInfo(ctx, &state); err != nil {
		resp.Diagnostics.AddError("获取账号信息失败", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunMongodbAccount) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state MongodbAccountConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 如果密码变更，更新密码
	if !plan.Password.Equal(state.Password) {
		err := c.updateAccountPassword(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError("更新MongoDB账号密码失败", err.Error())
			return
		}
	}

	// 如果权限变更，更新权限
	// 检查roles块或者database字段是否变更
	rolesChanged := len(plan.Roles) != len(state.Roles)
	if !rolesChanged {
		for i := range plan.Roles {
			if !plan.Roles[i].Database.Equal(state.Roles[i].Database) ||
				!plan.Roles[i].Role.Equal(state.Roles[i].Role) {
				rolesChanged = true
				break
			}
		}
	}

	if rolesChanged || !plan.Database.Equal(state.Database) {
		err := c.updateAccountPermission(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError("更新MongoDB账号权限失败", err.Error())
			return
		}
	}

	if err := c.readAccountInfo(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("获取账号信息失败", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *CtyunMongodbAccount) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MongodbAccountConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteReq := &mongodb.MongodbDeleteAccountRequest{
		AccountName:  state.Name.ValueString(),
		DatabaseName: state.Database.ValueString(),
	}

	headers := &mongodb.MongodbDeleteAccountRequestHeaders{
		RegionID:   state.RegionID.ValueString(),
		ProdInstId: state.InstanceID.ValueString(),
	}
	if !state.ProjectID.IsNull() {
		headers.ProjectID = state.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "删除MongoDB账号", map[string]interface{}{
		"instance_id":  state.InstanceID.ValueString(),
		"account_name": state.Name.ValueString(),
	})

	deleteResp, err := c.meta.Apis.SdkMongodbApis.MongodbDeleteAccountApi.Do(ctx, c.meta.Credential, deleteReq, headers)
	if err != nil {
		resp.Diagnostics.AddError("删除MongoDB账号失败", err.Error())
		return
	}

	if deleteResp.StatusCode != 800 {
		resp.Diagnostics.AddError("删除MongoDB账号失败", fmt.Sprintf("API返回错误: %s", *deleteResp.Message))
		return
	}
}

func (c *CtyunMongodbAccount) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// readAccountInfo 查询账号信息并更新到配置中
func (c *CtyunMongodbAccount) readAccountInfo(ctx context.Context, plan *MongodbAccountConfig) error {
	// 解析ID获取instance_id和account_name
	instanceID := plan.InstanceID.ValueString()

	describeReq := &mongodb.MongodbDescribeAccountsRequest{
		ProdInstId: instanceID,
		PageNow:    plan.PageNow.ValueInt32(),
		PageSize:   plan.PageSize.ValueInt32(),
	}

	headers := &mongodb.MongodbDescribeAccountsRequestHeaders{
		RegionID: plan.RegionID.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		headers.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	describeResp, err := c.meta.Apis.SdkMongodbApis.MongodbDescribeAccountsApi.Do(ctx, c.meta.Credential, describeReq, headers)
	if err != nil {
		return err
	}

	if describeResp.StatusCode != 800 {
		if describeResp.Message != nil {
			return fmt.Errorf("API返回错误: %s", *describeResp.Message)
		} else {
			return fmt.Errorf("API返回错误，状态码: %d", describeResp.StatusCode)
		}
	}

	if describeResp.ReturnObj == nil || len(describeResp.ReturnObj.List) == 0 {
		return fmt.Errorf("未找到账号信息")
	}

	// 这里获取账号信息, 并更新到配置中  循环list 获取账号信息user字段 ==plan.user
	for _, account := range describeResp.ReturnObj.List {
		if account.User == plan.Name.ValueString() {
			// 如果数据库没有设置，则使用API返回的值
			if plan.Database.IsNull() || plan.Database.IsUnknown() {
				plan.Database = types.StringValue(account.DB)
			}

			// 如果roles块为空，则根据API返回的角色信息填充
			if len(plan.Roles) == 0 && len(account.Roles) > 0 {
				plan.Roles = make([]MongodbAccountRole, len(account.Roles))
				for i, role := range account.Roles {
					plan.Roles[i] = MongodbAccountRole{
						Database: types.StringValue(role.DB),
						Role:     types.StringValue(role.Role),
					}
				}
			}

			return nil
		}
	}

	return fmt.Errorf("未找到指定账号: %s", plan.Name.ValueString())
}

// updateAccountPassword 更新账户密码
func (c *CtyunMongodbAccount) updateAccountPassword(ctx context.Context, plan *MongodbAccountConfig) error {
	// 对密码进行Base64编码
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(plan.Password.ValueString()))

	updatePasswordReq := &mongodb.MongodbUpdateAccountPasswordRequest{
		AccountName:     plan.Name.ValueString(),
		AccountPassword: encodedPassword,
		Database:        plan.Database.ValueString(),
	}

	headers := &mongodb.MongodbUpdateAccountPasswordRequestHeaders{
		RegionID:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		headers.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "更新MongoDB账号密码", map[string]interface{}{
		"instance_id":  plan.InstanceID.ValueString(),
		"account_name": plan.Name.ValueString(),
	})

	updatePasswordResp, err := c.meta.Apis.SdkMongodbApis.MongodbUpdateAccountPasswordApi.Do(ctx, c.meta.Credential, updatePasswordReq, headers)
	if err != nil {
		return err
	}

	if updatePasswordResp.StatusCode != 800 {
		return fmt.Errorf("API返回错误: %s", *updatePasswordResp.Message)
	}

	return nil
}

// updateAccountPermission 更新账户权限
func (c *CtyunMongodbAccount) updateAccountPermission(ctx context.Context, plan *MongodbAccountConfig) error {
	// 创建权限对象数组
	var roles []mongodb.MongodbAccountRole

	// 如果定义了roles块，则使用roles块中的权限
	for _, role := range plan.Roles {
		roles = append(roles, mongodb.MongodbAccountRole{
			DB:   role.Database.ValueString(),
			Role: role.Role.ValueString(),
		})
	}

	modifyPermissionReq := &mongodb.MongodbModifyAccountPermissionRequest{
		AccountName: plan.Name.ValueString(),
		Roles:       &roles,
	}

	// 只有当database字段被设置时才添加到请求中
	if !plan.Database.IsNull() && !plan.Database.IsUnknown() {
		modifyPermissionReq.DatabaseName = plan.Database.ValueStringPointer()
	}

	headers := &mongodb.MongodbModifyAccountPermissionRequestHeaders{
		RegionID:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		headers.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "更新MongoDB账号权限", map[string]interface{}{
		"instance_id":  plan.InstanceID.ValueString(),
		"account_name": plan.Name.ValueString(),
	})

	modifyPermissionResp, err := c.meta.Apis.SdkMongodbApis.MongodbModifyAccountPermissionApi.Do(ctx, c.meta.Credential, modifyPermissionReq, headers)
	if err != nil {
		return err
	}

	if modifyPermissionResp.StatusCode != 800 {
		return fmt.Errorf("API返回错误: %s", *modifyPermissionResp.Message)
	}

	return nil
}

type MongodbAccountRole struct {
	Database types.String `tfsdk:"db"`
	Role     types.String `tfsdk:"role"`
}

type MongodbAccountConfig struct {
	ID         types.String         `tfsdk:"id"`
	InstanceID types.String         `tfsdk:"instance_id"`
	RegionID   types.String         `tfsdk:"region_id"`
	ProjectID  types.String         `tfsdk:"project_id"`
	Name       types.String         `tfsdk:"name"`
	Password   types.String         `tfsdk:"password"`
	Database   types.String         `tfsdk:"database"`
	PageNow    types.Int32          `tfsdk:"page_now"`
	PageSize   types.Int32          `tfsdk:"page_size"`
	Roles      []MongodbAccountRole `tfsdk:"roles"`
}
