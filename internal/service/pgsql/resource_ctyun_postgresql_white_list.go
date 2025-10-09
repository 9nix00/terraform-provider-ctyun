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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CtyunPgsqlWhiteList{}
	_ resource.ResourceWithConfigure   = &CtyunPgsqlWhiteList{}
	_ resource.ResourceWithImportState = &CtyunPgsqlWhiteList{}
)

type CtyunPgsqlWhiteList struct {
	meta *common.CtyunMetadata
}

func (c *CtyunPgsqlWhiteList) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_white_list"
}
func NewCtyunPgsqlWhiteList() resource.Resource {
	return &CtyunPgsqlWhiteList{}
}

func (c *CtyunPgsqlWhiteList) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPgsqlWhiteList) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {

}

func (c *CtyunPgsqlWhiteList) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
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
			"inst_id": schema.StringAttribute{
				Required:    true,
				Description: "MySQL实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
			"mode": schema.StringAttribute{
				Required:    true,
				Description: "修改模式， cover(覆盖) ， append(追加) ， delete(删除,若分组下的ip被全部删除，则会将该分组也删除，默认分组(default)则会被设置成只允许本机访问，即只有127.0.0.1这个白名单ip)",
			},
			"ip_list": schema.SetAttribute{
				Required:    true,
				Description: "ip列表,数量限制：1-1000",
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 1000),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"ip_list_result": schema.SetAttribute{
				Computed:    true,
				Description: "变更后最终的ip列表,数量限制：1-1000",
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 1000),
				},
			},
		},
	}
}

func (c *CtyunPgsqlWhiteList) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPostgresqlWhiteListConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 开始创建备份集
	err = c.updateWhiteListRequest(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMergePostgresqlWhiteList(ctx, &plan)
	if err != nil {
		return
	}
	//plan.ID = types.StringValue(plan.BackupName.ValueString())
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPgsqlWhiteList) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPostgresqlWhiteListConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergePostgresqlWhiteList(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPgsqlWhiteList) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunPgsqlWhiteList) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *CtyunPgsqlWhiteList) updateWhiteListRequest(ctx context.Context, config *CtyunPostgresqlWhiteListConfig) error {
	var ips []string
	diags := config.IpList.ElementsAs(ctx, &ips, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	params := &pgsql.PgsqlUpdateWhiteListRequest{
		ProdInstId: config.InstID.ValueString(),
		Mode:       config.Mode.ValueString(),
		IpList:     ips,
	}
	header := &pgsql.PgsqlUpdateWhiteListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlUpdateWhiteListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例添加白名单ip失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	return nil
}

func (c *CtyunPgsqlWhiteList) getAndMergePostgresqlWhiteList(ctx context.Context, config *CtyunPostgresqlWhiteListConfig) error {
	resp, err := c.getWhiteIpList(ctx, config)
	if err != nil {
		return err
	}
	ips, diags := types.SetValueFrom(ctx, types.StringType, resp.ReturnObj)
	if diags.HasError() {
		err = fmt.Errorf(diags[0].Detail())
		return err
	}
	config.IpListResult = ips
	return nil
}

func (c *CtyunPgsqlWhiteList) getWhiteIpList(ctx context.Context, config *CtyunPostgresqlWhiteListConfig) (*pgsql.PgsqlGetWhiteListResponse, error) {
	params := &pgsql.PgsqlGetWhiteListRequest{
		ProdInstId: config.InstID.ValueString(),
	}
	header := &pgsql.PgsqlGetWhiteListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetWhiteListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例获取白名单ip失败，接口返回nil，请联系研发确认问题原因！")
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	}
	return resp, nil
}

type CtyunPostgresqlWhiteListConfig struct {
	InstID       types.String `tfsdk:"inst_id"`
	RegionID     types.String `tfsdk:"region_id"`
	ProjectID    types.String `tfsdk:"project_id"`
	Mode         types.String `tfsdk:"mode"`
	IpList       types.Set    `tfsdk:"ip_list"`
	IpListResult types.Set    `tfsdk:"ip_list_result"`
}
