package pgsql

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
)

var (
	_ resource.Resource                = &CtyunPgsqlAssociationEip{}
	_ resource.ResourceWithConfigure   = &CtyunPgsqlAssociationEip{}
	_ resource.ResourceWithImportState = &CtyunPgsqlAssociationEip{}
)

type CtyunPgsqlAssociationEip struct {
	meta *common.CtyunMetadata
}

func (c *CtyunPgsqlAssociationEip) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_association_eip"
}
func NewCtyunMysqlAssociationEip() resource.Resource {
	return &CtyunPgsqlAssociationEip{}
}

func (c *CtyunPgsqlAssociationEip) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"eip_id": schema.StringAttribute{
				Required:    true,
				Description: "弹性id",
			},
			"eip": schema.StringAttribute{
				Required:    true,
				Description: "弹性ip",
			},
			"inst_id": schema.StringAttribute{
				Required:    true,
				Description: "实例id",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "项目id",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池Id",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"eip_status": schema.Int32Attribute{
				Computed:    true,
				Description: " 弹性ip状态 0->unbind，1->bind,2->binding",
				Validators: []validator.Int32{
					int32validator.Between(0, 2),
				},
			},
		},
	}
}

func (c *CtyunPgsqlAssociationEip) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunPgsqlAssociationEipConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 实例绑定弹性IP
	err = c.MysqlBindEip(ctx, &plan)
	if err != nil {
		return
	}
}

func (c *CtyunPgsqlAssociationEip) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
}

func (c *CtyunPgsqlAssociationEip) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (c *CtyunPgsqlAssociationEip) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunPgsqlAssociationEipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	unbindParams := &pgsql.PgsqlUnBindEipRequest{
		EipID:  state.EipID.ValueString(),
		Eip:    state.Eip.ValueString(),
		InstID: state.InstID.ValueString(),
	}
	unbindHeader := &pgsql.PgsqlUnBindEipRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		unbindHeader.ProjectId = state.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlUnBindEipApi.Do(ctx, c.meta.Credential, unbindParams, unbindHeader)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
}

func (c *CtyunPgsqlAssociationEip) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
func (c *CtyunPgsqlAssociationEip) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunPgsqlAssociationEip) MysqlBindEip(ctx context.Context, config *CtyunPgsqlAssociationEipConfig) (err error) {
	params := &pgsql.PgsqlBindEipRequest{
		EipID:  config.EipID.ValueString(),
		Eip:    config.Eip.ValueString(),
		InstID: config.InstID.ValueString(),
	}
	header := &pgsql.PgsqlBindEipRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		header.ProjectId = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlBindEipApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s ", *resp.Message)
		return
	}
	return
}

type CtyunPgsqlAssociationEipConfig struct {
	EipID     types.String `tfsdk:"eip_id"`     //弹性id
	Eip       types.String `tfsdk:"eip"`        //弹性ip
	InstID    types.String `tfsdk:"inst_id"`    //实例id
	ProjectID types.String `tfsdk:"project_id"` //项目id
	RegionID  types.String `tfsdk:"region_id"`  //区域Id
	EipStatus types.Int32  `tfsdk:"eip_status"` //弹性ip状态 0->unbind，1->bind,2->binding
}
