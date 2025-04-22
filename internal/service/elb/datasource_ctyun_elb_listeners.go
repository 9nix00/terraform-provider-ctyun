package elb

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/common"
	ctelb "terraform-provider-ctyun/internal/core/ctelb"
)

var (
	_ datasource.DataSource              = &ctyunElbListeners{}
	_ datasource.DataSourceWithConfigure = &ctyunElbListeners{}
)

type ctyunElbListeners struct {
	meta *common.CtyunMetadata
}

func NewElbListeners() datasource.DataSource {
	return &ctyunElbListeners{}
}

func (c ctyunElbListeners) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c ctyunElbListeners) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_listeners"
}

func (c ctyunElbListeners) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (c ctyunElbListeners) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	//SDK hhttps://eop.ctyun.cn/ebp/ctapiDocument/search?sid=24&api=5654&data=88&isNormal=1&vid=82
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunElbListenersConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	params := &ctelb.CtelbListListenerRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
	}

	if !config.ProjectID.IsNull() {
		params.ProjectID = config.ProjectID.ValueString()
	}
	if !config.IDs.IsNull() {
		params.IDs = config.IDs.ValueString()
	}
	if !config.Name.IsNull() {
		params.Name = config.Name.ValueString()
	}
	if !config.LoadBalancerID.IsNull() {
		params.LoadBalancerID = config.LoadBalancerID.ValueString()
	}
	if !config.AccessControlID.IsNull() {
		params.AccessControlID = config.AccessControlID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbListListenerApi.Do(ctx, c.meta.SdkCredential, params)
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
	var listeners []CtyunElbListenersDetailModel
	for _, listenerItem := range resp.ReturnObj {
		var listener CtyunElbListenersDetailModel
		listener.RegionID = types.StringValue(listenerItem.RegionID)
		listener.AzName = types.StringValue(listenerItem.AzName)
		listener.ID = types.StringValue(listenerItem.ID)
		listener.Name = types.StringValue(listenerItem.Name)
		listener.Description = types.StringValue(listenerItem.Description)
		listener.LoadBalancerID = types.StringValue(listenerItem.LoadBalancerID)
		listener.Protocol = types.StringValue(listenerItem.Protocol)
		listener.ProtocolPort = types.Int32Value(listenerItem.ProtocolPort)
		listener.CertificateID = types.StringValue(listenerItem.CertificateID)
		listener.AccessControlID = types.StringValue(listenerItem.AccessControlID)
		listener.AccessControlType = types.StringValue(listenerItem.AccessControlType)
		listener.ForwardedForEnabled = types.BoolValue(*listenerItem.ForwardedForEnabled)
		listener.Status = types.StringValue(listenerItem.Status)
		listener.CreatedTime = types.StringValue(listenerItem.CreatedTime)
		listener.UpdatedTime = types.StringValue(listenerItem.UpdatedTime)
		// 处理defaultAction
		var defaultActions []CtyunDefaultActionModel
		if listenerItem.DefaultAction != nil && len(listenerItem.DefaultAction) > 0 {
			for _, defaultActionItem := range listener.DefaultAction {
				var defaultAction CtyunDefaultActionModel
				defaultAction.Type = defaultActionItem.Type
				defaultAction.ForwardConfig = defaultActionItem.ForwardConfig
				defaultAction.RedirectListenerID = defaultActionItem.RedirectListenerID
				defaultActions = append(defaultActions, defaultAction)
			}
		}
		listener.DefaultAction = defaultActions

		listeners = append(listeners, listener)
	}
	config.Listeners = listeners
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunElbListenersConfig struct {
	RegionID        types.String                   `tfsdk:"region_id"`
	ProjectID       types.String                   `tfsdk:"project_id"`
	IDs             types.String                   `tfsdk:"ids"`
	Name            types.String                   `tfsdk:"name"`
	LoadBalancerID  types.String                   `tfsdk:"load_balancer_id"`
	AccessControlID types.String                   `tfsdk:"access_control_id"`
	Listeners       []CtyunElbListenersDetailModel `tfsdk:"listeners"`
}

type CtyunElbListenersDetailModel struct {
	RegionID            types.String              `tfsdk:"region_id"`
	AzName              types.String              `tfsdk:"az_name"`
	ProjectID           types.String              `tfsdk:"project_id"`
	ID                  types.String              `tfsdk:"id"`
	Name                types.String              `tfsdk:"name"`
	Description         types.String              `tfsdk:"description"`
	LoadBalancerID      types.String              `tfsdk:"load_balancer_id"`
	Protocol            types.String              `tfsdk:"protocol"`
	ProtocolPort        types.Int32               `tfsdk:"protocol_port"`
	CertificateID       types.String              `tfsdk:"certificate_id"`
	CaEnabled           types.Bool                `tfsdk:"ca_enabled"` //是否开启双向认证
	ClientCertificateID types.String              `tfsdk:"client_certificate_id"`
	DefaultAction       []CtyunDefaultActionModel `tfsdk:"default_action"`
	AccessControlID     types.String              `tfsdk:"access_control_id"`
	AccessControlType   types.String              `tfsdk:"access_control_type"`
	ForwardedForEnabled types.Bool                `tfsdk:"forwarded_for_enabled"`
	Status              types.String              `tfsdk:"status"`
	CreatedTime         types.String              `tfsdk:"created_time"`
	UpdatedTime         types.String              `tfsdk:"updated_time"`
}

type CtyunDefaultActionModel struct {
	Type               types.String `tfsdk:"type"`
	ForwardConfig      types.List   `tfsdk:"forward_config"`
	RedirectListenerID types.String `tfsdk:"redirect_listener_id"`
}

type CtyunForwardConfigModel struct {
	TargetGroups types.List `cty:"target_groups"`
}

type CtyunTargetGroupModel struct {
	TargetGroupID types.String `cty:"target_group_id"`
	Weight        types.Int32  `cty:"weight"`
}
