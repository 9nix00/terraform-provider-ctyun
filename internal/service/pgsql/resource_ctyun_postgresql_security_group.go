package pgsql

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &CtyunPostgresqlSecurityGroup{}
	_ resource.ResourceWithConfigure   = &CtyunPostgresqlSecurityGroup{}
	_ resource.ResourceWithImportState = &CtyunPostgresqlSecurityGroup{}
)

func NewCtyunPostgresqlSecurityGroup() resource.Resource {
	return &CtyunPostgresqlSecurityGroup{}
}

type CtyunPostgresqlSecurityGroup struct {
	meta *common.CtyunMetadata
}

type CtyunPostgresqlSecurityGroupConfig struct {
	ID               types.String `tfsdk:"id"`
	InstanceID       types.String `tfsdk:"instance_id"`
	ProjectID        types.String `tfsdk:"project_id"`
	RegionID         types.String `tfsdk:"region_id"`
	SecurityGroupIds types.Set    `tfsdk:"security_group_ids"`
}

func (c *CtyunPostgresqlSecurityGroup) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_security_group"
}

func (c *CtyunPostgresqlSecurityGroup) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `**PostgreSQL 实例安全组资源**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识，格式为 instance_id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "PostgreSQL实例ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
			// 项目相关
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
			"security_group_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "安全组ID集合",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (c *CtyunPostgresqlSecurityGroup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPostgresqlSecurityGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPostgresqlSecurityGroupConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 创建安全组关联
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}

	// 设置ID
	plan.ID = plan.InstanceID

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunPostgresqlSecurityGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPostgresqlSecurityGroupConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 查询安全组信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunPostgresqlSecurityGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPostgresqlSecurityGroupConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 更新安全组关联
	err = c.update(ctx, &plan)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunPostgresqlSecurityGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPostgresqlSecurityGroupConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 删除安全组关联
	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunPostgresqlSecurityGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	return
}

func (c *CtyunPostgresqlSecurityGroup) create(ctx context.Context, plan *CtyunPostgresqlSecurityGroupConfig) (err error) {
	// 获取安全组ID集合
	var securityGroupIds []string
	diags := plan.SecurityGroupIds.ElementsAs(ctx, &securityGroupIds, false)
	if diags.HasError() {
		return fmt.Errorf("无法解析安全组ID集合: %v", err)
	}

	// 循环添加安全组
	for _, sgId := range securityGroupIds {
		request := &pgsql.PostgresqlAddSecurityGroupRequest{
			SecurityGroupId: sgId,
			InstanceId:      plan.InstanceID.ValueString(),
		}

		header := &pgsql.PostgresqlAddSecurityGroupRequestHeader{
			ProjectId: plan.ProjectID.ValueStringPointer(),
		}

		tflog.Info(ctx, "为PostgreSQL实例绑定安全组", map[string]interface{}{
			"instance_id":       plan.InstanceID.ValueString(),
			"security_group_id": sgId,
		})

		resp, err := c.meta.Apis.SdkCtPgsqlApis.PostgresqlAddSecurityGroupApi.Do(ctx, c.meta.Credential, request, header)
		if err != nil {
			return fmt.Errorf("绑定安全组 %s 失败: %v", sgId, err)
		}

		if resp.StatusCode != 200 && resp.StatusCode != 800 {
			return fmt.Errorf("API return error. Status Code: %d, Message: %s, Error: %s", resp.StatusCode, resp.Message, resp.Error)
		}
	}

	return nil
}

func (c *CtyunPostgresqlSecurityGroup) getAndMerge(ctx context.Context, state *CtyunPostgresqlSecurityGroupConfig) (err error) {
	// 查询实例的安全组信息
	request := &pgsql.PgsqlSecurityGroupListRequest{
		RegionID: state.RegionID.ValueString(), // 这个参数在API中可能不是必需的
		InstID:   state.InstanceID.ValueString(),
	}

	header := &pgsql.PgsqlSecurityGroupListRequestHeader{
		ProjectID: state.ProjectID.ValueStringPointer(),
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlSecurityGroupListApi.Do(ctx, c.meta.Credential, request, header)
	if err != nil {
		return err
	}

	if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Status Code: %d, Message: %v, Error: %v", resp.StatusCode, resp.Message, resp.Error)
	}

	if resp.ReturnObj == nil || len(resp.ReturnObj.Data) == 0 {
		// 如果没有安全组，设置为空集合
		state.SecurityGroupIds = types.SetNull(types.StringType)
		return nil
	}

	// 提取安全组ID集合
	var securityGroupIds []string
	for _, sg := range resp.ReturnObj.Data {
		securityGroupIds = append(securityGroupIds, sg.ID)
	}

	// 更新状态
	state.SecurityGroupIds, _ = types.SetValueFrom(ctx, types.StringType, securityGroupIds)
	// 设置ID
	state.ID = state.InstanceID

	return nil
}

func (c *CtyunPostgresqlSecurityGroup) update(ctx context.Context, plan *CtyunPostgresqlSecurityGroupConfig) (err error) {
	// 获取当前状态和计划中的安全组集合
	var planSecurityGroupIds []string
	diags := plan.SecurityGroupIds.ElementsAs(ctx, &planSecurityGroupIds, false)
	if diags.HasError() {
		return fmt.Errorf("无法解析计划中的安全组ID集合: %v", err)
	}

	// 获取当前状态中的安全组集合
	var currentSecurityGroupIds []string
	err = c.getSecurityGroupIdsFromState(ctx, plan, &currentSecurityGroupIds)
	if err != nil {
		return fmt.Errorf("无法获取当前安全组集合: %v", err)
	}

	// 计算需要添加和删除的安全组
	toAdd, toRemove := c.calculateDiff(currentSecurityGroupIds, planSecurityGroupIds)

	// 删除需要移除的安全组
	for _, sgId := range toRemove {
		err = c.removeSecurityGroup(ctx, plan, sgId)
		if err != nil {
			return fmt.Errorf("删除安全组 %s 失败: %v", sgId, err)
		}
	}

	// 添加需要新增的安全组
	for _, sgId := range toAdd {
		err = c.addSecurityGroup(ctx, plan, sgId)
		if err != nil {
			return fmt.Errorf("添加安全组 %s 失败: %v", sgId, err)
		}
	}

	return nil
}

func (c *CtyunPostgresqlSecurityGroup) delete(ctx context.Context, state *CtyunPostgresqlSecurityGroupConfig) (err error) {
	// 获取当前安全组集合
	var currentSecurityGroupIds []string
	err = c.getSecurityGroupIdsFromState(ctx, state, &currentSecurityGroupIds)
	if err != nil {
		return fmt.Errorf("无法获取当前安全组集合: %v", err)
	}

	// 删除所有安全组
	for _, sgId := range currentSecurityGroupIds {
		err = c.removeSecurityGroup(ctx, state, sgId)
		if err != nil {
			return fmt.Errorf("删除安全组 %s 失败: %v", sgId, err)
		}
	}

	return nil
}

// 获取当前状态中的安全组集合
func (c *CtyunPostgresqlSecurityGroup) getSecurityGroupIdsFromState(ctx context.Context, plan *CtyunPostgresqlSecurityGroupConfig, securityGroupIds *[]string) error {
	request := &pgsql.PgsqlSecurityGroupListRequest{
		RegionID: plan.RegionID.ValueString(), // 这个参数在API中可能不是必需的
		InstID:   plan.InstanceID.ValueString(),
	}

	header := &pgsql.PgsqlSecurityGroupListRequestHeader{
		ProjectID: plan.ProjectID.ValueStringPointer(),
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlSecurityGroupListApi.Do(ctx, c.meta.Credential, request, header)
	if err != nil {
		return err
	}

	if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Status Code: %d, Message: %v, Error: %v", resp.StatusCode, resp.Message, resp.Error)
	}

	if resp.ReturnObj != nil && len(resp.ReturnObj.Data) > 0 {
		for _, sg := range resp.ReturnObj.Data {
			*securityGroupIds = append(*securityGroupIds, sg.ID)
		}
	}

	return nil
}

// 计算安全组差异
func (c *CtyunPostgresqlSecurityGroup) calculateDiff(current, plan []string) (toAdd, toRemove []string) {
	// 创建映射便于查找
	currentMap := make(map[string]bool)
	planMap := make(map[string]bool)

	for _, id := range current {
		currentMap[id] = true
	}

	for _, id := range plan {
		planMap[id] = true
	}

	// 找出需要添加的安全组 (在计划中但不在当前状态中)
	for id := range planMap {
		if !currentMap[id] {
			toAdd = append(toAdd, id)
		}
	}

	// 找出需要删除的安全组 (在当前状态中但不在计划中)
	for id := range currentMap {
		if !planMap[id] {
			toRemove = append(toRemove, id)
		}
	}

	return toAdd, toRemove
}

// 添加安全组
func (c *CtyunPostgresqlSecurityGroup) addSecurityGroup(ctx context.Context, plan *CtyunPostgresqlSecurityGroupConfig, sgId string) error {
	request := &pgsql.PostgresqlAddSecurityGroupRequest{
		SecurityGroupId: sgId,
		InstanceId:      plan.InstanceID.ValueString(),
	}

	header := &pgsql.PostgresqlAddSecurityGroupRequestHeader{
		ProjectId: plan.ProjectID.ValueStringPointer(),
	}

	tflog.Info(ctx, "为PostgreSQL实例绑定安全组", map[string]interface{}{
		"instance_id":       plan.InstanceID.ValueString(),
		"security_group_id": sgId,
	})

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PostgresqlAddSecurityGroupApi.Do(ctx, c.meta.Credential, request, header)
	if err != nil {
		return fmt.Errorf("绑定安全组 %s 失败: %v", sgId, err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 800 {
		return fmt.Errorf("API return error. Status Code: %d, Message: %s, Error: %s", resp.StatusCode, resp.Message, resp.Error)
	}

	return nil
}

// 删除安全组
func (c *CtyunPostgresqlSecurityGroup) removeSecurityGroup(ctx context.Context, plan *CtyunPostgresqlSecurityGroupConfig, sgId string) error {
	request := &pgsql.PgsqlDeleteSecurityGroupRequest{
		SecurityGroupId: sgId,
		InstanceId:      plan.InstanceID.ValueString(),
	}

	header := &pgsql.PgsqlDeleteSecurityGroupRequestHeader{
		ProjectID: plan.ProjectID.ValueStringPointer(),
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlDeleteSecurityGroupApi.Do(ctx, c.meta.Credential, request, header)
	if err != nil {
		return fmt.Errorf("删除安全组 %s 失败: %v", sgId, err)
	}

	if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Status Code: %d, Message: %s, Error: %s", resp.StatusCode, resp.Message, resp.Error)
	}

	return nil
}
