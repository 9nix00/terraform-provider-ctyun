package pgsql

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
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
	"terraform-provider-ctyun/internal/utils"
	"time"
)

var (
	_ resource.Resource                = &CtyunPostgresqlInstance{}
	_ resource.ResourceWithConfigure   = &CtyunPostgresqlInstance{}
	_ resource.ResourceWithImportState = &CtyunPostgresqlInstance{}
)

type CtyunPostgresqlInstance struct {
	meta *common.CtyunMetadata
}

func (c *CtyunPostgresqlInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunPostgresqlInstance) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func NewCtyunPostgresqlInstance() resource.Resource {
	return &CtyunPostgresqlInstance{}
}

func (c *CtyunPostgresqlInstance) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_instance"
}

func (c *CtyunPostgresqlInstance) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "pgsql provider",
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
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id,如果不填这默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"host_type": schema.StringAttribute{
				Required:    true,
				Description: "主机类型 host type: S6 or S7等。可根据data.ctyun_postgresql_specs查询",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"prod_id": schema.StringAttribute{
				Required:    true,
				Description: "产品ID。Single1222-（单实例12.22版本）, MasterSlave1222（一主一备12.22版本）, Single1417（单实例14.17版本）, MasterSlave1417（一主一备14.17版本）, Single1320（单实例13.20版本）, MasterSlave1320（一主一备13.20版本）, ReadOnly1222（只读实例12.22版本）, ReadOnly1320（只读实例13.20版本）, ReadOnly1417（只读实例14.17版本）, Single1512（单实例15.12版本）, MasterSlave1512（一主一备15.12版本）, ReadOnly1512（只读实例15.12版本）, Master2Slave1222（一主两备12.22版本）, Master2Slave1417（一主两备14.17版本）, Master2Slave1320（一主两备13.20版本）, Master2Slave1512（一主两备15.12版本）, Single168（单实例16.8版本）, MasterSlave168（一主一备16.8版本）, Master2Slave168（一主两备16.8版本）, ReadOnly168（只读实例16.8版本）。注：扩容过程中，不支持磁盘、规格和实例扩容同时进行",
				Validators: []validator.String{
					stringvalidator.OneOf(business.PgsqlProdIds...),
				},
			},
			// 存储与备份
			"backup_storage_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "备份存储类型: SSD=超高IO, SATA=普通IO, SAS=高IO, SSD-genric=通用型SSD, FAST-SSD=极速型SSD",
				Validators: []validator.String{
					stringvalidator.OneOf(business.StorageType...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_type": schema.StringAttribute{
				Required:    true,
				Description: "主存储类型: SSD=超高IO, SATA=普通IO, SAS=高IO, SSD-genric=通用型SSD, FAST-SSD=极速型SSD",
				Validators: []validator.String{
					stringvalidator.OneOf(business.StorageType...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_space": schema.Int32Attribute{
				Required:    true,
				Description: "主存储空间(单位:G，范围100-32768)。扩容过程中不支持磁盘、规格和实例扩容同时进行",
				Validators: []validator.Int32{
					int32validator.Between(100, 32768),
				},
			},
			"backup_storage_space": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "备份存储空间大小",
			},
			// 网络配置
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "虚拟私有云Id",
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
			"appoint_vip": schema.StringAttribute{
				Optional:    true,
				Description: "指定VIP",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// 实例配置
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
				Description: "实例密码（8-32位由大写字母、小写字母、数字、特殊字符中的任意三种组成 特殊字符为!@#$%^&*()_+-=），RSA公钥加密存储",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 32),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// 订购选项
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
			// 高级配置
			"case_sensitive": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否区分大小写: true=区分, false=不区分。默认不区分",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
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
			// 项目相关
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
			},
			"is_mgr": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否开启MRG，默认false",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},

			// 节点配置 (重要: 嵌套对象列表)
			"instance_series": schema.StringAttribute{
				Required:    true,
				Description: "实例规格，取值范围:S(通用型)， C(计算增强型)，M(内存增强型)",
				Validators: []validator.String{
					stringvalidator.OneOf(business.MysqlInstanceSeries...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"prod_performance_spec": schema.StringAttribute{
				Required:    true,
				Description: "实例规格(例: 4C8G)。可根据data.ctyun_postgresql_specs获取。不支持规格和实例扩容同时进行：prod_id和prod_performance_spec不能同时与原配置不一致",
			},
			// 自动扩展配置
			"auto_scale": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "存储自动扩容: false=关闭, true=开启。默认关闭",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			// todo validator
			"max_scale": schema.Int32Attribute{
				Optional:    true,
				Description: "存储扩容上限(单位G)",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"active_scale_rate": schema.Int32Attribute{
				Optional:    true,
				Description: "触发扩容百分比(1-100)",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"availability_zone_info": schema.ListNestedAttribute{
				Required: true,
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
							Description: "节点类型(master/readNode)",
						},
						//"display_name": schema.StringAttribute{
						//	Optional:    true,
						//	Description: "可用区显示名",
						//},
						//"spec_id": schema.StringAttribute{
						//	Optional:    true,
						//	Description: "规格ID",
						//},
					},
				},
			},
			//"new_order_id": schema.StringAttribute{
			//	Computed:    true,
			//	Description: "订单id",
			//},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "pgsql 实例id",
			},
			"alive": schema.Int32Attribute{
				Computed:    true,
				Description: "实例是否存活,0:存活，-1:异常",
			},
			"disk_rated": schema.Int32Attribute{
				Computed:    true,
				Description: "磁盘使用率",
			},
			"outer_prod_inst_id": schema.StringAttribute{
				Computed:    true,
				Description: "对外的实例ID，对应PaaS平台",
			},
			"prod_db_engine": schema.StringAttribute{
				Computed:    true,
				Description: "数据库实例引擎",
			},
			"prod_order_status": schema.Int32Attribute{
				Computed:    true,
				Description: "订单状态，0：正常，1：冻结，2：删除，3：操作中，4：失败,2005:扩容中",
			},
			"prod_running_status": schema.Int32Attribute{
				Computed:    true,
				Description: "实例状态",
			},
			"prod_type": schema.Int32Attribute{
				Computed:    true,
				Description: "实例部署方式 0：单机部署,1：主备部署",
			},
			"read_port": schema.Int32Attribute{
				Computed:    true,
				Description: "读端口",
			},
			"write_port": schema.StringAttribute{
				Computed:    true,
				Description: "写端口",
			},
			"tool_type": schema.Int32Attribute{
				Computed:    true,
				Description: "备份工具类型，1：pg_baseback，2：pgbackrest，3：s3",
			},

			"running_control": schema.StringAttribute{
				Optional:    true,
				Description: "控制是否暂停，启用和重启实例，取值范围：stop, start, restart",
				Validators: []validator.String{
					stringvalidator.OneOf("stop", "start", "restart"),
				},
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订单id",
			},
		},
	}
}

func (c *CtyunPostgresqlInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPostgresqlInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 开始创建
	err = c.CreatePgsqlInstance(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取pgsql详情
	err = c.getAndMergePgsqlInstance(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPostgresqlInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPostgresqlInstanceConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergePgsqlInstance(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "未找到实例") {
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

func (c *CtyunPostgresqlInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunPostgresqlInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunPostgresqlInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.updatePgsqlInstance(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergePgsqlInstance(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPostgresqlInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunPostgresqlInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	deleteParams := &pgsql.PgsqlRefundRequest{
		InstId: state.ID.ValueString(),
	}
	deleteHeader := &pgsql.PgsqlRefundRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		deleteHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}

	// 确保订单已完成状态才能退订
	err = c.StartedOrderLoop(ctx, &state, business.MysqlOrderStatusStarted, 60)
	if err != nil {
		return
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlRefundApi.Do(ctx, c.meta.Credential, deleteParams, deleteHeader)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 轮询确认时候退订成功
	err = c.RefundLoop(ctx, &state, 60)
	if err != nil {
		return
	}
	response.Diagnostics.AddWarning("删除PostgreSql集群成功", "集群退订后，若立即删除子网或安全组可能会失败，需要等待底层资源释放")
}

func (c *CtyunPostgresqlInstance) CreatePgsqlInstance(ctx context.Context, config *CtyunPostgresqlInstanceConfig) (err error) {
	cycleType := config.CycleType.ValueString()
	isMgr := fmt.Sprintf("%t", config.IsMGR.ValueBool())

	// 构建创建参数
	params := &pgsql.PgsqlCreateRequest{
		RegionId:        config.RegionID.ValueString(),
		BillMode:        business.MysqlBillMode[cycleType],
		HostType:        config.HostType.ValueString(),
		ProdVersion:     business.PgsqlProdVersionDict[config.ProdID.ValueString()],
		ProdId:          business.PgsqlProdIDDict[config.ProdID.ValueString()],
		VpcId:           config.VpcID.ValueString(),
		SubnetId:        config.SubnetId.ValueString(),
		SecurityGroupId: config.SecurityGroupId.ValueString(),
		Name:            config.Name.ValueString(),
		Period:          config.CycleCount.ValueInt32(),
		Count:           1,
		IsMGR:           &isMgr,
	}
	// 处理密码，并对密码进行RSA加密
	encodePassword := business.Encode(config.Password.ValueString())
	params.Password = encodePassword
	//config.Password = types.StringValue(encodePassword)
	//params.Password = config.Password.ValueString()
	if config.CaseSensitive.ValueBool() {
		params.CaseSensitive = "0"
	} else {
		params.CaseSensitive = "1"
	}
	if config.CycleType.ValueString() == business.OnDemandCycleType {
		params.AutoRenewStatus = 0
	} else {
		params.AutoRenewStatus = map[bool]int32{true: 1, false: 0}[config.AutoRenew.ValueBool()]
	}

	header := &pgsql.PgsqlCreateRequestHeader{}
	if config.ProjectID.ValueString() != "0" {
		params.ProjectId = config.ProjectID.ValueStringPointer()
		header.ProjectId = config.ProjectID.ValueStringPointer()
	}
	// 处理backupStorage
	if config.BackupStorageType.ValueString() != "" {
		params.BackupStorageType = config.BackupStorageType.ValueStringPointer()
	}

	if config.AppointVip.ValueString() != "" {
		params.AppointVip = config.AppointVip.ValueStringPointer()
	}
	if config.OsType.ValueString() != "" {
		osType := business.MysqlOSTypeDict[config.OsType.ValueString()]
		params.OsType = &osType
	}
	if config.CpuType.ValueString() != "" {
		cpuType := business.MysqlCpuTypeDict[config.CpuType.ValueString()]
		params.CpuType = &cpuType
	}

	// 处理MysqlNodeInfoList
	mysqlNodeInfoList := []pgsql.PgsqlCreateRequestMysqlNodeInfoList{}
	mysqlNodeInfo := pgsql.PgsqlCreateRequestMysqlNodeInfoList{}
	mysqlNodeInfo.InstSpec = business.PgsqlInstanceSeriesDict[config.InstanceSeries.ValueString()]
	mysqlNodeInfo.StorageType = config.StorageType.ValueString()
	mysqlNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	mysqlNodeInfo.ProdPerformanceSpec = config.ProdPerformanceSpec.ValueString()
	mysqlNodeInfo.Disks = 1
	mysqlNodeInfo.NodeType = business.PgsqlNodeTypeDict[config.ProdID.ValueString()]

	if config.BackupStorageSpace.ValueInt32() != 0 {
		backupStorageSpace := fmt.Sprintf("%d", config.BackupStorageSpace.ValueInt32())
		mysqlNodeInfo.BackupStorageSpace = &backupStorageSpace
	}
	// 处理availabilityZoneInfo
	azModelList := []pgsql.PgsqlCreateRequestAvailabilityZoneInfo{}
	availabilityZoneModel := []AvailabilityZoneModel{}
	diag := config.AvailabilityZoneInfo.ElementsAs(ctx, &availabilityZoneModel, true)
	if diag.HasError() {
		err = errors.New("AvailabilityZoneInfo解析错误")
		return
	}
	for _, azModelItem := range availabilityZoneModel {
		azModel := pgsql.PgsqlCreateRequestAvailabilityZoneInfo{}
		azModel.AvailabilityZoneName = azModelItem.AvailabilityZoneName.ValueString()
		azModel.AvailabilityZoneCount = azModelItem.AvailabilityZoneCount.ValueInt32()
		azModel.NodeType = azModelItem.NodeType.ValueString()
		azModelList = append(azModelList, azModel)
	}
	mysqlNodeInfo.AvailabilityZoneInfo = azModelList
	mysqlNodeInfoList = append(mysqlNodeInfoList, mysqlNodeInfo)
	params.MysqlNodeInfoList = mysqlNodeInfoList
	// 处理AutoScaleParam
	autoScaleParams := pgsql.PgsqlCreateRequestAutoScaleParam{}
	autoScaleParams.AutoScale = fmt.Sprintf("%t", config.AutoScale.ValueBool())
	if config.MaxScale.ValueInt32() != 0 {
		autoScaleParams.MaxScale = int64(config.MaxScale.ValueInt32())
	}
	if config.ActiveScaleRate.ValueInt32() != 0 {
		autoScaleParams.ActiveScaleRate = fmt.Sprintf("%d", config.ActiveScaleRate.ValueInt32())
	}
	if autoScaleParams.AutoScale != "" || autoScaleParams.MaxScale != 0 || autoScaleParams.ActiveScaleRate != "" {
		params.AutoScaleParam = &autoScaleParams
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlCreateApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	//else if resp.ReturnObj == nil {
	//	err = common.InvalidReturnObjError
	//	return
	//}
	// 保存orderId
	//if resp.ReturnObj.NewOrderId == nil {
	//	err = errors.New("订单id为空，创建有误！")
	//	return
	//}
	config.MasterOrderID = utils.SecStringValue(resp.ReturnObj.Data.NewOrderId)
	return
}

func (c *CtyunPostgresqlInstance) getAndMergePgsqlInstance(ctx context.Context, config *CtyunPostgresqlInstanceConfig) (err error) {
	// 通过查询list获取instId
	// 若实例id为空，表示需要轮询列表查询inst id
	if config.ID.IsNull() || config.ID.ValueString() == "" {
		if config.Name.ValueStringPointer() == nil {
			err = errors.New("实例名为空，有误！")
			return
		}
		listParams := &pgsql.PgsqlListRequest{
			PageNum:      1,
			PageSize:     100,
			ProdInstName: config.Name.ValueStringPointer(),
		}
		listHeaders := &pgsql.PgsqlListRequestHeader{
			RegionID: config.RegionID.ValueString(),
		}
		if config.ProjectID.ValueString() != "" {
			listHeaders.ProjectID = config.ProjectID.ValueStringPointer()
		}
		resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlListApi.Do(ctx, c.meta.Credential, listParams, listHeaders)
		if err2 != nil {
			return err2
		} else if resp.StatusCode != 800 {
			err = fmt.Errorf("API return error. Message: %s", *resp.Message)
			return
		}
		if len(resp.ReturnObj.List) > 1 {
			err = errors.New("实例名称重复，存在异常！")
			return
		} else if len(resp.ReturnObj.List) == 0 {
			// 若列表为空，说明实例未创建成功，继续轮询查询
			err = c.ListLoop(ctx, config, listParams, listHeaders)
			if err != nil {
				return
			}
		}
	}
	if config.ID.ValueString() == "" {
		err = errors.New("实例id为空")
		return
	}
	if config.ID.ValueString() == "" {
		err = errors.New("在查询实例详情时，实例id为空")
	}
	// 获取pgsql详情
	detailParams := &pgsql.PgsqlDetailRequest{
		ProdInstId: config.ID.ValueString(),
	}
	detailHeaders := &pgsql.PgsqlDetailRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if config.ProjectID.ValueString() != "" {
		detailHeaders.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
	if err != nil {
		return err
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析pgsql 详情
	returnObj := resp.ReturnObj
	config.Alive = types.Int32Value(returnObj.Alive)
	config.StorageSpace = types.Int32Value(returnObj.DiskSize)
	config.StorageType = types.StringValue(returnObj.DiskType)
	config.ProdPerformanceSpec = types.StringValue(returnObj.MachineSpec)
	config.OuterProdInstId = types.StringValue(returnObj.OuterProdInstId)
	config.ProdDbEngine = types.StringValue(returnObj.ProdDbEngine)
	config.Name = types.StringValue(returnObj.ProdInstName)
	config.ProdType = types.Int32Value(returnObj.ProdType)
	prodId, err := strconv.ParseInt(returnObj.SpuCode, 10, 64)
	config.ProdID = types.StringValue(business.PgsqlProdIDRevDict[prodId])
	config.ReadPort = types.Int32Value(returnObj.ReadPort)
	config.WritePort = types.StringValue(returnObj.WritePort)
	config.CycleType = types.StringValue(map[int32]string{1: "month", 2: "on_demand"}[returnObj.BillMode])
	config.SecurityGroupId = types.StringValue(returnObj.SecurityGroupId)
	config.DiskRated = types.Int32Value(returnObj.DiskRated)
	config.ProdOrderStatus = types.Int32Value(returnObj.ProdOrderStatus)
	config.ProdRunninStatus = types.Int32Value(returnObj.ProdRunningStatus)
	config.ToolType = types.Int32Value(returnObj.ToolType)
	config.BackupStorageType = types.StringValue(returnObj.BackupDiskType)
	backupDiskSize := c.ParseStorageSize(&returnObj.BackupDiskSize)
	diskSize, err := strconv.ParseInt(backupDiskSize, 10, 32)
	if err != nil {
		return
	}
	config.BackupStorageSpace = types.Int32Value(int32(diskSize))
	return
}

func (c *CtyunPostgresqlInstance) updatePgsqlInstance(ctx context.Context, state *CtyunPostgresqlInstanceConfig, plan *CtyunPostgresqlInstanceConfig) (err error) {
	if state.ID.ValueString() == "" {
		err = errors.New("在变配实例时，实例id为空")
	}

	// 修改RDS实例名称
	// 当plan name不为空，且plan name 与 state name 不一致时，触发实例名称修改
	if plan.Name.ValueString() != "" && plan.Name.ValueString() != state.Name.ValueString() {
		// 确保操作时，实例处于running状态，避免更新失败
		err = c.RunningStatusLoop(ctx, state, business.MysqlRunningStatusStarted, business.MysqlOrderStatusStarted, 30)
		if err != nil {
			return
		}
		modifyNameParams := &pgsql.PgsqlUpdateInstanceNameRequest{
			ProdInstId:   state.ID.ValueString(),
			InstanceName: plan.Name.ValueString(),
		}
		modifyNameHeaders := &pgsql.PgsqlUpdateInstanceNameRequestHeader{
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			modifyNameHeaders.ProjectID = state.ProjectID.ValueStringPointer()
		}
		resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlUpdateInstanceNameApi.Do(ctx, c.meta.Credential, modifyNameParams, modifyNameHeaders)
		if err2 != nil {
			return err2
		} else if resp.StatusCode != 800 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
	}
	// 更变安全组
	if plan.SecurityGroupId.ValueString() != "" && plan.SecurityGroupId.ValueString() != state.SecurityGroupId.ValueString() {
		// 确保操作时，实例处于running状态，避免更新失败
		err = c.RunningStatusLoop(ctx, state, business.MysqlRunningStatusStarted, business.MysqlOrderStatusStarted, 30)
		if err != nil {
			return
		}

		deleteSgParams := &pgsql.PgsqlDeleteSecurityGroupRequest{
			SecurityGroupId: state.SecurityGroupId.ValueString(),
			InstanceId:      state.SecurityGroupId.ValueString(),
		}
		deleteSgHeader := &pgsql.PgsqlDeleteSecurityGroupRequestHeader{}
		updateSecurityGroupParams := &pgsql.PgsqlUpdateSecurityGroupRequest{
			SecurityGroupId:    state.SecurityGroupId.ValueString(),
			InstanceId:         state.ID.ValueString(),
			NewSecurityGroupId: plan.SecurityGroupId.ValueString(),
		}
		updatedSecurityGroupHeaders := &pgsql.PgsqlUpdateSecurityGroupRequestHeader{}
		if !state.ProjectID.IsNull() {
			updatedSecurityGroupHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			deleteSgHeader.ProjectID = state.ProjectID.ValueStringPointer()
		}

		// 先替换
		resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlUpdateSecurityGroupApi.Do(ctx, c.meta.Credential, updateSecurityGroupParams, updatedSecurityGroupHeaders)
		if err2 != nil {
			return err2
		} else if resp.StatusCode != 200 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
		// 再解绑原securityGroup
		// 先解绑原来的安全组
		deleteResp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlDeleteSecurityGroupApi.Do(ctx, c.meta.Credential, deleteSgParams, deleteSgHeader)
		if err2 != nil {
			return err
		} else if deleteResp.StatusCode != 200 {
			err = fmt.Errorf("API return error. Message: %s", deleteResp.Message)
			return
		}
	}
	// 轮询确认姓名和安全组修改成功
	err = c.InfoLoop(ctx, state, plan, 60)
	if err != nil {
		return
	}

	// 扩容云数据库实例
	// 扩容磁盘
	nodeType := business.PgsqlNodeTypeDict[plan.ProdID.ValueString()]
	upgradeParams := &pgsql.PgsqlUpgradeRequest{
		InstId:   state.ID.ValueString(),
		NodeType: &nodeType,
	}

	upgradeHeaders := &pgsql.PgsqlUpgradeRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		upgradeHeaders.ProjectID = state.ProjectID.ValueStringPointer()
	}
	// 若StorageSpace不为空，触发扩容主存储空间，且plan storageSpace与state不相同时
	if plan.StorageSpace.ValueInt32() != 0 && plan.StorageSpace.ValueInt32() != state.StorageSpace.ValueInt32() {
		upgradeParams.DiskVolume = plan.StorageSpace.ValueInt32Pointer()
	}
	// 若backupStorageSpace不为空，触发备用存储空间扩容，且plan backupStorageSpace与state不相同
	if plan.BackupStorageSpace.ValueInt32() != 0 && plan.BackupStorageSpace.ValueInt32() != state.BackupStorageSpace.ValueInt32() {
		storageSize32 := plan.BackupStorageSpace.ValueInt32()
		upgradeParams.DiskVolume = &storageSize32
		nodeType = "backup"
		upgradeParams.NodeType = &nodeType
	}
	// 规格扩容
	if plan.ProdPerformanceSpec.ValueString() != "" && plan.ProdPerformanceSpec.ValueString() != state.ProdPerformanceSpec.ValueString() {
		upgradeParams.ProdPerformanceSpec = plan.ProdPerformanceSpec.ValueStringPointer()
	}
	// 类型扩容,单机到一主一备， 单机到一主两备，一主一备到一主两备
	// 若plan.prodID不为空
	if plan.ProdID.ValueString() != "" && plan.ProdID.ValueString() != state.ProdID.ValueString() {
		prodId := business.PgsqlProdIDDict[plan.ProdID.ValueString()]
		upgradeParams.ProdId = &prodId
	}
	// 若规格不为空或者prodID不为空的情况下，需要配置azList参数
	if upgradeParams.ProdPerformanceSpec != nil || upgradeParams.ProdId != nil {
		var availabilityZoneList []AvailabilityZoneModel
		var azInfoList []pgsql.PgsqlUpgradeRequestAzList
		diags := plan.AvailabilityZoneInfo.ElementsAs(ctx, &availabilityZoneList, true)
		if diags.HasError() {
			return
		}
		for _, azInfoItem := range availabilityZoneList {
			var azInfo pgsql.PgsqlUpgradeRequestAzList
			azInfo.AvailabilityZoneName = azInfoItem.AvailabilityZoneName.ValueString()
			azInfo.AvailabilityZoneCount = azInfoItem.AvailabilityZoneCount.ValueInt32()
			azInfoList = append(azInfoList, azInfo)
		}
		upgradeParams.AzList = azInfoList
	}
	// 若更新参数不为空，触发更新
	if upgradeParams.DiskVolume != nil || upgradeParams.ProdPerformanceSpec != nil || upgradeParams.ProdId != nil {
		// 确保操作时，实例处于running状态，避免更新失败
		err = c.RunningStatusLoop(ctx, state, business.MysqlRunningStatusStarted, business.MysqlOrderStatusStarted, 30)
		if err != nil {
			return
		}
		// 若机器刚创建完成，需要同步实例远端状态
		err = c.getAndMergePgsqlInstance(ctx, state)
		if err != nil {
			return
		}
		// 判断prod_id 和spec和disks是否同时需要更新，若是，则返回报错
		if c.countSame(plan, state) > 1 {
			err = errors.New("不支持磁盘、规格和主备扩容同时进行！")
			return
		}

		detailParams := &pgsql.PgsqlDetailRequest{
			ProdInstId: state.ID.ValueString(),
		}
		detailHeaders := &pgsql.PgsqlDetailRequestHeader{
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
		}
		respt, errt := c.meta.Apis.SdkCtPgsqlApis.PgsqlDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
		fmt.Println(respt.ReturnObj.ProdRunningStatus, respt.ReturnObj.ProdOrderStatus)
		if errt != nil {
			err = errt
			return
		}
		resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlUpgradeApi.Do(ctx, c.meta.Credential, upgradeParams, upgradeHeaders)
		if err2 != nil {
			return err2
		} else if resp.StatusCode != 200 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
		// 更新后，轮询确认时候更新完成
		err = c.UpgradeLoop(ctx, state, plan, 100)
		if err != nil {
			return
		}
		// 扩容完成后，同步state AvailabilityZoneModel 状态
		if !plan.AvailabilityZoneInfo.IsNull() {
			state.AvailabilityZoneInfo = plan.AvailabilityZoneInfo
		}

	}
	// 启动实例
	if plan.RunningControl.ValueString() == "start" {
		err = c.startInstance(ctx, state, plan)
		if err != nil {
			return
		}
	}
	// 停用实例
	// 当state ProdRunningState = started时，才可以停用实例
	if plan.RunningControl.ValueString() == "stop" && state.ProdRunninStatus.ValueInt32() == business.PgsqlProdRunningStatusStarted {
		// 确保操作时，实例处于running状态，避免更新失败
		err = c.RunningStatusLoop(ctx, state, business.MysqlRunningStatusStarted, business.MysqlOrderStatusStarted, 30)
		if err != nil {
			return
		}
		stopParams := &pgsql.PgsqlStopRequest{
			ProdInstId: state.ID.ValueString(),
		}
		stopHeaders := &pgsql.PgsqlStopRequestHeader{
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			stopHeaders.ProjectID = state.ProjectID.ValueString()
		}
		resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlStopApi.Do(ctx, c.meta.Credential, stopParams, stopHeaders)
		if err2 != nil {
			return err2
		} else if resp.StatusCode != 800 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}

		err = c.RunningStatusLoop(ctx, state, business.PgsqlProdRunningStatusStopped, business.MysqlOrderStatusStarted)
		if err != nil {
			return
		}
	}
	// 重启实例
	if plan.RunningControl.ValueString() == "restart" {
		// 确保操作时，实例处于running状态，避免更新失败
		err = c.RunningStatusLoop(ctx, state, business.MysqlRunningStatusStarted, business.MysqlOrderStatusStarted, 30)
		if err != nil {
			return
		}
		restartParams := &pgsql.PgsqlRestartRequest{
			ProdInstId: state.ID.ValueString(),
		}
		restartHeaders := &pgsql.PgsqlRestartRequestHeader{
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			restartHeaders.ProjectID = state.ProjectID.ValueString()
		}
		resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlRestartApi.Do(ctx, c.meta.Credential, restartParams, restartHeaders)
		if err2 != nil {
			return err2
		} else if resp.StatusCode != 800 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
		err = c.RunningStatusLoop(ctx, state, business.PgsqlProdRunningStatusStarted, business.PgsqlProdOrderStatusRunning)
		if err != nil {
			return
		}
	}
	state.RunningControl = plan.RunningControl
	return
}

func (c *CtyunPostgresqlInstance) ListLoop(ctx context.Context, config *CtyunPostgresqlInstanceConfig, params *pgsql.PgsqlListRequest, headers *pgsql.PgsqlListRequestHeader, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	dealyCount := 2
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlListApi.Do(ctx, c.meta.Credential, params, headers)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			}
			if len(resp.ReturnObj.List) > 1 {
				err = errors.New("实例名称重复，存在异常！")
				return false
			} else if len(resp.ReturnObj.List) == 0 {
				return true
			} else if resp.ReturnObj.List[0].ProdInstId == "" {
				return true
			} else {
				if resp.ReturnObj.List[0].ProdOrderStatus != business.PgsqlProdOrderStatusRunning || resp.ReturnObj.List[0].ProdRunningStatus != business.PgsqlProdRunningStatusStarted {
					return true
				}
				// 确保与页面保持一致
				if dealyCount > 0 {
					dealyCount--
					return true
				}
				config.ID = types.StringValue(resp.ReturnObj.List[0].ProdInstId)

				return false
			}
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未创建或查询到！")
	}
	return
}

func (c *CtyunPostgresqlInstance) RunningStatusLoop(ctx context.Context, state *CtyunPostgresqlInstanceConfig, runningStatus int32, orderStatus int32, loopCount ...int) (err error) {

	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	// 设置一个容忍机制，根据pgsql定制，pgsql开通等操作可能出现报错回滚，此期间会存在查询不到实例情况
	tolerateCount := 20
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	// 因为pgsql console和openapi有一个同步误差时间，需要多轮询几轮，目前暂定4轮
	syncCount := 4
	result := retryer.Start(
		func(currentTime int) bool {
			detailParams := &pgsql.PgsqlDetailRequest{
				ProdInstId: state.ID.ValueString(),
			}
			detailHeaders := &pgsql.PgsqlDetailRequestHeader{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				tolerateCount--
				if tolerateCount < 0 {
					err = err2
					return false
				}
				return true
			} else if resp.StatusCode != 800 {
				tolerateCount--
				if tolerateCount < 0 {
					err = fmt.Errorf("API return error. Message: %s", resp.Message)
					return false
				}
				return true
			} else if resp.ReturnObj == nil {
				tolerateCount--
				if tolerateCount < 0 {
					err = common.InvalidReturnObjError
					return false
				}
				return true
			}
			detailRunningStatus := resp.ReturnObj.ProdRunningStatus
			detailOrderStatus := resp.ReturnObj.ProdOrderStatus
			// 判断是否被停用，如果被停用需要恢复使用
			if detailRunningStatus == business.PgsqlProdRunningStatusStopped && runningStatus != business.PgsqlProdRunningStatusStopped {
				err = c.startInstance(ctx, state, nil)
				if err != nil {
					return false
				}
				return true
			}
			if detailRunningStatus == runningStatus && detailOrderStatus == orderStatus {
				if syncCount <= 0 {
					return false
				} else {
					syncCount--
				}
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未启动成功！")
	}
	return
}

func (c *CtyunPostgresqlInstance) InfoLoop(ctx context.Context, state *CtyunPostgresqlInstanceConfig, plan *CtyunPostgresqlInstanceConfig, loopCount ...int) (err error) {
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
			detailParams := &pgsql.PgsqlDetailRequest{
				ProdInstId: state.ID.ValueString(),
			}
			detailHeaders := &pgsql.PgsqlDetailRequestHeader{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			// 更新成功，跳出轮询的条件：
			// 1) plan name不为空，且plan name = state name
			// 2) plan security group id 不为空， 且plan security_group_id = state security_group_id
			flagName := false
			flagSecurityGroup := false
			if plan.Name.ValueString() != "" {
				if plan.Name.ValueString() == resp.ReturnObj.ProdInstName {
					flagName = true
				}
			}
			if plan.SecurityGroupId.ValueString() != "" {
				if plan.SecurityGroupId.ValueString() == resp.ReturnObj.SecurityGroupId {
					flagSecurityGroup = true
				}
			}
			if flagName && flagSecurityGroup {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未更新成功！")
	}
	return
}

func (c *CtyunPostgresqlInstance) UpgradeLoop(ctx context.Context, state *CtyunPostgresqlInstanceConfig, plan *CtyunPostgresqlInstanceConfig, loopCount ...int) (err error) {
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
			detailParams := &pgsql.PgsqlDetailRequest{
				ProdInstId: state.ID.ValueString(),
			}
			detailHeaders := &pgsql.PgsqlDetailRequestHeader{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}

			returnObj := resp.ReturnObj
			runningStatus := returnObj.ProdRunningStatus
			orderStatus := returnObj.ProdOrderStatus
			if returnObj.DiskSize == plan.StorageSpace.ValueInt32() && returnObj.MachineSpec == plan.ProdPerformanceSpec.ValueString() && returnObj.SpuCode == fmt.Sprintf("%d", business.PgsqlProdIDDict[plan.ProdID.ValueString()]) {
				flag := false
				if plan.BackupStorageSpace.ValueInt32() == 0 {
					flag = true
				}
				if plan.BackupStorageSpace.ValueInt32() != 0 && c.ParseStorageSize(&returnObj.BackupDiskSize) == fmt.Sprintf("%d", plan.BackupStorageSpace.ValueInt32()) {
					flag = true
				}
				if runningStatus == business.MysqlRunningStatusStarted && orderStatus == business.MysqlOrderStatusStarted && flag {
					return false
				}
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未升配成功！")
	}
	return
}

func (c *CtyunPostgresqlInstance) StartedOrderLoop(ctx context.Context, state *CtyunPostgresqlInstanceConfig, orderStatus int32, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	// 设置一个容忍机制，根据pgsql定制，pgsql开通等操作可能出现报错回滚，此期间会存在查询不到实例情况
	tolerateCount := 30
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			detailParams := &pgsql.PgsqlDetailRequest{
				ProdInstId: state.ID.ValueString(),
			}
			detailHeaders := &pgsql.PgsqlDetailRequestHeader{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				tolerateCount--
				if tolerateCount < 0 {
					err = err2
					return false
				}
				return true
			} else if resp.StatusCode != 800 {
				tolerateCount--
				if tolerateCount < 0 {
					err = fmt.Errorf("API return error. Message: %s", resp.Message)
					return false
				}
				return true
			} else if resp.ReturnObj == nil {
				tolerateCount--
				if tolerateCount < 0 {
					err = common.InvalidReturnObjError
					return false
				}
				return true
			}
			detailOrderStatus := resp.ReturnObj.ProdOrderStatus
			if detailOrderStatus == orderStatus {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未启动成功！")
	}
	return
}
func (c *CtyunPostgresqlInstance) RefundLoop(ctx context.Context, config *CtyunPostgresqlInstanceConfig, loopCount ...int) (err error) {
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
			listParams := &pgsql.PgsqlListRequest{
				PageNum:      1,
				PageSize:     100,
				ProdInstName: config.Name.ValueStringPointer(),
			}
			listHeaders := &pgsql.PgsqlListRequestHeader{
				RegionID: config.RegionID.ValueString(),
			}
			if config.ProjectID.ValueString() != "" {
				listHeaders.ProjectID = config.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlListApi.Do(ctx, c.meta.Credential, listParams, listHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			}
			if len(resp.ReturnObj.List) > 1 {
				err = errors.New("实例名称重复，存在异常！")
				return false
			} else if len(resp.ReturnObj.List) == 0 {
				// 若列表为空，说明实例已删除
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未退订成功！")
	}
	return
}

func (c *CtyunPostgresqlInstance) ParseStorageSize(storageSize *string) string {
	str := *storageSize
	if str[len(str)-1] == 'G' {
		str = str[:len(str)-1]
	} else {
		str = "0"
	}
	return str
}

func (c *CtyunPostgresqlInstance) countSame(plan *CtyunPostgresqlInstanceConfig, state *CtyunPostgresqlInstanceConfig) int {
	count := 0
	if plan.BackupStorageSpace.ValueInt32() != 0 && state.BackupStorageSpace.ValueInt32() != plan.BackupStorageSpace.ValueInt32() {
		count += 1
	}
	if plan.StorageSpace.ValueInt32() != 0 && state.StorageSpace.ValueInt32() != plan.StorageSpace.ValueInt32() {
		count += 1
	}
	if plan.ProdID.ValueString() != "" && state.ProdID.ValueString() != plan.ProdID.ValueString() {
		count += 1
	}
	if plan.ProdPerformanceSpec.ValueString() != "" && state.ProdPerformanceSpec.ValueString() != plan.ProdPerformanceSpec.ValueString() {
		count += 1
	}
	return count
}

func (c *CtyunPostgresqlInstance) startInstance(ctx context.Context, state *CtyunPostgresqlInstanceConfig, plan *CtyunPostgresqlInstanceConfig) (err error) {
	startParams := &pgsql.PgsqlStartRequest{
		ProdInstId: state.ID.ValueString(),
	}
	startHeaders := &pgsql.PgsqlStartRequestHeader{
		RegionID: state.RegionID.ValueString(),
	}
	if state.ProjectID.ValueString() != "" {
		startHeaders.ProjectID = state.ProjectID.ValueString()
	}
	resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlStartApi.Do(ctx, c.meta.Credential, startParams, startHeaders)
	if err2 != nil {
		return err2
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 轮询确认是否启动完成
	err = c.RunningStatusLoop(ctx, state, business.PgsqlProdRunningStatusStarted, business.PgsqlProdOrderStatusRunning, 60)
	if err != nil {
		return
	}
	return
}

type CtyunPostgresqlInstanceConfig struct {
	CycleType            types.String `tfsdk:"cycle_type"`             // 计费模式： 1是包周期，2是按需
	RegionID             types.String `tfsdk:"region_id"`              // 目标资源池Id
	HostType             types.String `tfsdk:"host_type"`              //主机类型 host type: S6 or S7
	ProdID               types.String `tfsdk:"prod_id"`                // 产品id
	BackupStorageType    types.String `tfsdk:"backup_storage_type"`    // 备份存储类型, SSD=超高IO、SATA=普通IO、SAS=高IO、SSD-genric=通用型SSD、FAST-SSD=极速型SSD
	VpcID                types.String `tfsdk:"vpc_id"`                 // 虚拟私有云Id，，回收站恢复到新实例场景非必传则取原实例配置
	SubnetId             types.String `tfsdk:"subnet_id"`              // 子网Id，，回收站恢复到新实例场景非必传则取原实例配置
	SecurityGroupId      types.String `tfsdk:"security_group_id"`      // 安全组，回收站恢复到新实例场景非必传则取原实例配置
	AppointVip           types.String `tfsdk:"appoint_vip"`            // 指定vip
	Name                 types.String `tfsdk:"name"`                   // 集群名称(若开通只读实例，默认在主实例名称后面加"-read")
	Password             types.String `tfsdk:"password"`               // 管理员密码（RSA公钥加密）
	CycleCount           types.Int32  `tfsdk:"cycle_count"`            // 购买时长：单位月（范围：1-36）
	AutoRenew            types.Bool   `tfsdk:"auto_renew"`             // 自动续订状态 （0-不自动续订,1-自动续订）
	CaseSensitive        types.Bool   `tfsdk:"case_sensitive"`         // 是否区分大小写 0 区分 1 不区分 2待定
	OsType               types.String `tfsdk:"os_type"`                // 操作系统类型，默认2，0=裸机，1=windows，2=centos，3=ubuntu，4=android，5=redhat，6=kylin，7=uos，8=suse，9=asianux，10=open_euler，11=ctyunos，12=euler
	CpuType              types.String `tfsdk:"cpu_type"`               // cpu类型，默认30，10=鲲鹏，20=海光，30=intel，40=amd，50=飞腾，60=龙芯，70=兆芯
	ProjectID            types.String `tfsdk:"project_id"`             // 企业项目ID，默认0
	IsMGR                types.Bool   `tfsdk:"is_mgr"`                 // 是否开启MRG，默认false
	InstanceSeries       types.String `tfsdk:"instance_series"`        // 实例规格（默认：通用型=1） InstSpec
	StorageType          types.String `tfsdk:"storage_type"`           // 存储类型: SSD=超高IO、SATA=普通IO、SAS=高IO、SSD-genric=通用型SSD、FAST-SSD=极速型SSD
	StorageSpace         types.Int32  `tfsdk:"storage_space"`          // 存储空间(单位:G，范围100,32768)
	ProdPerformanceSpec  types.String `tfsdk:"prod_performance_spec"`  // 规格(例: 4C8G)
	AvailabilityZoneInfo types.List   `tfsdk:"availability_zone_info"` // 可用区信息
	BackupStorageSpace   types.Int32  `tfsdk:"backup_storage_space"`   // 备份存储空间大小
	AutoScale            types.Bool   `tfsdk:"auto_scale"`             // 0 不自动扩容 1 自动扩存储
	MaxScale             types.Int32  `tfsdk:"max_scale"`              // 存储扩容上限，单位G
	ActiveScaleRate      types.Int32  `tfsdk:"active_scale_rate"`      // 触发扩容百分比，取值范围1-100
	ID                   types.String `tfsdk:"id"`                     // 实例ID
	Alive                types.Int32  `tfsdk:"alive"`                  // 实例是否存活,0:存活，-1:异常
	DiskRated            types.Int32  `tfsdk:"disk_rated"`             // 磁盘使用率
	OuterProdInstId      types.String `tfsdk:"outer_prod_inst_id"`     // 对外的实例ID，对应PaaS平台
	ProdDbEngine         types.String `tfsdk:"prod_db_engine"`         // 数据库实例引擎
	ProdOrderStatus      types.Int32  `tfsdk:"prod_order_status"`      // 订单状态，0：正常，1：冻结，2：删除，3：操作中，4：失败,2005:扩容中
	ProdRunninStatus     types.Int32  `tfsdk:"prod_running_status"`    // 实例状态
	ProdType             types.Int32  `tfsdk:"prod_type"`              // 实例部署方式 0：单机部署,1：主备部署
	ReadPort             types.Int32  `tfsdk:"read_port"`              // 读端口
	WritePort            types.String `tfsdk:"write_port"`             // 写端口
	ToolType             types.Int32  `tfsdk:"tool_type"`              // 备份工具类型，1：pg_baseback，2：pgbackrest，3：s3
	RunningControl       types.String `tfsdk:"running_control"`        //
	MasterOrderID        types.String `tfsdk:"master_order_id"`        // 订单id
}

type AvailabilityZoneModel struct {
	AvailabilityZoneName  types.String `tfsdk:"availability_zone_name"`  // 资源池可用区名称
	AvailabilityZoneCount types.Int32  `tfsdk:"availability_zone_count"` // 资源池可用区总数
	NodeType              types.String `tfsdk:"node_type"`               // 分布AZ的节点类型
}
