package zos

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctzos"
	"terraform-provider-ctyun/internal/utils"
)

var (
	_ datasource.DataSource              = &ctyunZosBuckets{}
	_ datasource.DataSourceWithConfigure = &ctyunZosBuckets{}
)

type ctyunZosBuckets struct {
	meta *common.CtyunMetadata
}

func NewCtyunZosBuckets() datasource.DataSource {
	return &ctyunZosBuckets{}
}

func (c *ctyunZosBuckets) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_zos_buckets"
}

type CtyunZosBucketsModel struct {
	Bucket       types.String `tfsdk:"bucket"`
	RegionName   types.String `tfsdk:"region_name"`
	ProjectID    types.String `tfsdk:"project_id"`
	StorageType  types.String `tfsdk:"storage_type"`
	IsEncrypted  types.Bool   `tfsdk:"is_encrypted"`
	CmkUUID      types.String `tfsdk:"cmk_uuid"`
	AzPolicy     types.String `tfsdk:"az_policy"`
	CreationDate types.String `tfsdk:"creation_date"`
}

type CtyunZosBucketsConfig struct {
	RegionID  types.String `tfsdk:"region_id"`
	ProjectID types.String `tfsdk:"project_id"`
	PageNo    types.Int64  `tfsdk:"page_no"`
	PageSize  types.Int64  `tfsdk:"page_size"`

	CurrentCount types.Int64            `tfsdk:"current_count"`
	TotalCount   types.Int64            `tfsdk:"total_count"`
	Buckets      []CtyunZosBucketsModel `tfsdk:"buckets"`
}

func (c *ctyunZosBuckets) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026735/10181237**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"project_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "企业项目ID",
			},
			"page_no": schema.Int64Attribute{
				Optional:    true,
				Description: "列表的页码",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "每页数据量大小，取值1-50",
				Validators: []validator.Int64{
					int64validator.Between(1, 50),
				},
			},
			"current_count": schema.Int64Attribute{
				Computed:    true,
				Description: "分页查询时每页的行数",
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "总数",
			},
			"buckets": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"bucket": schema.StringAttribute{
							Computed:    true,
							Description: "桶名",
						},
						"region_name": schema.StringAttribute{
							Computed:    true,
							Description: "区域名称",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目ID",
						},
						"az_policy": schema.StringAttribute{
							Computed:    true,
							Description: "AZ策略，single-az或multi-az",
						},
						"storage_type": schema.StringAttribute{
							Computed:    true,
							Description: "存储类型，STANDARD、STANDARD_IA、GLACIER，分别表示标准、低频、归档，默认STANDARD",
						},
						"is_encrypted": schema.BoolAttribute{
							Computed:    true,
							Description: "是否加密",
						},
						"cmk_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "加密ID，若isEncrypted为false，此值为空字符串",
						},
						"creation_date": schema.StringAttribute{
							Computed:    true,
							Description: "创建日期，为ISO8601格式",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunZosBuckets) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunZosBucketsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	projectId := c.meta.GetExtraIfEmpty(config.ProjectID.ValueString(), common.ExtraProjectId)
	// 组装请求体
	params := &ctzos.ZosListBucketsRequest{
		RegionID: regionId,
	}
	if projectId != "" {
		params.ProjectID = projectId
	}
	pageNo := config.PageNo.ValueInt64()
	pageSize := config.PageSize.ValueInt64()
	if pageNo > 0 {
		params.PageNo = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}

	// 调用API
	resp, err := c.meta.Apis.SdkCtZosApis.ZosListBucketsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	config.Buckets = []CtyunZosBucketsModel{}
	config.TotalCount = types.Int64Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int64Value(resp.ReturnObj.CurrentCount)
	for _, b := range resp.ReturnObj.BucketList {
		item := CtyunZosBucketsModel{
			Bucket:       types.StringValue(b.Bucket),
			RegionName:   types.StringValue(b.RegionName),
			ProjectID:    types.StringValue(b.ProjectID),
			StorageType:  types.StringValue(b.StorageType),
			IsEncrypted:  utils.SecBoolValue(b.IsEncrypted),
			CmkUUID:      types.StringValue(b.CmkUUID),
			AzPolicy:     types.StringValue(b.AZPolicy),
			CreationDate: types.StringValue(b.CreationDate),
		}
		config.Buckets = append(config.Buckets, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunZosBuckets) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
