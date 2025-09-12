package mongodb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &CtyunMongodbParamTemplateResource{}
	_ resource.ResourceWithConfigure   = &CtyunMongodbParamTemplateResource{}
	_ resource.ResourceWithImportState = &CtyunMongodbParamTemplateResource{}
)

func NewCtyunMongodbParamTemplateResource() resource.Resource {
	return &CtyunMongodbParamTemplateResource{}
}

// CtyunMongodbParamTemplateResource defines the resource implementation.
type CtyunMongodbParamTemplateResource struct {
	meta *common.CtyunMetadata
}

// CtyunMongodbParamTemplateResourceModel describes the resource data model.
type CtyunMongodbParamTemplateResourceModel struct {
	ID            types.String `tfsdk:"id"`
	RegionID      types.String `tfsdk:"region_id"`
	ProjectID     types.String `tfsdk:"project_id"`
	TemplateName  types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	NodeType      types.String `tfsdk:"node_type"`
	EngineVersion types.String `tfsdk:"engine_version"`
	//TemplateType     types.String `tfsdk:"template_type"`
	//SourceTemplateId types.String `tfsdk:"source_template_id"`
	TemplateId  types.String `tfsdk:"template_id"`
	CreatedTime types.String `tfsdk:"created_time"`
	UpdatedTime types.String `tfsdk:"updated_time"`
}

func (r *CtyunMongodbParamTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_param_template"
}

func (r *CtyunMongodbParamTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "天翼云MongoDB参数组资源",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源ID，格式为region_id:template_id",
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "参数组名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "参数组描述",
			},
			"node_type": schema.StringAttribute{
				Required:    true,
				Description: "引擎类型",
				Validators: []validator.String{
					stringvalidator.OneOf("Mongod", "Shard", "Mongos", "Configsvr"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"engine_version": schema.StringAttribute{
				Required:    true,
				Description: "引擎版本",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			//"template_type": schema.StringAttribute{
			//	Required:    true,
			//	Description: "模板类型",
			//	Validators: []validator.String{
			//		stringvalidator.OneOf("system", "user"),
			//	},
			//	PlanModifiers: []planmodifier.String{
			//		stringplanmodifier.RequiresReplace(),
			//	},
			//},
			//"source_template_id": schema.StringAttribute{
			//	Optional:    true,
			//	Description: "源参数组ID",
			//	PlanModifiers: []planmodifier.String{
			//		stringplanmodifier.RequiresReplace(),
			//	},
			//},
			"template_id": schema.StringAttribute{
				Computed:    true,
				Description: "参数组ID",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间",
			},
			"updated_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间",
			},
		},
	}
}

func (r *CtyunMongodbParamTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	r.meta = meta
}

func (r *CtyunMongodbParamTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CtyunMongodbParamTemplateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.createParamTemplate(ctx, data, resp) {
		return
	}

	// 保存数据到Terraform状态
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CtyunMongodbParamTemplateResource) createParamTemplate(ctx context.Context, data CtyunMongodbParamTemplateResourceModel, resp *resource.CreateResponse) bool {
	// Create param template
	createReq := &mongodb.MongodbCreateParamTemplateRequest{
		TemplateName:  data.TemplateName.ValueString(),
		EngineVersion: data.EngineVersion.ValueString(),
		NodeType:      data.NodeType.ValueString(),
	}

	if !data.Description.IsNull() {
		createReq.TemplateDesc = data.Description.ValueString()
	}

	header := &mongodb.MongodbCreateParamTemplateRequestHeaders{
		RegionID: data.RegionID.ValueString(),
	}

	if !data.ProjectID.IsNull() {
		header.ProjectID = data.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "开始创建MongoDB参数组", map[string]interface{}{
		"template_name": data.TemplateName.ValueString(),
	})

	response, err := r.meta.Apis.SdkMongodbApis.MongodbCreateParamTemplateApi.Do(ctx, r.meta.Credential, createReq, header)
	if err != nil {
		resp.Diagnostics.AddError("创建MongoDB参数组失败", err.Error())
		return true
	}

	if response.StatusCode != 200 {
		resp.Diagnostics.AddError("创建MongoDB参数组失败", fmt.Sprintf("API返回错误，状态码: %d, 错误信息: %s", response.StatusCode, response.Error))
		return true
	}

	// 设置参数组ID
	templateId := *response.ReturnObj
	data.TemplateId = types.StringValue(templateId)

	// 设置资源ID
	data.ID = types.StringValue(templateId)
	return false
}

func (r *CtyunMongodbParamTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CtyunMongodbParamTemplateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 获取参数组信息
	describeReq := &mongodb.MongodbDescribeParamTemplatesRequest{
		PageNow:  1,
		PageSize: 100,
	}

	header := &mongodb.MongodbDescribeParamTemplatesRequestHeaders{
		RegionID: data.RegionID.ValueString(),
	}

	if !data.ProjectID.IsNull() {
		header.ProjectID = data.ProjectID.ValueStringPointer()
	}

	response, err := r.meta.Apis.SdkMongodbApis.MongodbDescribeParamTemplatesApi.Do(ctx, r.meta.Credential, describeReq, header)
	if err != nil {
		resp.Diagnostics.AddError("查询MongoDB参数组信息失败", err.Error())
		return
	}

	if response.StatusCode != 200 {
		resp.Diagnostics.AddError("查询MongoDB参数组信息失败", fmt.Sprintf("API返回错误，状态码: %d, 错误信息: %s", response.StatusCode, response.Error))
		return
	}

	// 查找参数组信息
	var templateInfo *mongodb.MongodbParamTemplateInfo
	for _, item := range response.ReturnObj.List {
		if item.TemplateId == data.TemplateId.ValueString() {
			templateInfo = &item
			break
		}
	}

	if templateInfo == nil {
		resp.Diagnostics.AddError("未找到MongoDB参数组信息", fmt.Sprintf("参数组ID: %s", data.TemplateId.ValueString()))
		return
	}

	// 更新数据
	data.TemplateName = types.StringValue(templateInfo.TemplateName)
	data.Description = types.StringValue(templateInfo.TemplateDesc)

	data.EngineVersion = types.StringValue(templateInfo.EngineVersion)

	data.CreatedTime = types.StringValue(templateInfo.CreatedTime)
	data.UpdatedTime = types.StringValue(templateInfo.UpdatedTime)

	// 保存数据到Terraform状态
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CtyunMongodbParamTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CtyunMongodbParamTemplateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 只有描述可以更新
	if !data.Description.Equal(data.Description) {
		// 更新参数组描述
		updateReq := &mongodb.MongodbUpdateParamTemplateDescRequest{
			TemplateId:   data.TemplateId.ValueString(),
			TemplateDesc: data.Description.ValueString(),
		}

		header := &mongodb.MongodbUpdateParamTemplateDescRequestHeaders{
			RegionID: data.RegionID.ValueString(),
		}

		if !data.ProjectID.IsNull() {
			header.ProjectID = data.ProjectID.ValueStringPointer()
		}

		tflog.Info(ctx, "开始更新MongoDB参数组描述", map[string]interface{}{
			"template_id": data.TemplateId.ValueString(),
		})

		response, err := r.meta.Apis.SdkMongodbApis.MongodbUpdateParamTemplateDescApi.Do(ctx, r.meta.Credential, updateReq, header)
		if err != nil {
			resp.Diagnostics.AddError("更新MongoDB参数组描述失败", err.Error())
			return
		}

		if response.StatusCode != 200 {
			resp.Diagnostics.AddError("更新MongoDB参数组描述失败", fmt.Sprintf("API返回错误，状态码: %d, 错误信息: %s", response.StatusCode, response.Error))
			return
		}
	}

	// 保存数据到Terraform状态
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CtyunMongodbParamTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CtyunMongodbParamTemplateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 删除参数组
	deleteReq := &mongodb.MongodbDeleteParamTemplateRequest{
		//TemplateId: data.TemplateId.ValueString(),
	}

	header := &mongodb.MongodbDeleteParamTemplateRequestHeaders{
		RegionID: data.RegionID.ValueString(),
	}

	if !data.ProjectID.IsNull() {
		header.ProjectID = data.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "开始删除MongoDB参数组", map[string]interface{}{
		"template_id": data.TemplateId.ValueString(),
	})

	response, err := r.meta.Apis.SdkMongodbApis.MongodbDeleteParamTemplateApi.Do(ctx, r.meta.Credential, deleteReq, header)
	if err != nil {
		resp.Diagnostics.AddError("删除MongoDB参数组失败", err.Error())
		return
	}

	if response.StatusCode != 200 {
		resp.Diagnostics.AddError("删除MongoDB参数组失败", fmt.Sprintf("API返回错误，状态码: %d, 错误信息: %s", response.StatusCode, response.Error))
		return
	}

	tflog.Info(ctx, "MongoDB参数组删除成功", map[string]interface{}{
		"template_id": data.TemplateId.ValueString(),
	})
}

func (r *CtyunMongodbParamTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
