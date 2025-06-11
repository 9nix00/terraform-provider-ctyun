package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
)

var (
	_ datasource.DataSource              = &CtyunPgsqlSpecs{}
	_ datasource.DataSourceWithConfigure = &CtyunPgsqlSpecs{}
)

type CtyunPgsqlSpecs struct {
	meta *common.CtyunMetadata
}

func NewCtyunPgsqlSpecs() *CtyunPgsqlSpecs {
	return &CtyunPgsqlSpecs{}
}
func (c *CtyunPgsqlSpecs) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPgsqlSpecs) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_specs"
}

func (c *CtyunPgsqlSpecs) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id",
			},
			"inst_id": schema.StringAttribute{
				Optional:    true,
				Description: "实例id",
			},
			"eip_id": schema.StringAttribute{
				Optional:    true,
				Description: "弹性ip唯一标识",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "项目id",
			},
			"eips": schema.ListNestedAttribute{
				Description: "eip 列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"eip_id": schema.StringAttribute{
							Computed:    true,
							Description: "弹性ip唯一标识",
						},
						"eip": schema.StringAttribute{
							Computed:    true,
							Description: "弹性IP",
						},
						"bind_status": schema.Int32Attribute{
							Computed:    true,
							Description: "0-未绑定，1-已绑定",
							Validators: []validator.Int32{
								int32validator.Between(0, 1),
							},
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态标识：ACTIVE=已使用，DOWN=未使用，ERROR=中间状态-异常，UPDATING=中间状态-更新中，BANDING_OR_UNBANGDING=中间状态-绑定或解绑中，DELETING=中间状态-删除中，DELETED=中间状态-已删除",
							Validators: []validator.String{
								stringvalidator.OneOf(business.PgsqlBindEipStatus...),
							},
						},
						"band_width": schema.Int32Attribute{
							Computed:    true,
							Description: "加入的共享带宽，单位M",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunPgsqlSpecs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPgsqlSpecsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	params := &pgsql.PgsqlBoundEipListRequest{
		RegionID: regionId,
	}
	if config.InstID.ValueString() != "" {
		params.InstID = config.InstID.ValueStringPointer()
	}
	if config.EipID.ValueString() != "" {
		params.EipID = config.EipID.ValueStringPointer()
	}
	header := &pgsql.PgsqlBoundEipListRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlBoundEipListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s ", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	returnObj := resp.ReturnObj.Data
	// 解析返回类型
	eips := []CtyunPgsqlSpecInfoModel{}
	for _, eipItem := range returnObj {
		var eip CtyunPgsqlSpecInfoModel
		eip.EipID = types.StringValue(eipItem.EipID)
		eip.Eip = types.StringValue(eipItem.Eip)
		eip.BindStatus = types.Int32Value(eipItem.BindStatus)
		eip.Status = types.StringValue(eipItem.Status)
		eip.BandWidth = types.Int32Value(eipItem.BandWidth)
		eips = append(eips, eip)
	}
	config.Eips = eips
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunPgsqlSpecsConfig struct {
	RegionID  types.String              `tfsdk:"region_id"`  //资源池id
	InstID    types.String              `tfsdk:"inst_id"`    //实例id
	EipID     types.String              `tfsdk:"eip_id"`     //实例id
	ProjectID types.String              `tfsdk:"project_id"` //项目id
	Eips      []CtyunPgsqlSpecInfoModel `tfsdk:"eips"`       // eip列表
}
type CtyunPgsqlSpecInfoModel struct {
	EipID      types.String `tfsdk:"eip_id"`      // 弹性ip唯一标识
	Eip        types.String `tfsdk:"eip"`         // 弹性IP
	BindStatus types.Int32  `tfsdk:"bind_status"` // 0-未绑定，1-已绑定
	Status     types.String `tfsdk:"status"`      // 状态标识：ACTIVE=已使用，DOWN=未使用，ERROR=中间状态-异常，UPDATING=中间状态-更新中，BANDING_OR_UNBANGDING=中间状态-绑定或解绑中，DELETING=中间状态-删除中，DELETED=中间状态-已删除
	BandWidth  types.Int32  `tfsdk:"band_width"`  // 加入的共享带宽，单位M
}
