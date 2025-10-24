package mongodb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &CtyunMongodbParamTemplatesDataSource{}
	_ datasource.DataSourceWithConfigure = &CtyunMongodbParamTemplatesDataSource{}
)

func NewCtyunMongodbParamTemplatesDataSource() datasource.DataSource {
	return &CtyunMongodbParamTemplatesDataSource{}
}

// CtyunMongodbParamTemplatesDataSource defines the data source implementation.
type CtyunMongodbParamTemplatesDataSource struct {
	meta *common.CtyunMetadata
}

// CtyunMongodbParamTemplatesDataSourceModel describes the data source data model.
type CtyunMongodbParamTemplatesDataSourceModel struct {
	ID            types.String                     `tfsdk:"id"`
	RegionID      types.String                     `tfsdk:"region_id"`
	ProjectID     types.String                     `tfsdk:"project_id"`
	EngineType    types.String                     `tfsdk:"engine_type"`
	EngineVersion types.String                     `tfsdk:"engine_version"`
	TemplateType  types.String                     `tfsdk:"template_type"`
	PageNo        types.Int32                      `tfsdk:"page_no"`
	PageSize      types.Int32                      `tfsdk:"page_size"`
	Templates     []CtyunMongodbParamTemplateModel `tfsdk:"templates"`
}

type CtyunMongodbParamTemplateModel struct {
	TemplateId       types.String `tfsdk:"template_id"`
	TemplateName     types.String `tfsdk:"template_name"`
	TemplateDesc     types.String `tfsdk:"template_desc"`
	EngineType       types.String `tfsdk:"engine_type"`
	EngineVersion    types.String `tfsdk:"engine_version"`
	TemplateType     types.String `tfsdk:"template_type"`
	SourceTemplateId types.String `tfsdk:"source_template_id"`
	CreatedTime      types.String `tfsdk:"created_time"`
	UpdatedTime      types.String `tfsdk:"updated_time"`
}

func (d *CtyunMongodbParamTemplatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_param_templates"
}

func (d *CtyunMongodbParamTemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "天翼云MongoDB参数组数据源",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "数据源ID",
			},
			"region_id": schema.StringAttribute{
				Required:    true,
				Description: "资源池ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "企业项目ID",
			},
			"engine_type": schema.StringAttribute{
				Optional:    true,
				Description: "引擎类型",
				Validators: []validator.String{
					stringvalidator.OneOf("mongodb"),
				},
			},
			"engine_version": schema.StringAttribute{
				Optional:    true,
				Description: "引擎版本",
			},
			"template_type": schema.StringAttribute{
				Optional:    true,
				Description: "模板类型",
				Validators: []validator.String{
					stringvalidator.OneOf("system", "user"),
				},
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "页码",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "每页条数",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
					int32validator.AtMost(100),
				},
			},
			"templates": schema.ListNestedAttribute{
				Computed:    true,
				Description: "参数组列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"template_id": schema.StringAttribute{
							Computed:    true,
							Description: "参数组ID",
						},
						"template_name": schema.StringAttribute{
							Computed:    true,
							Description: "参数组名称",
						},
						"template_desc": schema.StringAttribute{
							Computed:    true,
							Description: "参数组描述",
						},
						"engine_type": schema.StringAttribute{
							Computed:    true,
							Description: "引擎类型",
						},
						"engine_version": schema.StringAttribute{
							Computed:    true,
							Description: "引擎版本",
						},
						"template_type": schema.StringAttribute{
							Computed:    true,
							Description: "模板类型",
						},
						"source_template_id": schema.StringAttribute{
							Computed:    true,
							Description: "源参数组ID",
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
				},
			},
		},
	}
}

func (d *CtyunMongodbParamTemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	d.meta = meta
}

func (d *CtyunMongodbParamTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CtyunMongodbParamTemplatesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 设置默认分页参数
	pageNo := int32(1)
	if !data.PageNo.IsNull() {
		pageNo = data.PageNo.ValueInt32()
	}

	pageSize := int32(20)
	if !data.PageSize.IsNull() {
		pageSize = data.PageSize.ValueInt32()
	}

	// 查询参数组列表
	describeReq := &mongodb.MongodbDescribeParamTemplatesRequest{
		PageNow:  pageNo,
		PageSize: pageSize,
	}

	if !data.EngineType.IsNull() {
		describeReq.EngineType = data.EngineType.ValueStringPointer()
	}

	if !data.EngineVersion.IsNull() {
		describeReq.EngineVersion = data.EngineVersion.ValueStringPointer()
	}

	if !data.TemplateType.IsNull() {
		describeReq.TemplateType = data.TemplateType.ValueStringPointer()
	}

	header := &mongodb.MongodbDescribeParamTemplatesRequestHeaders{
		RegionID: data.RegionID.ValueString(),
	}

	if !data.ProjectID.IsNull() {
		header.ProjectID = data.ProjectID.ValueStringPointer()
	}

	response, err := d.meta.Apis.SdkMongodbApis.MongodbDescribeParamTemplatesApi.Do(ctx, d.meta.Credential, describeReq, header)
	if err != nil {
		resp.Diagnostics.AddError("查询MongoDB参数组列表失败", err.Error())
		return
	}

	if response.StatusCode != 200 {
		resp.Diagnostics.AddError("查询MongoDB参数组列表失败", fmt.Sprintf("API返回错误，状态码: %d, 错误信息: %s", response.StatusCode, response.Error))
		return
	}

	// 转换参数组信息
	var templates []CtyunMongodbParamTemplateModel
	for _, item := range response.ReturnObj.List {
		template := CtyunMongodbParamTemplateModel{
			TemplateId:    types.StringValue(item.TemplateId),
			TemplateName:  types.StringValue(item.TemplateName),
			TemplateDesc:  types.StringValue(item.TemplateDesc),
			EngineType:    types.StringValue(item.EngineType),
			EngineVersion: types.StringValue(item.EngineVersion),
			TemplateType:  types.StringValue(item.TemplateType),
			CreatedTime:   types.StringValue(item.CreatedTime),
			UpdatedTime:   types.StringValue(item.UpdatedTime),
		}

		if item.SourceTemplateId != nil {
			template.SourceTemplateId = types.StringValue(*item.SourceTemplateId)
		}

		templates = append(templates, template)
	}

	data.Templates = templates
	data.ID = types.StringValue(fmt.Sprintf("%s:%s:%s", data.RegionID.ValueString(), data.EngineType.ValueString(), data.EngineVersion.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
