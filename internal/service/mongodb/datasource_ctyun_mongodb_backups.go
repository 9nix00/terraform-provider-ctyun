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
	_ datasource.DataSource              = &CtyunMongodbBackupsDataSource{}
	_ datasource.DataSourceWithConfigure = &CtyunMongodbBackupsDataSource{}
)

func NewCtyunMongodbBackupsDataSource() datasource.DataSource {
	return &CtyunMongodbBackupsDataSource{}
}

// CtyunMongodbBackupsDataSource defines the data source implementation.
type CtyunMongodbBackupsDataSource struct {
	meta *common.CtyunMetadata
}

// CtyunMongodbBackupsDataSourceModel describes the data source data model.
type CtyunMongodbBackupsDataSourceModel struct {
	ID         types.String              `tfsdk:"id"`
	RegionID   types.String              `tfsdk:"region_id"`
	ProjectID  types.String              `tfsdk:"project_id"`
	InstanceID types.String              `tfsdk:"instance_id"`
	BackupID   types.String              `tfsdk:"backup_id"`
	BackupType types.String              `tfsdk:"backup_type"`
	PageNo     types.Int32               `tfsdk:"page_no"`
	PageSize   types.Int32               `tfsdk:"page_size"`
	Backups    []CtyunMongodbBackupModel `tfsdk:"backups"`
}

type CtyunMongodbBackupModel struct {
	BackupID          types.String `tfsdk:"backup_id"`
	BackupName        types.String `tfsdk:"backup_name"`
	BackupMethod      types.String `tfsdk:"backup_method"`
	BackupType        types.String `tfsdk:"backup_type"`
	BackupStatus      types.String `tfsdk:"backup_status"`
	BackupStartTime   types.String `tfsdk:"backup_start_time"`
	BackupEndTime     types.String `tfsdk:"backup_end_time"`
	BackupSize        types.Int64  `tfsdk:"backup_size"`
	Description       types.String `tfsdk:"description"`
	BackupTriggerType types.String `tfsdk:"backup_trigger_type"`
}

func (d *CtyunMongodbBackupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_backups"
}

func (d *CtyunMongodbBackupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "天翼云MongoDB备份数据源",

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
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID",
			},
			"backup_id": schema.StringAttribute{
				Optional:    true,
				Description: "备份ID",
			},
			"backup_type": schema.StringAttribute{
				Optional:    true,
				Description: "备份类型",
				Validators: []validator.String{
					stringvalidator.OneOf("full", "incremental"),
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
			"backups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "备份列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"backup_id": schema.StringAttribute{
							Computed:    true,
							Description: "备份ID",
						},
						"backup_name": schema.StringAttribute{
							Computed:    true,
							Description: "备份名称",
						},
						"backup_method": schema.StringAttribute{
							Computed:    true,
							Description: "备份方式",
						},
						"backup_type": schema.StringAttribute{
							Computed:    true,
							Description: "备份类型",
						},
						"backup_status": schema.StringAttribute{
							Computed:    true,
							Description: "备份状态",
						},
						"backup_start_time": schema.StringAttribute{
							Computed:    true,
							Description: "备份开始时间",
						},
						"backup_end_time": schema.StringAttribute{
							Computed:    true,
							Description: "备份结束时间",
						},
						"backup_size": schema.Int64Attribute{
							Computed:    true,
							Description: "备份大小",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "备份描述",
						},
						"backup_trigger_type": schema.StringAttribute{
							Computed:    true,
							Description: "备份触发类型",
						},
					},
				},
			},
		},
	}
}

func (d *CtyunMongodbBackupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	d.meta = meta
}

func (d *CtyunMongodbBackupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CtyunMongodbBackupsDataSourceModel

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

	// 查询备份列表
	describeReq := &mongodb.MongodbDescribeBackupsRequest{
		ProdInstId: data.InstanceID.ValueString(),
		PageNow:    pageNo,
		PageSize:   pageSize,
	}

	//if !data.BackupID.IsNull() {
	//	describeReq.BackupId = data.BackupID.ValueStringPointer()
	//}

	if !data.BackupType.IsNull() {
		describeReq.BackupType = data.BackupType.ValueStringPointer()
	}

	header := &mongodb.MongodbDescribeBackupsRequestHeaders{
		RegionID: data.RegionID.ValueString(),
	}

	if !data.ProjectID.IsNull() {
		header.ProjectID = data.ProjectID.ValueStringPointer()
	}

	response, err := d.meta.Apis.SdkMongodbApis.MongodbDescribeBackupsApi.Do(ctx, d.meta.Credential, describeReq, header)
	if err != nil {
		resp.Diagnostics.AddError("查询MongoDB备份列表失败", err.Error())
		return
	}

	if response.StatusCode != 200 {
		resp.Diagnostics.AddError("查询MongoDB备份列表失败", fmt.Sprintf("API返回错误，状态码: %d, 错误信息: %s", response.StatusCode, response.Error))
		return
	}

	// 转换备份信息
	var backups []CtyunMongodbBackupModel
	for _, item := range response.ReturnObj.List {
		backup := CtyunMongodbBackupModel{
			BackupID:          types.StringValue(item.BackupId),
			BackupName:        types.StringValue(item.BackupName),
			BackupMethod:      types.StringValue(item.BackupMethod),
			BackupType:        types.StringValue(item.BackupType),
			BackupStatus:      types.StringValue(item.BackupStatus),
			BackupStartTime:   types.StringValue(item.BackupStartTime),
			BackupEndTime:     types.StringValue(item.BackupEndTime),
			BackupSize:        types.Int64Value(item.BackupSize),
			BackupTriggerType: types.StringValue(item.BackupTriggerType),
		}

		if item.Description != nil {
			backup.Description = types.StringValue(*item.Description)
		}

		backups = append(backups, backup)
	}

	data.Backups = backups
	data.ID = types.StringValue(fmt.Sprintf("%s:%s", data.RegionID.ValueString(), data.InstanceID.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
