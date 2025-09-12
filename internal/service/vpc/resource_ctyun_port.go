package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

func NewCtyunNetworkInterface() resource.Resource {
	return &ctyunNetworkInterface{}
}

type ctyunNetworkInterface struct {
	meta *common.CtyunMetadata
}

func (c *ctyunNetworkInterface) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_port"
}

type CtyunNetworkInterfaceResource struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	RegionId                types.String `tfsdk:"region_id"`
	SubnetId                types.String `tfsdk:"subnet_id"`
	PrimaryIpAddress        types.String `tfsdk:"primary_ip_address"`
	SecurityGroupIds        types.Set    `tfsdk:"security_group_ids"`
	SecondaryPrivateIpCount types.Int32  `tfsdk:"secondary_private_ip_count"`
	SecondaryPrivateIps     types.Set    `tfsdk:"secondary_private_ips"`
	Ipv6AddressCount        types.Int32  `tfsdk:"ipv6_address_count"`
	Ipv6Addresses           types.List   `tfsdk:"ipv6_addresses"`
	NetworkInterfaceId      types.String `tfsdk:"network_interface_id"`
	MacAddress              types.String `tfsdk:"mac_address"`
	InstanceId              types.String `tfsdk:"instance_id"`
	InstanceType            types.String `tfsdk:"instance_type"`
	Status                  types.String `tfsdk:"status"`
}

func (c *ctyunNetworkInterface) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**弹性网卡资源，详细说明请见文档：https://www.ctyun.cn/document/10026730**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "网卡ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "弹性网卡名称。支持拉丁字母、中文、数字，下划线，连字符，中文/英文字母开头，不能以http:/https:开头，长度2-32  支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 32),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\u4e00-\u9fa5][0-9a-zA-Z_\u4e00-\u9fa5-]*[0-9a-zA-Z_\u4e00-\u9fa5]$"), "弹性网卡名称不符合规则"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "弹性网卡描述。支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}|《》？：“”【】、；‘'，。、，不能以http:/https:开头，长度0-128 支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(128),
					validator2.Desc(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"primary_ip_address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "主私有IP地址，如果不指定则自动分配",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					validator2.Ip(),
				},
			},
			"security_group_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "安全组ID列表，最多支持10个",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(10),
				},
			},
			"secondary_private_ip_count": schema.Int32Attribute{
				Optional:    true,
				Description: "辅助私有IP地址数量，指定私有IP地址数量，让系统为您自动创建IP地址，最多支持10个",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
					int32validator.AtMost(10),
				},
			},
			"secondary_private_ips": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "辅助私有IP地址列表，指定私有IP地址，不能和secondary_private_ip_count同时指定，最多支持10个",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(10),
				},
			},
			"ipv6_address_count": schema.Int32Attribute{
				Optional:    true,
				Description: "IPv6地址数量，指定IPv6地址数量，让系统为您自动创建IPv6地址，最多支持10个",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
					int32validator.AtMost(10),
				},
			},
			"ipv6_addresses": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "IPv6地址列表，指定IPv6地址，不能和ipv6_address_count同时指定，最多支持10个",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
					listplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(10),
				},
			},
			"network_interface_id": schema.StringAttribute{
				Computed:    true,
				Description: "网卡ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mac_address": schema.StringAttribute{
				Computed:    true,
				Description: "MAC地址",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Computed:    true,
				Description: "绑定的实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_type": schema.StringAttribute{
				Computed:    true,
				Description: "实例类型",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "网卡状态",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *ctyunNetworkInterface) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunNetworkInterfaceResource
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	resp, err := c.createNetworkInterface(ctx, plan)

	if err != nil {
		response.Diagnostics.AddError(
			"创建弹性网卡失败",
			fmt.Sprintf("创建弹性网卡时发生错误: %s", err.Error()),
		)
		return
	}

	if resp.ReturnObj == nil || resp.ReturnObj.NetworkInterfaceID == nil {
		// 打印详细的错误信息
		errorInfo := ""
		if resp.Message != nil {
			errorInfo += fmt.Sprintf("Message: %s ", *resp.Message)
		}
		if resp.ErrorCode != nil {
			errorInfo += fmt.Sprintf("ErrorCode: %s ", *resp.ErrorCode)
		}
		if resp.Description != nil {
			errorInfo += fmt.Sprintf("Description: %s ", *resp.Description)
		}

		response.Diagnostics.AddError(
			"创建弹性网卡失败",
			fmt.Sprintf("API返回数据: %s", errorInfo),
		)
		return
	}
	// 更新计划中的所有字段
	plan.Id = types.StringPointerValue(resp.ReturnObj.NetworkInterfaceID)
	plan.NetworkInterfaceId = types.StringPointerValue(resp.ReturnObj.NetworkInterfaceID)
	// 查询网卡详细信息并更新状态
	updatedPlan, err := c.getAndMergePort(ctx, &plan)
	if err != nil {
		response.Diagnostics.AddError(
			"获取弹性网卡信息失败",
			fmt.Sprintf("获取弹性网卡信息时发生错误: %s", err.Error()),
		)
		return
	}

	if updatedPlan != nil {
		plan = *updatedPlan
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunNetworkInterface) createNetworkInterface(ctx context.Context, plan CtyunNetworkInterfaceResource) (*ctvpc.CtvpcCreatePortResponse, error) {
	// 根据传入的不同参数确定创建方式
	regionId := plan.RegionId.ValueString()
	// 构造创建请求参数
	createReq := &ctvpc.CtvpcCreatePortRequest{
		ClientToken: uuid.NewString(),
		RegionID:    regionId,
		SubnetID:    plan.SubnetId.ValueString(),
	}

	// 处理主私有IP地址
	if !plan.PrimaryIpAddress.IsNull() && !plan.PrimaryIpAddress.IsUnknown() {
		primaryIp := plan.PrimaryIpAddress.ValueString()
		createReq.PrimaryPrivateIp = &primaryIp
	}

	// 处理安全组ID列表
	if !plan.SecurityGroupIds.IsNull() && len(plan.SecurityGroupIds.Elements()) > 0 {
		var sgIds []string
		plan.SecurityGroupIds.ElementsAs(ctx, &sgIds, false)
		sgIdPtrs := make([]*string, len(sgIds))
		for i, sgId := range sgIds {
			sgIdPtrs[i] = &sgId
		}
		createReq.SecurityGroupIds = sgIdPtrs
	}

	// 处理辅助私有IP地址数量
	if !plan.SecondaryPrivateIpCount.IsNull() {
		createReq.SecondaryPrivateIpCount = plan.SecondaryPrivateIpCount.ValueInt32()
	}

	// 处理辅助私有IP地址列表
	if !plan.SecondaryPrivateIps.IsNull() && len(plan.SecondaryPrivateIps.Elements()) > 0 {
		var secondaryIps []string
		plan.SecondaryPrivateIps.ElementsAs(ctx, &secondaryIps, false)
		secondaryIpPtrs := make([]*string, len(secondaryIps))
		for i, ip := range secondaryIps {
			secondaryIpPtrs[i] = &ip
		}
		createReq.SecondaryPrivateIps = secondaryIpPtrs
	}

	// 处理IPv6地址数量
	if !plan.Ipv6AddressCount.IsNull() {
		ipv6Count := plan.Ipv6AddressCount.ValueInt32()
		if ipv6Count > 0 {
			ipv6Addresses := make([]*string, ipv6Count)
			for i := int32(0); i < ipv6Count; i++ {
				ipv6Addresses[i] = nil // 表示自动分配
			}
			createReq.Ipv6Addresses = ipv6Addresses
		}
	}

	// 处理IPv6地址列表
	if !plan.Ipv6Addresses.IsNull() && len(plan.Ipv6Addresses.Elements()) > 0 {
		var ipv6Addresses []string
		plan.Ipv6Addresses.ElementsAs(ctx, &ipv6Addresses, false)
		ipv6AddrPtrs := make([]*string, len(ipv6Addresses))
		for i, addr := range ipv6Addresses {
			ipv6AddrPtrs[i] = &addr
		}
		createReq.Ipv6Addresses = ipv6AddrPtrs
	}

	// 处理名称和描述
	if !plan.Name.IsNull() {
		name := plan.Name.ValueString()
		createReq.Name = &name
	}
	if !plan.Description.IsNull() {
		description := plan.Description.ValueString()
		createReq.Description = &description
	}

	// 调用API创建弹性网卡
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreatePortApi.Do(ctx, c.meta.SdkCredential, createReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ctyunNetworkInterface) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunNetworkInterfaceResource
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 更新状态
	updatedState, err := c.getAndMergePort(ctx, &state)
	if err != nil {
		response.Diagnostics.AddError(
			"获取弹性网卡信息失败",
			fmt.Sprintf("获取弹性网卡信息时发生错误: %s", err.Error()),
		)
		return
	}

	// 如果返回为nil，表示资源不存在，需要从状态中移除
	if updatedState == nil {
		response.State.RemoveResource(ctx)
		return
	}

	state = *updatedState

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (c *ctyunNetworkInterface) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state CtyunNetworkInterfaceResource
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	_, done := c.updateNetworkInterface(ctx, plan, state, response)
	if done {
		return
	}

	// 查询网卡详细信息并更新状态
	updatedPlan, err := c.getAndMergePort(ctx, &plan)
	if err != nil {
		response.Diagnostics.AddError(
			"获取弹性网卡信息失败",
			fmt.Sprintf("获取弹性网卡信息时发生错误: %s", err.Error()),
		)
		return
	}

	if updatedPlan != nil {
		plan = *updatedPlan
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunNetworkInterface) updateNetworkInterface(ctx context.Context, plan CtyunNetworkInterfaceResource, state CtyunNetworkInterfaceResource, response *resource.UpdateResponse) (string, bool) {
	// 检查必要参数
	networkInterfaceId := plan.NetworkInterfaceId.ValueString()
	if networkInterfaceId == "" {
		networkInterfaceId = state.NetworkInterfaceId.ValueString()
	}

	if networkInterfaceId == "" {
		response.Diagnostics.AddError(
			"更新弹性网卡失败",
			"network_interface_id不能为空",
		)
		return "", true
	}
	// 检查是否需要更新名称或描述
	if !plan.Name.Equal(state.Name) || !plan.Description.Equal(state.Description) {
		updateReq := &ctvpc.CtvpcUpdatePortRequest{
			ClientToken:        uuid.NewString(),
			RegionID:           plan.RegionId.ValueString(),
			NetworkInterfaceID: networkInterfaceId,
		}

		// 处理名称
		if !plan.Name.IsNull() {
			name := plan.Name.ValueString()
			updateReq.Name = &name
		}

		// 处理描述
		if !plan.Description.IsNull() {
			description := plan.Description.ValueString()
			updateReq.Description = &description
		}
		// 处理安全组ID列表
		if !plan.SecurityGroupIds.IsNull() && len(plan.SecurityGroupIds.Elements()) > 0 {
			var sgIds []string
			plan.SecurityGroupIds.ElementsAs(ctx, &sgIds, false)
			sgIdPtrs := make([]*string, len(sgIds))
			for i, sgId := range sgIds {
				sgIdPtrs[i] = &sgId
			}
			updateReq.SecurityGroupIDs = sgIdPtrs
		} else {
			// 如果安全组为空，则传递空数组
			updateReq.SecurityGroupIDs = []*string{}
		}
		// 调用API更新网卡属性
		_, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdatePortApi.Do(ctx, c.meta.SdkCredential, updateReq)
		if err != nil {
			response.Diagnostics.AddError(
				"更新弹性网卡属性失败",
				fmt.Sprintf("更新弹性网卡属性时发生错误: %s", err.Error()),
			)
			return "", true
		}
	}
	return networkInterfaceId, false
}

func (c *ctyunNetworkInterface) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunNetworkInterfaceResource
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 检查必要参数
	var networkInterfaceId string
	if !state.NetworkInterfaceId.IsNull() {
		networkInterfaceId = state.NetworkInterfaceId.ValueString()
	}

	if networkInterfaceId == "" {
		response.Diagnostics.AddError(
			"删除弹性网卡失败",
			"network_interface_id不能为空",
		)
		return
	}

	var regionId string
	if !state.RegionId.IsNull() {
		regionId = state.RegionId.ValueString()
	}

	// 构造删除请求参数
	deleteReq := &ctvpc.CtvpcDeletePortRequest{
		ClientToken:        uuid.NewString(),
		RegionID:           regionId,
		NetworkInterfaceID: networkInterfaceId,
	}

	// 调用API删除弹性网卡
	_, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeletePortApi.Do(ctx, c.meta.SdkCredential, deleteReq)
	if err != nil {
		response.Diagnostics.AddError(
			"删除弹性网卡失败",
			fmt.Sprintf("删除弹性网卡时发生错误: %s", err.Error()),
		)
		return
	}

}

// ImportState 导入资源状态
func (c *ctyunNetworkInterface) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var id, regionId string
	err := terraform_extend.Split(request.ID, &id, &regionId)
	if err != nil {
		response.Diagnostics.AddError(
			"导入弹性网卡失败",
			fmt.Sprintf("导入弹性网卡时发生错误: %s", err.Error()),
		)
		return
	}

	// 设置导入的属性
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), id)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("network_interface_id"), id)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("region_id"), regionId)...)
}

// Configure 配置资源
func (c *ctyunNetworkInterface) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// getNetworkInterface 获取网卡详细信息
func (c *ctyunNetworkInterface) getNetworkInterface(ctx context.Context, regionId, networkInterfaceId string) (*ctvpc.CtvpcShowPortReturnObjResponse, error) {
	// 检查networkInterfaceId是否为空
	if networkInterfaceId == "" {
		return nil, fmt.Errorf("networkInterfaceId不能为空")
	}

	req := &ctvpc.CtvpcShowPortRequest{
		RegionID:           regionId,
		NetworkInterfaceID: networkInterfaceId,
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowPortApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return nil, err
	}

	if resp.ReturnObj == nil {
		// 打印详细的错误信息
		errorInfo := ""
		if resp.Message != nil {
			errorInfo += fmt.Sprintf("Message: %s ", *resp.Message)
		}
		if resp.ErrorCode != nil {
			errorInfo += fmt.Sprintf("ErrorCode: %s ", *resp.ErrorCode)
		}
		if resp.Description != nil {
			errorInfo += fmt.Sprintf("Description: %s ", *resp.Description)
		}
		if errorInfo != "" {
			return nil, fmt.Errorf("API返回数据为空 (%s)", errorInfo)
		}
		return nil, fmt.Errorf("API返回数据为空")
	}

	return resp.ReturnObj, nil
}

// getAndMergePort 查询网卡信息并合并到资源配置中
func (c *ctyunNetworkInterface) getAndMergePort(ctx context.Context, plan *CtyunNetworkInterfaceResource) (*CtyunNetworkInterfaceResource, error) {
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowPortApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowPortRequest{
		RegionID:           plan.RegionId.ValueString(),
		NetworkInterfaceID: plan.NetworkInterfaceId.ValueString(),
	})
	if err != nil {
		// 检查是否是因为网卡不存在导致的错误
		if err.Error() == common.OpenapiVpcPortNotFound {
			return nil, nil
		}
		return nil, err
	}

	if resp.ReturnObj == nil {
		// 打印详细的错误信息
		errorInfo := ""
		if resp.Message != nil {
			errorInfo += fmt.Sprintf("Message: %s ", *resp.Message)
		}
		if resp.ErrorCode != nil {
			errorInfo += fmt.Sprintf("ErrorCode: %s ", *resp.ErrorCode)
		}
		if resp.Description != nil {
			errorInfo += fmt.Sprintf("Description: %s ", *resp.Description)
		}
		if errorInfo != "" {
			return nil, fmt.Errorf("API返回数据为空 (%s)", errorInfo)
		}
		return nil, fmt.Errorf("API返回数据为空")
	}
	// 更新计划中的所有字段
	plan.Id = types.StringPointerValue(resp.ReturnObj.NetworkInterfaceID)
	plan.NetworkInterfaceId = types.StringPointerValue(resp.ReturnObj.NetworkInterfaceID)
	plan.Name = types.StringPointerValue(resp.ReturnObj.NetworkInterfaceName)
	plan.Description = types.StringPointerValue(resp.ReturnObj.Description)
	plan.MacAddress = types.StringPointerValue(resp.ReturnObj.MacAddress)
	plan.SubnetId = types.StringPointerValue(resp.ReturnObj.SubnetID)
	plan.PrimaryIpAddress = types.StringPointerValue(resp.ReturnObj.PrimaryPrivateIp)
	plan.InstanceId = types.StringPointerValue(resp.ReturnObj.InstanceID)
	plan.InstanceType = types.StringPointerValue(resp.ReturnObj.InstanceType)

	// 设置状态
	if resp.ReturnObj.AdminStatus != nil {
		plan.Status = types.StringValue(*resp.ReturnObj.AdminStatus)
	} else {
		plan.Status = types.StringValue("UNKNOWN")
	}

	// 设置安全组ID
	if resp.ReturnObj.SecurityGroupIds != nil {
		sgIds := make([]attr.Value, len(resp.ReturnObj.SecurityGroupIds))
		for i, sgId := range resp.ReturnObj.SecurityGroupIds {
			if sgId != nil {
				sgIds[i] = types.StringValue(*sgId)
			}
		}
		plan.SecurityGroupIds, _ = types.SetValue(types.StringType, sgIds)
	}

	// 设置辅助私有IP
	if resp.ReturnObj.SecondaryPrivateIps != nil {
		secondaryIps := make([]attr.Value, len(resp.ReturnObj.SecondaryPrivateIps))
		for i, ip := range resp.ReturnObj.SecondaryPrivateIps {
			if ip != nil {
				secondaryIps[i] = types.StringValue(*ip)
			}
		}
		plan.SecondaryPrivateIps, _ = types.SetValue(types.StringType, secondaryIps)
	}

	// 设置IPv6地址
	if resp.ReturnObj.Ipv6Addresses != nil {
		ipv6Addrs := make([]attr.Value, len(resp.ReturnObj.Ipv6Addresses))
		for i, addr := range resp.ReturnObj.Ipv6Addresses {
			if addr != nil {
				ipv6Addrs[i] = types.StringValue(*addr)
			}
		}
		plan.Ipv6Addresses, _ = types.ListValue(types.StringType, ipv6Addrs)
	} else {
		// 如果没有IPv6地址，确保字段被正确初始化为空列表
		plan.Ipv6Addresses, _ = types.ListValue(types.StringType, []attr.Value{})
	}

	return plan, nil
}
