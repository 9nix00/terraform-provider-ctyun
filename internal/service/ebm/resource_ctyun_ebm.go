package ebm

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-ctyun/internal/core/ctebm"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
	"terraform-provider-ctyun/internal/utils"
	"time"

	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

var (
	_ resource.Resource                = &ctyunEbm{}
	_ resource.ResourceWithConfigure   = &ctyunEbm{}
	_ resource.ResourceWithImportState = &ctyunEbm{}
)

type ctyunEbm struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbm() resource.Resource {
	return &ctyunEbm{}
}

func (c *ctyunEbm) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebm"
}

type CtyunEbmConfig struct {
	ID                   types.String `tfsdk:"id"`
	InstanceID           types.String `tfsdk:"instance_id"`
	RegionID             types.String `tfsdk:"region_id"`
	AzName               types.String `tfsdk:"az_name"`
	DeviceType           types.String `tfsdk:"device_type"`
	InstanceName         types.String `tfsdk:"instance_name"`
	Hostname             types.String `tfsdk:"hostname"`
	ImageUUID            types.String `tfsdk:"image_uuid"`
	Password             types.String `tfsdk:"password"`
	ProjectID            types.String `tfsdk:"project_id"`
	SystemVolumeRaidUUID types.String `tfsdk:"system_volume_raid_uuid"`
	DataVolumeRaidUUID   types.String `tfsdk:"data_volume_raid_uuid"`
	VpcID                types.String `tfsdk:"vpc_id"`
	ExtIP                types.String `tfsdk:"ext_ip"`
	IpType               types.String `tfsdk:"ip_type"`
	BandWidth            types.Int32  `tfsdk:"band_width"`
	PublicIP             types.String `tfsdk:"public_ip"`
	SecurityGroupIDs     types.Set    `tfsdk:"security_group_ids"`
	DiskList             types.List   `tfsdk:"disk_list"`
	NetworkCardList      types.List   `tfsdk:"network_card_list"`
	UserData             types.String `tfsdk:"user_data"`
	KeyPairName          types.String `tfsdk:"key_pair_name"`
	AutoRenew            types.Bool   `tfsdk:"auto_renew"`
	CycleCount           types.Int32  `tfsdk:"cycle_count"`
	CycleType            types.String `tfsdk:"cycle_type"`
	MasterOrderID        types.String `tfsdk:"master_order_id"`
	Status               types.String `tfsdk:"status"`
}

type CtyunEbmDiskList struct {
	DiskType types.String `tfsdk:"disk_type"`
	Title    types.String `tfsdk:"title"`
	Type     types.String `tfsdk:"type"`
	Size     types.Int64  `tfsdk:"size"`
}

type CtyunEbmNetworkCardList struct {
	FixedIP     types.String `tfsdk:"fixed_ip"`
	Master      types.Bool   `tfsdk:"master"`
	Ipv6        types.String `tfsdk:"ipv6"`
	SubnetID    types.String `tfsdk:"subnet_id"`
	PortID      types.String `tfsdk:"port_id"`
	InterfaceID types.String `tfsdk:"interface_id"`
}

func (c *ctyunEbm) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027724**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID",
			},
			"instance_id": schema.StringAttribute{
				Computed:    true,
				Description: "物理机UUID",
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订购的受理单id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称",
				Default:     defaults.AcquireFromGlobalString(common.ExtraAzName, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device_type": schema.StringAttribute{
				Required:    true,
				Description: "物理机套餐类型",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_name": schema.StringAttribute{
				Required:    true,
				Description: "物理机名称，长度为2-31位",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 31),
				},
			},
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "hostname，linux系统2到63位长度；windows系统2-15位长度；<br/>允许使用大小写字母、数字、连字符'-'、点号'.'，不能连续使用'-'或者'.'，'-'和'.'不能用于开头或结尾，不能仅使用数字",
			},
			"image_uuid": schema.StringAttribute{
				Required:    true,
				Description: "物理机镜像id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Sensitive:   true,
				Optional:    true,
				Computed:    true,
				Description: "密码(必须包含大小写字母和（一个数字或者特殊字符）长度8到30位)，未传入有效的keyName时必须传入password",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("key_pair_name"),
					}...),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
			},
			"system_volume_raid_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "本地系统盘raid类型，如果有本地盘则必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data_volume_raid_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "本地数据盘raid类型，如果有本地盘则必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "主网卡网络ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ext_ip": schema.StringAttribute{
				Required:    true,
				Description: "是否使用弹性公网IP，取值范围:[自动分配:auto_assign,不使用:not_use,使用已有:use_exist]",
				Validators: []validator.String{
					stringvalidator.OneOf(business.EbmExtIp...),
				},
			},
			"ip_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "弹性IP版本，取值范围:[ipv4=v4地址,ipv6=v6地址]，默认值:ipv4",
				Default:     stringdefault.StaticString("ipv4"),
			},
			"band_width": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "带宽，取值范围:[1~2000]，默认值:100",
				Default:     int32default.StaticInt32(0),
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("ext_ip"),
						types.StringValue(business.EbmExtIpAutoAssign),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("ext_ip"),
						types.StringValue(business.EbmExtIpNotUse),
						types.StringValue(business.EbmExtIpUseExist),
					),
					int32validator.Between(1, 2000),
				},
			},
			"public_ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "弹性公网IP的id",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("ext_ip"),
						types.StringValue(business.EbmExtIpUseExist),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("ext_ip"),
						types.StringValue(business.EbmExtIpNotUse),
						types.StringValue(business.EbmExtIpAutoAssign),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_group_ids": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				Description: "安全组ID，套餐smartNicExist为true可支持安全组。创建弹性裸金属必须传入安全组ID，标准裸金属不支持传入安全组ID",
				ElementType: types.StringType,
			},
			"disk_list": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "云盘信息列表，套餐中supportCloud为true表示支持云盘",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
					listplanmodifier.RequiresReplace(),
				},
				Default: listdefault.StaticValue(types.ListValueMust(
					utils.StructToTFObjectTypes(CtyunEbmDiskList{}),
					[]attr.Value{},
				)),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"disk_type": schema.StringAttribute{
							Required:    true,
							Description: "磁盘类型system或data，套餐中cloudBoot为true表示支持云盘系统盘",
							Validators: []validator.String{
								stringvalidator.OneOf(business.EbmSystemDiskType, "data"),
							},
						},
						"title": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "磁盘名称，长度2~64,不支持中文",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "磁盘分类，取值范围:[SAS=SAS盘,SATA=SATA盘,SSD-genric=SSD-genric盘,SSD=SSD盘]",
							Validators: []validator.String{
								stringvalidator.OneOf(business.EbsDiskTypes...),
							},
						},
						"size": schema.Int64Attribute{
							Required:    true,
							Description: "磁盘容量",
						},
					},
				},
			},
			"network_card_list": schema.ListNestedAttribute{
				Required:    true,
				Description: "网卡",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"port_id": schema.StringAttribute{
							Computed:    true,
							Description: "PORT UUID",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"interface_id": schema.StringAttribute{
							Computed:    true,
							Description: "网卡UUID",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"fixed_ip": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "内网IPv4地址",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"master": schema.BoolAttribute{
							Required:    true,
							Description: "是否主节点(True代表主节点)",
						},
						"ipv6": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "内网IPv6地址",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"subnet_id": schema.StringAttribute{
							Required:    true,
							Description: "子网id",
						},
					},
				},
			},
			"user_data": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "用户自定义数据,需要以Base64方式编码,Base64编码后的长度限制为1-16384字符",
				Default:     stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 16384),
				},
			},
			"key_pair_name": schema.StringAttribute{
				Optional:    true,
				Description: "密钥对名词",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("password"),
					}...),
				},
			},
			"auto_renew": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否自动续订，默认非自动续订。",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围:[on_demand=按需,month=按月,year=按年]，cycleType与cycleCount一起填写",
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypeOnDemand, business.OrderCycleTypeYear, business.OrderCycleTypeMonth),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int32Attribute{
				Optional:    true,
				Description: "订购时长，最长订购周期为60个月（5年）；非按需时必填",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeYear),
						types.StringValue(business.OrderCycleTypeMonth),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
				},
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "物理机状态",
				Validators: []validator.String{
					stringvalidator.OneOf(
						business.EbmStatusRunning,
						business.EbmStatusStopped,
					),
				},
				Default: stringdefault.StaticString(business.EbmStatusRunning),
			},
		},
	}
}

func (c *ctyunEbm) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbmConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建前检查
	err = c.checkBeforeCreateInstance(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	returnObj, err := c.createInstance(ctx, plan)
	if err != nil {
		return
	}

	// 先保存订单号
	masterOrderId := *returnObj.MasterOrderID
	plan.MasterOrderID = types.StringValue(masterOrderId)
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 根据订单号轮询查资源的uuid
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	loop, err := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId, 600)
	if err != nil {
		return
	}
	plan.InstanceID = types.StringValue(loop.Uuid[0])
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建机器后状态默认为启动状态，可根据用户要求的状态，去执行对应的操作，比如关机
	err = c.handleInstance(ctx, plan, business.EbmStatusRunning, plan.Status.ValueString())
	if err != nil {
		return
	}
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *ctyunEbm) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbmConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 通过订单号同步
	if !c.acquireAndSetIdIfOrderNotFinished(ctx, &state, response) {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "instance is not found") {
			// 查下主网卡是否存在
			var exist bool
			portID := c.getMasterPortID(ctx, state)
			exist, err = business.NewPortService(c.meta).Exist(ctx, portID, state.RegionID.ValueString())
			if err != nil {
				return
			}
			// 主网卡存在，则监听到主网卡删除为止
			if exist {
				err = c.checkAfterDelete(ctx, state)
				if err != nil {
					return
				}
			}
			// 主网卡不存在，清理state
			response.State.RemoveResource(ctx)
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEbm) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunEbmConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunEbmConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	if !plan.NetworkCardList.Equal(state.NetworkCardList) {
		err = fmt.Errorf("请使用其他能力修改网卡")
		return
	}
	// 处理开关机
	err = c.handleInstance(ctx, state, state.Status.ValueString(), plan.Status.ValueString())
	if err != nil {
		return
	}
	// 修改实例名称
	err = c.updateInstanceName(ctx, state, plan)
	if err != nil {
		return
	}
	// 修改密码或主机名
	err = c.updatePasswordOrHostname(ctx, state, plan)
	if err != nil {
		return
	}
	state.Password = plan.Password
	if err != nil {
		return
	}
	state.CycleCount = plan.CycleCount
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEbm) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbmConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 关机
	err = c.handleInstance(ctx, state, state.Status.ValueString(), business.EbmStatusStopped)
	if err != nil {
		return
	}
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDelete(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunEbm) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [instanceID],[regionID],[azName]
func (c *ctyunEbm) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbmConfig
	var instanceUUID, regionID, azName string
	err = terraform_extend.Split(request.ID, &instanceUUID, &regionID, &azName)
	if err != nil {
		return
	}

	plan.InstanceID = types.StringValue(instanceUUID)
	plan.AzName = types.StringValue(azName)
	plan.RegionID = types.StringValue(regionID)
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

// createInstance 创建物理机
func (c *ctyunEbm) createInstance(ctx context.Context, plan CtyunEbmConfig) (returnObj ctebm.EbmCreateInstanceV4plusReturnObjResponse, err error) {
	regionID := plan.RegionID.ValueString()
	projectID := plan.ProjectID.ValueString()
	azName := plan.AzName.ValueString()
	publicIP := plan.PublicIP.ValueString()
	password := plan.Password.ValueString()
	systemVolumeRaidUUID := plan.SystemVolumeRaidUUID.ValueString()
	dataVolumeRaidUUID := plan.DataVolumeRaidUUID.ValueString()
	ipType := plan.IpType.ValueString()
	userData := plan.UserData.ValueString()
	keyName := plan.KeyPairName.ValueString()
	bandwidth := plan.BandWidth.ValueInt32()
	securityGroupIDs, _ := c.buildSecGroupList(ctx, plan)
	securityGroupStr := strings.Join(securityGroupIDs, ",")
	diskList, _ := c.buildDiskList(ctx, plan)
	networkCardList, _ := c.buildNetworkCardList(ctx, plan)
	extIp, _ := business.EbmExtIpMap.FromOriginalScene(plan.ExtIP.ValueString(), business.EbmExtIpMapScene1)
	params := &ctebm.EbmCreateInstanceV4plusRequest{
		RegionID:        regionID,
		AzName:          azName,
		DeviceType:      plan.DeviceType.ValueString(),
		InstanceName:    plan.InstanceName.ValueString(),
		Hostname:        plan.Hostname.ValueString(),
		ImageUUID:       plan.ImageUUID.ValueString(),
		VpcID:           plan.VpcID.ValueString(),
		ExtIP:           extIp.(string),
		ProjectID:       &projectID,
		IpType:          &ipType,
		DiskList:        diskList,
		NetworkCardList: networkCardList,
		AutoRenewStatus: map[bool]int32{false: 0, true: 1}[plan.AutoRenew.ValueBool()],
		ClientToken:     uuid.NewString(),
		OrderCount:      1,
		SecurityGroupID: &securityGroupStr,
	}
	if password != "" {
		params.Password = &password
	} else if keyName != "" {
		params.KeyName = &keyName
	} else {
		err = fmt.Errorf("password or keyname is empty")
	}
	if userData != "" {
		params.UserData = &userData
	}
	if systemVolumeRaidUUID != "" {
		params.SystemVolumeRaidUUID = &systemVolumeRaidUUID
	}
	if dataVolumeRaidUUID != "" {
		params.DataVolumeRaidUUID = &dataVolumeRaidUUID
	}

	if bandwidth > 0 {
		params.BandWidth = bandwidth
	}

	if publicIP != "" {
		params.PublicIP = &publicIP
	}

	switch plan.CycleType.ValueString() {
	case business.OrderCycleTypeOnDemand:
		params.InstanceChargeType = business.EbmOrderOnDemand
	case business.OrderCycleTypeMonth, business.OrderCycleTypeYear:
		params.InstanceChargeType = business.EbmOrderOnCycle
		params.CycleType = strings.ToUpper(plan.CycleType.ValueString())
		params.CycleCount = plan.CycleCount.ValueInt32()
	}

	resp, err := c.meta.Apis.CtEbmApis.EbmCreateInstanceV4plusApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	returnObj = *resp.ReturnObj
	return
}

// checkBeforeCreateInstance 创建前检查
func (c *ctyunEbm) checkBeforeCreateInstance(ctx context.Context, plan CtyunEbmConfig) error {
	cycleCount := plan.CycleCount.ValueInt32()
	cycleType := plan.CycleType.ValueString()
	if cycleType == business.OrderCycleTypeMonth && cycleCount > 11 ||
		cycleType == business.OrderCycleTypeYear && cycleCount > 5 {
		return fmt.Errorf("创建包周期物理机时，以月为单位，最长支持11月；以年为单位，最长支持5年")
	}

	// 确保当前虚拟私有云存在，且子网与虚拟私有云存在对应关系
	vpc := plan.VpcID.ValueString()
	subnets, err := business.NewVpcService(c.meta).GetVpcSubnet(ctx, vpc, plan.RegionID.ValueString(), plan.ProjectID.ValueString())
	if err != nil {
		return err
	}
	// 查询套餐
	deviceTypeConfig, err := c.getDeviceTypeConfig(ctx, plan)
	if err != nil {
		return err
	}

	networkCardList, err := c.buildNetworkCardList(ctx, plan)
	if err != nil {
		return err
	}
	diskList, err := c.buildDiskList(ctx, plan)
	if err != nil {
		return err
	}
	for _, card := range networkCardList {
		if subnet, ok := subnets[card.SubnetID]; !ok {
			return fmt.Errorf("子网 %s 不属于 %s", card.SubnetID, vpc)
		} else if *deviceTypeConfig.SmartNicExist && subnet.Type != business.SubnetTypeCommonInt {
			return fmt.Errorf("该套餐 %s 为弹性裸金属, 必须使用普通子网", plan.DeviceType.ValueString())
		} else if !*deviceTypeConfig.SmartNicExist && subnet.Type != business.SubnetTypeEbmInt {
			return fmt.Errorf("该套餐 %s 为标准裸金属, 必须使用裸金属子网", plan.DeviceType.ValueString())
		}

	}
	// 弹性裸金属必须有安全组id，标准裸金属一定不能有安全组id
	secGroup, err := c.buildSecGroupList(ctx, plan)
	if err != nil {
		return err
	}
	if *deviceTypeConfig.SmartNicExist && len(secGroup) == 0 {
		return fmt.Errorf("该套餐 %s 为弹性裸金属，必须传递安全组ID", plan.DeviceType.ValueString())
	}
	if !*deviceTypeConfig.SmartNicExist && len(secGroup) != 0 {
		return fmt.Errorf("该套餐 %s 为标准裸金属，不能传递安全组ID", plan.DeviceType.ValueString())
	}
	// 安全组必须存在
	for _, g := range secGroup {
		err = business.NewSecurityGroupService(c.meta).MustExist(ctx, g, plan.RegionID.ValueString())
		if err != nil {
			return err
		}
	}

	// 校验eip
	if plan.PublicIP.ValueString() != "" {
		err = business.NewEipService(c.meta).MustExist(ctx, plan.PublicIP.ValueString(), plan.RegionID.ValueString())
		if err != nil {
			return err
		}
	}

	// 高级版必须关联云硬盘
	if !*deviceTypeConfig.SupportCloud && len(diskList) > 0 {
		return fmt.Errorf("该套餐 %s 不支持关联云硬盘", plan.DeviceType.ValueString())
	}
	if *deviceTypeConfig.CloudBoot && len(diskList) == 0 {
		return fmt.Errorf("该套餐 %s 需要从云硬盘启动，必须关联云硬盘", plan.DeviceType.ValueString())
	}
	var extSys bool
	for _, disk := range diskList {
		if disk.DiskType == business.EbmSystemDiskType {
			extSys = true
			if disk.Size < 100 || disk.Size > 2048 {
				return fmt.Errorf("云盘系统盘容量取值范围：[100, 2048]，单位GB")
			}
		} else if disk.DiskType == business.EbmDataDiskType {
			if disk.Size < 10 || disk.Size > 32768 {
				return fmt.Errorf("云盘数据盘容量取值范围：[10, 32768]，单位GB")
			}
		}
	}
	if !extSys && *deviceTypeConfig.CloudBoot {
		return fmt.Errorf("该套餐 %s 需要从云硬盘启动，必须设置云盘系统盘", plan.DeviceType.ValueString())
	}
	if deviceTypeConfig.SystemVolumeAmount > 0 && plan.SystemVolumeRaidUUID.ValueString() == "" {
		return fmt.Errorf("该套餐 %s 必须传递本地系统盘ID", plan.DeviceType.ValueString())
	}
	if deviceTypeConfig.DataVolumeAmount > 0 && plan.DataVolumeRaidUUID.ValueString() == "" {
		return fmt.Errorf("该套餐 %s 必须传递本地数据盘ID", plan.DeviceType.ValueString())
	}

	// 检查库存
	enough, err := c.checkStock(ctx, plan)
	if err != nil {
		return err
	} else if !enough {
		return fmt.Errorf("该套餐 %s 库存不足", plan.DeviceType.ValueString())
	}

	// 检查镜像
	available, err := c.checkImage(ctx, plan)
	if err != nil {
		return err
	} else if !available {
		return fmt.Errorf("该套餐 %s 不能使用镜像 %s", plan.DeviceType.ValueString(), plan.ImageUUID.ValueString())
	}

	return nil
}

// checkImage 检查镜像是否可用
func (c *ctyunEbm) checkImage(ctx context.Context, plan CtyunEbmConfig) (available bool, err error) {
	deviceType := plan.DeviceType.ValueString()
	imageUUID := plan.ImageUUID.ValueString()
	params := &ctebm.EbmImageListRequest{
		RegionID:   plan.RegionID.ValueString(),
		AzName:     plan.AzName.ValueString(),
		DeviceType: deviceType,
		ImageUUID:  &imageUUID,
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmImageListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	available = len(resp.ReturnObj.Results) > 0
	return
}

// checkStock 获取库存
func (c *ctyunEbm) checkStock(ctx context.Context, plan CtyunEbmConfig) (enough bool, err error) {
	deviceType := plan.DeviceType.ValueString()
	params := &ctebm.EbmDeviceStockListRequest{
		RegionID:   plan.RegionID.ValueString(),
		AzName:     plan.AzName.ValueString(),
		DeviceType: &deviceType,
		Count:      1,
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmDeviceStockListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	enough = *resp.ReturnObj.Results[0].Success
	return
}

// getDeviceTypeConfig 查询套餐详情
func (c *ctyunEbm) getDeviceTypeConfig(ctx context.Context, plan CtyunEbmConfig) (result ctebm.EbmDeviceTypeListReturnObjResultsResponse, err error) {
	deviceType := plan.DeviceType.ValueString()
	params := &ctebm.EbmDeviceTypeListRequest{
		RegionID:   plan.RegionID.ValueString(),
		AzName:     plan.AzName.ValueString(),
		DeviceType: &deviceType,
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmDeviceTypeListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	if len(resp.ReturnObj.Results) == 0 {
		err = fmt.Errorf("未查询到该套餐 %s", deviceType)
		return
	}
	result = *resp.ReturnObj.Results[0]
	return
}

// buildSecGroup 构建安全组列表
func (c *ctyunEbm) buildSecGroupList(ctx context.Context, plan CtyunEbmConfig) (secGroupIDs []string, err error) {
	if plan.SecurityGroupIDs.IsNull() {
		return
	}
	// 处理安全组集合
	diags := plan.SecurityGroupIDs.ElementsAs(ctx, &secGroupIDs, true) // 第二个参数为是否忽略未知值
	if diags.HasError() {
		err = fmt.Errorf("invalid security group ids")
		return
	}
	return
}

// buildDiskList 构建创建物理机时的云硬盘列表结构
func (c *ctyunEbm) buildDiskList(ctx context.Context, plan CtyunEbmConfig) (diskListReq []*ctebm.EbmCreateInstanceV4plusDiskListRequest, err error) {
	if plan.DiskList.IsNull() {
		return
	}
	var diskList []CtyunEbmDiskList
	diags := plan.DiskList.ElementsAs(ctx, &diskList, false)
	if diags.HasError() {
		err = fmt.Errorf("invalid disk list")
		return
	}
	for _, disk := range diskList {
		title := disk.Title.ValueString()
		item := &ctebm.EbmCreateInstanceV4plusDiskListRequest{
			DiskType: disk.DiskType.ValueString(),
			Size:     int32(disk.Size.ValueInt64()),
			RawType:  strings.ToUpper(disk.Type.ValueString()),
		}
		if title != "" {
			item.Title = &title
		}
		diskListReq = append(diskListReq, item)
	}
	return
}

// buildNetworkCardList 构建创建物理机时的网卡结构
func (c *ctyunEbm) buildNetworkCardList(ctx context.Context, plan CtyunEbmConfig) (
	networkCardListReq []*ctebm.EbmCreateInstanceV4plusNetworkCardListRequest,
	err error) {
	var networkCardList []CtyunEbmNetworkCardList
	if plan.NetworkCardList.IsNull() {
		return
	}
	diags := plan.NetworkCardList.ElementsAs(ctx, &networkCardList, false)
	if diags.HasError() {
		err = fmt.Errorf("invalid network card list")
		return
	}
	for _, card := range networkCardList {
		fixedIP := card.FixedIP.ValueString()
		ipv6 := card.FixedIP.ValueString()
		params := &ctebm.EbmCreateInstanceV4plusNetworkCardListRequest{
			Master:   card.Master.ValueBool(),
			SubnetID: card.SubnetID.ValueString(),
		}
		if fixedIP != "" {
			params.FixedIP = &fixedIP
		}
		if ipv6 != "" {
			params.Ipv6 = &ipv6
		}
		networkCardListReq = append(networkCardListReq, params)
	}
	return
}

// handleInstance 操作机器，开机或关机
func (c *ctyunEbm) handleInstance(ctx context.Context, plan CtyunEbmConfig, currentStatus string, targetStatus string) (err error) {
	if currentStatus == targetStatus {
		return
	}
	switch targetStatus {
	case business.EbmStatusStopped:
		return c.stopInstance(ctx, plan)
	case business.EbmStatusRunning:
		return c.startInstance(ctx, plan)
	}
	return errors.New("操作机器状态失败，请检查实例状态")
}

// startInstance 启动物理机
func (c *ctyunEbm) startInstance(ctx context.Context, plan CtyunEbmConfig) (err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmStartInstanceApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmStartInstanceRequest{
		RegionID:     plan.RegionID.ValueString(),
		AzName:       plan.AzName.ValueString(),
		InstanceUUID: plan.InstanceID.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var executeSuccessFlag bool
	var status string
	retryer, _ := business.NewRetryer(time.Second*10, 60)
	retryer.Start(
		func(currentTime int) bool {
			status, err = c.getInstanceStatus(ctx, plan)
			if err != nil {
				return false
			}
			switch status {
			case business.EbmStatusRunning:
				// 执行成功
				executeSuccessFlag = true
				return false
			default:
				return true
			}
		},
	)
	if err != nil {
		return err
	}
	if !executeSuccessFlag {
		return errors.New("执行开启ebm动作时，ebm状态异常：status")
	}
	return
}

// stopInstance 关闭物理机
func (c *ctyunEbm) stopInstance(ctx context.Context, plan CtyunEbmConfig) (err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmStopInstanceApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmStopInstanceRequest{
		RegionID:     plan.RegionID.ValueString(),
		AzName:       plan.AzName.ValueString(),
		InstanceUUID: plan.InstanceID.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var executeSuccessFlag bool
	var status string
	retryer, _ := business.NewRetryer(time.Second*10, 60)
	retryer.Start(
		func(currentTime int) bool {
			status, err = c.getInstanceStatus(ctx, plan)
			if err != nil {
				return false
			}
			switch status {
			case business.EbmStatusStopped:
				// 执行成功
				executeSuccessFlag = true
				return false
			default: // 其他状态持续轮询
				return true
			}
		})
	if err != nil {
		return err
	}
	if !executeSuccessFlag {
		return errors.New("执行关闭ebm动作时，ebm状态异常，当前状态：" + status)
	}
	return
}

// getAndMerge 查询ebm
func (c *ctyunEbm) getAndMerge(ctx context.Context, cfg *CtyunEbmConfig) (err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmDescribeInstanceV4plusApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmDescribeInstanceV4plusRequest{
		RegionID:     cfg.RegionID.ValueString(),
		InstanceUUID: cfg.InstanceID.ValueString(),
		AzName:       cfg.AzName.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	instance := resp.ReturnObj
	cfg.InstanceID = utils.SecStringValue(instance.InstanceUUID)
	cfg.RegionID = utils.SecStringValue(instance.RegionID)
	cfg.AzName = utils.SecStringValue(instance.AzName)
	cfg.DeviceType = utils.SecStringValue(instance.DeviceType)
	cfg.InstanceName = utils.SecStringValue(instance.DisplayName)
	cfg.Hostname = utils.SecStringValue(instance.InstanceName)
	cfg.ImageUUID = utils.SecStringValue(instance.ImageID)
	cfg.SystemVolumeRaidUUID = utils.SecStringValue(instance.SystemVolumeRaidID)
	cfg.DataVolumeRaidUUID = utils.SecStringValue(instance.DataVolumeRaidID)
	cfg.VpcID = utils.SecStringValue(instance.VpcID)
	cfg.Status = utils.SecLowerStringValue(instance.EbmState)

	cfg.PublicIP = utils.SecStringValue(instance.PublicIP)
	cardList := []CtyunEbmNetworkCardList{}
	diskList := []CtyunEbmDiskList{}
	cardObj := utils.StructToTFObjectTypes(CtyunEbmNetworkCardList{})
	diskObj := utils.StructToTFObjectTypes(CtyunEbmDiskList{})
	for _, card := range instance.Interfaces {
		master := utils.SecBoolValue(card.Master)
		if master.ValueBool() && len(card.SecurityGroups) > 0 {
			var secGroups []string
			for _, g := range card.SecurityGroups {
				secGroups = append(secGroups, utils.SecString(g.SecurityGroupID))
			}
			cfg.SecurityGroupIDs, _ = types.SetValueFrom(ctx, types.StringType, secGroups)
		}
		item := CtyunEbmNetworkCardList{
			FixedIP:     utils.SecStringValue(card.Ipv4),
			Master:      master,
			Ipv6:        utils.SecStringValue(card.Ipv6),
			SubnetID:    utils.SecStringValue(card.SubnetUUID),
			PortID:      utils.SecStringValue(card.PortUUID),
			InterfaceID: utils.SecStringValue(card.InterfaceUUID),
		}
		cardList = append(cardList, item)
	}
	cfg.NetworkCardList, _ = types.ListValueFrom(ctx, cardObj, cardList)

	for _, diskId := range instance.AttachedVolumes {
		diskInfo, err := business.NewEbsService(c.meta).GetEbsInfo(ctx, *diskId, cfg.RegionID.ValueString())
		if err != nil {
			return err
		}
		item := CtyunEbmDiskList{
			Type:  types.StringValue(strings.ToLower(diskInfo.DiskType)),
			Title: types.StringValue(diskInfo.DiskName),
			Size:  types.Int64Value(diskInfo.DiskSize),
		}
		if diskInfo.IsSystemVolume {
			item.DiskType = types.StringValue(business.EbmSystemDiskType)
		} else {
			item.DiskType = types.StringValue(business.EbmDataDiskType)
		}
		diskList = append(diskList, item)
	}
	cfg.DiskList, _ = types.ListValueFrom(ctx, diskObj, diskList)
	cfg.ID = cfg.InstanceID

	return nil
}

// getMasterOrderIdIfOrderInProgress 获取masterOrderId
func (c *ctyunEbm) getMasterOrderIdIfOrderInProgress(err ctyunsdk.CtyunRequestError) (string, error) {
	resp := struct {
		MasterOrderId string `json:"masterOrderID"`
		MasterOrderNo string `json:"masterOrderNO"`
	}{}
	if err.CtyunResponse() == nil {
		return "", err
	}
	_, err = err.CtyunResponse().ParseByStandardModel(&resp)
	if err != nil {
		return "", err
	}
	return resp.MasterOrderId, err
}

// acquireIdIfOrderNotFinished 重新获取id，如果前订单状态有问题需要重新轮询
// 返回值：数据是否有效
func (c *ctyunEbm) acquireAndSetIdIfOrderNotFinished(ctx context.Context, state *CtyunEbmConfig, response *resource.ReadResponse) bool {
	id := state.InstanceID.ValueString()
	masterOrderId := state.MasterOrderID.ValueString()
	if id != "" {
		// 数据是完整的，无需处理
		return true
	}
	if state.MasterOrderID.ValueString() == "" {
		// 没有受理的订购单id，数据是不可恢复的，直接把当前状态移除并且返回
		response.State.RemoveResource(ctx)
		return false
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	resp, err := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
	if err != nil || len(resp.Uuid) == 0 {
		// 报错了，或者受理没有返回数据的情况，那么意思是这个单子并没有开通出来，此时数据无法恢复
		response.State.RemoveResource(ctx)
		return false
	}

	// 成功把id恢复出来
	state.InstanceID = types.StringValue(resp.Uuid[0])
	response.State.Set(ctx, state)
	return true
}

// updateInstanceName 更新实例名称
func (c *ctyunEbm) updateInstanceName(ctx context.Context, state CtyunEbmConfig, plan CtyunEbmConfig) (err error) {
	// 判断名字是否相同
	if plan.InstanceName.Equal(state.InstanceName) {
		return
	}

	name := plan.InstanceName.ValueString()
	resp, err := c.meta.Apis.CtEbmApis.EbmUpdateInstanceApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmUpdateInstanceRequest{
		RegionID:     state.RegionID.ValueString(),
		AzName:       state.AzName.ValueString(),
		DisplayName:  &name,
		InstanceUUID: state.InstanceID.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
	}
	return
}

// updatePasswordOrHostname 修改密码或主机名
func (c *ctyunEbm) updatePasswordOrHostname(ctx context.Context, state CtyunEbmConfig, plan CtyunEbmConfig) (err error) {
	if state.Password.Equal(plan.Password) && state.Hostname.Equal(plan.Hostname) {
		return
	}
	// 修改前需要检查机器状态是否是关机
	err = c.checkBeforeUpdatePasswordOrHostname(ctx, state)
	if err != nil {
		return
	}
	// 修改密码
	err = c.updatePassword(ctx, state, plan)
	if err != nil {
		return
	}
	// 修改主机名
	err = c.updateHostname(ctx, state, plan)
	if err != nil {
		return
	}

	return
}

// updatePassword 修改密码
func (c *ctyunEbm) updatePassword(ctx context.Context, state CtyunEbmConfig, plan CtyunEbmConfig) (err error) {
	if state.Password.Equal(plan.Password) {
		return
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmResetPasswordApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmResetPasswordRequest{
		RegionID:     state.RegionID.ValueString(),
		AzName:       state.AzName.ValueString(),
		InstanceUUID: state.InstanceID.ValueString(),
		NewPassword:  plan.Password.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
	}
	// 通过状态检查是否修改完成
	err = c.checkAfterUpdatePasswordOrHostname(ctx, state)
	if err != nil {
		return
	}
	// 关机
	err = c.stopInstance(ctx, state)
	if err != nil {
		return
	}
	return
}

// updatePassword 修改主机名
func (c *ctyunEbm) updateHostname(ctx context.Context, state CtyunEbmConfig, plan CtyunEbmConfig) (err error) {
	if state.Hostname.Equal(plan.Hostname) {
		return
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmResetHostnameApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmResetHostnameRequest{
		RegionID:     state.RegionID.ValueString(),
		AzName:       state.AzName.ValueString(),
		InstanceUUID: state.InstanceID.ValueString(),
		Hostname:     plan.Hostname.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
	}
	// 通过状态检查是否修改完成
	err = c.checkAfterUpdatePasswordOrHostname(ctx, state)
	if err != nil {
		return
	}
	// 关机
	err = c.stopInstance(ctx, state)
	if err != nil {
		return
	}
	return
}

// checkBeforeUpdatePasswordOrHostname 修改密码或主机名前对机器状态做检查
func (c *ctyunEbm) checkBeforeUpdatePasswordOrHostname(ctx context.Context, state CtyunEbmConfig) error {
	var executeSuccessFlag bool
	var status string
	var err error
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			status, err = c.getInstanceStatus(ctx, state)
			if err != nil {
				return false
			}
			switch status {
			case business.EbmStatusStopping, business.EbmStatusResettingPassword, business.EbmStatusResettingHostname:
				return true
			case business.EbmStatusStopped:
				executeSuccessFlag = true
				return false
			default:
				return false
			}
		})
	if err != nil {
		return err
	}
	if !executeSuccessFlag {
		return errors.New("修改物理机密码或hostname失败，请确认物理机状态，修改密码或hostname必须先关机，当前状态：" + status)
	}
	return nil
}

// checkAfterUpdatePasswordOrHostname 修改后检查机器状态
func (c *ctyunEbm) checkAfterUpdatePasswordOrHostname(ctx context.Context, state CtyunEbmConfig) error {
	var executeSuccessFlag bool
	var status string
	var err error
	var cnt int
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			status, err = c.getInstanceStatus(ctx, state)
			if err != nil {
				return false
			}
			switch status {
			case business.EbmStatusResettingPassword, business.EbmStatusResettingHostname:
				return true
			case business.EbmStatusRunning:
				executeSuccessFlag = true
				return false
			case business.EbmStatusStopped:
				cnt++
				if cnt > 3 {
					return false
				}
				return true
			default:
				return false
			}
		})
	if err != nil {
		return err
	}
	if !executeSuccessFlag {
		return errors.New("修改物理机密码或hostname后，检查失败，请确认物理机状态：" + status)
	}
	return nil
}

// getInstanceStatus 查询物理机状态
func (c *ctyunEbm) getInstanceStatus(ctx context.Context, state CtyunEbmConfig) (status string, err error) {
	return business.NewEbmService(c.meta).GetEbmStatus(
		ctx,
		state.InstanceID.ValueString(),
		state.RegionID.ValueString(),
		state.AzName.ValueString(),
	)
}

// delete 删除物理机
func (c *ctyunEbm) delete(ctx context.Context, state CtyunEbmConfig) (err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmDeleteInstanceApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmDeleteInstanceRequest{
		RegionID:     state.RegionID.ValueString(),
		AzName:       state.AzName.ValueString(),
		InstanceUUID: state.InstanceID.ValueString(),
		ClientToken:  uuid.NewString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	err = helper.RefundLoop(ctx, c.meta.Credential, *resp.ReturnObj.MasterOrderID)
	if err != nil {
		return
	}
	return
}

// checkAfterDelete 删除后检查
func (c *ctyunEbm) checkAfterDelete(ctx context.Context, state CtyunEbmConfig) (err error) {
	portID := c.getMasterPortID(ctx, state)
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	var exist bool
	retryer.Start(
		func(currentTime int) bool {
			exist, err = business.NewPortService(c.meta).Exist(ctx, portID, state.RegionID.ValueString())
			if err != nil {
				return false
			}
			if !exist {
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if exist {
		err = fmt.Errorf("裸金属 %s 的主网卡 %s 残留", state.InstanceID.ValueString(), portID)
	}
	return
}

// getMasterPortID 获取主网卡id
func (c *ctyunEbm) getMasterPortID(ctx context.Context, state CtyunEbmConfig) (portID string) {
	if state.NetworkCardList.IsNull() {
		return
	}
	var networkCardList []CtyunEbmNetworkCardList
	diags := state.NetworkCardList.ElementsAs(ctx, &networkCardList, false)
	if diags.HasError() { // 无需处理
		return
	}
	for _, card := range networkCardList {
		if card.Master.ValueBool() {
			portID = card.PortID.ValueString()
		}
	}
	return
}
