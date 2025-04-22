package elb

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	ctelb "terraform-provider-ctyun/internal/core/ctelb"
	"terraform-provider-ctyun/internal/utils"
)

var (
	_ resource.Resource                = &CtyunElbLoadBalancerResource{}
	_ resource.ResourceWithConfigure   = &CtyunElbLoadBalancerResource{}
	_ resource.ResourceWithImportState = &CtyunElbLoadBalancerResource{}
)

type CtyunElbLoadBalancerResource struct {
	meta *common.CtyunMetadata
}

func NewCtyunElbBalancerResource() resource.Resource {
	return &CtyunElbLoadBalancerResource{}
}
func (c *CtyunElbLoadBalancerResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_loadbalancer"
}

func (c *CtyunElbLoadBalancerResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "**文档详情：https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=24&api=5643&data=88&isNormal=1&vid=82",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Description: "",
			},
			"client_token": schema.StringAttribute{
				Description: "",
			},
			"project_id": schema.StringAttribute{
				Description: "",
			},
			"vpc_id": schema.StringAttribute{
				Description: "",
			},
			"subnet_id": schema.StringAttribute{
				Description: "",
			},
			"name": schema.StringAttribute{
				Description: "",
			},
			"description": schema.StringAttribute{
				Description: "",
			},
			"eip_id": schema.StringAttribute{
				Description: "",
			},
			"sla_name": schema.StringAttribute{
				Description: "",
			},
			"resource_type": schema.StringAttribute{
				Description: "",
			},
			"private_ip_address": schema.StringAttribute{
				Description: "",
			},
			"delete_protection": schema.StringAttribute{
				Description: "",
			},
			"id": schema.StringAttribute{
				Description: "",
			},
		},
	}
}

func (c *CtyunElbLoadBalancerResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunElbLoadBalancerConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	//创建前检查
	err = c.checkBeforeCreateElb(ctx, plan)
	if err != nil {
		return
	}

	// 创建
	returnObj, err := c.createElb(ctx, &plan)
	if err != nil {
		return
	}
	// todo 没有masterOrderIp，不确认是否为异步请求
	plan.ID = types.StringValue(returnObj.ID)
	// 创建后反查创建后的nat信息
	err = c.getAndMergeElb(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbLoadBalancerResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunElbLoadBalancerConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeElb(ctx, &state)
	if err != nil {
		// 有待确定
		if strings.Contains(err.Error(), "is not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunElbLoadBalancerResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取tf文件中配置
	var plan CtyunElbLoadBalancerConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunElbLoadBalancerConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
	}

	// 更新基本信息
	err = c.updateElbInfo(ctx, state, plan)
	if err != nil {
		return
	}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeElb(ctx, &state)
	if err != nil {
		return
	}

	//todo
	//升级为保障型负载均衡实例
	//保障型负载均衡实例创建
	//保障型负载均衡实例变配
	//保障型负载均衡实例续订
	//保障型负载均衡实例退订
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunElbLoadBalancerResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunElbLoadBalancerConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	params := &ctelb.CtelbDeleteLoadBalancerRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionID.String(),
	}
	if !state.ProjectID.IsNull() {
		params.ProjectID = state.ProjectID.ValueString()
	}
	if !state.ID.IsNull() {
		params.ID = state.ID.ValueString()
		params.ElbID = state.ID.ValueString()
	}

	// SDK ctelb_delete_load_balancer_api.go
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbDeleteLoadBalancerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj != nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}
func (c *CtyunElbLoadBalancerResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {

}

func (c *CtyunElbLoadBalancerResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunElbLoadBalancerResource) createElb(ctx context.Context, plan *CtyunElbLoadBalancerConfig) (returnObj ctelb.CtelbCreateLoadBalancerReturnObjResponse, err error) {
	params := &ctelb.CtelbCreateLoadBalancerRequest{
		ClientToken:  uuid.NewString(),
		RegionID:     plan.RegionID.ValueString(),
		SubnetID:     plan.SubnetID.ValueString(),
		Name:         plan.Name.ValueString(),
		SlaName:      plan.SlaName.ValueString(),
		ResourceType: plan.ResourceType.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		params.ProjectID = plan.ProjectID.ValueString()
	}
	if !plan.VpcID.IsNull() {
		params.VpcID = plan.VpcID.ValueString()
	}
	if !plan.Description.IsNull() {
		params.Description = plan.Description.ValueString()
	}

	if plan.ResourceType.ValueString() == business.LbResourceTypeExternal || !plan.EipID.IsNull() {
		params.EipID = plan.EipID.ValueString()
	}
	if !plan.PrivateIpAddress.IsNull() {
		params.PrivateIpAddress = plan.PrivateIpAddress.ValueString()
	}
	if !plan.DeleteProtection.IsNull() {
		params.DeleteProtection = plan.DeleteProtection.ValueBoolPointer()
	}

	//SDK ctelb_create_load_balancer_api.go
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbCreateLoadBalancerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}

	returnObj = *resp.ReturnObj
	return
}

func (c *CtyunElbLoadBalancerResource) checkBeforeCreateElb(_ context.Context, plan CtyunElbLoadBalancerConfig) error {
	// regionid不能为空，subnetID	(子网id)不能为空,name不能为空，slaName不能为空，resourceType不能为空
	regionId := plan.RegionID
	subnetId := plan.SubnetID
	slaName := plan.SlaName
	resourceType := plan.ResourceType
	name := plan.Name
	eipId := plan.EipID
	if regionId.IsNull() {
		return fmt.Errorf("regionID不能为空!")
	}
	if subnetId.IsNull() {
		return fmt.Errorf("subnetId-子网的ID不能为空!")
	}
	if slaName.IsNull() {
		return fmt.Errorf("slaName-lb的规格名称不能为空！")
	}
	if resourceType.IsNull() {
		return fmt.Errorf("resourceType-资源类型不能为空！")
	}
	if !c.isContains(resourceType.ValueString(), business.LbResourceType) {
		return fmt.Errorf("resourceType资源类型取值存在问题，resourceType取值范围为{internal：内网负载均衡，external：公网负载均衡}")
	}
	//当resourceType=external为必填, eipID不能为空
	if resourceType.ValueString() == business.LbResourceTypeExternal && eipId.IsNull() {
		return fmt.Errorf("当resourceType=external为必填, eipID不能为空")
	}

	if name.IsNull() {
		return fmt.Errorf("name不能为空")
	}
	return nil
}

func (c *CtyunElbLoadBalancerResource) getAndMergeElb(ctx context.Context, plan *CtyunElbLoadBalancerConfig) (err error) {
	//查看ELB详情： ctelb_show_load_balancer_api.go
	params := &ctelb.CtelbShowLoadBalancerRequest{
		RegionID: plan.RegionID.ValueString(),
		ElbID:    plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbShowLoadBalancerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析resp.ReturnObj,将最新的elb信息同步到state中
	if len(resp.ReturnObj) > 1 {
		err = fmt.Errorf("ReturnObj长度>1")
		return
	}
	elbObj := resp.ReturnObj[0]
	// todo 我认为这里返回list是不合理的，id应该一一对应，我这里写成取第1个对象
	if plan.RegionID.ValueString() != elbObj.RegionID {
		err = fmt.Errorf("elb详情regionid(%s)与plan的reigonid(%s)不一致！", elbObj.RegionID, plan.RegionID.ValueString())
		return
	}
	if plan.ID.ValueString() != elbObj.ID {
		err = fmt.Errorf("详情elb id(%s)与plan的elb id(%s)不一致！", elbObj.RegionID, plan.RegionID.ValueString())
		return
	}
	plan.AzName = types.StringValue(elbObj.AzName)
	plan.ProjectID = types.StringValue(elbObj.ProjectID)
	plan.Name = types.StringValue(elbObj.Name)
	plan.Description = types.StringValue(elbObj.Description)
	plan.VpcID = types.StringValue(elbObj.VpcID)
	plan.SubnetID = types.StringValue(elbObj.SubnetID)
	plan.PortID = types.StringValue(elbObj.PortID)
	plan.PrivateIpAddress = types.StringValue(elbObj.PrivateIpAddress)
	plan.Ipv6Address = types.StringValue(elbObj.Ipv6Address)
	plan.SlaName = types.StringValue(elbObj.SlaName)
	plan.DeleteProtection = types.BoolValue(*elbObj.DeleteProtection)
	plan.AdminStatus = types.StringValue(elbObj.AdminStatus)
	plan.Status = types.StringValue(elbObj.Status)
	plan.ResourceType = types.StringValue(elbObj.ResourceType)
	plan.CreatedTime = types.StringValue(elbObj.CreatedTime)
	plan.UpdatedTime = types.StringValue(elbObj.UpdatedTime)
	EipInfoList := elbObj.EipInfo
	var eipInfos []EipInfoModel
	if EipInfoList != nil && len(EipInfoList) > 0 {
		for _, eipItem := range EipInfoList {
			var eipInfo EipInfoModel
			eipInfo.ResourceID = types.StringValue(eipItem.ResourceID)
			eipInfo.EipID = types.StringValue(eipItem.EipID)
			eipInfo.Bandwidth = types.Int32Value(eipItem.Bandwidth)
			if eipItem.IsTalkOrder != nil {
				eipInfo.IsTalkOrder = types.BoolValue(*eipItem.IsTalkOrder)
			}
			eipInfos = append(eipInfos, eipInfo)
		}
	}
	eipInfoType := utils.StructToTFObjectTypes(EipInfoModel{})
	plan.eipInfo, _ = types.ListValueFrom(ctx, eipInfoType, eipInfos)
	return
}

func (c *CtyunElbLoadBalancerResource) updateElbInfo(ctx context.Context, state CtyunElbLoadBalancerConfig, plan CtyunElbLoadBalancerConfig) (err error) {
	// SDK ctelb_update_load_balancer_api.go
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateLoadBalancerApi.Do(ctx, c.meta.SdkCredential, &ctelb.CtelbUpdateLoadBalancerRequest{
		ClientToken:      uuid.NewString(),
		RegionID:         plan.RegionID.ValueString(),
		ID:               plan.ID.ValueString(),
		ElbID:            plan.ID.ValueString(),
		SlaName:          plan.SlaName.ValueString(),
		Name:             plan.Name.ValueString(),
		Description:      plan.Description.ValueString(),
		DeleteProtection: plan.DeleteProtection.ValueBoolPointer(),
	})
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

func (c *CtyunElbLoadBalancerResource) isContains(value string, collect []string) bool {
	for _, v := range collect {
		if v == value {
			return true
		}
	}
	return false
}

type CtyunElbLoadBalancerConfig struct {
	RegionID         types.String `tfsdk:"region_id"`          //区域ID
	ClientToken      types.String `tfsdk:"client_token"`       //客户端存根，用于保证订单幂等性, 长度 1 - 64
	ProjectID        types.String `tfsdk:"project_id"`         //企业项目 ID，默认为'0'
	VpcID            types.String `tfsdk:"vpc_id"`             //vpc的ID
	SubnetID         types.String `tfsdk:"subnet_id"`          //子网的ID
	Name             types.String `tfsdk:"name"`               //唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32
	Description      types.String `tfsdk:"description"`        //支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128
	EipID            types.String `tfsdk:"eip_id"`             //弹性公网IP的ID。当resourceType=external为必填
	SlaName          types.String `tfsdk:"sla_name"`           //lb的规格名称,支持elb.s1.small和elb.default，默认为elb.default
	ResourceType     types.String `tfsdk:"resource_type"`      //资源类型。internal：内网负载均衡，external：公网负载均衡
	PrivateIpAddress types.String `tfsdk:"private_ip_address"` //负载均衡的私有IP地址，不指定则自动分配
	DeleteProtection types.Bool   `tfsdk:"delete_protection"`  //删除保护。false（不开启）、true（开）。 默认：不开启
	ID               types.String `tfsdk:"id"`                 //负载均衡ID
	AzName           types.String `tfsdk:"az_name"`
	PortID           types.String `tfsdk:"port_id"`
	Ipv6Address      types.String `tfsdk:"ipv6_address"`
	eipInfo          types.List   `tfsdk:"eip_info"`
	dminStatus       types.String `tfsdk:"admin_status"`
	AdminStatus      types.String `tfsdk:"admin_status"`
	Status           types.String `tfsdk:"status"`
	CreatedTime      types.String `tfsdk:"created_time"`
	UpdatedTime      types.String `tfsdk:"created_time"`
	//Elbs             types.List   `tfsdk:"elbs"`

}

//type CtyunElbDetailModel struct {
//	// 详情信息
//	RegionID         types.String `tfsdk:"region_id"` //区域ID
//	AzName           types.String `tfsdk:"az_name"`
//	ID               types.String `tfsdk:"id"`                 //负载均衡ID
//	ProjectID        types.String `tfsdk:"project_id"`         //企业项目 ID，默认为'0'
//	Name             types.String `tfsdk:"name"`               //唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32
//	Description      types.String `tfsdk:"description"`        //支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128
//	VpcID            types.String `tfsdk:"vpc_id"`             //vpc的ID
//	SubnetID         types.String `tfsdk:"subnet_id"`          //子网的ID
//	PortID           types.String `tfsdk:"port_id"`            //负载均衡实例默认创建port ID
//	PrivateIpAddress types.String `tfsdk:"private_ip_address"` //负载均衡的私有IP地址，不指定则自动分配
//	Ipv6Address      types.String `tfsdk:"ipv6_address"`       //负载均衡实例的IPv6地址
//	SlaName          types.String `tfsdk:"sla_name"`           //lb的规格名称,支持elb.s1.small和elb.default，默认为elb.default
//	eipInfo          types.List   `tfsdk:"eip_info"`           //[]EipInfoModel
//	DeleteProtection types.Bool   `tfsdk:"delete_protection"`  //删除保护。false（不开启）、true（开）。 默认：不开启
//	AdminStatus      types.String `tfsdk:"admin_status"`       //管理状态: DOWN / ACTIVE
//	Status           types.String `tfsdk:"status"`             //负载均衡状态: DOWN / ACTIVE
//	ResourceType     types.String `tfsdk:"resource_type"`      //资源类型。internal：内网负载均衡，external：公网负载均衡
//	CreatedTime      types.String `tfsdk:"created_time"`       //创建时间，为UTC格式
//	UpdatedTime      types.String `tfsdk:"updated_time"`       //更新时间，为UTC格式
//}
