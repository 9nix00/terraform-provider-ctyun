package ec

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CtyunEcCloudGatewaySdwanAssociation{}
	_ resource.ResourceWithConfigure   = &CtyunEcCloudGatewaySdwanAssociation{}
	_ resource.ResourceWithImportState = &CtyunEcCloudGatewaySdwanAssociation{}
)

func NewCtyunEcCloudGatewaySdwanAssociation() resource.Resource {
	return &CtyunEcCloudGatewaySdwanAssociation{}
}

type CtyunEcCloudGatewaySdwanAssociation struct {
	meta *common.CtyunMetadata
}

type CtyunEcCloudGatewaySdwanAssociationConfig struct {
	ID      types.String `tfsdk:"id"`
	EcID    types.String `tfsdk:"ec_id"`
	SdwanID types.String `tfsdk:"sdwan_id"`
	CgwList []CgwInfo    `tfsdk:"cgw_list"`
}

type CgwInfo struct {
	RtbID types.String `tfsdk:"rtb_id"`
	CgwID types.String `tfsdk:"cgw_id"`
}

func (c *CtyunEcCloudGatewaySdwanAssociation) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ec_cloud_gateway_sdwan_association"
}

func (c *CtyunEcCloudGatewaySdwanAssociation) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10038220`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识符，格式为 {ec_id}/{sdwan_id}",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"sdwan_id": schema.StringAttribute{
				Required:    true,
				Description: "SDWAN ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cgw_list": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "需要绑定的云网关列表，如果全部解绑则传空列表 支持更新",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rtb_id": schema.StringAttribute{
							Required:    true,
							Description: "云网关默认路由表ID 支持更新",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"cgw_id": schema.StringAttribute{
							Required:    true,
							Description: "云网关ID 支持更新",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
		},
	}
}

func (c *CtyunEcCloudGatewaySdwanAssociation) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunEcCloudGatewaySdwanAssociation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunEcCloudGatewaySdwanAssociationConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}

	// 设置ID为 ec_id/sdwan_id 的格式
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s", plan.EcID.ValueString(), plan.SdwanID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunEcCloudGatewaySdwanAssociation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunEcCloudGatewaySdwanAssociationConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 对于关联资源，我们无法直接读取关联状态，因此只验证资源是否存在
	// 在实际应用中，可能需要通过其他API查询绑定状态
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunEcCloudGatewaySdwanAssociation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunEcCloudGatewaySdwanAssociationConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}

	// 设置ID为 ec_id/sdwan_id 的格式
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s", plan.EcID.ValueString(), plan.SdwanID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunEcCloudGatewaySdwanAssociation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunEcCloudGatewaySdwanAssociationConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 解绑操作，传入空的 cgwList
	state.CgwList = []CgwInfo{}
	err = c.create(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunEcCloudGatewaySdwanAssociation) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// 暂不支持导入
	resp.Diagnostics.AddError(
		"Import not supported",
		"This resource does not support import.",
	)
}

func (c *CtyunEcCloudGatewaySdwanAssociation) create(ctx context.Context, plan *CtyunEcCloudGatewaySdwanAssociationConfig) (err error) {
	// 构造请求参数
	req := &ec.EcEcBindCloudGatewayRequest{
		EcID:    plan.EcID.ValueString(),
		SdwanID: plan.SdwanID.ValueString(),
	}

	// 添加云网关列表
	cgwList := make([]*ec.EcEcBindCloudGatewayCgwListRequest, len(plan.CgwList))
	for i, cgw := range plan.CgwList {
		cgwList[i] = &ec.EcEcBindCloudGatewayCgwListRequest{
			RtbID: cgw.RtbID.ValueString(),
			CgwID: cgw.CgwID.ValueString(),
		}
	}
	req.CgwList = cgwList

	resp, err := c.meta.Apis.SdkEcApis.EcEcBindCloudGatewayApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp == nil {
		return fmt.Errorf("API return error. StatusCode is nil")
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return nil
}
