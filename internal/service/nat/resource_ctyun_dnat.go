package nat

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
	"terraform-provider-ctyun/internal/utils"
	"time"
)

var (
	_ resource.Resource                = &ctyunDnatResource{}
	_ resource.ResourceWithConfigure   = &ctyunDnatResource{}
	_ resource.ResourceWithImportState = &ctyunDnatResource{}
)

type ctyunDnatResource struct {
	meta *common.CtyunMetadata
}

func (c *ctyunDnatResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_nat_dnat"
}

func NewCtyunDnatResource() resource.Resource {
	return &ctyunDnatResource{}
}

func (c *ctyunDnatResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `详细说明请见文档：`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"nat_gateway_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "NAT网关Id",
			},
			"dnat_id": schema.StringAttribute{
				Computed:    true,
				Description: "Dnat规则的id",
			},
			"external_id": schema.StringAttribute{
				Required:    true,
				Description: "弹性公网id",
			},
			"external_ip": schema.StringAttribute{
				Computed:    true,
				Description: "弹性公网ip",
			},
			"external_port": schema.Int32Attribute{
				Required:    true,
				Description: "弹性IP公网端口, 1 - 1024",
				Validators: []validator.Int32{
					int32validator.Between(1, 1024),
				},
			},
			"internal_ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "内部 IP,virtual_machine_type=2(自定义),必填",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("virtual_machine_type"),
						types.Int32Value(business.VirtualMachineTypeCustom),
					),
				},
			},
			"internal_port": schema.Int32Attribute{
				Required:    true,
				Description: "主机内网端口，1-65535",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"port_id": schema.StringAttribute{
				Optional:    true,
				Description: "对应网卡id",
			},
			"port_name": schema.StringAttribute{
				Optional:    true,
				Description: "网卡名称",
			},
			"device_id": schema.StringAttribute{
				Optional:    true,
				Description: "网卡对应的设备ID",
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "协议：tcp/udp",
				Validators: []validator.String{
					stringvalidator.OneOf(business.DNatProtocols...),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "运行状态: ACTIVE / FREEZING / CREATING",
				Validators: []validator.String{
					stringvalidator.OneOf(business.DNatStatus...),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "创建时间",
			},
			"ip_expire_time": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "ip到期时间",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "描述",
			},
			"virtual_machine_id": schema.StringAttribute{
				Optional:    true,
				Description: "虚拟机id",
			},
			"virtual_machine_type": schema.Int32Attribute{
				Optional:    true,
				Description: "云主机类型，1-选择云主机，serverType字段必传 2-自定义，internalIp必传",
				//Validators: []validator.Int32{
				//	int32validator.Between(1, 2),
				//},
			},
			"server_type": schema.StringAttribute{
				Optional:    true,
				Description: "当 virtual_machine_type 为 1 时，serverType 必传，支持: VM / BM （仅支持大写）",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("virtual_machine_type"),
						types.Int32Value(business.VirtualMachineTypeCloud),
					),
				},
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "绑定状态，取值 in_progress / done",
				Validators: []validator.String{
					stringvalidator.OneOf(business.DnatStatus...),
				},
			},
		},
	}
}

func (c *ctyunDnatResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunDnatConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建前检查
	err = c.checkBeforeCreateDnat(ctx, plan)
	if err != nil {
		return
	}

	// 创建DNAT
	params, err := c.createDnat(ctx, plan)
	if err != nil {
		return
	}

	// 轮询查询dnat规则是否创建，dnat和snat的轮询策略是重复请求即可
	loopResponse, err := c.CreateLoop(ctx, params, 600)
	if err != nil {
		return
	}

	if loopResponse == nil {
		return
	}
	// 如果loopResponse不为空，则表示创建成功,保存dnat id和状态
	plan.DNatID = loopResponse.DNatID
	plan.Status = loopResponse.Status

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 反查信息
	err = c.getAndMergeDnat(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)

}

func (c *ctyunDnatResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunDnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}
	// 通过DnatId同步
	params := &ctvpc.CtvpcCreateDnatEntryRequest{
		RegionID:           state.RegionID.ValueString(),
		NatGatewayID:       state.NatGatewayID.ValueString(),
		ExternalID:         state.ExternalID.String(),
		ExternalPort:       state.ExternalPort.ValueInt32(),
		VirtualMachineID:   state.VirtualMachineID.ValueStringPointer(),
		VirtualMachineType: state.VirtualMachineType.ValueInt32(),
		InternalPort:       state.InternalPort.ValueInt32(),
		Protocol:           state.Protocol.ValueString(),
		ClientToken:        uuid.NewString(),
		Description:        state.Description.ValueStringPointer(),
	}
	if state.VirtualMachineType.ValueInt32() == business.VirtualMachineTypeCloud {
		// 如果云主机类型为1-选择云主机，传参数serverType
		params.ServerType = state.ServerType.ValueStringPointer()
	} else if state.VirtualMachineType.ValueInt32() == business.VirtualMachineTypeCustom {
		// 如果云主机类型为2 - 自定义，internalIp必传
		params.InternalIp = state.InternalIP.ValueStringPointer()
	}
	// 通过创建的参数进行同步
	if !c.acquireAndSetIdIfCreateNotFinished(ctx, &state, response, params) {
		return
	}

	// 查询远端
	err = c.getAndMergeDnat(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(request.State.Set(ctx, &state)...)
}

func (c *ctyunDnatResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf 文件中的
	var plan CtyunDnatConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// state中的
	var state CtyunDnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新dnat规则
	err = c.updateDNat(ctx, state, plan)
	if err != nil {
		return
	}
	// 查询远端信息
	err = c.getAndMergeDnat(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunDnatResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunDnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 删除操作
	// 1. 定义删除参数
	params := &ctvpc.CtvpcDeleteDnatEntryRequest{
		RegionID:    state.RegionID.ValueString(),
		DNatID:      state.DNatID.ValueString(),
		ClientToken: uuid.NewString(),
	}
	// 2.调用删除方法
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteDnatEntryApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}

	// 轮询检测时候彻底删除
	err = c.DeleteLoop(ctx, state)
	if err != nil {
		return
	}

}

func (c *ctyunDnatResource) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)

	c.meta = meta
}

func (c *ctyunDnatResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var cfg CtyunDnatConfig
	var id string
	err = terraform_extend.Split(request.ID, &id)
	if err != nil {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(cfg.RegionID.ValueString(), common.ExtraRegionId)
	cfg.RegionID = types.StringValue(regionId)

	natGatewayId := cfg.NatGatewayID.ValueString()
	cfg.NatGatewayID = types.StringValue(natGatewayId)
	err = c.getAndMergeDnat(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *ctyunDnatResource) getAndMergeDnat(ctx context.Context, cfg *CtyunDnatConfig) (err error) {
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowDnatEntryApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowDnatEntryRequest{
		RegionID:     cfg.RegionID.ValueString(),
		NatGatewayID: cfg.NatGatewayID.ValueString(),
		DNatID:       cfg.DNatID.ValueString(),
	})
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	dnat := resp.ReturnObj
	cfg.DNatID = utils.SecStringValue(dnat.DNatID)
	cfg.CreatedAt = utils.SecStringValue(dnat.CreationTime)
	cfg.Description = utils.SecStringValue(dnat.Description)
	cfg.IpExpireTime = utils.SecStringValue(dnat.IpExpireTime)
	cfg.ExternalIP = utils.SecStringValue(dnat.ExternalIp)
	cfg.ExternalID = utils.SecStringValue(dnat.ExternalID)
	//cfg.ExternalPort = types.Int32Value(dnat.ExternalPort)
	cfg.InternalIP = utils.SecStringValue(dnat.InternalIp)
	cfg.Protocol = utils.SecStringValue(dnat.Protocol)
	cfg.State = utils.SecStringValue(dnat.State)
	externalPort := types.Int32Value(dnat.ExternalPort)
	if c.isPort(externalPort, "external") {
		cfg.ExternalPort = externalPort
	}
	internalPort := types.Int32Value(dnat.InternalPort)
	if c.isPort(internalPort, "internal") {
		cfg.InternalPort = internalPort
	}
	return nil
}

func (c *ctyunDnatResource) isPort(port types.Int32, flag string) bool {
	if port.IsNull() {
		return false
	}
	if flag == "internal" {
		if port.ValueInt32() > 0 && port.ValueInt32() <= 65535 {
			return true
		}
	} else if flag == "external" {
		if port.ValueInt32() > 0 && port.ValueInt32() <= 1024 {
		}
		return true
	}
	return false
}

// checkBeforeCreateDnat 创建dnat之前进行检查
func (c *ctyunDnatResource) checkBeforeCreateDnat(_ context.Context, plan CtyunDnatConfig) (err error) {
	// dnat创建前，需要先创建nat网关
	if plan.NatGatewayID.ValueString() == "" {
		return fmt.Errorf("nat Gateway ID is empty")
	}
	//当 virtualMachineType 为 1 时，serverType 必传，支持: VM / BM （仅支持大写）
	virtualMachineType := plan.VirtualMachineType.ValueInt32()
	if virtualMachineType == business.VirtualMachineTypeCloud {
		if plan.ServerType.IsNull() || plan.ServerType.ValueString() == "" {
			return fmt.Errorf("server Type is empty")
		}
		// 判断传值,仅可以为VM/BM
		if !c.contains(plan.ServerType.ValueString(), business.ServerTypes) {
			return fmt.Errorf("server Type is invalid, serverType:%s", plan.ServerType.ValueString())
		}
	} else if virtualMachineType == business.VirtualMachineTypeCustom {
		//自定义，internalIp必传
		if plan.InternalIP.IsNull() || plan.InternalIP.ValueString() == "" {
			return fmt.Errorf("internal IP is empty")
		}

	} else {
		return fmt.Errorf("virtual machine type is invalid, virtualMachineType:%d", virtualMachineType)
	}

	// 端口判断
	if !c.isPort(plan.InternalPort, "internal") {
		return fmt.Errorf("internal port is invalid, internalPort:%s", plan.InternalPort)
	}
	if !c.isPort(plan.ExternalPort, "external") {
		return fmt.Errorf("external port is invalid, externalPort:%s", plan.ExternalPort)
	}
	// 协议判断
	if plan.Protocol.IsNull() || plan.Protocol.ValueString() == "" {
		return fmt.Errorf("protocol is empty")
	} else if !c.contains(plan.Protocol.ValueString(), business.DNatProtocols) {
		return fmt.Errorf("protocol is invalid, protocol:%s", plan.Protocol.ValueString())
	}
	return nil
}

// createDnat 创建dnat规则
func (c *ctyunDnatResource) createDnat(ctx context.Context, plan CtyunDnatConfig) (createParams *ctvpc.CtvpcCreateDnatEntryRequest, err error) {
	// 定义创建dnat规则的请求参数
	params := &ctvpc.CtvpcCreateDnatEntryRequest{
		RegionID:           plan.RegionID.ValueString(),
		NatGatewayID:       plan.NatGatewayID.ValueString(),
		ExternalID:         plan.ExternalID.ValueString(),
		ExternalPort:       plan.ExternalPort.ValueInt32(),
		VirtualMachineID:   plan.VirtualMachineID.ValueStringPointer(),
		VirtualMachineType: plan.VirtualMachineType.ValueInt32(),
		InternalPort:       plan.InternalPort.ValueInt32(),
		Protocol:           plan.Protocol.ValueString(),
		ClientToken:        uuid.NewString(),
		Description:        plan.Description.ValueStringPointer(),
	}

	if plan.VirtualMachineType.ValueInt32() == business.VirtualMachineTypeCloud {
		// 如果云主机类型为1-选择云主机，传参数serverType
		params.ServerType = plan.ServerType.ValueStringPointer()
	} else if plan.VirtualMachineType.ValueInt32() == business.VirtualMachineTypeCustom {
		// 如果云主机类型为2 - 自定义，internalIp必传
		params.InternalIp = plan.InternalIP.ValueStringPointer()
	}

	// SDK接口：ctvpc_create_dnat_entry_api.go
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateDnatEntryApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	createParams = params
	return
}

func (c *ctyunDnatResource) CreateLoop(ctx context.Context, params *ctvpc.CtvpcCreateDnatEntryRequest, loopCount ...int) (loopResponse *DnatLoopCreateResponse, err error) {
	var respError error = nil
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}

	retryer, _ := business.NewRetryer(time.Second*5, count)
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateDnatEntryApi.Do(ctx, c.meta.SdkCredential, params)
			if err != nil {
				respError = err
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				respError = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
				return false
			}
			status := *resp.ReturnObj.Dnat.Status
			switch status {
			case business.NatCreateStatusING:
				//仍在开通
				return true
			case business.NatCreateStatusDone:
				DNatId := resp.ReturnObj.Dnat.DnatID
				loopResponse = &DnatLoopCreateResponse{
					DNatID: utils.SecStringValue(DNatId),
					Status: types.StringValue(status),
				}
				return false
			default:
				// 其他状态
				respError = errors.New("invalid status, status: " + status)
				return false
			}
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return nil, errors.New("轮询已达最大次数，资源仍未创建成功！")
	}
	return loopResponse, respError
}

func (c *ctyunDnatResource) updateLoop(ctx context.Context, state *CtyunDnatConfig, updatedParams *ctvpc.CtvpcUpdateDnatEntryAttributeRequest, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*5, count)
	if err != nil {
		return
	}
	result := retryer.Start(func(currentTime int) bool {
		resp, err2 := c.meta.Apis.SdkCtVpcApis.CtvpcShowDnatEntryApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowDnatEntryRequest{
			RegionID:     state.RegionID.ValueString(),
			NatGatewayID: state.NatGatewayID.ValueString(),
			DNatID:       state.DNatID.ValueString(),
		})
		if err2 != nil {
			err = err2
			return false
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return false
		} else if resp.ReturnObj == nil {
			err = common.InvalidReturnObjError
			return false
		}

		dnatInfo := resp.ReturnObj
		if updatedParams.ExternalID != nil && *dnatInfo.ExternalID != *updatedParams.ExternalID {
			return true
		}
		if updatedParams.ExternalPort != 0 && dnatInfo.ExternalPort != updatedParams.ExternalPort {
			return true
		}
		if updatedParams.InternalIp != nil && *dnatInfo.InternalIp != *updatedParams.InternalIp {
			return true
		}
		if updatedParams.InternalPort != 0 && dnatInfo.InternalPort != updatedParams.InternalPort {
			return true
		}
		if updatedParams.Description != nil && *dnatInfo.Description != *updatedParams.Description {
			return true
		}
		return false
	})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未更新!Dnat: " + state.DNatID.ValueString())
	}
	return
}

func (c *ctyunDnatResource) DeleteLoop(ctx context.Context, state CtyunDnatConfig) (err error) {
	var respErr error
	retryer, err := business.NewRetryer(time.Second*5, 60)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowDnatEntryApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowDnatEntryRequest{
				RegionID:     state.RegionID.ValueString(),
				NatGatewayID: state.NatGatewayID.ValueString(),
				DNatID:       state.DNatID.ValueString(),
			})
			if err != nil {
				respErr = err
				return false
			}
			// 如果返回为空了，说明已经删除成功
			if resp.ReturnObj == nil {
				return false
			} else {
				// 如果仍能查询到dnat信息，说明仍未删除完成
				return true
			}
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未删除!Dnat: " + state.DNatID.ValueString())
	}
	return respErr
}

// acquireIdIfOrderNotFinished 重新获取id，如果前订单状态有问题需要重新轮询
// 返回值：数据是否有效
func (c *ctyunDnatResource) acquireAndSetIdIfCreateNotFinished(ctx context.Context, state *CtyunDnatConfig, response *resource.ReadResponse, params *ctvpc.CtvpcCreateDnatEntryRequest) bool {
	DNatId := state.DNatID.ValueString()
	if DNatId != "" {
		// 数据无需处理，是完整的
		return true
	}

	loopResponse, err := c.CreateLoop(ctx, params)
	if err != nil || loopResponse == nil {
		response.State.RemoveResource(ctx)
		return false
	}
	// 如果请求成功的话，则把id取出
	state.DNatID = loopResponse.DNatID
	state.Status = loopResponse.Status
	response.State.Set(ctx, state)
	return true

}

func (c *ctyunDnatResource) updateDNat(ctx context.Context, state CtyunDnatConfig, plan CtyunDnatConfig) (err error) {
	if state.RegionID.ValueString() != plan.RegionID.ValueString() {
		err = fmt.Errorf("when updating Dnat Information, the Planned Dnat regionId needs to remain the same as the original")
		return
	}

	params := &ctvpc.CtvpcUpdateDnatEntryAttributeRequest{
		RegionID:           plan.RegionID.ValueString(),
		DNatID:             state.DNatID.ValueString(),
		Protocol:           plan.Protocol.ValueString(),
		VirtualMachineType: plan.VirtualMachineType.ValueInt32(),
		ClientToken:        uuid.NewString(),
		ServerType:         plan.Description.ValueStringPointer(),
	}
	if plan.ExternalID.ValueString() != "" {
		params.ExternalID = plan.ExternalID.ValueStringPointer()
	}
	// 判断弹性IP公网端口和主机内网端口是否符合标准，如果需要更新的端口不符合标准，则不更新
	if c.isPort(plan.ExternalPort, "external") {
		params.ExternalPort = plan.ExternalPort.ValueInt32()
	}
	if c.isPort(plan.InternalPort, "internal") {
		params.InternalPort = plan.InternalPort.ValueInt32()
	}
	if plan.InternalIP.ValueString() != "" {
		params.InternalIp = plan.InternalIP.ValueStringPointer()
	}
	if plan.Description.ValueString() != "" {
		params.Description = plan.Description.ValueStringPointer()
	}
	if plan.VirtualMachineType.ValueInt32() == business.VirtualMachineTypeCloud && plan.ServerType.ValueString() != "" {
		params.ServerType = plan.ServerType.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdateDnatEntryAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	}
	// 轮询确认已经更新完毕
	err = c.updateLoop(ctx, &state, params, 50)
	if err != nil {
		return
	}
	return
}

//func (c *ctyunDnatResource) deleteDNat(ctx context.Context, state CtyunDnatConfig)

// contains 方法用于判断value时候被包含在list中，区分大小写
func (c *ctyunDnatResource) contains(value string, list []string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

type CtyunDnatConfig struct {
	RegionID           types.String `tfsdk:"region_id"`            //区域id
	NatGatewayID       types.String `tfsdk:"nat_gateway_id"`       //要查询的私网NAT的ID
	DNatID             types.String `tfsdk:"dnat_id"`              //DNAT规则的ID
	ExternalID         types.String `tfsdk:"external_id"`          //中转IP ID
	ExternalIP         types.String `tfsdk:"external_ip"`          //中转IP
	ExternalPort       types.Int32  `tfsdk:"external_port"`        //外部端口
	InternalIP         types.String `tfsdk:"internal_ip"`          //内部IP
	InternalPort       types.Int32  `tfsdk:"internal_port"`        //内部端口
	PortID             types.String `tfsdk:"port_id"`              //对应的网卡ID
	PortName           types.String `tfsdk:"port_name"`            //网卡名称
	DeviceID           types.String `tfsdk:"device_id"`            //网卡对应的设备ID
	Protocol           types.String `tfsdk:"protocol"`             //协议: tcp/udp
	State              types.String `tfsdk:"state"`                //DNAT状态: running代表运行中, freeze代表已冻结, expired代表已到期
	CreatedAt          types.String `tfsdk:"created_at"`           //创建时间
	IpExpireTime       types.String `tfsdk:"ip_expire_time"`       //ip到期时间
	Description        types.String `tfsdk:"description"`          //描述
	VirtualMachineID   types.String `tfsdk:"virtual_machine_id"`   //云主机
	VirtualMachineType types.Int32  `tfsdk:"virtual_machine_type"` //云主机类型1-选择云主机，serverType字段必传 ;2-自定义，internalIp必传
	ServerType         types.String `tfsdk:"server_type"`          //当 virtualMachineType 为 1 时，serverType 必传，支持: VM / BM （仅支持大写）
	Status             types.String `tfsdk:"status"`               //绑定状态，取值 in_progress / done
}

type DnatLoopCreateResponse struct {
	DNatID types.String
	Status types.String
}
