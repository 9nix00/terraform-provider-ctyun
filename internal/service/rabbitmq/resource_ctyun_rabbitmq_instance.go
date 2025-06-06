package rabbitmq

//
//import (
//	"context"
//	"fmt"
//	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
//	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
//	"github.com/hashicorp/terraform-plugin-framework/path"
//	"github.com/hashicorp/terraform-plugin-framework/resource"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
//	"github.com/hashicorp/terraform-plugin-framework/types"
//	"regexp"
//	"strings"
//	"terraform-provider-ctyun/internal/business"
//	"terraform-provider-ctyun/internal/common"
//
//	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
//	"terraform-provider-ctyun/internal/extend/terraform/defaults"
//	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
//	"terraform-provider-ctyun/internal/utils"
//	"time"
//)
//
//var (
//	_ resource.Resource                = &ctyunRabbitmqInstance{}
//	_ resource.ResourceWithConfigure   = &ctyunRabbitmqInstance{}
//	_ resource.ResourceWithImportState = &ctyunRabbitmqInstance{}
//)
//
//type ctyunRabbitmqInstance struct {
//	meta       *common.CtyunMetadata
//	vpcService *business.VpcService
//	sgService  *business.SecurityGroupService
//}
//
//func NewCtyunRabbitmqInstance() resource.Resource {
//	return &ctyunRabbitmqInstance{}
//}
//
//func (c *ctyunRabbitmqInstance) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
//	response.TypeName = request.ProviderTypeName + "_rabbitmq_instance"
//}
//
//type CtyunRabbitmqInstanceConfig struct {
//	ID              types.String `tfsdk:"id"`
//	MasterOrderID   types.String `tfsdk:"master_order_id"`
//	RegionID        types.String `tfsdk:"region_id"`
//	InstanceName    types.String `tfsdk:"instance_name"`
//	HostType        types.String `tfsdk:"host_type"`
//	DiskType        types.String `tfsdk:"disk_type"`
//	DiskSize        types.Int32  `tfsdk:"disk_size"`
//	CpuNum          types.Int32  `tfsdk:"cpu_num"`
//	MemSize         types.Int32  `tfsdk:"mem_size"`
//	NodeNum         types.Int32  `tfsdk:"node_num"`
//	EngineType      types.String `tfsdk:"engine_type"`
//	VpcID           types.String `tfsdk:"vpc_id"`
//	SubnetID        types.String `tfsdk:"subnet_id"`
//	SecurityGroupID types.String `tfsdk:"security_group_id"`
//	AzInfo          types.String `tfsdk:"az_info"`
//	CycleType       types.String `tfsdk:"cycle_type"`
//	CycleCount      types.Int32  `tfsdk:"cycle_cnt"`
//}
//
//func (c *ctyunRabbitmqInstance) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
//	response.Schema = schema.Schema{
//		MarkdownDescription: `**详细说明请见文档：**`,
//		Attributes: map[string]schema.Attribute{
//			"id": schema.StringAttribute{
//				Computed:    true,
//				Description: "ID",
//			},
//			"master_order_id": schema.StringAttribute{
//				Computed:    true,
//				Description: "主订单号",
//			},
//			"region_id": schema.StringAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "资源池ID",
//				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"instance_name": schema.StringAttribute{
//				Required:    true,
//				Description: "实例名称",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"host_type": schema.StringAttribute{
//				Required:    true,
//				Description: "主机类型，例如：S6",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"disk_type": schema.StringAttribute{
//				Required:    true,
//				Description: "磁盘类型",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"disk_size": schema.Int32Attribute{
//				Required:    true,
//				Description: "单个节点的磁盘存储空间，单位为GB，存储空间取值范围100GB ~ 10000，并且为100的倍数。实例总存储空间为diskSize * nodeNum",
//				Validators: []validator.Int32{
//					int32validator.Between(100, 10000),
//				},
//			},
//			"node_num": schema.Int32Attribute{
//				Required:    true,
//				Description: "节点数。单机版为1个，集群版3~50个",
//				Validators: []validator.Int32{
//					int32validator.Between(1, 50),
//				},
//			},
//			"zone_list": schema.SetAttribute{
//				Required:    true,
//				ElementType: types.StringType,
//				Description: "实例所在可用区信息",
//				PlanModifiers: []planmodifier.Set{
//					setplanmodifier.RequiresReplace(),
//				},
//			},
//
//			"vpc_id": schema.StringAttribute{
//				Required:    true,
//				Description: "关联的vpcID",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"subnet_id": schema.StringAttribute{
//				Required:    true,
//				Description: "子网ID",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"security_group_id": schema.StringAttribute{
//				Required:    true,
//				Description: "安全组ID",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//					stringplanmodifier.UseStateForUnknown(),
//				},
//			},
//			"enable_ipv6": schema.BoolAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "是否启用IPv6，默认为false",
//				Default:     booldefault.StaticBool(false),
//				PlanModifiers: []planmodifier.Bool{
//					boolplanmodifier.RequiresReplace(),
//				},
//			},
//			"plain_port": schema.Int32Attribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "公共接入点(PLAINTEXT)端口，范围在8000到9100之间，默认为8090",
//				Validators: []validator.Int32{
//					int32validator.Between(8000, 9100),
//				},
//				Default: int32default.StaticInt32(8090),
//				PlanModifiers: []planmodifier.Int32{
//					int32planmodifier.RequiresReplace(),
//				},
//			},
//			"sasl_port": schema.Int32Attribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "安全接入点(SASL_PLAINTEXT)端口，范围在8000到9100之间，默认为8092",
//				Validators: []validator.Int32{
//					int32validator.Between(8000, 9100),
//				},
//				Default: int32default.StaticInt32(8092),
//				PlanModifiers: []planmodifier.Int32{
//					int32planmodifier.RequiresReplace(),
//				},
//			},
//			"ssl_port": schema.Int32Attribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "SSL接入点(SASL_SSL)端口，范围在8000到9100之间，默认为8098。",
//				Validators: []validator.Int32{
//					int32validator.Between(8000, 9100),
//				},
//				Default: int32default.StaticInt32(8098),
//				PlanModifiers: []planmodifier.Int32{
//					int32planmodifier.RequiresReplace(),
//				},
//			},
//			"http_port": schema.Int32Attribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "HTTP接入点端口，范围在8000到9100之间，默认为8082",
//				Validators: []validator.Int32{
//					int32validator.Between(8000, 9100),
//				},
//				Default: int32default.StaticInt32(8082),
//				PlanModifiers: []planmodifier.Int32{
//					int32planmodifier.RequiresReplace(),
//				},
//			},
//			"retention_hours": schema.Int32Attribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "实例消息保留时长，默认为72小时，可选1~10000小时",
//				Validators: []validator.Int32{
//					int32validator.Between(1, 10000),
//				},
//				Default: int32default.StaticInt32(72),
//			},
//			"cycle_type": schema.StringAttribute{
//				Required:    true,
//				Description: "订购周期类型，取值范围：month：按月，on_demand：按需。当此值为month时，cycle_count为必填",
//				Validators: []validator.String{
//					stringvalidator.OneOf("month", "on_demand"),
//				},
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.UseStateForUnknown(),
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"cycle_count": schema.Int32Attribute{
//				Optional:    true,
//				Description: "订购时长，该参数在cycle_type为month时才生效，当cycleType=month，支持传递1、2、3、4、5、6、12、24、36",
//				Validators: []validator.Int32{
//					validator2.AlsoRequiresEqualInt32(
//						path.MatchRoot("cycle_type"),
//						types.StringValue(business.OrderCycleTypeMonth),
//					),
//					validator2.ConflictsWithEqualInt32(
//						path.MatchRoot("cycle_type"),
//						types.StringValue(business.OrderCycleTypeOnDemand),
//					),
//					int32validator.OneOf(1, 2, 3, 5, 6, 7, 12, 24, 36),
//				},
//				PlanModifiers: []planmodifier.Int32{
//					int32planmodifier.RequiresReplace(),
//				},
//			},
//			"auto_renew": schema.BoolAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "是否自动续订",
//				Default:     booldefault.StaticBool(false),
//				Validators: []validator.Bool{
//					validator2.ConflictsWithEqualBool(
//						path.MatchRoot("cycle_type"),
//						types.StringValue(business.OrderCycleTypeOnDemand),
//					),
//				},
//				PlanModifiers: []planmodifier.Bool{
//					boolplanmodifier.UseStateForUnknown(),
//					boolplanmodifier.RequiresReplace(),
//				},
//			},
//			"auto_renew_cycle_count": schema.Int32Attribute{
//				Optional:    true,
//				Description: "自动续订时长，支持自动续订范围：1-6月",
//				Validators: []validator.Int32{
//					validator2.AlsoRequiresEqualInt32(
//						path.MatchRoot("auto_renew"),
//						types.BoolValue(true),
//					),
//					validator2.ConflictsWithEqualInt32(
//						path.MatchRoot("auto_renew"),
//						types.BoolValue(false),
//					),
//					int32validator.Between(1, 6),
//				},
//				PlanModifiers: []planmodifier.Int32{
//					int32planmodifier.RequiresReplace(),
//				},
//			},
//		},
//	}
//}
//
//func (c *ctyunRabbitmqInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var plan CtyunRabbitmqInstanceConfig
//	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// 创建前检查
//	err = c.checkBeforeCreate(ctx, plan)
//	if err != nil {
//		return
//	}
//	// 创建
//	masterOrderID, err := c.create(ctx, plan)
//	if err != nil {
//		return
//	}
//	plan.MasterOrderID = types.StringValue(masterOrderID)
//	// 创建后检查
//	id, err := c.checkAfterCreate(ctx, plan)
//	if err != nil {
//		return
//	}
//	plan.ID = types.StringValue(id)
//
//	// 反查信息
//	err = c.getAndMerge(ctx, &plan)
//	if err != nil {
//		return
//	}
//
//	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
//}
//
//func (c *ctyunRabbitmqInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var state CtyunRabbitmqInstanceConfig
//	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// 查询远端
//	err = c.getAndMerge(ctx, &state)
//	if err != nil {
//		return
//	}
//
//	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
//}
//
//func (c *ctyunRabbitmqInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	// tf文件中的
//	var plan CtyunRabbitmqInstanceConfig
//	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// state中的
//	var state CtyunRabbitmqInstanceConfig
//	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	err = c.checkBeforeUpdate(ctx, plan, state)
//	if err != nil {
//		return
//	}
//	// 更新
//	err = c.update(ctx, plan, state)
//	if err != nil {
//		return
//	}
//	// 查询远端信息
//	err = c.getAndMerge(ctx, &state)
//	if err != nil {
//		return
//	}
//
//	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
//}
//
//func (c *ctyunRabbitmqInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var state CtyunRabbitmqInstanceConfig
//	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// 删除
//	err = c.delete(ctx, state)
//	if err != nil {
//		return
//	}
//	//response.State.RemoveResource(ctx)
//}
//
//func (c *ctyunRabbitmqInstance) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
//	if request.ProviderData == nil {
//		return
//	}
//	meta := request.ProviderData.(*common.CtyunMetadata)
//	c.meta = meta
//	c.vpcService = business.NewVpcService(meta)
//	c.sgService = business.NewSecurityGroupService(meta)
//}
//
//// 导入命令：terraform import [配置标识].[导入配置名称] [id],[regionID]
//func (c *ctyunRabbitmqInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var cfg CtyunRabbitmqInstanceConfig
//	var id, regionID string
//	err = terraform_extend.Split(request.ID, &id, &regionID)
//	if err != nil {
//		return
//	}
//	cfg.RegionID = types.StringValue(regionID)
//	cfg.ID = types.StringValue(id)
//	// 查询远端
//	err = c.getAndMerge(ctx, &cfg)
//	if err != nil {
//		return
//	}
//	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
//}
//
//// checkBeforeCreate 创建前检查
//func (c *ctyunRabbitmqInstance) checkBeforeCreate(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (err error) {
//	regionID, projectID := plan.RegionID.ValueString(), plan.ProjectID.ValueString()
//	vpc, subnetID, sgID := plan.VpcID.ValueString(), plan.SubnetID.ValueString(), plan.SecurityGroupID.ValueString()
//	subnets, err := c.vpcService.GetVpcSubnet(ctx, vpc, regionID, projectID)
//	if err != nil {
//		return err
//	}
//	_, exist := subnets[subnetID]
//	if !exist {
//		err = fmt.Errorf("子网不存在")
//		return err
//	}
//	err = c.sgService.MustExistInVpc(ctx, vpc, sgID, regionID)
//	if err != nil {
//		return err
//	}
//	err = c.checkSpecParams(ctx, plan)
//	if err != nil {
//		return err
//	}
//	return
//}
//
//// checkSpecParams 检查规格参数
//func (c *ctyunRabbitmqInstance) checkSpecParams(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (err error) {
//	nodeNum := plan.NodeNum.ValueInt32()
//	specName := plan.SpecName.ValueString()
//	diskType := plan.DiskType.ValueString()
//
//	if strings.HasSuffix(specName, "single") && nodeNum != 1 {
//		return fmt.Errorf("单机版实例节点数必须为1")
//	} else if strings.HasSuffix(specName, "cluster") && nodeNum < 3 {
//		return fmt.Errorf("集群版实例节点数必须大于等于3")
//	}
//	// 组装请求体
//	params := &ctgrabbitmq.CtgrabbitmqProdDetailRequest{
//		RegionId: plan.RegionID.ValueString(),
//	}
//	// 调用API
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqProdDetailApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//
//	var skuRes ctgrabbitmq.CtgrabbitmqProdDetailReturnObjResponseSkuResItem
//	var skuDisk ctgrabbitmq.CtgrabbitmqProdDetailReturnObjResponseSkuDiskItem
//	for _, s := range resp.ReturnObj.Data.Series {
//		for _, p := range s.Sku {
//			if p.ProdName == "集群版" && plan.NodeNum.ValueInt32() >= 3 {
//				skuRes = p.ResItem
//				skuDisk = p.DiskItem
//				break
//			} else if p.ProdName == "单机版" && plan.NodeNum.ValueInt32() == 1 {
//				skuRes = p.ResItem
//				skuDisk = p.DiskItem
//				break
//			}
//		}
//	}
//
//	var specAvailable bool
//	for _, r := range skuRes.ResItems {
//		for _, s := range r.Spec {
//			if s.SpecName == specName {
//				specAvailable = true
//				break
//			}
//		}
//		if specAvailable {
//			break
//		}
//	}
//	if !specAvailable {
//		return fmt.Errorf("本资源池不支持 %s", specName)
//	}
//
//	var diskAvailable bool
//	for _, d := range skuDisk.ResItems {
//		if d == diskType {
//			diskAvailable = true
//			break
//		}
//	}
//	if !diskAvailable {
//		return fmt.Errorf("本资源池不支持 %s", diskType)
//	}
//
//	return
//}
//
//// create 创建
//func (c *ctyunRabbitmqInstance) create(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (masterOrderID string, err error) {
//	switch plan.CycleType.ValueString() {
//	case business.OrderCycleTypeMonth:
//		return c.createPrePayOrder(ctx, plan)
//	case business.OrderCycleTypeOnDemand:
//		return c.createPostPayOrder(ctx, plan)
//	}
//	return
//}
//
//// createPrePayOrder 创建包年包月
//func (c *ctyunRabbitmqInstance) createPrePayOrder(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (masterOrderID string, err error) {
//	params := &ctgrabbitmq.CtgrabbitmqCreateOrderRequest{
//		RegionId:            plan.RegionID.ValueString(),
//		ProjectId:           plan.ProjectID.ValueString(),
//		CycleCnt:            plan.CycleCount.ValueInt32(),
//		ClusterName:         plan.InstanceName.ValueString(),
//		EngineVersion:       plan.EngineVersion.ValueString(),
//		SpecName:            plan.SpecName.ValueString(),
//		NodeNum:             plan.NodeNum.ValueInt32(),
//		DiskType:            plan.DiskType.ValueString(),
//		DiskSize:            plan.DiskSize.ValueInt32(),
//		VpcId:               plan.VpcID.ValueString(),
//		SubnetId:            plan.SubnetID.ValueString(),
//		SecurityGroupId:     plan.SecurityGroupID.ValueString(),
//		EnableIpv6:          plan.EnableIpv6.ValueBoolPointer(),
//		PlainPort:           plan.PlainPort.ValueInt32(),
//		SaslPort:            plan.SaslPort.ValueInt32(),
//		SslPort:             plan.SslPort.ValueInt32(),
//		HttpPort:            plan.HttpPort.ValueInt32(),
//		RetentionHours:      plan.RetentionHours.ValueInt32(),
//		AutoRenewStatus:     plan.AutoRenew.ValueBoolPointer(),
//		AutoRenewCycleCount: plan.AutoRenewCycleCount.ValueInt32(),
//	}
//
//	var zoneList []string
//	var str []types.String
//	plan.ZoneList.ElementsAs(ctx, &str, true)
//	for _, s := range str {
//		zoneList = append(zoneList, s.ValueString())
//	}
//	params.ZoneList = zoneList
//
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqCreateOrderApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	masterOrderID = resp.ReturnObj.Data.NewOrderId
//	return
//}
//
//// createPostPayOrder 创建按需
//func (c *ctyunRabbitmqInstance) createPostPayOrder(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (masterOrderID string, err error) {
//	params := &ctgrabbitmq.CtgrabbitmqCreatePostPayOrderRequest{
//		RegionId:        plan.RegionID.ValueString(),
//		ProjectId:       plan.ProjectID.ValueString(),
//		ClusterName:     plan.InstanceName.ValueString(),
//		EngineVersion:   plan.EngineVersion.ValueString(),
//		SpecName:        plan.SpecName.ValueString(),
//		NodeNum:         plan.NodeNum.ValueInt32(),
//		DiskType:        plan.DiskType.ValueString(),
//		DiskSize:        plan.DiskSize.ValueInt32(),
//		VpcId:           plan.VpcID.ValueString(),
//		SubnetId:        plan.SubnetID.ValueString(),
//		SecurityGroupId: plan.SecurityGroupID.ValueString(),
//		EnableIpv6:      plan.EnableIpv6.ValueBoolPointer(),
//		PlainPort:       plan.PlainPort.ValueInt32(),
//		SaslPort:        plan.SaslPort.ValueInt32(),
//		SslPort:         plan.SslPort.ValueInt32(),
//		HttpPort:        plan.HttpPort.ValueInt32(),
//		RetentionHours:  plan.RetentionHours.ValueInt32(),
//	}
//
//	var zoneList []string
//	var strings []types.String
//	plan.ZoneList.ElementsAs(ctx, &strings, true)
//	for _, s := range strings {
//		zoneList = append(zoneList, s.ValueString())
//	}
//	params.ZoneList = zoneList
//
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqCreatePostPayOrderApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	masterOrderID = resp.ReturnObj.Data.NewOrderId
//	return
//}
//
//// getAndMerge 从远端查询
//func (c *ctyunRabbitmqInstance) getAndMerge(ctx context.Context, plan *CtyunRabbitmqInstanceConfig) (err error) {
//	instance, err := c.getByNameOrID(ctx, *plan)
//	if err != nil {
//		return
//	}
//	plan.InstanceName = types.StringValue(instance.InstanceName)
//	if len(instance.Version) >= 3 {
//		plan.EngineVersion = types.StringValue(instance.Version[:3])
//	}
//	plan.SpecName = types.StringValue(instance.Specifications)
//	plan.NodeNum = types.Int32Value(int32(len(instance.NodeList)))
//
//	plan.DiskType = types.StringValue(instance.DiskType)
//	plan.DiskSize = types.Int32Value(utils.StringToInt32Must(instance.Space))
//	plan.VpcID = types.StringValue(instance.VpcId)
//	plan.SubnetID = types.StringValue(instance.SubnetId)
//
//	plan.EnableIpv6 = types.BoolValue(map[int32]bool{1: true, 0: false}[instance.Ipv6Enable])
//	if len(instance.NodeList) > 0 {
//		plan.PlainPort = types.Int32Value(utils.StringToInt32Must(instance.NodeList[0].VpcPort))
//		plan.SaslPort = types.Int32Value(utils.StringToInt32Must(instance.NodeList[0].SaslPort))
//		plan.SslPort = types.Int32Value(utils.StringToInt32Must(instance.NodeList[0].ListenNodePort))
//		plan.HttpPort = types.Int32Value(utils.StringToInt32Must(instance.NodeList[0].HttpPort))
//	}
//
//	config, err := c.getInstanceConfig(ctx, *plan)
//	if err != nil {
//		return
//	}
//	plan.RetentionHours = types.Int32Value(utils.StringToInt32Must(config["log.retention.hours"].Value))
//	if plan.ZoneList.IsNull() {
//		plan.ZoneList = types.SetNull(types.StringType)
//	}
//
//	return
//	// 下列字段没有地方获取
//	//CycleType
//	//CycleCount
//	//AutoRenew
//	//AutoRenewCycleCount
//	//SecurityGroupID
//}
//
//func (c *ctyunRabbitmqInstance) checkBeforeUpdate(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	err = c.checkSpecParams(ctx, plan)
//	if err != nil {
//		return
//	}
//	if state.NodeNum.ValueInt32() == 1 && !plan.NodeNum.Equal(state.NodeNum) {
//		return fmt.Errorf("单机版实例不能进行节点扩缩容操作")
//	}
//	if strings.HasSuffix(plan.SpecName.ValueString(), "single") && strings.HasSuffix(state.SpecName.ValueString(), "cluster") ||
//		strings.HasSuffix(plan.SpecName.ValueString(), "cluster") && strings.HasSuffix(state.SpecName.ValueString(), "single") {
//		return fmt.Errorf("不支持单机版和集群版互相变更")
//	}
//	instance, err := c.getByNameOrID(ctx, state)
//	if err != nil {
//		return
//	}
//	if instance.Status != 1 {
//		return fmt.Errorf("请在实例处于运行中状态时再进行更新操作")
//	}
//
//	return nil
//}
//
//// update 更新
//func (c *ctyunRabbitmqInstance) update(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	err = c.updateRetentionHours(ctx, plan, state)
//	if err != nil {
//		return
//	}
//	err = c.updateName(ctx, plan, state)
//	if err != nil {
//		return
//	}
//	err = c.updateDiskSize(ctx, plan, state)
//	if err != nil {
//		return
//	}
//	err = c.updateNodeNum(ctx, plan, state)
//	if err != nil {
//		return
//	}
//	err = c.updateSpec(ctx, plan, state)
//	if err != nil {
//		return
//	}
//	return
//}
//
//// updateDiskSize 更新磁盘大小
//func (c *ctyunRabbitmqInstance) updateDiskSize(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	if plan.DiskSize.Equal(state.DiskSize) {
//		return
//	}
//	if plan.DiskSize.ValueInt32() > state.DiskSize.ValueInt32() {
//		err = c.diskExtend(ctx, plan, state)
//	} else {
//		err = fmt.Errorf("目前不支持磁盘缩容")
//		//err = c.diskShrink(ctx, plan, state)
//	}
//	if err != nil {
//		return
//	}
//	return c.checkAfterUpdateDiskSize(ctx, plan, state)
//}
//
//// diskExtend 磁盘缩容
//func (c *ctyunRabbitmqInstance) diskShrink(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	params := &ctgrabbitmq.CtgrabbitmqDiskShrinkRequest{
//		RegionId:   state.RegionID.ValueString(),
//		ProdInstId: state.ID.ValueString(),
//		DiskSize:   plan.DiskSize.String(),
//	}
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqDiskShrinkApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	return
//}
//
//// diskExtend 磁盘扩容
//func (c *ctyunRabbitmqInstance) diskExtend(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	autoPay := true
//	params := &ctgrabbitmq.CtgrabbitmqDiskExtendRequest{
//		RegionId:       state.RegionID.ValueString(),
//		ProdInstId:     state.ID.ValueString(),
//		DiskExtendSize: plan.DiskSize.String(),
//		AutoPay:        &autoPay,
//	}
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqDiskExtendApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	return
//}
//
//// checkAfterUpdateDiskSize 检查磁盘大小是否变更成功
//func (c *ctyunRabbitmqInstance) checkAfterUpdateDiskSize(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	var executeSuccessFlag bool
//	retryer, _ := business.NewRetryer(time.Second*10, 180)
//	retryer.Start(
//		func(currentTime int) bool {
//			var instance *ctgrabbitmq.CtgrabbitmqInstQueryReturnObjDataResponse
//			instance, err = c.getByNameOrID(ctx, state)
//			if err != nil {
//				return false
//			}
//			if utils.StringToInt32Must(instance.Space) != plan.DiskSize.ValueInt32() || instance.Status != 1 {
//				return true
//			}
//			time.Sleep(30 * time.Second)
//			executeSuccessFlag = true
//			return false
//		})
//	if err != nil {
//		return
//	}
//	if !executeSuccessFlag {
//		err = fmt.Errorf("磁盘变配时间过长")
//	}
//	return
//}
//
//// updateNodeNum 更新节点数量
//func (c *ctyunRabbitmqInstance) updateNodeNum(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	if plan.NodeNum.Equal(state.NodeNum) {
//		return
//	}
//	if plan.NodeNum.ValueInt32() > state.NodeNum.ValueInt32() {
//		err = c.nodeExtend(ctx, plan, state)
//	} else {
//		err = fmt.Errorf("目前不支持节点缩容")
//		//err = c.nodeShrink(ctx, plan, state)
//	}
//	if err != nil {
//		return
//	}
//	return c.checkAfterUpdateNodeNum(ctx, plan, state)
//}
//
//// nodeShrink 节点缩容
//func (c *ctyunRabbitmqInstance) nodeShrink(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	params := &ctgrabbitmq.CtgrabbitmqNodeShrinkRequest{
//		RegionId:   state.RegionID.ValueString(),
//		ProdInstId: state.ID.ValueString(),
//		NodeNum:    plan.DiskSize.String(),
//	}
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqNodeShrinkApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	return
//}
//
//// nodeExtend 节点扩容
//func (c *ctyunRabbitmqInstance) nodeExtend(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	autoPay := true
//	params := &ctgrabbitmq.CtgrabbitmqNodeExtendRequest{
//		RegionId:      state.RegionID.ValueString(),
//		ProdInstId:    state.ID.ValueString(),
//		ExtendNodeNum: plan.NodeNum.ValueInt32(),
//		AutoPay:       &autoPay,
//	}
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqNodeExtendApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	return
//}
//
//// checkAfterUpdateNodeNum 检查节点数量是否变更成功
//func (c *ctyunRabbitmqInstance) checkAfterUpdateNodeNum(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	var executeSuccessFlag bool
//	retryer, _ := business.NewRetryer(time.Second*10, 180)
//	retryer.Start(
//		func(currentTime int) bool {
//			var instance *ctgrabbitmq.CtgrabbitmqInstQueryReturnObjDataResponse
//			instance, err = c.getByNameOrID(ctx, state)
//			if err != nil {
//				return false
//			}
//			if len(instance.NodeList) != int(plan.NodeNum.ValueInt32()) || instance.Status != 1 {
//				return true
//			}
//			time.Sleep(30 * time.Second)
//			executeSuccessFlag = true
//			return false
//		})
//	if err != nil {
//		return
//	}
//	if !executeSuccessFlag {
//		err = fmt.Errorf("节点数量变配时间过长")
//	}
//	return
//}
//
//// updateSpec 更新规格
//func (c *ctyunRabbitmqInstance) updateSpec(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	if plan.SpecName.Equal(state.SpecName) {
//		return
//	}
//	planU, err := c.parseSpec(plan.SpecName.ValueString())
//	if err != nil {
//		return
//	}
//	stateU, err := c.parseSpec(state.SpecName.ValueString())
//	if err != nil {
//		return
//	}
//	if planU > stateU {
//		err = c.specExtend(ctx, plan, state)
//	} else {
//		err = c.specShrink(ctx, plan, state)
//	}
//	if err != nil {
//		return
//	}
//	return c.checkAfterUpdateSpec(ctx, plan, state)
//}
//
//// parseSpec 从规格名称解析cpu
//func (c *ctyunRabbitmqInstance) parseSpec(s string) (u int, err error) {
//	re := regexp.MustCompile(`(\d+)u(\d+)g`)
//	matches := re.FindStringSubmatch(s)
//	if len(matches) != 3 {
//		err = fmt.Errorf("invalid format: %s", s)
//		return
//	}
//
//	if _, err = fmt.Sscanf(matches[1], "%d", &u); err != nil {
//		return
//	}
//	return
//}
//
//// specShrink 规格缩容
//func (c *ctyunRabbitmqInstance) specShrink(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	params := &ctgrabbitmq.CtgrabbitmqSpecShrinkRequest{
//		RegionId:   state.RegionID.ValueString(),
//		ProdInstId: state.ID.ValueString(),
//		SpecName:   plan.SpecName.ValueString(),
//	}
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqSpecShrinkApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	return
//}
//
//// specExtend 规格扩容
//func (c *ctyunRabbitmqInstance) specExtend(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	autoPay := true
//	params := &ctgrabbitmq.CtgrabbitmqSpecExtendRequest{
//		RegionId:   state.RegionID.ValueString(),
//		ProdInstId: state.ID.ValueString(),
//		SpecName:   plan.SpecName.ValueString(),
//		AutoPay:    &autoPay,
//	}
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqSpecExtendApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	return
//}
//
//// checkAfterUpdateSpec 检查规格是否变更成功
//func (c *ctyunRabbitmqInstance) checkAfterUpdateSpec(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	var executeSuccessFlag bool
//	retryer, _ := business.NewRetryer(time.Second*10, 180)
//	retryer.Start(
//		func(currentTime int) bool {
//			var instance *ctgrabbitmq.CtgrabbitmqInstQueryReturnObjDataResponse
//			instance, err = c.getByNameOrID(ctx, state)
//			if err != nil {
//				return false
//			}
//			if instance.Specifications != plan.SpecName.ValueString() || instance.Status != 1 {
//				return true
//			}
//			time.Sleep(30 * time.Second)
//			executeSuccessFlag = true
//			return false
//		})
//	if err != nil {
//		return
//	}
//	if !executeSuccessFlag {
//		err = fmt.Errorf("规格变配时间过长")
//	}
//	return
//}
//
//// updateName 更新实例名称
//func (c *ctyunRabbitmqInstance) updateName(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	if plan.InstanceName.Equal(state.InstanceName) {
//		return
//	}
//	params := &ctgrabbitmq.CtgrabbitmqInstancesModifyNameV3Request{
//		RegionId:     state.RegionID.ValueString(),
//		ProdInstId:   state.ID.ValueString(),
//		InstanceName: plan.InstanceName.ValueString(),
//	}
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqInstancesModifyNameV3Api.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		return fmt.Errorf("API return error. Message: %s", resp.Message)
//	} else if resp.ReturnObj == nil {
//		return common.InvalidReturnObjError
//	} else if resp.ReturnObj.Data != "modify success" {
//		return fmt.Errorf("API return error. Data: %s", resp.ReturnObj.Data)
//	}
//	return
//}
//
//// updateRetentionHours 更新实例消息保留时长
//func (c *ctyunRabbitmqInstance) updateRetentionHours(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	if plan.RetentionHours.Equal(state.RetentionHours) {
//		return
//	}
//	params := &ctgrabbitmq.CtgrabbitmqUpdateInstanceConfigRequest{
//		RegionId:   state.RegionID.ValueString(),
//		ProdInstId: state.ID.ValueString(),
//		StaticConfigs: []*ctgrabbitmq.CtgrabbitmqUpdateInstanceConfigStaticConfigsRequest{
//			{Name: "log.retention.hours", Value: plan.RetentionHours.String()},
//		},
//	}
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqUpdateInstanceConfigApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		return fmt.Errorf("API return error. Message: %s", resp.Message)
//	} else if resp.ReturnObj == nil {
//		return common.InvalidReturnObjError
//	} else if resp.ReturnObj.Data != "modify success" {
//		return fmt.Errorf("API return error. Data: %s", resp.ReturnObj.Data)
//	}
//	// 更新后需要重启
//	return c.reboot(ctx, plan, state)
//}
//
//// reboot 重启实例
//func (c *ctyunRabbitmqInstance) reboot(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
//	params := &ctgrabbitmq.CtgrabbitmqInstancesRestartV3Request{
//		RegionId:   state.RegionID.ValueString(),
//		ProdInstId: state.ID.ValueString(),
//	}
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqInstancesRestartV3Api.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	// 等待重启完成
//	var executeSuccessFlag bool
//	retryer, _ := business.NewRetryer(time.Second*10, 60)
//	retryer.Start(
//		func(currentTime int) bool {
//			var instance *ctgrabbitmq.CtgrabbitmqInstQueryReturnObjDataResponse
//			instance, err = c.getByNameOrID(ctx, state)
//			if err != nil {
//				return false
//			}
//			if instance.Status != 1 {
//				return true
//			}
//			executeSuccessFlag = true
//			return false
//		})
//	if err != nil {
//		return
//	}
//	if !executeSuccessFlag {
//		err = fmt.Errorf("重启时间过长")
//	}
//	return
//}
//
//// delete 删除
//func (c *ctyunRabbitmqInstance) delete(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (err error) {
//	params := &ctgrabbitmq.CtgrabbitmqUnsubscribeInstV3Request{
//		RegionId:   plan.RegionID.ValueString(),
//		ProdInstId: plan.ID.ValueString(),
//	}
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqUnsubscribeInstV3Api.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	return
//}
//
//// checkAfterCreate 创建后检查
//func (c *ctyunRabbitmqInstance) checkAfterCreate(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (id string, err error) {
//	var executeSuccessFlag bool
//	retryer, _ := business.NewRetryer(time.Second*10, 180)
//	retryer.Start(
//		func(currentTime int) bool {
//			var instance *ctgrabbitmq.CtgrabbitmqInstQueryReturnObjDataResponse
//			instance, err = c.getByNameOrID(ctx, plan)
//			if err != nil {
//				return false
//			}
//			if instance == nil || instance.Status != 1 || instance.ProdInstId == "" {
//				return true
//			}
//			// 等待订单完成
//			time.Sleep(30 * time.Second)
//			id = instance.ProdInstId
//			executeSuccessFlag = true
//			return false
//		})
//	if err != nil {
//		return
//	}
//	if !executeSuccessFlag {
//		err = fmt.Errorf("创建时间过长")
//	}
//	return
//}
//
//// checkAfterDelete 删除后检查
//func (c *ctyunRabbitmqInstance) checkAfterDelete(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (err error) {
//	var executeSuccessFlag bool
//	retryer, _ := business.NewRetryer(time.Second*10, 180)
//	retryer.Start(
//		func(currentTime int) bool {
//			var instance *ctgrabbitmq.CtgrabbitmqInstQueryReturnObjDataResponse
//			instance, err = c.getByNameOrID(ctx, plan)
//			if err != nil {
//				return false
//			}
//			if instance != nil && instance.Status != 5 {
//				return true
//			}
//			executeSuccessFlag = true
//			return false
//		})
//	if err != nil {
//		return
//	}
//	if !executeSuccessFlag {
//		err = fmt.Errorf("删除时间过长")
//	}
//	return
//}
//
//// getByNameOrID 根据ID或名称查询集群
//func (c *ctyunRabbitmqInstance) getByNameOrID(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (instance *ctgrabbitmq.CtgrabbitmqInstQueryReturnObjDataResponse, err error) {
//	params := &ctgrabbitmq.CtgrabbitmqInstQueryRequest{
//		RegionId:       plan.RegionID.ValueString(),
//		OuterProjectId: plan.ProjectID.ValueString(),
//	}
//
//	if plan.ID.ValueString() != "" {
//		params.ProdInstId = plan.ID.ValueString()
//	} else if plan.InstanceName.ValueString() != "" {
//		params.Name = plan.InstanceName.ValueString()
//		e := true
//		params.ExactMatchName = &e
//	}
//
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqInstQueryApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	if len(resp.ReturnObj.Data) > 0 {
//		instance = resp.ReturnObj.Data[0]
//		if instance == nil {
//			err = common.InvalidReturnObjResultsError
//		}
//	}
//	return
//}
//
//// getInstanceConfig 获取实例配置
//func (c *ctyunRabbitmqInstance) getInstanceConfig(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (attr map[string]*ctgrabbitmq.CtgrabbitmqGetInstanceConfigReturnObjDataResponse, err error) {
//	params := &ctgrabbitmq.CtgrabbitmqGetInstanceConfigRequest{
//		RegionId:   plan.RegionID.ValueString(),
//		ProdInstId: plan.ID.ValueString(),
//	}
//
//	resp, err := c.meta.Apis.SdkRabbitmqApis.CtgrabbitmqGetInstanceConfigApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	attr = map[string]*ctgrabbitmq.CtgrabbitmqGetInstanceConfigReturnObjDataResponse{}
//	for _, d := range resp.ReturnObj.Data {
//		attr[d.Name] = d
//	}
//	return
//}
