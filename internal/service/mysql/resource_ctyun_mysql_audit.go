package mysql

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CtyunMysqlAudit{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlAudit{}
	_ resource.ResourceWithImportState = &CtyunMysqlAudit{}
)

type CtyunMysqlAudit struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMysqlAudit) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_audit"
}
func NewCtyunMysqlAudit() resource.Resource {
	return &CtyunMysqlAudit{}
}

func (c *CtyunMysqlAudit) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMysqlAudit) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunMysqlAuditConfig
	var ID, regionId, projectId, instId string
	err = terraform_extend.Split(request.ID, &ID, &regionId, &projectId, &instId)
	if err != nil {
		return
	}

	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *CtyunMysqlAudit) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10133568",
		Attributes: map[string]schema.Attribute{
			"inst_id": schema.StringAttribute{
				Computed:    true,
				Description: "mysql实例Id",
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
			"audit_switch": schema.BoolAttribute{
				Required:    true,
				Description: "sql审计开关状态。false：关闭状态；true：开启状态。",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (c *CtyunMysqlAudit) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMysqlAuditConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 开启/关闭sql审计

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlAudit) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	return
}

func (c *CtyunMysqlAudit) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunMysqlAudit) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *CtyunMysqlAudit) create(ctx context.Context, config *CtyunMysqlAuditConfig) error {
	params := &mysql.TeledbStartAuditRequest{
		OuterProdInstId: config.InstID.ValueString(),
		AuditSwitch:     config.AuditSwitch.ValueBool(),
	}
	header := &mysql.TeledbStartAuditRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbStartAuditApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("开始/关闭数据库（id=%s）SQL审计失败，接口返回nil", config.InstID.ValueString())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	status, err := c.getMysqlAuditStatus(ctx, config)
	if err != nil || status == nil {
		return err
	}

	if *status != config.AuditSwitch.ValueBool() {
		err = fmt.Errorf("开始/关闭数据库（id=%s）SQL审计失败，当前状态为：%t", config.InstID.ValueString(), *status)
		return nil
	}
	return nil
}

func (c *CtyunMysqlAudit) getAndMerge(ctx context.Context, config *CtyunMysqlAuditConfig) error {
	return nil
}

func (c *CtyunMysqlAudit) getMysqlAuditStatus(ctx context.Context, config *CtyunMysqlAuditConfig) (*bool, error) {
	params := &mysql.TeledbGetAuditStatusRequest{
		OuterProdInstId: config.InstID.ValueString(),
	}
	header := &mysql.TeledbGetAuditStatusRequestHeader{
		RegionID: config.RegionID.ValueString(),
		InstID:   config.InstID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetAuditStatusApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取数据库（id=%s）SQL审计状态失败，接口返回nil", config.InstID.ValueString())
		return nil, err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return &resp.ReturnObj.AuditLogSwitch, nil
}

type CtyunMysqlAuditConfig struct {
	InstID      types.String `tfsdk:"inst_id"`
	ProjectID   types.String `tfsdk:"project_id"`
	RegionID    types.String `tfsdk:"region_id"`
	AuditSwitch types.Bool   `tfsdk:"audit_switch"`
}
