package mongodb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &CtyunMongodbBackupResource{}
	_ resource.ResourceWithConfigure   = &CtyunMongodbBackupResource{}
	_ resource.ResourceWithImportState = &CtyunMongodbBackupResource{}
)

func NewCtyunMongodbBackupResource() resource.Resource {
	return &CtyunMongodbBackupResource{}
}

// CtyunMongodbBackupResource defines the resource implementation.
type CtyunMongodbBackupResource struct {
	meta *common.CtyunMetadata
}

// CtyunMongodbBackupResourceModel describes the resource data model.
type CtyunMongodbBackupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	RegionID    types.String `tfsdk:"region_id"`
	ProjectID   types.String `tfsdk:"project_id"`
	InstanceID  types.String `tfsdk:"instance_id"`
	BackupName  types.String `tfsdk:"backup_name"`
	Description types.String `tfsdk:"description"`
}

func (r *CtyunMongodbBackupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_backup"
}

func (r *CtyunMongodbBackupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "天翼云MongoDB备份资源",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源ID，格式为region_id:instance_id:backup_name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"backup_name": schema.StringAttribute{
				Required:    true,
				Description: "备份名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "备份描述",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *CtyunMongodbBackupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	r.meta = meta
}

func (r *CtyunMongodbBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CtyunMongodbBackupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create backup
	createReq := &mongodb.MongodbCreateBackupRequest{
		ProdInstId: data.InstanceID.ValueString(),
	}

	createReq.BackupName = data.BackupName.ValueStringPointer()

	if !data.Description.IsNull() {
		createReq.Description = data.Description.ValueStringPointer()
	}

	header := &mongodb.MongodbCreateBackupRequestHeaders{
		RegionID: data.RegionID.ValueString(),
	}

	if !data.ProjectID.IsNull() {
		header.ProjectID = data.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "开始创建MongoDB备份", map[string]interface{}{
		"instance_id": data.InstanceID.ValueString(),
	})

	response, err := r.meta.Apis.SdkMongodbApis.MongodbCreateBackupApi.Do(ctx, r.meta.Credential, createReq, header)
	if err != nil {
		resp.Diagnostics.AddError("创建MongoDB备份失败", err.Error())
		return
	}

	if response.StatusCode != 800 {
		resp.Diagnostics.AddError("创建MongoDB备份失败", fmt.Sprintf("API返回错误，状态码: %d, 错误信息: %s", response.StatusCode, response.Error))
		return
	}

	// 设置资源ID，使用region_id:instance_id:backup_name组合作为唯一标识
	data.ID = types.StringValue(fmt.Sprintf("%s:%s:%s", data.RegionID.ValueString(), data.InstanceID.ValueString(), *createReq.BackupName))
	// 保存数据到Terraform状态
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CtyunMongodbBackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CtyunMongodbBackupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 获取备份信息
	describeReq := &mongodb.MongodbDescribeBackupsRequest{
		ProdInstId: data.InstanceID.ValueString(),
		BackupName: data.BackupName.ValueStringPointer(),
		PageNow:    1,
		PageSize:   100,
	}

	header := &mongodb.MongodbDescribeBackupsRequestHeaders{
		RegionID: data.RegionID.ValueString(),
	}

	if !data.ProjectID.IsNull() {
		header.ProjectID = data.ProjectID.ValueStringPointer()
	}

	response, err := r.meta.Apis.SdkMongodbApis.MongodbDescribeBackupsApi.Do(ctx, r.meta.Credential, describeReq, header)
	if err != nil {
		resp.Diagnostics.AddError("查询MongoDB备份信息失败", err.Error())
		return
	}

	if response.StatusCode != 800 {
		resp.Diagnostics.AddError("查询MongoDB备份信息失败", fmt.Sprintf("API返回错误，状态码: %d, 错误信息: %s", response.StatusCode, response.Error))
		return
	}

	// 查找备份信息
	var backupInfo *mongodb.MongodbBackupInfo
	for _, item := range response.ReturnObj.List {
		if item.BackupName == data.BackupName.ValueString() {
			backupInfo = &item
			break
		}
	}

	if backupInfo == nil {
		resp.Diagnostics.AddError("未找到MongoDB备份信息", fmt.Sprintf("备份名称: %s", data.BackupName.ValueString()))
		return
	}

	// 更新数据
	data.BackupName = types.StringValue(backupInfo.BackupName)

	if backupInfo.Description != nil {
		data.Description = types.StringValue(*backupInfo.Description)
	}

	// 保存数据到Terraform状态
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CtyunMongodbBackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// 备份资源不支持更新
	var data CtyunMongodbBackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CtyunMongodbBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CtyunMongodbBackupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 删除备份
	deleteReq := &mongodb.MongodbDeleteBackupRequest{
		ProdInstId: data.InstanceID.ValueString(),
		BackupId:   data.BackupName.ValueString(),
	}

	header := &mongodb.MongodbDeleteBackupRequestHeaders{
		RegionID: data.RegionID.ValueString(),
	}

	if !data.ProjectID.IsNull() {
		header.ProjectID = data.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "开始删除MongoDB备份", map[string]interface{}{
		"backup_id": data.BackupName.ValueString(),
	})

	response, err := r.meta.Apis.SdkMongodbApis.MongodbDeleteBackupApi.Do(ctx, r.meta.Credential, deleteReq, header)
	if err != nil {
		resp.Diagnostics.AddError("删除MongoDB备份失败", err.Error())
		return
	}

	if response.StatusCode != 200 {
		resp.Diagnostics.AddError("删除MongoDB备份失败", fmt.Sprintf("API返回错误，状态码: %d, 错误信息: %s", response.StatusCode, response.Error))
		return
	}

	tflog.Info(ctx, "MongoDB备份删除成功", map[string]interface{}{
		"backup_id": data.BackupName.ValueString(),
	})
}

func (r *CtyunMongodbBackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
