package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunPostgresqlCharacterSet{}
	_ datasource.DataSourceWithConfigure = &ctyunPostgresqlCharacterSet{}
)

type ctyunPostgresqlCharacterSet struct {
	meta *common.CtyunMetadata
}

func NewCtyunPostgresqlCharacterSet() datasource.DataSource {
	return &ctyunPostgresqlCharacterSet{}
}
func (c *ctyunPostgresqlCharacterSet) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunPostgresqlCharacterSet) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_character_set"
}

func (c *ctyunPostgresqlCharacterSet) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，默认使用provider配置",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "项目ID",
			},
			"character_set": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "字符集列表",
			},
		},
	}
}

func (c *ctyunPostgresqlCharacterSet) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPostgresqlCharacterSetConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}

	params := &pgsql.PgsqlGetCharacterSetRequest{
		Engine: "postgreSQL",
	}
	header := &pgsql.PgsqlGetCharacterSetRequestHeader{
		RegionID: regionId,
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetCharacterSetApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询postgresql实例数据库支持的字符集失败，接口返回nil，请联系研发确认问题原因！")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	characterSetResp := resp.ReturnObj.CharacterSetNameItems
	characterSet, diagnostics := types.SetValueFrom(ctx, types.StringType, characterSetResp)
	if diagnostics.HasError() {
		err = fmt.Errorf(diagnostics[0].Detail())
		return
	}
	config.CharacterSet = characterSet
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunPostgresqlCharacterSetConfig struct {
	RegionID     types.String `tfsdk:"region_id"`
	ProjectID    types.String `tfsdk:"project_id"`
	CharacterSet types.Set    `tfsdk:"character_set"`
}
