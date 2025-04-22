package elb

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/common"
	ctelb "terraform-provider-ctyun/internal/core/ctelb"
)

var (
	_ resource.Resource                = &CtyunElbListenerResource{}
	_ resource.ResourceWithConfigure   = &CtyunElbListenerResource{}
	_ resource.ResourceWithImportState = &CtyunElbListenerResource{}
)

type CtyunElbListenerResource struct {
	meta *common.CtyunMetadata
}

func NewCtyunElbListenerResource() resource.Resource {
	return &CtyunElbListenerResource{}
}

func (c CtyunElbListenerResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c CtyunElbListenerResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_listener"

}

func (c CtyunElbListenerResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (c CtyunElbListenerResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunElbListenerConfig
	response.Diagnostics.Append(response.State.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.CheckBeforeCreateElbListener(ctx, plan)
	if err != nil {
		return
	}

	// 创建
	err = c.CreateElbListener(ctx, plan)
	if err != nil {
		return
	}
	// 创建后反查创建后的nat信息
	err = c.getAndMergeListener(ctx, plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbListenerResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunElbListenerResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunElbListenerResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunElbListenerResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunElbListenerResource) CreateElbListener(ctx context.Context, plan CtyunElbListenerConfig) (err error) {
	//SDK ctelb_create_listener_api.go
	params := &ctelb.CtelbCreateListenerRequest{
		ClientToken:         uuid.NewString(),
		RegionID:            plan.RegionId.ValueString(),
		LoadBalancerID:      plan.LoadBalancerId.ValueString(),
		Name:                plan.Name.ValueString(),
		Protocol:            plan.Protocol.ValueString(),
		ProtocolPort:        plan.ProtocolPort.ValueInt32(),
		DefaultAction:       nil,
		AccessControlID:     "",
		AccessControlType:   "",
		ForwardedForEnabled: nil,
	}
	if !plan.Description.IsNull() {
		params.Description = plan.Description.ValueString()
	}
	if !plan.CertificateId.IsNull() {
		params.CertificateID = plan.CertificateId.ValueString()
	}
	if !plan.CaEnabled.IsNull() {
		params.CaEnabled = plan.CaEnabled.ValueBoolPointer()
	}
	if !plan.ClientCertificateId.IsNull() {
		params.ClientCertificateID = plan.ClientCertificateId.ValueString()
	}
	var defaultAction CtyunDefaultActionModel
	diags := plan.DefaultAction.ElementsAs(ctx, &defaultAction, false)
	if diags.HasError() {
		return
	}
	var forwardConfig CtyunForwardConfigModel
	as := defaultAction.ForwardConfig.ElementsAs(ctx, &forwardConfig, false)
	if as.HasError() {
		return
	}
	var targetGroups []CtyunTargetGroupModel
	elementsAs := forwardConfig.TargetGroups.ElementsAs(ctx, &targetGroups, false)
	if elementsAs.HasError() {
		return
	}
	var targetGroupsRequest []*ctelb.CtelbCreateListenerDefaultActionForwardConfigTargetGroupsRequest
	for _, targetGroupItem := range targetGroups {
		var targetGroup ctelb.CtelbCreateListenerDefaultActionForwardConfigTargetGroupsRequest
		targetGroup.TargetGroupID = targetGroupItem.TargetGroupID.ValueString()
		targetGroup.Weight = targetGroupItem.Weight.ValueInt32()
		targetGroupsRequest = append(targetGroupsRequest, &targetGroup)
	}

	params.DefaultAction = &ctelb.CtelbCreateListenerDefaultActionRequest{
		RawType: defaultAction.Type.ValueString(),
		ForwardConfig: &ctelb.CtelbCreateListenerDefaultActionForwardConfigRequest{
			TargetGroups: targetGroupsRequest,
		},
		RedirectListenerID: defaultAction.RedirectListenerID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbCreateListenerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return nil
}

func (c *CtyunElbListenerResource) getAndMergeListener(ctx context.Context, plan CtyunElbListenerConfig) (err error) {
	return nil
}

func (c *CtyunElbListenerResource) CheckBeforeCreateElbListener(ctx context.Context, plan CtyunElbListenerConfig) (err error) {
	//todo
	return nil
}

type CtyunElbListenerConfig struct {
	RegionId            types.String `cty:"region_id"`
	LoadBalancerId      types.String `cty:"loadbalancer_id"`
	Name                types.String `cty:"name"`
	Description         types.String `cty:"description"`
	Protocol            types.String `cty:"protocol"`
	ProtocolPort        types.Int32  `cty:"protocol_port"`
	CertificateId       types.String `cty:"certificate_id"`
	CaEnabled           types.Bool   `cty:"ca_enabled"`
	ClientCertificateId types.String `cty:"client_certificate_id"`
	DefaultAction       types.Set    `cty:"default_action"` //CtyunDefaultActionModel
	AccessControlId     types.String `cty:"access_control_id"`
	AccessControlType   types.String `cty:"access_control_type"`
	ForwardedForEnabled types.Bool   `cty:"forwarded_for_enabled"`
}

type DefaultActionModel struct {
	Type               types.String `tfsdk:"type"`
	ForwardConfig      types.Set    `tfsdk:"forward_config"`
	RedirectListenerID types.String `tfsdk:"redirect_listener_id"`
}
