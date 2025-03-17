package ebm

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-ctyun/internal/core/ctebm"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
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
	ID                   types.String        `tfsdk:"id"`
	RegionID             types.String        `tfsdk:"region_id"`
	AzName               types.String        `tfsdk:"az_name"`
	DeviceType           types.String        `tfsdk:"device_type"`
	InstanceName         types.String        `tfsdk:"instance_name"`
	Hostname             types.String        `tfsdk:"hostname"`
	ImageUUID            types.String        `tfsdk:"image_uuid"`
	Password             types.String        `tfsdk:"password"`
	ProjectID            types.String        `tfsdk:"project_id"`
	SystemVolumeRaidUUID types.String        `tfsdk:"system_volume_raid_uuid"`
	DataVolumeRaidUUID   types.String        `tfsdk:"data_volume_raid_uuid"`
	VpcID                types.String        `tfsdk:"vpc_id"`
	ExtIP                types.String        `tfsdk:"ext_ip"`
	IpType               types.String        `tfsdk:"ip_type"`
	BandWidth            types.Int64         `tfsdk:"band_width"`
	PublicIP             types.String        `tfsdk:"public_ip"`
	SecurityGroupID      types.String        `tfsdk:"security_group_id"`
	DiskList             basetypes.ListValue `tfsdk:"disk_list"`
	NetworkCardList      basetypes.ListValue `tfsdk:"network_card_list"`
	UserData             types.String        `tfsdk:"user_data"`
	KeyName              types.String        `tfsdk:"key_name"`
	PayVoucherPrice      types.Float64       `tfsdk:"pay_voucher_price"`
	AutoRenewStatus      types.Int64         `tfsdk:"auto_renew_status"`
	InstanceChargeType   types.String        `tfsdk:"instance_charge_type"`
	CycleCount           types.Int64         `tfsdk:"cycle_count"`
	CycleType            types.String        `tfsdk:"cycle_type"`
	MasterOrderID        types.String        `tfsdk:"master_order_id"`
	Status               types.String        `tfsdk:"status"`
}

type CtyunEbmDiskList struct {
	DiskType types.String `tfsdk:"disk_type"`
	Title    types.String `tfsdk:"title"`
	Type     types.String `tfsdk:"type"`
	Size     types.Int64  `tfsdk:"size"`
}

type CtyunEbmNetworkCardList struct {
	Title    types.String `tfsdk:"title"`
	FixedIP  types.String `tfsdk:"fixed_ip"`
	Master   types.Bool   `tfsdk:"master"`
	Ipv6     types.String `tfsdk:"ipv6"`
	SubnetID types.String `tfsdk:"subnet_id"`
}

func (c *ctyunEbm) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027724**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "id",
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订购的受理单id",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "区域ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称",
				Default:     defaults.AcquireFromGlobalString(common.ExtraAzName, true),
			},
			"device_type": schema.StringAttribute{
				Required:    true,
				Description: "物理机套餐类型",
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
			},
			"password": schema.StringAttribute{
				Sensitive:   true,
				Optional:    true,
				Computed:    true,
				Description: "密码(必须包含大小写字母和（一个数字或者特殊字符）长度8到30位)，未传入有效的keyName时必须传入password",
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
			},
			"data_volume_raid_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "本地数据盘raid类型，如果有本地盘则必填",
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "主网卡网络ID",
			},
			"ext_ip": schema.StringAttribute{
				Required:    true,
				Description: "是否使用弹性公网IP，取值范围:[1=自动分配,0=不使用,2=使用已有]",
			},
			"ip_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "弹性IP版本，取值范围:[ipv4=v4地址,ipv6=v6地址]，默认值:ipv4",
				Default:     stringdefault.StaticString("ipv4"),
			},
			"band_width": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "带宽，取值范围:[1~2000]，默认值:100",
				Default:     int64default.StaticInt64(0),
			},
			"public_ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "弹性公网IP的id",
				Default:     stringdefault.StaticString(""),
			},
			"security_group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "安全组ID，套餐smartNicExist为true可支持安全组。创建弹性裸金属必须传入安全组ID，标准裸金属不支持传入安全组ID",
			},
			"disk_list": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "云盘信息列表，套餐中supportCloud为true表示支持云盘",
				Default: listdefault.StaticValue(types.ListValueMust(
					types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"disk_type": types.StringType,
							"size":      types.Int64Type,
							"title":     types.StringType,
							"type":      types.StringType,
						},
					},
					[]attr.Value{},
				)),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"disk_type": schema.StringAttribute{
							Required:    true,
							Description: "磁盘类型system或data，套餐中cloudBoot为true表示支持云盘系统盘",
						},
						"title": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "磁盘名称，长度2~64,不支持中文",
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "磁盘分类，取值范围:[SAS=SAS盘,SATA=SATA盘,SSD-genric=SSD-genric盘,SSD=SSD盘]",
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
						"title": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "网卡名称",
							Default:     stringdefault.StaticString(""),
						},
						"fixed_ip": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "内网IPv4地址",
						},
						"master": schema.BoolAttribute{
							Required:    true,
							Description: "是否主节点(True代表主节点)",
						},
						"ipv6": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "内网IPv6地址",
							Default:     stringdefault.StaticString(""),
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
			},
			"key_name": schema.StringAttribute{
				Optional:    true,
				Description: "密钥对名词",
			},
			"pay_voucher_price": schema.Float64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "代金券，满足以下规则：两位小数，不足两位自动补0，超过两位小数无效；不可为负数；字段为0时表示不使用代金券",
				Default:     float64default.StaticFloat64(0.00),
			},
			"auto_renew_status": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "是否自动续订，默认非自动续订。取值范围：<br/>0（不续费），<br/>1（自动续费），<br/>注：按月购买，自动续订周期为1个月；按年购买，自动续订周期为1年",
				Default:     int64default.StaticInt64(0),
			},
			"instance_charge_type": schema.StringAttribute{
				Required:    true,
				Description: "实例计费类型 <br/>*ORDER_ON_CYCLE：包年包月<br/>*ORDER_ON_DEMAND：按量付费",
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "订购时长，该参数需要与cycleType一同使用<br/>注：最长订购周期为60个月（5年）；cycleType与cycleCount一起填写；按量付费（即instanceChargeType为ORDER_ON_DEMAND）时，无需填写该参数（填写无效）",
			},
			"cycle_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "订购周期类型，取值范围:[MONTH=按月,YEAR=按年]<br/>注：cycleType与cycleCount一起填写；按量付费（即instanceChargeType为ORDER_ON_DEMAND）时，无需填写该参数（填写无效）",
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
	var plan CtyunEbmConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建
	returnObj, err := c.createInstance(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 先保存订单号
	masterOrderId := *returnObj.MasterOrderID
	plan.MasterOrderID = types.StringValue(masterOrderId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 根据订单号轮询查资源的uuid
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	loop, err := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId, 600)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	plan.ID = types.StringValue(loop.Uuid[0])
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建机器后状态默认为启动状态，可根据用户要求的状态，去执行对应的操作，比如关机
	err = c.handleInstance(ctx, plan, business.EbmStatusRunning, plan.Status.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	// 反查信息
	err = c.getAndMergeEbm(ctx, &plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunEbm) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunEbmConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	if !c.acquireAndSetIdIfOrderNotFinished(ctx, &state, response) {
		return
	}
	err := c.getAndMergeEbm(ctx, &state)
	if errors.Is(err, common.InvalidReturnObjError) {
		response.State.RemoveResource(ctx)
		return
	} else if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEbm) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
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

	// 处理开关机
	err := c.handleInstance(ctx, state, state.Status.ValueString(), plan.Status.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	// 修改基础信息
	err = c.updateInstanceInfo(ctx, state, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	// 修改密码
	err = c.updatePassword(ctx, state, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	state.Password = plan.Password

	err = c.getAndMergeEbm(ctx, &state)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEbm) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunEbmConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	resp, err := c.meta.Apis.CtEbmApis.EbmDeleteInstanceApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmDeleteInstanceRequest{
		RegionID:     state.RegionID.ValueString(),
		AzName:       state.AzName.ValueString(),
		InstanceUUID: state.ID.ValueString(),
		ClientToken:  uuid.NewString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf(*resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	err = helper.RefundLoop(ctx, c.meta.Credential, *resp.ReturnObj.MasterOrderID)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
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

// 导入命令：terraform import [配置标识].[导入配置名称] [uuid]
func (c *ctyunEbm) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunEbmConfig
	var id string
	err := terraform_extend.Split(request.ID, &id)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	regionId := c.meta.GetExtraIfEmpty(cfg.RegionID.ValueString(), common.ExtraRegionId)
	cfg.RegionID = types.StringValue(regionId)
	azName := c.meta.GetExtraIfEmpty(cfg.AzName.ValueString(), common.ExtraAzName)
	cfg.AzName = types.StringValue(azName)

	cfg.ID = types.StringValue(id)
	err = c.getAndMergeEbm(ctx, &cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
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
	securityGroupID := plan.SecurityGroupID.ValueString()
	userData := plan.UserData.ValueString()
	keyName := plan.KeyName.ValueString()
	instanceChargeType := plan.InstanceChargeType.ValueString()
	cycleType := plan.CycleType.ValueString()
	bandwidth := plan.BandWidth.ValueInt64()
	diskList := c.buildDiskList(ctx, plan)
	networkCardList := c.buildNetworkCardList(ctx, plan)
	// 需要校验的很多，比如弹性裸金属需要安全组id、自选ip时需要传ip，自动分配需要传带宽，密码和密钥对必须有一个，raidid是否传递等等
	params := &ctebm.EbmCreateInstanceV4plusRequest{
		RegionID:             regionID,
		AzName:               azName,
		DeviceType:           plan.DeviceType.ValueString(),
		InstanceName:         plan.InstanceName.ValueString(),
		Hostname:             plan.Hostname.ValueString(),
		ImageUUID:            plan.ImageUUID.ValueString(),
		Password:             &password,
		SystemVolumeRaidUUID: &systemVolumeRaidUUID,
		DataVolumeRaidUUID:   &dataVolumeRaidUUID,
		VpcID:                plan.VpcID.ValueString(),
		ExtIP:                plan.ExtIP.ValueString(),
		PayVoucherPrice:      float32(plan.PayVoucherPrice.ValueFloat64()),
		ProjectID:            &projectID,
		IpType:               &ipType,
		DiskList:             diskList,
		NetworkCardList:      networkCardList,
		UserData:             &userData,
		KeyName:              &keyName,
		AutoRenewStatus:      int32(plan.AutoRenewStatus.ValueInt64()),
		InstanceChargeType:   &instanceChargeType,
		CycleCount:           int32(plan.CycleCount.ValueInt64()),
		CycleType:            &cycleType,
		ClientToken:          uuid.NewString(),
		OrderCount:           1,
	}

	if bandwidth > 0 {
		params.BandWidth = int32(bandwidth)
	}
	if securityGroupID != "" {
		params.SecurityGroupID = &securityGroupID
	}
	if publicIP != "" {
		params.PublicIP = &publicIP
	}

	resp, err := c.meta.Apis.CtEbmApis.EbmCreateInstanceV4plusApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf(*resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	returnObj = *resp.ReturnObj
	return
}

// buildDiskList 构建创建物理机时的云硬盘列表结构
func (c *ctyunEbm) buildDiskList(ctx context.Context, plan CtyunEbmConfig) (diskListReq []*ctebm.EbmCreateInstanceV4plusDiskListRequest) {
	if plan.DiskList.IsNull() {
		return
	}
	var diskList []CtyunEbmDiskList
	diags := plan.DiskList.ElementsAs(ctx, &diskList, false)
	if diags.HasError() {
		return
	}
	for _, disk := range diskList {
		title := disk.Title.ValueString()
		diskListReq = append(diskListReq, &ctebm.EbmCreateInstanceV4plusDiskListRequest{
			DiskType: disk.DiskType.ValueString(),
			Size:     int32(disk.Size.ValueInt64()),
			Title:    &title,
			RawType:  disk.DiskType.ValueString(),
		})
	}
	return
}

// buildNetworkCardList 构建创建物理机时的网卡结构
func (c *ctyunEbm) buildNetworkCardList(ctx context.Context, plan CtyunEbmConfig) (networkCardListReq []*ctebm.EbmCreateInstanceV4plusNetworkCardListRequest) {
	var networkCardList []CtyunEbmNetworkCardList
	if plan.NetworkCardList.IsNull() {
		return
	}
	diags := plan.NetworkCardList.ElementsAs(ctx, &networkCardList, false)
	if diags.HasError() {
		return
	}
	for _, card := range networkCardList {
		title := card.Title.ValueString()
		fixedIP := card.FixedIP.ValueString()
		ipv6 := card.FixedIP.ValueString()
		params := &ctebm.EbmCreateInstanceV4plusNetworkCardListRequest{
			Master:   card.Master.ValueBool(),
			SubnetID: card.SubnetID.ValueString(),
		}
		if title != "" {
			params.Title = &title
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
		InstanceUUID: plan.ID.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf(*resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*5, 20)
	retryer.Start(
		func(currentTime int) bool {
			status, err := c.getInstanceStatus(ctx, plan)
			if err != nil {
				return false
			}
			switch status {
			case business.EbmStatusStarting:
				// 执行中
				return true
			case business.EbmStatusRunning:
				// 执行成功
				executeSuccessFlag = true
				return false
			default:
				// 默认为执行失败
				return false
			}
		},
	)
	if !executeSuccessFlag {
		return errors.New("执行开启ebm动作时，ebm状态异常")
	}
	return
}

// stopInstance 关闭物理机
func (c *ctyunEbm) stopInstance(ctx context.Context, plan CtyunEbmConfig) (err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmStopInstanceApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmStopInstanceRequest{
		RegionID:     plan.RegionID.ValueString(),
		AzName:       plan.AzName.ValueString(),
		InstanceUUID: plan.ID.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf(*resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*5, 20)
	retryer.Start(
		func(currentTime int) bool {
			status, err := c.getInstanceStatus(ctx, plan)
			if err != nil {
				return false
			}
			switch status {
			case business.EbmStatusStopping:
				// 执行中
				return true
			case business.EbmStatusStopped:
				// 执行成功
				executeSuccessFlag = true
				return false
			default:
				// 默认为执行失败
				return false
			}
		})
	if !executeSuccessFlag {
		return errors.New("执行关闭ebm动作时，ebm状态异常")
	}
	return
}

// getAndMergeEbm 查询ebm
func (c *ctyunEbm) getAndMergeEbm(ctx context.Context, cfg *CtyunEbmConfig) (err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmDescribeInstanceV4plusApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmDescribeInstanceV4plusRequest{
		RegionID:     cfg.RegionID.ValueString(),
		InstanceUUID: cfg.ID.ValueString(),
		AzName:       cfg.AzName.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf(*resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	instance := resp.ReturnObj
	cfg.ID = types.StringValue(*instance.InstanceUUID)
	cfg.RegionID = types.StringValue(*instance.RegionID)
	cfg.AzName = types.StringValue(*instance.AzName)
	cfg.DeviceType = types.StringValue(*instance.DeviceType)
	cfg.InstanceName = types.StringValue(*instance.DisplayName)
	cfg.Hostname = types.StringValue(*instance.InstanceName)
	cfg.ImageUUID = types.StringValue(*instance.ImageID)
	cfg.SystemVolumeRaidUUID = types.StringValue(*instance.SystemVolumeRaidID)
	cfg.DataVolumeRaidUUID = types.StringValue(*instance.DataVolumeRaidID)
	cfg.VpcID = types.StringValue(*instance.VpcID)
	cfg.Status = types.StringValue(*instance.EbmState)
	//cfg.PublicIP = types.StringValue(*instance.PublicIP)
	if cfg.InstanceChargeType.ValueString() == "ORDER_ON_DEMAND" {
		cfg.CycleType = types.StringValue("")
		cfg.CycleCount = types.Int64Value(0)
	}
	networkCard := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"title":     types.StringType,
			"fixed_ip":  types.StringType,
			"master":    types.BoolType,
			"ipv6":      types.StringType,
			"subnet_id": types.StringType,
		},
	}
	cards := []attr.Value{}
	for _, card := range instance.Interfaces {
		c, _ := types.ObjectValue(
			networkCard.AttrTypes,
			map[string]attr.Value{
				"title":     types.StringValue(""),
				"fixed_ip":  types.StringValue(*card.Ipv4),
				"ipv6":      types.StringValue(*card.Ipv6),
				"master":    types.BoolValue(*card.Master),
				"subnet_id": types.StringValue(*card.SubnetUUID),
			},
		)
		cards = append(cards, c)
	}
	cfg.NetworkCardList, _ = types.ListValue(
		networkCard,
		cards,
	)

	// 密码、企业项目、ExtIP、IpType、
	//bindwith、SecurityGroupID、disklist、NetworkCardList、UserData、KeyName, PayVoucherPrice
	// AutoRenewStatus InstanceChargeType\CycleCount\CycleType\OrderCount没返回
	// 需要填充

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
	id := state.ID.ValueString()
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
	state.ID = types.StringValue(resp.Uuid[0])
	response.State.Set(ctx, state)
	return true
}

// updateInstanceInfo 更新主机的部分信息
func (c *ctyunEbm) updateInstanceInfo(ctx context.Context, state CtyunEbmConfig, plan CtyunEbmConfig) (err error) {
	// 判断名字是否相同
	if plan.InstanceName.Equal(state.InstanceName) {
		return
	}

	name := plan.InstanceName.ValueString()
	resp, err := c.meta.Apis.CtEbmApis.EbmUpdateInstanceApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmUpdateInstanceRequest{
		RegionID:     state.RegionID.ValueString(),
		AzName:       state.AzName.ValueString(),
		DisplayName:  &name,
		InstanceUUID: state.ID.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf(*resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// updatePassword 修改密码
func (c *ctyunEbm) updatePassword(ctx context.Context, state CtyunEbmConfig, plan CtyunEbmConfig) (err error) {
	if state.Password.Equal(plan.Password) {
		return
	}
	// 修改前需要检查机器状态
	err = c.checkBeforeUpdatePassword(ctx, state)
	if err != nil {
		return
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmResetPasswordApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmResetPasswordRequest{
		RegionID:     state.RegionID.ValueString(),
		AzName:       state.AzName.ValueString(),
		InstanceUUID: state.ID.ValueString(),
		NewPassword:  plan.Password.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf(*resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	err = c.checkAfterUpdatePassword(ctx, state)
	if err != nil {
		return
	}
	// 因为改完密码会默认开机，而改密码前一定是关机状态，这里要手动恢复到关机态
	err = c.stopInstance(ctx, state)
	if err != nil {
		return
	}
	return
}

// checkBeforeUpdatePassword 修改密码前对机器状态做检查
func (c *ctyunEbm) checkBeforeUpdatePassword(ctx context.Context, state CtyunEbmConfig) error {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 120) // 20min
	retryer.Start(
		func(currentTime int) bool {
			status, err := c.getInstanceStatus(ctx, state)
			if err != nil {
				return false
			}
			switch status {
			case business.EbmStatusStopping, business.EbmStatusResettingPassword:
				return true
			case business.EbmStatusStopped:
				executeSuccessFlag = true
				return false
			default:
				return false
			}
		})
	if !executeSuccessFlag {
		return errors.New("修改物理机密码前置检查失败，请确认物理机状态")
	}
	return nil
}

// checkAfterUpdatePassword 修改密码后检查是否需要恢复机器状态
func (c *ctyunEbm) checkAfterUpdatePassword(ctx context.Context, state CtyunEbmConfig) error {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 120)
	retryer.Start(
		func(currentTime int) bool {
			status, err := c.getInstanceStatus(ctx, state)
			if err != nil {
				return false
			}
			switch status {
			case business.EbmStatusResettingPassword:
				return true
			case business.EbmStatusRunning:
				executeSuccessFlag = true
				return false
			default:
				return false
			}
		})
	if !executeSuccessFlag {
		return errors.New("修改物理机密码后置检查失败，请确认物理机状态")
	}
	return nil
}

// getInstanceStatus 查询物理机状态
func (c *ctyunEbm) getInstanceStatus(ctx context.Context, state CtyunEbmConfig) (status string, err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmDescribeInstanceV4plusApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmDescribeInstanceV4plusRequest{
		RegionID:     state.RegionID.ValueString(),
		InstanceUUID: state.ID.ValueString(),
		AzName:       state.AzName.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf(*resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return *resp.ReturnObj.EbmState, err
}
