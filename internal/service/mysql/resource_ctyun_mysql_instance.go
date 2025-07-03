package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
	"time"
)

var (
	_ resource.Resource                = &CtyunMysqlInstance{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlInstance{}
	_ resource.ResourceWithImportState = &CtyunMysqlInstance{}
)

type CtyunMysqlInstance struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMysqlInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunMysqlInstance) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func NewCtyunMysqlInstance() resource.Resource {
	return &CtyunMysqlInstance{}
}

func (c *CtyunMysqlInstance) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_instance"
}

func (c *CtyunMysqlInstance) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month：按月，on_demand：按需。当此值为month时，cycle_count为必填",
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypeOnDemand, business.OrderCycleTypeMonth),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int32Attribute{
				Optional:    true,
				Description: "订购时长，该参数当且仅当在cycle_type为month时填写，支持传递1-36",
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					int32validator.Between(1, 36),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"auto_renew": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否自动续订，默认非自动续订，当cycle_type不等于on_demand时才可填写，当cycle_count<12，到期自动续订1个月，当cycle_count>=12，到期自动续订12个月",
				Default:     booldefault.StaticBool(false),
				Validators: []validator.Bool{
					validator2.ConflictsWithEqualBool(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
				},
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
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
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "虚拟私有云Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"host_type": schema.StringAttribute{ //host_type
				Required:    true,
				Description: "主机类型host_type: S6 or S7等。可根据data.ctyun_mysql_specs获取",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"security_group_id": schema.StringAttribute{
				Required:    true,
				Description: "安全组Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "实例名称（长度在 4 到 64个字符，必须以字母开头，不区分大小写，可以包含字母、数字、中划线或下划线，不能包含其他特殊字符）",
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 64),
				},
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "实例密码为8-26位，需为字母、数字和特殊字符~!@#%^*_-+:,.?/{[]}的组合，区分大小写。",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 26),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"prod_id": schema.StringAttribute{
				Required:    true,
				Description: "产品id。在扩容过程中，不支持规格和实例扩容同时进行，ProdID和prod_performance_spec不能同时与原配置不一致。prod_id取值范围：Single57（单实例5.7版本）, Single80（单实例8.0版本）, MasterSlave57（一主一备5.7版本）, MasterSlave80（一主一备8.0版本）, Master2Slave57（一主两备5.7版本）, Master2Slave80（一主两备8.0版本）",
				Validators: []validator.String{
					stringvalidator.OneOf(business.MysqlProdIds...),
				},
			},
			"instance_series": schema.StringAttribute{
				Required:    true,
				Description: "实例规格，取值范围：S(通用型)，C(计算增强型)，M(内存增强型)",
				Validators: []validator.String{
					stringvalidator.OneOf(business.MysqlInstanceSeries...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_type": schema.StringAttribute{
				Required:    true,
				Description: "存储类型: SSD=超高IO、SATA=普通IO、SAS=高IO、SSD-genric=通用型SSD、FAST-SSD=极速型SSD",
				Validators: []validator.String{
					stringvalidator.OneOf(business.StorageType...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_space": schema.Int32Attribute{
				Required:    true,
				Description: "存储空间(单位:G，范围100,32768)",
				Validators: []validator.Int32{
					int32validator.Between(100, 32768),
				},
			},
			"backup_storage_space": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "备份节点存储空间(单位:G，范围100,32768)，不支持主备磁盘空间同时升配。若storage_space和backup_storage_space都不为空，优先升配备份节点存储空间",
				Validators: []validator.Int32{
					int32validator.Between(100, 32768),
				},
			},
			"prod_performance_spec": schema.StringAttribute{
				Required:    true,
				Description: "规格(例: 4C8G),可根据data.ctyun_mysql_specs获取。不支持规格和实例扩容同时进行：ProdID和prod_performance_spec不能同时与原配置不一致",
			},
			"availability_zone_info": schema.ListNestedAttribute{
				Required:    true,
				Description: "可用区信息",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"availability_zone_name": schema.StringAttribute{
							Required:    true,
							Description: "资源池可用区名称",
						},
						"availability_zone_count": schema.Int32Attribute{
							Required:    true,
							Description: "资源池可用区总数",
						},
						"node_type": schema.StringAttribute{
							Required:    true,
							Description: "表示分布AZ的节点类型，master/slave",
						},
					},
				},
			},
			"cpu_type": schema.StringAttribute{
				Required:    true,
				Description: "cpu类型：KunPeng(鲲鹏)，Hygon(海光)，Intel(intel)，AMD(amd),Phytium(飞腾)，Loongson(龙芯)",
				Validators: []validator.String{
					stringvalidator.OneOf(business.MysqlCpuType...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"os_type": schema.StringAttribute{
				Required:    true,
				Description: "系统类型：nil(裸机)，windows，centos，ubuntu，android，redhat，kylin，uos，suse，asianux，open_euler，ctyunos，euler",
				Validators: []validator.String{
					stringvalidator.OneOf(business.MysqlOSType...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订单id",
			},
			"inst_id": schema.StringAttribute{
				Computed:    true,
				Description: "实例Id",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
			},
			"prod_running_status": schema.Int32Attribute{
				Computed:    true,
				Description: "0.正常 1.重启中 2.备份中 3.恢复中 4.修改参数中 5.应用参数组中 6.扩容预处理中 7.扩容预处理完成 8.修改端口中 9.迁移中 10.重置密码中 11.修改数据复制方式中 12.缩容预处理中 13.缩容预处理完成 15.内核小版本升级 17.迁移可用区中 18.修改备份配置中 20.停止中 21.已停止 22.启动中 26.白名单配置中",
				Validators: []validator.Int32{
					int32validator.Between(0, 26),
				},
			},
			"vip": schema.StringAttribute{
				Computed:    true,
				Description: "虚拟IP地址",
			},
			"write_port": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "写数据端口",
			},
			"read_port": schema.StringAttribute{
				Computed:    true,
				Description: "读端口",
			},
			"prod_db_engine": schema.StringAttribute{
				Computed:    true,
				Description: "数据库引擎",
			},
			"eip": schema.StringAttribute{
				Computed:    true,
				Description: "弹性ip",
			},
			"eip_status": schema.Int32Attribute{
				Computed:    true,
				Description: "弹性ip状态 0->unbind，1->bind,2->binding",
			},
			"ssl_status": schema.Int32Attribute{
				Computed:    true,
				Description: "Ssl状态 0->off，1->on",
			},
			"new_mysql_version": schema.StringAttribute{
				Computed:    true,
				Description: "mysql版本",
			},
			"audit_log_status": schema.Int32Attribute{
				Computed:    true,
				Description: "日志审计开关",
			},
			"inst_release_protection_status": schema.Int32Attribute{
				Computed:    true,
				Description: "实例释放保护开关 1:on,0:off",
			},
			"pause_enable": schema.BoolAttribute{
				Computed:    true,
				Description: "是否允许暂停",
			},
			"mysql_port": schema.StringAttribute{
				Computed:    true,
				Description: "数据库端口",
			},
			"security_group_status": schema.Int32Attribute{
				Computed:    true,
				Description: "安全组状态 0->normal, 1->changing, 2->deleted",
			},
			"running_control": schema.StringAttribute{
				Optional:    true,
				Description: "控制是否暂停，启用和重启实例，取值范围：freeze, unfreeze, restart",
				Validators: []validator.String{
					stringvalidator.OneOf("freeze", "unfreeze", "restart"),
				},
			},
			"prod_order_status": schema.Int32Attribute{
				Computed:    true,
				Description: "0.正常 1.欠费暂停 2.已注销 3.创建中 4.施工失败 5.到期退订状态 6.新增的状态-openApi暂停 7.创建完成等待变更单 8.待注销 9.手动暂停 10.手动退订",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "实例Id",
			},
		},
	}
}

func (c *CtyunMysqlInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMysqlInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 开始创建
	err = c.CreateMysqlInstance(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMergeMysqlInstance(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMysqlInstanceConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergeMysqlInstance(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunMysqlInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunMysqlInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.updateMysqlInstance(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeMysqlInstance(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunMysqlInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	deleteParams := &mysql.TeledbRefundRequest{
		InstId: state.InstID.ValueString(),
	}
	deleteHeader := &mysql.TeledbRefundRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		deleteHeader.ProjectID = state.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbRefundApi.Do(ctx, c.meta.Credential, deleteParams, deleteHeader)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 轮询确认时候退订成功
	err = c.DeleteLoop(ctx, &state, 60)
	if err != nil {
		return
	}
}

// CreateMysqlInstance 创建mysql实例
func (c *CtyunMysqlInstance) CreateMysqlInstance(ctx context.Context, config *CtyunMysqlInstanceConfig) (err error) {
	cycleType := config.CycleType.ValueString()
	params := &mysql.TeledbCreateRequest{
		BillMode:        business.MysqlBillMode[cycleType],
		RegionId:        config.RegionID.ValueString(),
		ProdVersion:     business.MysqlProdVersionDict[config.ProdID.ValueString()],
		VpcId:           config.VpcID.ValueString(),
		HostType:        config.HostType.ValueString(),
		SubnetId:        config.SubnetID.ValueString(),
		SecurityGroupId: config.SecurityGroupID.ValueString(),
		Name:            config.Name.ValueString(),
		Period:          config.CycleCount.ValueInt32(),
		Count:           1,
		ProdId:          business.MysqlProdIdDict[config.ProdID.ValueString()],
		CpuType:         business.MysqlCpuTypeDict[config.CpuType.ValueString()],
		OsType:          business.MysqlOSTypeDict[config.OsType.ValueString()],
	}
	if !config.Password.IsNull() && !config.Password.IsUnknown() {
		password := business.Encode(config.Password.ValueString())
		params.Password = password
	}
	if cycleType == business.OnDemandCycleType {
		params.AutoRenewStatus = 0
	} else {
		params.AutoRenewStatus = map[bool]int32{true: 1, false: 0}[config.AutoRenew.ValueBool()]
	}

	header := &mysql.TeledbCreateRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}

	var MysqlNodeInfos []mysql.MysqlNodeInfoListRequest

	mysqlNodeInfo := mysql.MysqlNodeInfoListRequest{}
	mysqlNodeInfo.NodeType = business.NodeTypeDict[config.ProdID.ValueString()]
	mysqlNodeInfo.InstSpec = business.MysqlInstanceSeriesDict[config.InstanceSeries.ValueString()]
	mysqlNodeInfo.StorageType = config.StorageType.ValueString()
	mysqlNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	mysqlNodeInfo.ProdPerformanceSpec = config.ProdPerformanceSpec.ValueString()
	mysqlNodeInfo.Disks = 1
	// 处理availabilityZoneInfo可用区信息
	var availabilityZoneInfos []mysql.AvailabilityZoneInfoRequest
	var availabilityZoneInfoList []AvailabilityZoneModel
	diag := config.AvailabilityZoneInfo.ElementsAs(ctx, &availabilityZoneInfoList, true)
	if diag.HasError() {
		return
	}
	for _, availabilityZoneInfoItem := range availabilityZoneInfoList {
		availabilityZoneInfo := mysql.AvailabilityZoneInfoRequest{}
		availabilityZoneInfo.AvailabilityZoneName = availabilityZoneInfoItem.AvailabilityZoneName.ValueString()
		availabilityZoneInfo.AvailabilityZoneCount = availabilityZoneInfoItem.AvailabilityZoneCount.ValueInt32()
		availabilityZoneInfo.NodeType = availabilityZoneInfoItem.NodeType.ValueString()
		availabilityZoneInfos = append(availabilityZoneInfos, availabilityZoneInfo)
	}
	mysqlNodeInfo.AvailabilityZoneInfo = availabilityZoneInfos
	MysqlNodeInfos = append(MysqlNodeInfos, mysqlNodeInfo)
	params.MysqlNodeInfoList = MysqlNodeInfos

	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbCreateApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 保存orderId
	if resp.ReturnObj.Data.NewOrderId == nil {
		err = errors.New("订单id为空，创建有误！")
		return
	}

	config.MasterOrderID = types.StringValue(*resp.ReturnObj.Data.NewOrderId)
	return
}

func (c *CtyunMysqlInstance) getAndMergeMysqlInstance(ctx context.Context, config *CtyunMysqlInstanceConfig) (err error) {
	// 若实例id为空，可能是因为实例刚创建，需要通过查询列表获取
	if config.InstID.ValueString() == "" {
		mysqlListParams := &mysql.TeledbGetListRequest{
			PageNow:      1,
			PageSize:     100,
			ProdInstName: config.Name.ValueStringPointer(),
		}
		mysqlListHeaders := &mysql.TeledbGetListHeaders{
			RegionID: config.RegionID.ValueString(),
		}
		if config.ProjectID.ValueString() != "" {
			mysqlListHeaders.ProjectID = config.ProjectID.ValueStringPointer()
		}

		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetListApi.Do(ctx, c.meta.Credential, mysqlListParams, mysqlListHeaders)
		if err2 != nil {
			err = err2
			return
		}
		if len(resp.ReturnObj.List) > 1 {
			err = errors.New("实例名重复！")
			return
		} else if len(resp.ReturnObj.List) < 1 {
			//若根据name查询不到机器，可能存在还未创建好的情况，需要轮询
			resp, err = c.ListLoop(ctx, mysqlListParams, mysqlListHeaders, 60)
			if err != nil {
				return
			}
			if len(resp.ReturnObj.List) != 1 {
				err = errors.New("未查询该实例mysql，mysql name:" + config.Name.ValueString())
				return
			}
		}
		config.InstID = types.StringValue(resp.ReturnObj.List[0].OuterProdInstId)
		config.ID = config.InstID
		// 确认资源是否开通完成
		// 若暂未开通完成，需要轮询等待
		if resp.ReturnObj.List[0].ProdOrderStatus != business.MysqlOrderStatusStarted {
			err = c.CreateLoop(ctx, mysqlListParams, mysqlListHeaders)
			if err != nil {
				return err
			}
		}
	}
	// 获取实例详情
	if config.InstID.ValueString() == "" {
		err = errors.New("查询实例详情时，实例 ID为空")
		return err
	}
	detailParams := &mysql.TeledbQueryDetailRequest{
		OuterProdInstId: config.InstID.ValueString(),
	}
	detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if config.ProjectID.ValueString() != "" {
		detailHeaders.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
	if err != nil {
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 处理实例详情
	returnOjb := resp.ReturnObj
	config.ProdRunningStatus = types.Int32Value(returnOjb.ProdRunningStatus)
	config.ProdOrderStatus = types.Int32Value(returnOjb.ProdOrderStatus)
	config.Vip = types.StringValue(returnOjb.Vip)
	config.ReadPort = types.StringValue(returnOjb.ReadPort)
	config.ProdDbEngine = types.StringValue(returnOjb.ProdDbEngine)
	config.EIP = types.StringValue(returnOjb.EIP)
	config.EipStatus = types.Int32Value(returnOjb.EIPStatus)
	config.SSlStatus = types.Int32Value(returnOjb.SSlStatus)
	config.NewMysqlVersion = types.StringValue(returnOjb.NewMysqlVersion)
	config.AuditLogStatus = types.Int32Value(returnOjb.AuditLogStatus)
	config.InstReleaseProtectionStatus = types.Int32Value(returnOjb.InstReleaseProtectionStatus)
	config.PauseEnable = types.BoolValue(returnOjb.PauseEnable)
	config.MysqlPort = types.StringValue(returnOjb.MysqlPort)
	config.SecurityGroupStatus = types.Int32Value(returnOjb.SecurityGroupStatus)
	config.Name = types.StringValue(returnOjb.ProdInstName)
	writePort, err := strconv.ParseInt(returnOjb.WritePort, 10, 32)
	if err != nil {
		return
	}
	config.WritePort = types.Int32Value(int32(writePort))

	// 更新disk， 主机配置相关信息
	config.ProdID = types.StringValue(business.MysqlProdIdRevDict[returnOjb.ProdId])

	config.StorageSpace = types.Int32Value(returnOjb.DiskSize)
	config.BackupStorageSpace = types.Int32Value(returnOjb.BackupDiskSize)
	config.ProdPerformanceSpec = types.StringValue(returnOjb.MachineSpec)
	return
}

func (c *CtyunMysqlInstance) CreateLoop(ctx context.Context, ListParams *mysql.TeledbGetListRequest, ListHeaders *mysql.TeledbGetListHeaders, loopCount ...int) (err error) {

	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetListApi.Do(ctx, c.meta.Credential, ListParams, ListHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}

			status := resp.ReturnObj.List[0].ProdOrderStatus
			switch status {
			case business.MysqlOrderStatusStarted:
				return false
			case business.MysqlOrderStatusCreating:
				return true
			case business.MysqlOrderStatusWaiting:
				return true
			case business.MysqlRunningStatusBackup:
				return true
			default:
				// 在开通的时候，其他状态是异常的，因此抛出异常，并跳出轮询
				err = errors.New("mysql创建状态有误： " + fmt.Sprintf("%d", status))
				return false
			}
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未创建成功！")
	}
	return
}

func (c *CtyunMysqlInstance) ListLoop(ctx context.Context, params *mysql.TeledbGetListRequest, headers *mysql.TeledbGetListHeaders, loopCount ...int) (*mysql.TeledbGetListResponse, error) {
	var err error
	var response *mysql.TeledbGetListResponse
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return nil, err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetListApi.Do(ctx, c.meta.Credential, params, headers)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}

			if len(resp.ReturnObj.List) > 1 {
				err = fmt.Errorf("查询到多条为名为%s的记录！", *params.ProdInstName)
				return false
			}
			if len(resp.ReturnObj.List) == 1 {
				response = resp
				return false
			}
			// 未查询到，继续轮询
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return nil, errors.New("轮询已达最大次数，资源仍未创建或查询到！")
	}
	return response, nil
}

func (c *CtyunMysqlInstance) UpgradeLoop(ctx context.Context, state *CtyunMysqlInstanceConfig, plan *CtyunMysqlInstanceConfig, loopCount ...int) (err error) {

	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			// 获取实例详情
			detailParams := &mysql.TeledbQueryDetailRequest{
				OuterProdInstId: state.InstID.ValueString(),
			}
			detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
				InstID:   state.InstID.ValueString(),
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			runningStatus := resp.ReturnObj.ProdRunningStatus
			orderStatus := resp.ReturnObj.ProdOrderStatus
			// 若符合预期，跳出循环，扩容成功
			if resp.ReturnObj.ProdId == business.MysqlProdIdDict[plan.ProdID.ValueString()] && resp.ReturnObj.DiskSize == plan.StorageSpace.ValueInt32() && resp.ReturnObj.MachineSpec == plan.ProdPerformanceSpec.ValueString() {
				//若备份磁盘空间不为空，且预期的分配磁盘空间与远端磁盘备份空间不相同，则继续轮询
				if plan.BackupStorageSpace.ValueInt32() != 0 && plan.BackupStorageSpace.ValueInt32() != resp.ReturnObj.BackupDiskSize {
					return true
				}
				if runningStatus == business.MysqlRunningStatusStarted && orderStatus == business.MysqlOrderStatusStarted {
					return false
				} else {
					return true
				}
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未升级成功！")
	}
	return
}

func (c *CtyunMysqlInstance) RunningStatusLoop(ctx context.Context, config *CtyunMysqlInstanceConfig, runningStatus int32, orderStatus int32, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			mysqlListParams := &mysql.TeledbGetListRequest{
				PageNow:      1,
				PageSize:     100,
				ProdInstName: config.Name.ValueStringPointer(),
			}
			mysqlListHeaders := &mysql.TeledbGetListHeaders{
				RegionID: config.RegionID.ValueString(),
			}
			if config.ProjectID.ValueString() != "" {
				mysqlListHeaders.ProjectID = config.ProjectID.ValueStringPointer()
			}

			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetListApi.Do(ctx, c.meta.Credential, mysqlListParams, mysqlListHeaders)
			if err2 != nil {
				err = err2
				return false
			}
			if len(resp.ReturnObj.List) > 1 {
				err = errors.New("实例名重复！")
				return false
			} else if len(resp.ReturnObj.List) < 1 {
				//若根据name查询不到机器，可能存在还未创建好的情况，需要轮询
				resp, err = c.ListLoop(ctx, mysqlListParams, mysqlListHeaders, 60)
				if err != nil {
					return false
				}
				if len(resp.ReturnObj.List) != 1 {
					err = errors.New("未查询该实例mysql，mysql name:" + config.Name.ValueString())
					return false
				}
			}

			currentRunningStatus := resp.ReturnObj.List[0].ProdRunningStatus
			currentOrderStatus := resp.ReturnObj.List[0].ProdOrderStatus
			if currentOrderStatus == orderStatus && currentRunningStatus == runningStatus {
				return false
			}
			return true

		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍完成状态更新！")
	}
	return
}

func (c *CtyunMysqlInstance) updateInfoLoop(ctx context.Context, state *CtyunMysqlInstanceConfig, plan *CtyunMysqlInstanceConfig, loopCount ...int) (err error) {

	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			// 获取实例详情
			detailParams := &mysql.TeledbQueryDetailRequest{
				OuterProdInstId: state.InstID.ValueString(),
			}
			detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
				InstID:   state.InstID.ValueString(),
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			//status := resp.ReturnObj.ProdRunningStatus
			// 跳出轮询条件如下：
			// 当state.name = plan.name，并且write_port无须更新时
			// 当state.name = plan.name，且write_port符合预期时
			if resp.ReturnObj.ProdInstName == plan.Name.ValueString() {
				if plan.WritePort.ValueInt32() == 0 {
					return false
				} else {
					if resp.ReturnObj.WritePort == fmt.Sprintf("%d", plan.WritePort.ValueInt32()) {
						return false
					} else {
						return true
					}
				}
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未更新成功！")
	}
	return
}

func (c *CtyunMysqlInstance) StartedLoop(ctx context.Context, state *CtyunMysqlInstanceConfig, loopCount ...int) (err error) {
	count := 30
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			// 获取实例详情
			detailParams := &mysql.TeledbQueryDetailRequest{
				OuterProdInstId: state.InstID.ValueString(),
			}
			detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
				InstID:   state.InstID.ValueString(),
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			runningStatus := resp.ReturnObj.ProdRunningStatus
			orderStatus := resp.ReturnObj.ProdOrderStatus
			if runningStatus == business.MysqlRunningStatusStarted && orderStatus == business.MysqlRunningStatusStarted {
				return false
			}
			if orderStatus == business.MysqlOrderStatusPause {
				err = errors.New("订单处于暂停状态，不可进行变更操作")
				return false
			}
			if runningStatus == business.MysqlRunningStatusStopping || runningStatus == business.MysqlRunningStatusStopped {
				err = errors.New("主机处于关机状态，不可进行变更操作")
				return false
			}

			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未到达启动状态！")
	}
	return
}

func (c *CtyunMysqlInstance) DeleteLoop(ctx context.Context, state *CtyunMysqlInstanceConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			params := &mysql.TeledbGetListRequest{
				PageNow:      1,
				PageSize:     100,
				ProdInstName: state.Name.ValueStringPointer(),
			}
			headers := &mysql.TeledbGetListHeaders{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				headers.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetListApi.Do(ctx, c.meta.Credential, params, headers)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			// 若查询列表已经查询不到，资源已经销毁
			if len(resp.ReturnObj.List) == 0 {
				return false
			}
			status := resp.ReturnObj.List[0].ProdOrderStatus
			switch status {
			case business.MysqlOrderStatusDestroy:
				return false
			case business.MysqlOrderStatusDestroyed:
				return false
			case business.MysqlOrderStatusStarted:
				return true
			case business.MysqlOrderStatusPause:
				return true
			default:
				err = errors.New("退订状态有误，当前状态为：" + fmt.Sprintf("%d", status))
				return false
			}
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未退订成功！")
	}
	return
}

func (c *CtyunMysqlInstance) updateMysqlInstance(ctx context.Context, state *CtyunMysqlInstanceConfig, plan *CtyunMysqlInstanceConfig) (err error) {
	if state.InstID.ValueString() == "" {
		err = errors.New("变配实例时，实例ID为空！")
		return err
	}
	// 修改实例名称
	if plan.Name.ValueString() != "" && state.Name.ValueString() != plan.Name.ValueString() {
		updateNameParams := &mysql.TeledbUpdateInstanceNameRequest{
			OuterProdInstID:     state.InstID.ValueString(),
			InstanceDescription: plan.Name.ValueString(),
		}
		updatedNameHeaders := &mysql.TeledbUpdateInstanceNameRequestHeader{
			InstID:   state.InstID.ValueString(),
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			updatedNameHeaders.ProjectID = state.ProjectID.ValueStringPointer()
		}
		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbUpdateInstanceNameApi.Do(ctx, c.meta.Credential, updateNameParams, updatedNameHeaders)
		if err2 != nil {
			err = err2
			return
		} else if resp.StatusCode != 0 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
	}

	// 修改实例写端口
	if plan.WritePort.ValueInt32() != 0 && state.WritePort.ValueInt32() != plan.WritePort.ValueInt32() {
		// 更新之前需要确定主机状态必须为started
		err = c.StartedLoop(ctx, state)
		if err != nil {
			return
		}
		updateWritePortParams := &mysql.TeledbUpdateWritePortRequest{
			OuterProdInstId: state.InstID.ValueString(),
			WritePort:       fmt.Sprintf("%d", plan.WritePort.ValueInt32()),
		}
		updateWritePortHeaders := &mysql.TeledbUpdateWritePortRequestHeader{
			InstID:   state.InstID.ValueString(),
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			updateWritePortHeaders.ProjectID = state.ProjectID.ValueString()
		}
		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbUpdateWritePortApi.Do(ctx, c.meta.Credential, updateWritePortParams, updateWritePortHeaders)
		if err2 != nil {
			return err2
		} else if resp.StatusCode != 0 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
	}
	// 轮询基础信息是否修改成功
	err = c.updateInfoLoop(ctx, state, plan)
	if err != nil {
		return
	}
	nodeType := business.NodeTypeDict[plan.ProdID.ValueString()]
	upgradeParams := &mysql.TeledbUpgradeRequest{
		InstId:   state.InstID.ValueString(),
		NodeType: &nodeType,
	}

	// 若StorageSpace不为空，触发主节点扩容存储空间
	if plan.StorageSpace.ValueInt32() != 0 && state.StorageSpace.ValueInt32() != plan.StorageSpace.ValueInt32() {
		upgradeParams.DiskVolume = plan.StorageSpace.ValueInt32Pointer()
	}
	// 若BackupStorageSpace不为空，触发备节点扩容存储空间
	if plan.BackupStorageSpace.ValueInt32() != 0 && state.BackupStorageSpace.ValueInt32() != plan.BackupStorageSpace.ValueInt32() {
		upgradeParams.DiskVolume = plan.BackupStorageSpace.ValueInt32Pointer()
		nodeType = business.PgsqlStorageTypeBackUp
		upgradeParams.NodeType = &nodeType
	}
	upgradeHeader := &mysql.TeledbUpgradeRequestHeader{}
	if plan.ProjectID.ValueString() != "" {
		upgradeHeader.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	// 扩容云数据库实例
	// 若plan.ProdPerformanceSpec不为空,且state和plan的ProdPerformanceSpec不一致，触发规格扩容
	if plan.ProdPerformanceSpec.ValueString() != "" && state.ProdPerformanceSpec.ValueString() != plan.ProdPerformanceSpec.ValueString() {
		upgradeParams.ProdPerformanceSpec = plan.ProdPerformanceSpec.ValueStringPointer()
	}
	// 若plan.prodId不为空,且state和plan的prodId不一致，触发实例类型扩容
	if plan.ProdID.ValueString() != "" && state.ProdID.ValueString() != plan.ProdID.ValueString() {
		prodId := business.MysqlProdIdDict[plan.ProdID.ValueString()]
		upgradeParams.ProdId = &prodId
	}
	// 若实例扩容或更新ProdID---从单节点升级至，一主一备、一主两备。需要补充AZ信息
	if upgradeParams.ProdPerformanceSpec != nil || upgradeParams.ProdId != nil {
		var azInfoList []AvailabilityZoneModel
		var upgradeAzList []mysql.AvailabilityZoneInfo

		diag := plan.AvailabilityZoneInfo.ElementsAs(ctx, &azInfoList, true)
		if diag.HasError() {
			return
		}

		for _, azInfoItem := range azInfoList {
			azInfo := mysql.AvailabilityZoneInfo{
				AvailabilityZoneName:  azInfoItem.AvailabilityZoneName.ValueString(),
				AvailabilityZoneCount: azInfoItem.AvailabilityZoneCount.ValueInt32(),
			}
			upgradeAzList = append(upgradeAzList, azInfo)
		}
		upgradeParams.AzList = upgradeAzList
	} else if !plan.AvailabilityZoneInfo.Equal(state.AvailabilityZoneInfo) {
		err = errors.New("未变配实例规格或者实例节点时，az info不可修改！")
		return err
	}

	// 若ProdPerformanceSpec, DiskVolume或者ProdId不为空时候，触发变配
	if upgradeParams.ProdPerformanceSpec != nil || upgradeParams.DiskVolume != nil || upgradeParams.ProdId != nil {
		// 更新之前需要确定主机状态必须为started
		err = c.StartedLoop(ctx, state)
		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbUpgradeApi.Do(ctx, c.meta.Credential, upgradeParams, upgradeHeader)
		if err2 != nil {
			err = err2
			return
		} else if resp.StatusCode != 200 {
			err = errors.New("扩容失败！")
			return
		}
		// 扩容后，轮循请求实例详情，确认已经完成升配
		err = c.UpgradeLoop(ctx, state, plan)
		if err != nil {
			return
		}
		// 扩容完成后，同步state AvailabilityZoneModel 状态
		if !plan.AvailabilityZoneInfo.IsNull() {
			state.AvailabilityZoneInfo = plan.AvailabilityZoneInfo
		}
	}

	// 启动实例
	if state.ProdOrderStatus.ValueInt32() == business.MysqlOrderStatusPause && plan.RunningControl.ValueString() == "unfreeze" {
		startParams := &mysql.TeledbStartRequest{
			OuterProdInstId: state.InstID.ValueString(),
		}
		startHeaders := &mysql.TeledbStartRequestHeader{
			InstID:   state.InstID.ValueString(),
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			startHeaders.ProjectID = state.ProjectID.ValueString()
		}
		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbStartApi.Do(ctx, c.meta.Credential, startParams, startHeaders)
		if err2 != nil {
			err = err2
			return
		} else if resp.StatusCode != 0 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
		// 轮询验证，是否已启动
		err = c.RunningStatusLoop(ctx, state, business.MysqlRunningStatusStarted, business.MysqlOrderStatusStarted, 60)
		if err != nil {
			return
		}
	}

	// 停止实例
	if plan.RunningControl.ValueString() == "freeze" {
		// 进行重启、停止实例时，确保实例处于started状态
		err = c.StartedLoop(ctx, state)
		if err != nil {
			return
		}
		pauseParams := &mysql.TeledbStopRequest{
			OuterProdInstId: state.InstID.ValueString(),
		}
		pauseHeader := &mysql.TeledbStopRequestHeader{
			InstID:   state.InstID.ValueString(),
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			pauseHeader.ProjectID = state.ProjectID.ValueString()
		}
		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbStopApi.Do(ctx, c.meta.Credential, pauseParams, pauseHeader)
		if err2 != nil {
			err = err2
			return err
		} else if resp.StatusCode != 0 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
		// 轮询验证，是否已停止，停止状态下，验证订单状态，预期=6
		err = c.RunningStatusLoop(ctx, state, business.MysqlRunningStatusStarted, business.MysqlOrderStatusPause, 60)
		if err != nil {
			return
		}
	}

	// 重启实例
	if plan.RunningControl.ValueString() == "restart" {
		// 进行重启、关机实例时，确保实例处于started状态
		err = c.StartedLoop(ctx, state)
		if err != nil {
			return
		}
		restartParams := &mysql.TeledbRestartRequest{
			OuterProdInstId: state.InstID.ValueString(),
		}
		restartHeader := &mysql.TeledbRestartRequestHeader{
			InstID:   state.InstID.ValueString(),
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			restartHeader.ProjectID = state.ProjectID.ValueString()
		}
		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbRestartApi.Do(ctx, c.meta.Credential, restartParams, restartHeader)
		if err2 != nil {
			err = err2
			return err
		} else if resp.StatusCode != 0 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
		//轮询验证，是否已完成重启
		err = c.RunningStatusLoop(ctx, state, business.MysqlRunningStatusStarted, business.MysqlOrderStatusStarted, 60)
		if err != nil {
			return
		}
	}
	state.RunningControl = plan.RunningControl
	return
}

type CtyunMysqlInstanceConfig struct {
	CycleType                   types.String `tfsdk:"cycle_type"`                     // 计费模式： 支持on_demand和month
	RegionID                    types.String `tfsdk:"region_id"`                      // 资源池Id
	VpcID                       types.String `tfsdk:"vpc_id"`                         // 虚拟私有云Id
	HostType                    types.String `tfsdk:"host_type"`                      // 主机类型 host type: S6 or S7
	SubnetID                    types.String `tfsdk:"subnet_id"`                      // 子网Id
	SecurityGroupID             types.String `tfsdk:"security_group_id"`              // 安全组
	Name                        types.String `tfsdk:"name"`                           // 集群名称
	Password                    types.String `tfsdk:"password"`                       // 管理员密码（RSA公钥加密）
	CycleCount                  types.Int32  `tfsdk:"cycle_count"`                    // 购买时长：单位月（范围：1-12，24，36）
	AutoRenew                   types.Bool   `tfsdk:"auto_renew"`                     // 自动续订状态
	ProdID                      types.String `tfsdk:"prod_id"`                        // 产品id
	CpuType                     types.String `tfsdk:"cpu_type"`                       // cpu类型：10是鲲鹏，20是海光，30是intel,40是amd,50是飞腾，60是龙芯，70是兆芯
	OsType                      types.String `tfsdk:"os_type"`                        // 系统类型：0是裸机，1是windows，2是centos，3是ubuntu，4是android，5是redhat，6是kylin，7是uos,8是suse，9是asianux，10是open_euler，11是ctyunos，12是euler
	MasterOrderID               types.String `tfsdk:"master_order_id"`                // 订单id
	InstID                      types.String `tfsdk:"inst_id"`                        // 实例id
	ProjectID                   types.String `tfsdk:"project_id"`                     // 项目id
	ProdRunningStatus           types.Int32  `tfsdk:"prod_running_status"`            // 以查询实例列表为主，0.正常 1.重启中 2.备份中 3.恢复中 4.修改参数中 5.应用参数组中 6&7.扩容中 8.修改端口中 9.迁移中 10.重置密码中
	Vip                         types.String `tfsdk:"vip"`                            // 虚拟IP地址
	WritePort                   types.Int32  `tfsdk:"write_port"`                     // 写数据端口
	ReadPort                    types.String `tfsdk:"read_port"`                      // 读端口
	ProdDbEngine                types.String `tfsdk:"prod_db_engine"`                 // 数据库引擎
	EIP                         types.String `tfsdk:"eip"`                            // 弹性ip
	EipStatus                   types.Int32  `tfsdk:"eip_status"`                     // 弹性ip状态 0->unbind，1->bind,2->binding
	SSlStatus                   types.Int32  `tfsdk:"ssl_status"`                     // Ssl状态 0->off，1->on
	NewMysqlVersion             types.String `tfsdk:"new_mysql_version"`              // mysql版本
	AuditLogStatus              types.Int32  `tfsdk:"audit_log_status"`               // 日志审计开关
	InstReleaseProtectionStatus types.Int32  `tfsdk:"inst_release_protection_status"` // 实例释放保护开关 1:on,0:off
	PauseEnable                 types.Bool   `tfsdk:"pause_enable"`                   // 是否允许暂停
	MysqlPort                   types.String `tfsdk:"mysql_port"`                     // 数据库端口
	SecurityGroupStatus         types.Int32  `tfsdk:"security_group_status"`          // 安全组状态 0->normal, 1->changing, 2->deleted
	InstanceSeries              types.String `tfsdk:"instance_series"`                // 实例规格（默认：通用型=1） InstSpec
	StorageType                 types.String `tfsdk:"storage_type"`                   // 存储类型：SSD, SATA, SAS, SSD-genric, FAST-SSD
	StorageSpace                types.Int32  `tfsdk:"storage_space"`                  // 存储空间（单位：GB，范围100到32768）
	BackupStorageSpace          types.Int32  `tfsdk:"backup_storage_space"`           // 备份节点，存储空间扩容使用
	ProdPerformanceSpec         types.String `tfsdk:"prod_performance_spec"`          // 规格（例：4C8G）
	AvailabilityZoneInfo        types.List   `tfsdk:"availability_zone_info"`         // 可用区信息
	RunningControl              types.String `tfsdk:"running_control"`                //
	ProdOrderStatus             types.Int32  `tfsdk:"prod_order_status"`
	ID                          types.String `tfsdk:"id"` // 实例id
}

type AvailabilityZoneModel struct {
	AvailabilityZoneName  types.String `tfsdk:"availability_zone_name"`  // 资源池可用区名称
	AvailabilityZoneCount types.Int32  `tfsdk:"availability_zone_count"` // 资源池可用区总数
	NodeType              types.String `tfsdk:"node_type"`               // 表示分布AZ的节点类型，master/slave
}

type UpdatedAZModel struct {
	AvailabilityZoneName  types.String `tfsdk:"availability_zone_name"`  // 资源池可用区名称
	AvailabilityZoneCount types.Int32  `tfsdk:"availability_zone_count"` // 资源池可用区总数
}
