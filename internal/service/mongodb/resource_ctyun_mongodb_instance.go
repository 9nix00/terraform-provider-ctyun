package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
	"terraform-provider-ctyun/internal/utils"
	"time"
)

var (
	_ resource.Resource                = &CtyunMongodbInstance{}
	_ resource.ResourceWithConfigure   = &CtyunMongodbInstance{}
	_ resource.ResourceWithImportState = &CtyunMongodbInstance{}
)

type CtyunMongodbInstance struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMongodbInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunMongodbInstance) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func NewCtyunMongodbInstance() resource.Resource {
	return &CtyunMongodbInstance{}
}

func (c *CtyunMongodbInstance) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mongodb_instance"
}

func (c *CtyunMongodbInstance) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "mongodb provider",
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
				Description: "区域id,如果不填这默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"prod_version": schema.StringAttribute{
				Optional:    true,
				Description: "版本",
			},
			"prod_spec_name": schema.StringAttribute{
				Optional:    true,
				Description: "产品名称规格名称",
			},
			"availability_zone": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "可用区名称",
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "虚拟私有云Id",
			},
			"host_type": schema.StringAttribute{
				Required:    true,
				Description: "主机类型 host type: S6 or S7",
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网Id",
			},
			"security_group_id": schema.StringAttribute{
				Required:    true,
				Description: "安全组Id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "集群名称(若开通只读实例，默认在主实例名称后面加-read)",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Description: "管理员密码（RSA公钥加密）",
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
			"purchase_count": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(1),
				Description: "购买数量(范围:1-50)",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
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
			"prod_id": schema.Int64Attribute{
				Required:    true,
				Description: "产品id",
				Validators: []validator.Int64{
					int64validator.OneOf(business.MongodbProdID...),
				},
			},
			"prod_performance_specs": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "该产品下面的单节点规格",
			},
			"host_ip": schema.StringAttribute{
				Computed:    true,
				Description: "主机ip",
			},
			"prod_performance_spec": schema.StringAttribute{
				Computed:    true,
				Description: "mongodb实例主机配置",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "项目id",
			},
			"new_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订单id",
			},
			"read_port": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "读端口,创建阶段不可填写，更新阶段可填",
			},
			"innodb_buffer_pool_size": schema.StringAttribute{
				Computed:    true,
				Description: "缓存池大小",
			},
			"innodb_thread_concurrency": schema.Int64Attribute{
				Computed:    true,
				Description: "线程数",
			},
			"prod_running_status": schema.Int32Attribute{
				Computed:    true,
				Description: "实例运行状态: 0->运行正常, 1->重启中, 2-备份操作中,3->恢复操作中,4->转换ssl,5->异常,6->修改参数组中,7->已冻结,8->已注销,9->施工中,10->施工失败,11->扩容中,12->主备切换中",
				Validators: []validator.Int32{
					int32validator.OneOf(business.MongodbRunningStatus...),
				},
			},
			"eip_id": schema.StringAttribute{
				Computed:    true,
				Description: "eip Id",
			},
			"allow_be_master": schema.BoolAttribute{
				Computed:    true,
				Description: "允许切换成为备用节点",
			},
			"is_upgrade_back_up": schema.BoolAttribute{
				Computed:    true,
				Description: "磁盘扩容时候会使用,是否主磁盘与备磁盘一起扩容",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "mongodb实例id",
			},
			"node_info_list": schema.ListNestedAttribute{
				Required:    true,
				Description: "DDS节点",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"node_type": schema.StringAttribute{
							Required:    true,
							Description: "节点类型 ：mongos=mongos节点；shard=分片节点；config=config节点；readonly=只读节点；ms=副本集；s=单机版；backup=备份机",
						},
						"inst_spec": schema.StringAttribute{
							Required:    true,
							Description: "实例类型，1=通用型，2=计算增强型，3=内存优化型，4=直通（未用到）",
						},
						"storage_type": schema.StringAttribute{
							Required:    true,
							Description: "存储类型: SSD=超高IO, SAS=高IO, SATA=普通IO，SSD-genric=通用型SSD",
							Validators: []validator.String{
								stringvalidator.OneOf(business.MongodbStorageType...),
							},
						},
						"storage_space": schema.Int32Attribute{
							Required:    true,
							Description: "存储空间(单位:G) 单机版和副本集必传：范围100-32768 、集群版shard和bckup节点必传：单个shard:范围100-2024，backup为单个shard的容量乘以shard的个数（注意：每一个shard对应3个availabilityZoneCount，参考下面字段的描述或者请求样例）",
							Validators: []validator.Int32{
								int32validator.Between(100, 32768),
							},
						},
						"prod_performance_spec": schema.StringAttribute{
							Optional:    true,
							Description: "规格: 4C8G 当 nodeType为backup类型 可不传",
						},
						"disks": schema.Int32Attribute{
							Required:    true,
							Description: "磁盘（默认为1）,2为Hbase，暂不支持",
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
										Description: "资源池可用区总数（开通集群版--nodeType为mongos时范围为[2,16]，nodeType为shard时,shard数量取值范围[2,16]，每一个shard对应3个availabilityZoneCount, 例：nodeType: shard且要开通shard数 量为3时，availabilityZoneCount:9 ；nodeType为config时节点默认为3即availabilityZoneCount: 3）",
									},
									"node_type": schema.StringAttribute{
										Required:    true,
										Description: "master:主节点、mongos:mongos节点、shard:shard节点 、config:config节点（存储类型storageType与shard节点一致，存储空间storageSpace为单个shard的storageSpace）、 backup:备份机(存储类型storageType与shard 节点一致，存储空间storageSpace为shard节点数量乘以单个shard的storageSpace)",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *CtyunMongodbInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMongodbInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 开始创建
	err = c.CreateMongodbInstance(ctx, &plan)
	if err != nil {
		return
	}
	// 创建完成后，同步云端信息
	err = c.getAndMergeMongodbInstance(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMongodbInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMongodbInstanceConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergeMongodbInstance(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
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

func (c *CtyunMongodbInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunMongodbInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunMongodbInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.updateMongodbInstance(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeMongodbInstance(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMongodbInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunMongodbInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	deleteParams := &mongodb.MongodbRefundRequest{
		InstId: state.ID.ValueString(),
	}
	deleteHeader := &mongodb.MongodbRefundRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		deleteHeader.ProjectID = state.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkMongodbApis.MongodbRefundApi.Do(ctx, c.meta.Credential, deleteParams, deleteHeader)
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

func (c *CtyunMongodbInstance) CreateMongodbInstance(ctx context.Context, config *CtyunMongodbInstanceConfig) (err error) {
	cycleType := config.CycleType.ValueString()
	params := &mongodb.MongodbCreateRequest{
		BillMode:          business.MysqlBillMode[cycleType],
		RegionId:          config.RegionID.ValueString(),
		VpcId:             config.VpcID.ValueString(),
		HostType:          config.HostType.ValueString(),
		SubnetId:          config.SubnetID.ValueString(),
		SecurityGroupId:   config.SecurityGroupID.ValueString(),
		Name:              config.Name.ValueString(),
		Password:          config.Password.ValueString(),
		Period:            config.CycleCount.ValueInt32(),
		Count:             config.PurchaseCount.ValueInt32(),
		ProdId:            config.ProdID.ValueInt64(),
		MysqlNodeInfoList: nil,
	}

	if cycleType == business.OnDemandCycleType {
		params.AutoRenewStatus = 0
	} else {
		params.AutoRenewStatus = map[bool]int32{true: 1, false: 0}[config.AutoRenew.ValueBool()]
	}
	if config.ProdVersion.ValueString() != "" {
		params.ProdVersion = config.ProdVersion.ValueStringPointer()
	}
	if config.ProdSpecName.ValueString() != "" {
		params.ProdSpecName = config.ProdSpecName.ValueStringPointer()
	}
	if !config.AvailabilityZone.IsNull() {
		var azZone []string
		diag := config.AvailabilityZone.ElementsAs(ctx, &azZone, true)
		if diag.HasError() {
			return
		}
		params.AvailabilityZone = azZone
	}
	if !config.ProdPerformanceSpecs.IsNull() {
		var specs []string
		diag := config.ProdPerformanceSpecs.ElementsAs(ctx, &specs, true)
		if diag.HasError() {
			return
		}
		params.ProdPerformanceSpecs = specs
	}
	// 处理nodeInfoList
	var nodeInfoList []NodeInfoListModel
	var mongodbNodeInfoListRequest []mongodb.MongodbNodeInfoListRequest
	diag := config.NodeInfoList.ElementsAs(ctx, &nodeInfoList, true)
	if diag.HasError() {
		return
	}
	for _, nodeInfoItem := range nodeInfoList {
		nodeInfo := mongodb.MongodbNodeInfoListRequest{
			NodeType:            nodeInfoItem.NodeType.ValueString(),
			InstSpec:            nodeInfoItem.InstSpec.ValueString(),
			StorageType:         nodeInfoItem.StorageType.ValueString(),
			StorageSpace:        nodeInfoItem.StorageSpace.ValueInt32(),
			ProdPerformanceSpec: nodeInfoItem.ProdPerformanceSpec.ValueString(),
			Disks:               nodeInfoItem.Disks.ValueInt32(),
		}
		// 处理AvailabilityZoneInfo
		var azZoneInfoList []AvailabilityZoneModel
		var azZoneInfo []mongodb.AvailabilityZoneInfoRequest

		diag = nodeInfoItem.AvailabilityZoneInfo.ElementsAs(ctx, &azZoneInfoList, true)
		if diag.HasError() {
			return
		}
		for _, azZoneInfoItem := range azZoneInfoList {
			azZone := mongodb.AvailabilityZoneInfoRequest{
				AvailabilityZoneName:  azZoneInfoItem.AvailabilityZoneName.ValueString(),
				AvailabilityZoneCount: azZoneInfoItem.AvailabilityZoneCount.ValueInt32(),
				NodeType:              azZoneInfoItem.NodeType.ValueString(),
			}
			azZoneInfo = append(azZoneInfo, azZone)
		}
		nodeInfo.AvailabilityZoneInfo = azZoneInfo
		mongodbNodeInfoListRequest = append(mongodbNodeInfoListRequest, nodeInfo)
	}
	params.MysqlNodeInfoList = mongodbNodeInfoListRequest

	header := &mongodb.MongodbCreateRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkMongodbApis.MongodbCreateApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
	if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 保存newOrderId
	config.NewOrderID = utils.SecStringValue(resp.ReturnObj.Data.NewOrderId)
	return
}

func (c *CtyunMongodbInstance) getAndMergeMongodbInstance(ctx context.Context, config *CtyunMongodbInstanceConfig) (err error) {

	listParams := &mongodb.MongodbGetListRequest{
		PageNow:      1,
		PageSize:     100,
		ProdInstName: config.Name.ValueStringPointer(),
	}
	listHeader := &mongodb.MongodbGetListHeaders{
		RegionID: config.RegionID.ValueString(),
	}
	if config.ProjectID.ValueString() != "" {
		listHeader.ProjectID = config.ProjectID.ValueStringPointer()
	}
	// 若实例id为空，实例刚刚创建，还未查询到id，需要轮询列表获取实例信息
	if config.ID.ValueString() == "" {
		resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, listParams, listHeader)
		if err2 != nil {
			err = err2
			return
		}
		if len(resp.ReturnObj.List) > 1 {
			err = errors.New("实例名重复！")
			return
		} else if len(resp.ReturnObj.List) < 1 {
			//若根据name查询不到机器，可能存在还未创建好的情况，需要轮询
			resp, err = c.ListLoop(ctx, listParams, listHeader, 60)
			if err != nil {
				return
			}
			if len(resp.ReturnObj.List) != 1 {
				err = errors.New("未查询该实例mysql，mysql name:" + config.Name.ValueString())
				return
			}
			// 查询到实例后，保存id
			if resp.ReturnObj.List[0].ProdInstId == "" {
				err = errors.New("实例创建后，实例id仍为空")
				return
			}
			config.ID = types.StringValue(resp.ReturnObj.List[0].ProdInstId)
		}
	}
	// 确认实例id不为空后,分两步查询实例详情
	// 1）轮询实例状态，确认已经正常运行，并获取实例部门详情：读端口、缓冲池信息和安全组信息
	listResp, err := c.RunningLoop(ctx, listParams, listHeader, 60)
	if err != nil {
		return err
	} else if listResp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s", *listResp.Message)
		return
	} else if listResp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	listReturnObj := listResp.ReturnObj.List[0]
	// 2）查询实例详情，获取allowBeMaster信息和eip id信息
	detailParams := &mongodb.MongodbQueryDetailRequest{
		ProdInstId: config.ID.ValueString(),
	}
	detailHeader := &mongodb.MongodbQueryDetailRequestHeaders{
		RegionID: config.RegionID.ValueString(),
	}
	if config.ProjectID.ValueString() != "" {
		detailHeader.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkMongodbApis.MongodbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeader)
	if err != nil {
		return err
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	detailReturnObj := resp.ReturnObj
	config.ReadPort = types.StringValue(detailReturnObj.Port)
	config.InnodbBufferPoolSize = types.StringValue(listReturnObj.InnodbBufferPoolSize)
	config.InnodbThreadConcurrency = types.Int64Value(listReturnObj.InnodbThreadConcurrency)
	config.ProdRunningStatus = types.Int32Value(listReturnObj.ProdRunningStatus)
	config.EipID = types.StringValue(detailReturnObj.NodeInfoVOS[0].OuterElasticIpId)
	config.AllowBeMaster = types.BoolValue(detailReturnObj.NodeInfoVOS[0].AllowBeMaster)
	config.Name = types.StringValue(listReturnObj.ProdInstName)
	config.SecurityGroupID = types.StringValue(listReturnObj.SecurityGroupId)
	prodID, err := strconv.ParseInt(listReturnObj.ProdId, 10, 64)
	if err != nil {
		return
	}
	config.ProdID = types.Int64Value(prodID)
	config.HostIp = types.StringValue(detailReturnObj.Host)
	config.ProdPerformanceSpec = types.StringValue(listReturnObj.MachineSpec)
	return
}

func (c *CtyunMongodbInstance) updateMongodbInstance(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig) (err error) {

	// 修改实例名称
	if plan.Name.ValueString() != "" && state.Name.ValueString() != plan.Name.ValueString() {
		// 修改实例前，确定实例状态为running
		err = c.PreCheckUpdateLoop(ctx, state)
		if err != nil {
			return
		}
		updateNameParams := &mongodb.MongodbUpdateInstanceNameRequest{
			ProdInstId:   state.ID.ValueString(),
			ProdInstName: plan.Name.ValueString(),
		}
		updateNameHeader := &mongodb.MongodbUpdateInstanceNameRequestHeader{
			RegionID: state.RegionID.ValueString(),
		}
		resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbUpdateInstanceNameApi.Do(ctx, c.meta.Credential, updateNameParams, updateNameHeader)
		if err2 != nil {
			err = err2
			return
		} else if resp.StatusCode != 800 {
			err = fmt.Errorf("API return error. Message: %s", *resp.Message)
			return
		}
	}

	// 修改安全组
	if plan.SecurityGroupID.ValueString() != "" && state.SecurityGroupID.ValueString() != plan.SecurityGroupID.ValueString() {
		// 修改实例前，确定实例状态为running
		err = c.PreCheckUpdateLoop(ctx, state)
		if err != nil {
			return
		}
		updateSecurityGroupParams := &mongodb.MongodbUpdateSecurityGroupRequest{
			SecurityGroupId:    state.SecurityGroupID.ValueString(),
			InstanceId:         state.ID.ValueString(),
			NewSecurityGroupId: plan.SecurityGroupID.ValueString(),
		}
		updateSecurityGroupHeader := &mongodb.MongodbUpdateSecurityGroupRequestHeader{}
		if state.ProjectID.ValueString() != "" {
			updateSecurityGroupHeader.ProjectID = state.ProjectID.ValueStringPointer()
		}
		resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbUpdateSecurityGroupApi.Do(ctx, c.meta.Credential, updateSecurityGroupParams, updateSecurityGroupHeader)
		if err2 != nil {
			err = err2
			return
		} else if resp.StatusCode != 200 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
	}
	// 变更完name和security group，确认是否变更完成
	err = c.PostCheckUpdateNameAndSecurityGroupLoop(ctx, state, plan, 60)
	if err != nil {
		return
	}
	// 修改实例端口
	if plan.ReadPort.ValueString() != "" && state.ReadPort.ValueString() != plan.ReadPort.ValueString() {
		// 修改实例前，确定实例状态为running
		err = c.PreCheckUpdateLoop(ctx, state, 60)
		if err != nil {
			return
		}
		updateParams := &mongodb.MongodbUpdatePortRequest{
			ProdInstId: state.ID.ValueString(),
			NewPort:    plan.ReadPort.ValueString(),
		}
		updateHeader := &mongodb.MongodbUpdatePortRequestHeader{
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			updateHeader.ProjectID = state.ProjectID.ValueStringPointer()
		}
		resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbUpdatePortApi.Do(ctx, c.meta.Credential, updateParams, updateHeader)
		if err2 != nil {
			err = err2
			return
		} else if resp.StatusCode != 800 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
		// 轮询确认端口更新完成
		err = c.UpdatePortLoop(ctx, state, plan, 60)
		if err != nil {
			return
		}
	}

	// 实例扩容
	// 解析出nodeInfoList
	if !plan.NodeInfoList.IsNull() {
		updateParams := &mongodb.MongodbUpgradeRequest{
			InstId: state.ID.ValueString(),
		}
		updateHeader := &mongodb.MongodbUpgradeRequestHeader{}
		if state.ProjectID.ValueString() != "" {
			updateHeader.ProjectID = state.ProjectID.ValueStringPointer()
		}
		var planNodeInfoList []NodeInfoListModel
		diag := plan.NodeInfoList.ElementsAs(ctx, &planNodeInfoList, true)
		if diag.HasError() {
			return
		}
		// 在更新阶段，nodeInfoList默认长度都为1
		if len(planNodeInfoList) != 1 {
			err = errors.New("在更新阶段，nodeInfoList输入有误！")
		}
		// 处理磁盘扩容
		planNodeInfo := planNodeInfoList[0]
		// 若plan.storageSpace不为0，触发扩容操作
		// 构建磁盘扩容请求接口
		if planNodeInfo.StorageSpace.ValueInt32() != 0 {
			updateParams.DiskVolume = planNodeInfo.StorageSpace.ValueInt32Pointer()
			updateParams.IsUpgradeBackup = plan.IsUpgradeBackUp.ValueBool()
			updateParams.NodeType = planNodeInfo.NodeType.ValueStringPointer()
		}

		// 规格升级
		if planNodeInfo.ProdPerformanceSpec.ValueString() != "" {
			updateParams.ProdPerformanceSpec = planNodeInfo.ProdPerformanceSpec.ValueStringPointer()
			updateParams.NodeType = planNodeInfo.NodeType.ValueStringPointer()
			if planNodeInfo.AvailabilityZoneInfo.IsNull() {
				err = errors.New("规格升级，azInfo不得为空！")
				return
			}
		}
		// 类型升级（例：DDS三副本集扩容到7副本）
		if plan.ProdID.ValueInt64() != 0 && state.ProdID.ValueInt64() == plan.ProdID.ValueInt64() {
			updateParams.ProdId = plan.ProdID.ValueInt64Pointer()
			updateParams.NodeType = planNodeInfo.NodeType.ValueStringPointer()
			if planNodeInfo.AvailabilityZoneInfo.IsNull() {
				err = errors.New("规格升级，azInfo不得为空！")
				return
			}
		}
		// 处理azInfo,如果azInfo不为空，三种情况：1）规格升级； 2）节点增加；3）类型升级
		if !planNodeInfo.AvailabilityZoneInfo.IsNull() {
			var availabilityZone []AvailabilityZoneModel
			var azList []mongodb.AvailabilityZoneInfo
			diag = planNodeInfo.AvailabilityZoneInfo.ElementsAs(ctx, &availabilityZone, true)
			if diag.HasError() {
				return
			}
			for _, azItem := range availabilityZone {
				var az mongodb.AvailabilityZoneInfo
				az.AvailabilityZoneName = azItem.AvailabilityZoneName.ValueString()
				az.AvailabilityZoneCount = azItem.AvailabilityZoneCount.ValueInt32()
				if azItem.NodeType.ValueString() != "" {
					az.NodeType = azItem.NodeType.ValueStringPointer()
				}
				azList = append(azList, az)
			}
			updateParams.AzList = azList
		}
		if updateParams.DiskVolume != nil || updateParams.ProdPerformanceSpec != nil || updateParams.AzList != nil {
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbUpgradeApi.Do(ctx, c.meta.Credential, updateParams, updateHeader)
			if err2 != nil {
				err = err2
				return
			} else if resp.StatusCode != 200 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return
			}
			// 轮询确认是否已扩容完成
			err = c.UpgradeLoop(ctx, state, plan, planNodeInfoList, 60)
			if err != nil {
				return
			}
		}
	}
	return
}

func (c *CtyunMongodbInstance) PreCheckUpdateLoop(ctx context.Context, state *CtyunMongodbInstanceConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	listParams := &mongodb.MongodbGetListRequest{
		PageNow:      1,
		PageSize:     100,
		ProdInstName: state.Name.ValueStringPointer(),
	}
	listHeader := &mongodb.MongodbGetListHeaders{
		RegionID: state.RegionID.ValueString(),
	}
	if state.ProjectID.ValueString() != "" {
		listHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}

	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, listParams, listHeader)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			if len(resp.ReturnObj.List) != 1 {
				err = errors.New("实例name不唯一，有误！")
				return false
			}
			if resp.ReturnObj.List[0].ProdRunningStatus == business.MongodbRunningStatusStarted {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，实例仍未运行成功！")
	}
	return nil
}

func (c *CtyunMongodbInstance) ListLoop(ctx context.Context, params *mongodb.MongodbGetListRequest, header *mongodb.MongodbGetListHeaders, loopCount ...int) (*mongodb.MongodbGetListResponse, error) {
	var err error
	var response *mongodb.MongodbGetListResponse
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
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, params, header)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
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
		return nil, errors.New("轮询已达最大次数，实例仍未创建或查询到！")
	}
	return response, nil
}

func (c *CtyunMongodbInstance) RunningLoop(ctx context.Context, params *mongodb.MongodbGetListRequest, header *mongodb.MongodbGetListHeaders, loopCount ...int) (*mongodb.MongodbGetListResponse, error) {
	var err error
	var response *mongodb.MongodbGetListResponse
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
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, params, header)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			runningStatus := resp.ReturnObj.List[0].ProdRunningStatus
			// 若实例状态已经运行正常，跳出轮询
			if runningStatus == business.MongodbRunningStatusStarted {
				response = resp
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return nil, errors.New("轮询已达最大次数，实例仍未启动！")
	}
	return response, nil
}

func (c *CtyunMongodbInstance) DeleteLoop(ctx context.Context, config *CtyunMongodbInstanceConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	listParams := &mongodb.MongodbGetListRequest{
		PageNow:      1,
		PageSize:     100,
		ProdInstName: config.Name.ValueStringPointer(),
	}
	listHeader := &mongodb.MongodbGetListHeaders{
		RegionID: config.RegionID.ValueString(),
	}
	if config.ProjectID.ValueString() != "" {
		listHeader.ProjectID = config.ProjectID.ValueStringPointer()
	}

	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, listParams, listHeader)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			}
			if len(resp.ReturnObj.List) == 0 {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，实例仍未删除成功！")
	}
	return
}

func (c *CtyunMongodbInstance) PostCheckUpdateNameAndSecurityGroupLoop(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	listParams := &mongodb.MongodbGetListRequest{
		PageNow:      1,
		PageSize:     100,
		ProdInstName: state.Name.ValueStringPointer(),
	}
	listHeader := &mongodb.MongodbGetListHeaders{
		RegionID: state.RegionID.ValueString(),
	}
	if state.ProjectID.ValueString() != "" {
		listHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}

	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, listParams, listHeader)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			}
			if len(resp.ReturnObj.List) != 1 {
				err = errors.New("根据name查询实例数量有误！")
				return false
			}
			updatedName := resp.ReturnObj.List[0].ProdInstName
			updatedSecurityGroupID := resp.ReturnObj.List[0].SecurityGroupId
			flagName := true
			flagSecurityGroup := true
			if plan.Name.ValueString() != "" {
				if updatedName != plan.Name.ValueString() {
					flagName = false
				}
			}
			if plan.SecurityGroupID.ValueString() != "" {
				if updatedSecurityGroupID != plan.SecurityGroupID.ValueString() {
					flagSecurityGroup = false
				}
			}
			if flagName && flagSecurityGroup {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，实例仍未删除成功！")
	}
	return
}

func (c *CtyunMongodbInstance) UpdatePortLoop(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			detailParams := &mongodb.MongodbQueryDetailRequest{
				ProdInstId: state.ID.ValueString(),
			}
			detailHeader := &mongodb.MongodbQueryDetailRequestHeaders{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeader.ProjectID = state.ProjectID.ValueStringPointer()
			}

			detailResp, err2 := c.meta.Apis.SdkMongodbApis.MongodbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeader)
			if err2 != nil {
				err2 = err
				return false
			} else if detailResp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *detailResp.Message)
				return false
			} else if detailResp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			// 若云端port信息与预期相符，退出轮询
			if detailResp.ReturnObj.Port == plan.ReadPort.ValueString() {
				return false
			}
			// 继续轮询
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，实例端口仍未更新成功！")
	}
	return
}

func (c *CtyunMongodbInstance) UpgradeLoop(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig, planNodeInfoList []NodeInfoListModel, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			detailParams := &mongodb.MongodbQueryDetailRequest{
				ProdInstId: state.ID.ValueString(),
			}
			detailHeader := &mongodb.MongodbQueryDetailRequestHeaders{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeader.ProjectID = state.ProjectID.ValueStringPointer()
			}

			listParams := &mongodb.MongodbGetListRequest{
				PageNow:      1,
				PageSize:     100,
				ProdInstName: state.Name.ValueStringPointer(),
			}
			listHeader := &mongodb.MongodbGetListHeaders{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				listHeader.ProjectID = state.ProjectID.ValueStringPointer()
			}

			listResp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, listParams, listHeader)
			if err2 != nil {
				err = err2
				return false
			} else if listResp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *listResp.Message)
				return false
			} else if listResp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}

			detailResp, err2 := c.meta.Apis.SdkMongodbApis.MongodbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeader)
			if err2 != nil {
				err2 = err
				return false
			} else if detailResp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *detailResp.Message)
				return false
			} else if detailResp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			// 验证扩容结果（磁盘空间，备份磁盘空间，配置，shard数和prodId）
			masterDiskFlag := true
			backupDiskFlag := true
			// 验证master磁盘空间
			if planNodeInfoList[0].NodeType.ValueString() == "master" {
				diskSize := detailResp.ReturnObj.DiskSize
				if planNodeInfoList[0].StorageSpace.ValueInt32() != 0 && planNodeInfoList[0].StorageSpace.ValueInt32() != diskSize {
					masterDiskFlag = false
				}
			} else if planNodeInfoList[0].NodeType.ValueString() == "backup" {
				diskSize := detailResp.ReturnObj.Backup.Size
				if planNodeInfoList[0].StorageSpace.ValueInt32() != 0 && fmt.Sprintf("%d", planNodeInfoList[0].StorageSpace.ValueInt32()) != diskSize {
					backupDiskFlag = false
				}
			}
			// 验证配置
			specFlag := true
			machineSpec := detailResp.ReturnObj.MachineSpec
			if planNodeInfoList[0].ProdPerformanceSpec.ValueString() != machineSpec {
				specFlag = false
			}

			// 验证prodID
			prodIDFlag := true
			prodID := listResp.ReturnObj.List[0].ProdId
			if plan.ProdID.ValueInt64() != 0 && prodID != fmt.Sprintf("%d", plan.ProdID.ValueInt64()) {
				prodIDFlag = false
			}

			if masterDiskFlag && backupDiskFlag && specFlag && prodIDFlag {
				return false
			}
			// 继续轮询
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，实例端口仍未更新成功！")
	}
	return

}

type CtyunMongodbInstanceConfig struct {
	CycleType               types.String `tfsdk:"cycle_type"`                // 计费模式： 1是包周期，2是按需
	RegionID                types.String `tfsdk:"region_id"`                 // 资源池Id
	ProdVersion             types.String `tfsdk:"prod_version"`              // 版本
	ProdSpecName            types.String `tfsdk:"prod_spec_name"`            // 产品名称规格名称
	AvailabilityZone        types.Set    `tfsdk:"availability_zone"`         // 可用区名称
	VpcID                   types.String `tfsdk:"vpc_id"`                    // 虚拟私有云Id
	HostType                types.String `tfsdk:"host_type"`                 // 主机类型 host type: S6 or S7
	SubnetID                types.String `tfsdk:"subnet_id"`                 // 子网Id
	SecurityGroupID         types.String `tfsdk:"security_group_id"`         // 安全组
	Name                    types.String `tfsdk:"name"`                      // 集群名称
	Password                types.String `tfsdk:"password"`                  // 管理员密码（RSA公钥加密）
	CycleCount              types.Int32  `tfsdk:"cycle_count"`               // 购买时长：单位月（范围：1-36）
	PurchaseCount           types.Int32  `tfsdk:"purchase_count"`            // 购买数量(范围:1-50)
	AutoRenew               types.Bool   `tfsdk:"auto_renew"`                // 自动续订状态（0-不自动续订，1-自动续订）
	ProdID                  types.Int64  `tfsdk:"prod_id"`                   // 产品id
	ProdPerformanceSpecs    types.Set    `tfsdk:"prod_performance_specs"`    // 该产品下面的单节点规格
	NodeInfoList            types.List   `tfsdk:"node_info_list"`            //
	ProjectID               types.String `tfsdk:"project_id"`                // 项目ID
	NewOrderID              types.String `tfsdk:"new_order_id"`              // 订单ID
	ID                      types.String `tfsdk:"id"`                        // 实例ID
	ReadPort                types.String `tfsdk:"read_port"`                 // 读端口
	InnodbBufferPoolSize    types.String `tfsdk:"innodb_buffer_pool_size"`   // 缓存池大小
	InnodbThreadConcurrency types.Int64  `tfsdk:"innodb_thread_concurrency"` // 线程数
	ProdRunningStatus       types.Int32  `tfsdk:"prod_running_status"`       // 实例运行状态: 0->运行正常, 1->重启中, 2-备份操作中,3->恢复操作中,4->转换ssl,5->异常,6->修改参数组中,7->已冻结,8->已注销,9->施工中,10->施工失败,11->扩容中,12->主备切换中
	EipID                   types.String `tfsdk:"eip_id"`                    // eip id
	AllowBeMaster           types.Bool   `tfsdk:"allow_be_master"`           // 允许切换成为备用节点
	IsUpgradeBackUp         types.Bool   `tfsdk:"is_upgrade_back_up"`        // DDS模块磁盘扩容时候会使用 是否主磁盘与备磁盘一起扩容
	HostIp                  types.String `tfsdk:"host_ip"`                   // 主机ip
	ProdPerformanceSpec     types.String `tfsdk:"prod_performance_spec"`     // 主机配置
}

type NodeInfoListModel struct {
	NodeType             types.String `tfsdk:"node_type"`              // 实例类型：master 或 readNode
	InstSpec             types.String `tfsdk:"inst_spec"`              // 实例规格（实例类型，1=通用型，2=计算增强型，3=内存优化型，4=直通（未用到））
	StorageType          types.String `tfsdk:"storage_type"`           // 存储类型：SSD, SATA, SAS, SSD-genric, FAST-SSD
	StorageSpace         types.Int32  `tfsdk:"storage_space"`          // 存储空间（单位：GB，范围100到32768）
	ProdPerformanceSpec  types.String `tfsdk:"prod_performance_spec"`  // 规格（例：4C8G）
	Disks                types.Int32  `tfsdk:"disks"`                  // 磁盘（默认为1）
	AvailabilityZoneInfo types.List   `tfsdk:"availability_zone_info"` // 可用区信息
}

type AvailabilityZoneModel struct {
	AvailabilityZoneName  types.String `tfsdk:"availability_zone_name"`  // 资源池可用区名称
	AvailabilityZoneCount types.Int32  `tfsdk:"availability_zone_count"` // 资源池可用区总数
	NodeType              types.String `tfsdk:"node_type"`               // 表示分布AZ的节点类型，master/slave/readNode
}
