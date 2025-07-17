package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strconv"
	"strings"
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
				Description: "区域id,如果不填这默认使用provider ctyun总region_id 或者环境变量",
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
			"host_type": schema.StringAttribute{
				Required:    true,
				Description: "主机类型 host type: S6 or S7等。可根据data.ctyun_mongodb_specs获取",
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
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]*$"), "实例名称不符合规范"),
				},
			},
			// 实现一个validator方法
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "实例密码（8-32位由大写字母、小写字母、数字、特殊字符中的任意三种组成 特殊字符为!@#$%^&*()_+-=），RSA公钥加密存储",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 32),
					validator2.MongodbPassword(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"prod_id": schema.StringAttribute{
				Required:    true,
				Description: "产品id，开通时用于确定开通单机/集群版/副本集和版本，取值范围包括：Single34（3.4单机版）,Single40（4.0单机版）,Replica3R34（3.4副本集三副本）,Replica3R40（4.0副本集三副本）,Replica5R34（3.4副本集五副本）,Replica5R40（4.0副本集五副本）,Replica7R34（3.4副本集七副本）,Replica7R40（4.0副本集七副本）,Cluster34（3.4集群版）,Cluster40（4.0集群版）,Single42（4.2单机版）,Replica3R42（4.2副本集三副本）,Replica5R42（4.2副本集五副本）,Replica7R42（4.2副本集七副本）,Cluster42（4.2集群版）,Single50（5.0单机版）,Replica3R50（5.0副本集三副本）,Replica5R50（5.0副本集五副本）,Replica7R50（5.0副本集七副本）,Cluster50（5.0集群版）,Cluster60（6.0集群版）,Replica3R60（6.0副本集三副本）,Replica5R60（6.0副本集五副本）,Replica7R60（6.0副本集七副本）,Single60（6.0单机版）",
				Validators: []validator.String{
					stringvalidator.OneOf(business.MongodbProdIDs...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订单id",
			},
			"read_port": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "读端口,创建阶段不可填写。若需要更新读取端口时可填，取值范围：1~65535",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"host_ip": schema.StringAttribute{
				Computed:    true,
				Description: "主机ip",
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
				Description: "实例运行状态: 0->运行正常, 1->重启中, 2-备份操作中, 3->恢复操作中,4->转换ssl,5->异常,6->修改参数组中,7->已冻结,8->已注销,9->施工中,10->施工失败,11->扩容中,12->主备切换中",
			},
			"prod_running_status_desc": schema.StringAttribute{
				Computed:    true,
				Description: "实例运行状态解释字段",
			},
			"eip_id": schema.StringAttribute{
				Computed:    true,
				Description: "eip Id",
			},
			"is_upgrade_back_up": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "磁盘扩容时候会使用,是否主磁盘与备磁盘一起扩容。默认true(主备一起扩容)",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "mongodb实例id",
			},
			"instance_series": schema.StringAttribute{
				Required:    true,
				Description: "实例规格，取值范围：S(通用型)，C(计算增强型)，M(内存增强型)",
				Validators: []validator.String{
					stringvalidator.OneOf(business.MysqlInstanceSeries...),
				},
			},
			"prod_performance_spec": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "实例规格，例如：4C8G",
			},
			"storage_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("SSD"),
				Description: "存储类型，默认为SSD。取值范围：SSD=超高IO, SAS=高IO, SATA=普通IO，SSD-genric=通用型SSD",
				Validators: []validator.String{
					stringvalidator.OneOf(business.MongodbStorageType...),
				},
			},
			"storage_space": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(100),
				Description: "存储空间(单位:G)，默认为100GB。取值范围：10-6144，backup节点为单个shard的容量乘以shard的个数",
				Validators: []validator.Int32{
					int32validator.Between(10, 6144),
				},
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
			// todo 必须在集群版才可填写
			"shard_num": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "shard节点数量，mongodb为集群版需填写，默认为2，取值范围：2~32",
				Default:     int32default.StaticInt32(2),
				Validators: []validator.Int32{
					int32validator.Between(2, 32),
				},
			},
			// todo 必须在集群版才可填写
			"mongos_num": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "mongos节点数量，mongodb为集群版需填写，默认为2，取值范围：2~32",
				Default:     int32default.StaticInt32(2),
				Validators: []validator.Int32{
					int32validator.Between(2, 32),
				},
			},
			"replica_num": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "副本集数量，mongodb为副本集需填写，默认为3，取值范围：[3, 5, 7]",
				Default:     int32default.StaticInt32(3),
				Validators: []validator.Int32{
					int32validator.OneOf(3, 5, 7),
				},
			},
			"backup_storage_space": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "backup节点磁盘空间，升配时用于区分节点升配",
			},
			"upgrade_node_type": schema.StringAttribute{
				Optional:    true,
				Description: "当实例为集群版，若升配mongos、shard节点规格时可填写。取值范围：shard, mongos",
				Validators: []validator.String{
					stringvalidator.OneOf("shard", "mongos"),
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
	// 确保实例创建成功后，判断port是否需要指定
	if !plan.ReadPort.IsNull() && !plan.ReadPort.IsUnknown() {
		err = c.updateReadPort(ctx, &plan, &plan)
		if err != nil {
			return
		}
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
		response.State.RemoveResource(ctx)
		err = nil
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
	response.Diagnostics.AddWarning("删除MongoDB集群成功", "集群退订后，若立即删除子网或安全组可能会失败，需要等待底层资源释放")
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
		Period:            config.CycleCount.ValueInt32(),
		Count:             1,
		ProdId:            business.MongodbProdIDDict[config.ProdID.ValueString()],
		MysqlNodeInfoList: nil,
	}

	password := business.Encode(config.Password.ValueString())
	params.Password = password
	//config.Password = types.StringValue(password)
	if cycleType == business.OnDemandCycleType {
		params.AutoRenewStatus = 0
	} else {
		params.AutoRenewStatus = map[bool]int32{true: 1, false: 0}[config.AutoRenew.ValueBool()]
	}

	var mongodbNodeInfoListRequest []mongodb.MongodbNodeInfoListRequest
	// 获取az信息
	if strings.Contains(config.ProdID.ValueString(), "single") {
		// 处理单节点nodeInfoList
		err2 := c.getSingleNodeInfo(ctx, config, &mongodbNodeInfoListRequest)
		if err2 != nil {
			err = err2
			return
		}
	} else if strings.Contains(config.ProdID.ValueString(), "replica") {
		// 处理副本级nodeInfoList
		err2 := c.getReplicaNodeInfo(ctx, config, &mongodbNodeInfoListRequest)
		if err2 != nil {
			err = err2
			return
		}
	} else if strings.Contains(config.ProdID.ValueString(), "cluster") {
		// 处理集群版本nodeInfoList
		err2 := c.getClusterNodeInfo(ctx, config, &mongodbNodeInfoListRequest)
		if err2 != nil {
			err = err2
			return
		}
	}
	// 处理nodeInfoList
	//var nodeInfoList []NodeInfoListModel

	//err = c.processNodeInfoList(ctx, &mongodbNodeInfoListRequest, config)

	//diag := config.NodeInfoList.ElementsAs(ctx, &nodeInfoList, true)
	//if diag.HasError() {
	//	return
	//}
	//for _, nodeInfoItem := range nodeInfoList {
	//	nodeInfo := mongodb.MongodbNodeInfoListRequest{
	//		NodeType:            nodeInfoItem.NodeType.ValueString(),
	//		InstSpec:            business.MongodbInstanceSeriesDict[nodeInfoItem.InstanceSeries.ValueString()],
	//		StorageType:         nodeInfoItem.StorageType.ValueString(),
	//		StorageSpace:        nodeInfoItem.StorageSpace.ValueInt32(),
	//		ProdPerformanceSpec: nodeInfoItem.ProdPerformanceSpec.ValueString(),
	//		Disks:               1,
	//	}
	//	// 处理AvailabilityZoneInfo
	//	var azZoneInfoList []AvailabilityZoneModel
	//	var azZoneInfo []mongodb.AvailabilityZoneInfoRequest
	//
	//	diag = nodeInfoItem.AvailabilityZoneInfo.ElementsAs(ctx, &azZoneInfoList, true)
	//	if diag.HasError() {
	//		return
	//	}
	//	for _, azZoneInfoItem := range azZoneInfoList {
	//		azZone := mongodb.AvailabilityZoneInfoRequest{
	//			AvailabilityZoneName:  azZoneInfoItem.AvailabilityZoneName.ValueString(),
	//			AvailabilityZoneCount: azZoneInfoItem.AvailabilityZoneCount.ValueInt32(),
	//			NodeType:              azZoneInfoItem.NodeType.ValueString(),
	//		}
	//		azZoneInfo = append(azZoneInfo, azZone)
	//	}
	//	nodeInfo.AvailabilityZoneInfo = azZoneInfo
	//	mongodbNodeInfoListRequest = append(mongodbNodeInfoListRequest, nodeInfo)
	//}
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
	config.MasterOrderID = utils.SecStringValue(resp.ReturnObj.Data.NewOrderId)
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
			} else if resp == nil {
				err = errors.New("获取mongodb列表信息，返回Nil")
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
	} else if listResp == nil {
		err = fmt.Errorf("列表查询返回为nil")
		return
	} else if listResp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s", *listResp.Message)
		return
	} else if listResp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	listReturnObj := listResp.ReturnObj.List[0]

	if config.ID.ValueString() == "" {
		err = errors.New("查询实例详情时，实例id为空")
	}
	// 2）查询实例详情，获取allowBeMaster信息和eip id信息

	detailReturnObj, err := c.getMongoDetailInfo(ctx, config)
	if err != nil {
		return
	}

	port, err := strconv.ParseInt(detailReturnObj.Port, 10, 32)
	if err != nil {
		return
	}
	config.ReadPort = types.Int32Value(int32(port))
	config.InnodbBufferPoolSize = types.StringValue(listReturnObj.InnodbBufferPoolSize)
	config.InnodbThreadConcurrency = types.Int64Value(listReturnObj.InnodbThreadConcurrency)
	config.ProdRunningStatus = types.Int32Value(listReturnObj.ProdRunningStatus)
	config.ProdRunningStatusDesc = types.StringValue(business.MongodbStatusDescDict[listReturnObj.ProdRunningStatus])
	config.EipID = types.StringValue(detailReturnObj.NodeInfoVOS[0].OuterElasticIpId)
	config.Name = types.StringValue(listReturnObj.ProdInstName)
	config.SecurityGroupID = types.StringValue(listReturnObj.SecurityGroupId)
	prodID, err := strconv.ParseInt(listReturnObj.ProdId, 10, 64)
	if err != nil {
		return
	}
	config.ProdID = types.StringValue(business.MongodbProdIDRevDict[prodID])
	config.HostIp = types.StringValue(detailReturnObj.Host)
	config.ProdPerformanceSpec = types.StringValue(listReturnObj.MachineSpec)

	return
}

func (c *CtyunMongodbInstance) updateMongodbInstance(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig) (err error) {
	if state.ID.ValueString() == "" {
		err = errors.New("在变配实例过程中， 实例id为空")
		return
	}
	// 修改实例名称
	if plan.Name.ValueString() != "" && state.Name.ValueString() != plan.Name.ValueString() {
		// 修改实例前，确定实例状态为running
		_, err = c.PreCheckUpdateLoop(ctx, state)
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
		_, err = c.PreCheckUpdateLoop(ctx, state)
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
	// 更新name，因为read-merge阶段，需要查询列表，查询列表通过name，如果name更新了，查询不到导致异常
	state.Name = types.StringValue(plan.Name.ValueString())
	// 修改实例端口
	if plan.ReadPort.ValueInt32() != 0 && state.ReadPort.ValueInt32() != plan.ReadPort.ValueInt32() {
		// 修改实例前，确定实例状态为running
		err = c.updateReadPort(ctx, state, plan)
		if err != nil {
			return
		}
	}

	// 实例扩容
	// 扩容磁盘
	err = c.upgradeStorage(ctx, state, plan)
	if err != nil {
		return
	}

	// 扩容规格
	// 若spec规格不为空，且plan阶段spec规格和state阶段spec不相同的时候，触发变配
	if !plan.ProdPerformanceSpec.IsNull() && !plan.ProdPerformanceSpec.Equal(state.ProdPerformanceSpec) {
		err = c.upgradeMongoSpec(ctx, state, plan)
		if err != nil {
			return
		}
	}

	// 扩容节点数
	// 处理副本集数量扩容

	// 处理扩容shard扩容

	// 处理扩容mongos扩容

	// 解析出nodeInfoList
	//if !plan.NodeInfoList.IsNull() && !plan.NodeInfoList.Equal(state.NodeInfoList) {
	//	updateParams := &mongodb.MongodbUpgradeRequest{
	//		InstId: state.ID.ValueString(),
	//	}
	//	updateHeader := &mongodb.MongodbUpgradeRequestHeader{}
	//	if state.ProjectID.ValueString() != "" {
	//		updateHeader.ProjectID = state.ProjectID.ValueStringPointer()
	//	}
	//	var planNodeInfoList []NodeInfoListModel
	//	diag := plan.NodeInfoList.ElementsAs(ctx, &planNodeInfoList, true)
	//	if diag.HasError() {
	//		return
	//	}
	//	// 在更新阶段，nodeInfoList默认长度都为1
	//	if len(planNodeInfoList) != 1 {
	//		err = errors.New("在更新阶段，nodeInfoList输入有误！")
	//		return
	//	}
	//	// 处理磁盘扩容
	//	planNodeInfo := planNodeInfoList[0]
	//	// 若plan.storageSpace不为0，触发扩容操作
	//	// 构建磁盘扩容请求接口
	//	diskUpgradeFlag, err2 := c.isDiskUpgrade(ctx, planNodeInfo, state)
	//	if err2 != nil {
	//		return
	//	}
	//	if planNodeInfo.StorageSpace.ValueInt32() != 0 && diskUpgradeFlag {
	//		// 确定实例处于running状态
	//		_, err = c.PreCheckUpdateLoop(ctx, state, 60)
	//		updateParams.DiskVolume = planNodeInfo.StorageSpace.ValueInt32Pointer()
	//		updateParams.IsUpgradeBackup = plan.IsUpgradeBackUp.ValueBoolPointer()
	//		updateParams.NodeType = planNodeInfo.NodeType.ValueStringPointer()
	//
	//		// 升配磁盘
	//		resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbUpgradeApi.Do(ctx, c.meta.Credential, updateParams, updateHeader)
	//		if err2 != nil {
	//			err = err2
	//			return
	//		} else if resp.StatusCode != 200 {
	//			err = fmt.Errorf("API return error. Message: %s", resp.Message)
	//			return
	//		}
	//		// 轮询确认是否已扩容完成
	//		err = c.UpgradeStorageLoop(ctx, state, plan, planNodeInfoList[0], 60)
	//		if err != nil {
	//			return
	//		}
	//		updateParams.DiskVolume = nil
	//		updateParams.IsUpgradeBackup = nil
	//	}
	//
	//	// 规格升级
	//	isUpgrade, err2 := c.isSpecUpgrade(ctx, state, planNodeInfo.ProdPerformanceSpec.ValueString())
	//	if err2 != nil {
	//		err = err2
	//		return
	//	}
	//	if planNodeInfo.ProdPerformanceSpec.ValueString() != "" && isUpgrade {
	//		updateParams.ProdPerformanceSpec = planNodeInfo.ProdPerformanceSpec.ValueStringPointer()
	//		updateParams.NodeType = planNodeInfo.NodeType.ValueStringPointer()
	//		if planNodeInfo.AvailabilityZoneInfo.IsNull() {
	//			err = errors.New("规格升级，azInfo不得为空！")
	//			return
	//		}
	//	}
	//	// 类型升级（例：DDS三副本集扩容到7副本）
	//	if !plan.ProdID.IsNull() && !plan.ProdID.IsUnknown() && state.ProdID.ValueString() != plan.ProdID.ValueString() {
	//		prodId := business.MongodbProdIDDict[plan.ProdID.ValueString()]
	//		updateParams.ProdId = &prodId
	//		updateParams.NodeType = planNodeInfo.NodeType.ValueStringPointer()
	//		if planNodeInfo.AvailabilityZoneInfo.IsNull() {
	//			err = errors.New("规格升级，azInfo不得为空！")
	//			return
	//		}
	//	}
	//	// 处理azInfo,如果azInfo不为空，三种情况：1）规格升级； 2）节点增加；3）类型升级
	//	if !planNodeInfo.AvailabilityZoneInfo.IsNull() {
	//		var availabilityZone []AvailabilityZoneModel
	//		var azList []mongodb.AvailabilityZoneInfo
	//		diag = planNodeInfo.AvailabilityZoneInfo.ElementsAs(ctx, &availabilityZone, true)
	//		if diag.HasError() {
	//			return
	//		}
	//		for _, azItem := range availabilityZone {
	//			var az mongodb.AvailabilityZoneInfo
	//			az.AvailabilityZoneName = azItem.AvailabilityZoneName.ValueString()
	//			az.AvailabilityZoneCount = azItem.AvailabilityZoneCount.ValueInt32()
	//			if azItem.NodeType.ValueString() != "" {
	//				az.NodeType = azItem.NodeType.ValueStringPointer()
	//			}
	//			azList = append(azList, az)
	//		}
	//		updateParams.AzList = azList
	//	}
	//	if updateParams.ProdPerformanceSpec != nil || updateParams.ProdId != nil {
	//		// 修改实例前，确定实例状态为running
	//		_, err = c.PreCheckUpdateLoop(ctx, state, 60)
	//		if err != nil {
	//			return
	//		}
	//		resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbUpgradeApi.Do(ctx, c.meta.Credential, updateParams, updateHeader)
	//		if err2 != nil {
	//			err = err2
	//			return
	//		} else if resp.StatusCode != 200 {
	//			err = fmt.Errorf("API return error. Message: %s", resp.Message)
	//			return
	//		}
	//		// 轮询确认是否已扩容完成
	//		err = c.UpgradeLoop(ctx, state, plan, planNodeInfoList, 60)
	//		if err != nil {
	//			return
	//		}
	//	}
	//
	//	// 更新完成后，将plan.NodeInfoList同步给state.NodeInfoList
	//	state.NodeInfoList = plan.NodeInfoList
	//	// 将state.upGradeBackup同步给plan
	//	state.IsUpgradeBackUp = plan.IsUpgradeBackUp
	//
	//}
	return
}

func (c *CtyunMongodbInstance) isDiskUpgrade(ctx context.Context, nodeInfoList NodeInfoListModel, state *CtyunMongodbInstanceConfig) (flag bool, err error) {
	flag = false

	detailResp, err := c.getMongoDetailInfo(ctx, state)
	if err != nil {
		return
	}

	if nodeInfoList.NodeType.ValueString() == "master" {
		masterStorageSpace := detailResp.DiskSize
		if nodeInfoList.StorageSpace.ValueInt32() != masterStorageSpace {
			flag = true
			return
		}
	} else if nodeInfoList.NodeType.ValueString() == "backup" {
		backupStorageSpace := detailResp.Backup.Size[:len(detailResp.Backup.Size)-1]
		if fmt.Sprintf("%d", nodeInfoList.StorageSpace.ValueInt32()) != backupStorageSpace {
			flag = true
			return
		}

	}
	return

}

// 查询详情，确认spec是否更改
func (c *CtyunMongodbInstance) isSpecUpgrade(ctx context.Context, state *CtyunMongodbInstanceConfig, spec string) (isUpgrade bool, err error) {
	isUpgrade = true
	err = nil
	if state.ProdPerformanceSpec.ValueString() != "" {
		if state.ProdPerformanceSpec.ValueString() == spec {
			return false, nil
		} else {
			return true, nil
		}
	}

	detailReturnObj, err := c.getMongoDetailInfo(ctx, state)
	if err != nil {
		return
	}

	if detailReturnObj.MachineSpec == spec {
		isUpgrade = false
		return
	}
	return
}

func (c *CtyunMongodbInstance) PreCheckUpdateLoop(ctx context.Context, state *CtyunMongodbInstanceConfig, loopCount ...int) (ListResp *mongodb.MongodbGetListResponse, err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	syncCount := 2
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return nil, err
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
			if resp.ReturnObj.List[0].ProdRunningStatus == business.MongodbRunningStatusStarted && resp.ReturnObj.List[0].ProdOrderStatus == business.MongodbOrderStatusStarted {
				if syncCount > 0 {
					syncCount--
					return true
				}
				ListResp = resp
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return nil, errors.New("轮询已达最大次数，实例仍未运行成功！")
	}
	return
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
		ProdInstName: plan.Name.ValueStringPointer(),
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
	tolerateCount := 5
	result := retryer.Start(
		func(currentTime int) bool {
			// 查询详情
			detailResp, err2 := c.getMongoDetailInfo(ctx, state)
			if err2 != nil {
				if tolerateCount <= 0 {
					err = err2
					return false
				}
				tolerateCount--
			}
			// 若云端port信息与预期相符，退出轮询
			if detailResp.Port == fmt.Sprintf("%d", plan.ReadPort.ValueInt32()) {
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
	tolerateCount := 30

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
				if tolerateCount <= 0 {
					err2 = err
					return false
				}
				tolerateCount--
			} else if detailResp.StatusCode != 800 {
				if tolerateCount <= 0 {
					err = fmt.Errorf("API return error. Message: %s", *detailResp.Message)
					return false
				}
				tolerateCount--
			} else if detailResp.ReturnObj == nil {
				if tolerateCount <= 0 {
					err = common.InvalidReturnObjError
					return false
				}
				tolerateCount--
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
			if plan.ProdID.ValueString() != "" && prodID != fmt.Sprintf("%d", business.MongodbProdIDDict[plan.ProdID.ValueString()]) {
				prodIDFlag = false
			}
			if specFlag && prodIDFlag {
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

func (c *CtyunMongodbInstance) UpgradeStorageLoop(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig, planNodeInfoList NodeInfoListModel, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	tolerateCount := 5
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
			// 若实例在升级过程中，直接继续轮询
			if listResp.ReturnObj.List[0].ProdRunningStatus != business.MongodbRunningStatusStarted || listResp.ReturnObj.List[0].ProdOrderStatus != business.MongodbOrderStatusStarted {
				return true
			}

			detailResp, err2 := c.meta.Apis.SdkMongodbApis.MongodbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeader)
			if err2 != nil {
				if tolerateCount <= 0 {
					err2 = err
					return false
				}
				tolerateCount--
				return true
			} else if detailResp.StatusCode != 800 {
				if tolerateCount <= 0 {
					err = fmt.Errorf("API return error. Message: %s", *detailResp.Message)
					return false
				}
				tolerateCount--
				return true
			} else if detailResp.ReturnObj == nil {
				if tolerateCount <= 0 {
					err = common.InvalidReturnObjError
					return false
				}
				tolerateCount--
				return true
			}
			// 验证扩容结果（磁盘空间，备份磁盘空间，配置，shard数和prodId）
			masterDiskFlag := true
			backupDiskFlag := true
			// 验证master磁盘空间
			if planNodeInfoList.NodeType.ValueString() == "master" {
				diskSize := detailResp.ReturnObj.DiskSize
				if planNodeInfoList.StorageSpace.ValueInt32() != 0 && planNodeInfoList.StorageSpace.ValueInt32() != diskSize {
					masterDiskFlag = false
				}
			}
			if plan.IsUpgradeBackUp.ValueBool() || planNodeInfoList.NodeType.ValueString() == "backup" {
				diskSize := detailResp.ReturnObj.Backup.Size[:len(detailResp.ReturnObj.Backup.Size)-1]
				if planNodeInfoList.StorageSpace.ValueInt32() != 0 && fmt.Sprintf("%d", planNodeInfoList.StorageSpace.ValueInt32()) != diskSize {
					backupDiskFlag = false
				}
			}
			if masterDiskFlag && backupDiskFlag {
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

func (c *CtyunMongodbInstance) updateReadPort(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig) (err error) {
	listResp, err2 := c.PreCheckUpdateLoop(ctx, state, 60)
	if err2 != nil {
		err = err2
		return
	}
	updateParams := &mongodb.MongodbUpdatePortRequest{
		ProdInstId: state.ID.ValueString(),
		NewPort:    fmt.Sprintf("%d", plan.ReadPort.ValueInt32()),
	}
	updateHeader := &mongodb.MongodbUpdatePortRequestHeader{
		RegionID: state.RegionID.ValueString(),
	}
	if state.ProjectID.ValueString() != "" {
		updateHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}
	fmt.Println(listResp.ReturnObj.List[0].ProdOrderStatus, listResp.ReturnObj.List[0].ProdRunningStatus)
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
	return
}

//func (c *CtyunMongodbInstance) processNodeInfoList(ctx context.Context, nodeInfoList *[]mongodb.MongodbNodeInfoListRequest, config *CtyunMongodbInstanceConfig) (err error) {
//	// 1.判断开通产品为单机版， 副本集和集群版
//	prodType := business.MongodbProdTypeDict[config.ProdID.ValueString()]
//	if prodType == business.MongodbProdTypeSingle {
//		var singleNodeInfo mongodb.MongodbNodeInfoListRequest
//		singleNodeInfo.NodeType = "s"
//		singleNodeInfo.InstSpec = config.HostType.ValueString()[0:1]
//		singleNodeInfo.StorageType = config.SingleStorageType.ValueString()
//		singleNodeInfo.StorageSpace = config.SingleStorageSpace.ValueInt32()
//		singleNodeInfo.ProdPerformanceSpec = config.SingleProdPerformanceSpec.ValueString()
//		singleNodeInfo.Disks = 1
//		// 若az信息不为空，则处理用户填写的az信息
//		if !config.SingleAzInfo.IsNull() && !config.SingleAzInfo.IsUnknown() {
//			var azZoneInfoList []AvailabilityZoneModel
//			var singleAzInfoList []mongodb.AvailabilityZoneInfoRequest
//			diags := config.SingleAzInfo.ElementsAs(ctx, &azZoneInfoList, true)
//			if diags.HasError() {
//				err = errors.New(diags[0].Detail())
//				return
//			}
//			for _, azInfo := range azZoneInfoList {
//				var singleAzInfo mongodb.AvailabilityZoneInfoRequest
//				singleAzInfo.AvailabilityZoneName = azInfo.AvailabilityZoneName.ValueString()
//				singleAzInfo.NodeType = azInfo.NodeType.ValueString()
//				singleAzInfo.AvailabilityZoneCount = azInfo.AvailabilityZoneCount.ValueInt32()
//				singleAzInfoList = append(singleAzInfoList, singleAzInfo)
//			}
//			singleNodeInfo.AvailabilityZoneInfo = singleAzInfoList
//		} else {
//			// 若用户az信息未填写，自动分配az信息
//			singleAzInfoList, err2 := c.generateAzInfo(ctx, config, "master")
//			if err2 != nil {
//				return
//			}
//			singleNodeInfo.AvailabilityZoneInfo = singleAzInfoList
//		}
//	} else if prodType == business.MongodbProdTypeReplica {
//
//	} else if prodType == business.MongodbProdTypeCluster {
//
//	}
//}

func (c *CtyunMongodbInstance) generateAzInfo(ctx context.Context, config *CtyunMongodbInstanceConfig, prodType string, nodeType string) (AzInfoList []mongodb.AvailabilityZoneInfoRequest, err error) {
	params := &mongodb.TeledbGetAvailabilityZoneRequest{
		RegionId: config.RegionID.ValueString(),
	}
	header := &mongodb.TeledbGetAvailabilityZoneRequestHeader{}
	// 1. 获取az信息
	resp, err := c.meta.Apis.SdkMongodbApis.TeledbGetAvailabilityZone.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = errors.New("查询az信息时返回为nil，请稍后再试")
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj.Data == nil {
		err = common.InvalidReturnObjError
		return
	}
	azList := resp.ReturnObj.Data
	azNum := len(azList)

	if azNum <= 0 {
		err = errors.New("未查询到该资源池az信息，请稍后重试，或者手动填写az信息进行创建")
		return
	}
	// 因此1个az和2个az分配规则相同，只有3个及以上az的资源池有区别
	// mongodb分配节点规则：主节点与备用节点1、2需完全相同，或者完全不相同
	// 副本集和集群节点，还需要生成backup

	if prodType == "single" || prodType == "backup" {
		var azInfo mongodb.AvailabilityZoneInfoRequest
		azInfo.NodeType = nodeType
		azInfo.AvailabilityZoneName = azList[0].AvailabilityZoneName
		azInfo.AvailabilityZoneCount = 1
		AzInfoList = append(AzInfoList, azInfo)
		return
	} else if prodType == "replica" {
		if azNum >= 3 {
			nodeDist := business.MongodbReplicaNodeDistMap[config.replicaNum.ValueInt32()]
			// 有3个az，节点可以平均分摊在各个az下
			for _, azItem := range azList {
				if len(AzInfoList) >= 3 {
					break
				}
				var azInfo mongodb.AvailabilityZoneInfoRequest
				azInfo.NodeType = nodeType
				azInfo.AvailabilityZoneCount = nodeDist % 10
				nodeDist = nodeDist / 10
				azInfo.AvailabilityZoneName = azItem.AvailabilityZoneName
				AzInfoList = append(AzInfoList, azInfo)
			}
		} else {
			var azInfo mongodb.AvailabilityZoneInfoRequest
			azInfo.NodeType = nodeType
			azInfo.AvailabilityZoneName = azList[0].AvailabilityZoneName
			azInfo.AvailabilityZoneCount = config.replicaNum.ValueInt32()
			AzInfoList = append(AzInfoList, azInfo)
		}
		return
	} else if prodType == "cluster" {
		// 默认为config的数量
		var nodeNum int32
		nodeNum = 3
		if nodeType == "mongos" {
			nodeNum = business.MongodbClusterNodeBaseNumMap[nodeType] * config.MongosNum.ValueInt32()
		} else if nodeType == "shard" {
			nodeNum = business.MongodbClusterNodeBaseNumMap[nodeType] * config.ShardNum.ValueInt32()
		}
		if azNum >= 3 {
			// 先计算每个AZ的节点数量
			distNodeNum := [3]int32{
				(int32(nodeNum) + 2) / 3,
				(int32(nodeNum) + 1) / 3,
				int32(nodeNum) / 3,
			}
			for idx, azItem := range azList {
				if len(AzInfoList) >= 3 {
					break
				}
				var azInfo mongodb.AvailabilityZoneInfoRequest
				azInfo.NodeType = nodeType
				azInfo.AvailabilityZoneName = azItem.AvailabilityZoneName
				azInfo.AvailabilityZoneCount = distNodeNum[idx]
				AzInfoList = append(AzInfoList, azInfo)
			}
		} else {
			// 处理单节点情况
			var azInfo mongodb.AvailabilityZoneInfoRequest
			azInfo.NodeType = nodeType
			azInfo.AvailabilityZoneName = azList[0].AvailabilityZoneName
			azInfo.AvailabilityZoneCount = nodeNum
			AzInfoList = append(AzInfoList, azInfo)
		}
		return
	} else {
		err = errors.New("mongodb数据库类型有误")
		return
	}
}

func (c *CtyunMongodbInstance) getSingleNodeInfo(ctx context.Context, config *CtyunMongodbInstanceConfig, mongoNodeInfoList *[]mongodb.MongodbNodeInfoListRequest) (err error) {
	var mongoMasterNodeInfo mongodb.MongodbNodeInfoListRequest
	mongoMasterNodeInfo.NodeType = "s"
	mongoMasterNodeInfo.InstSpec = "1"
	mongoMasterNodeInfo.StorageType = config.StorageType.ValueString()
	mongoMasterNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	mongoMasterNodeInfo.ProdPerformanceSpec = config.ProdPerformanceSpec.ValueString()
	// 处理azInfo，若azInfo不为空，用用户输入的azInfo
	if !config.AvailabilityZoneInfo.IsNull() && !config.AvailabilityZoneInfo.IsUnknown() {
		var azZoneInfoList []AvailabilityZoneModel
		var azZoneInfo []mongodb.AvailabilityZoneInfoRequest
		diags := config.AvailabilityZoneInfo.ElementsAs(ctx, &azZoneInfoList, true)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		for _, azInfoItem := range azZoneInfoList {
			azZone := mongodb.AvailabilityZoneInfoRequest{
				AvailabilityZoneName:  azInfoItem.AvailabilityZoneName.ValueString(),
				AvailabilityZoneCount: azInfoItem.AvailabilityZoneCount.ValueInt32(),
				NodeType:              azInfoItem.NodeType.ValueString(),
			}
			azZoneInfo = append(azZoneInfo, azZone)
		}
		mongoMasterNodeInfo.AvailabilityZoneInfo = azZoneInfo
	} else {
		// 若azInfo为空，则生成az信息
		azInfo, err2 := c.generateAzInfo(ctx, config, "single", "master")
		if err2 != nil {
			err = err2
			return
		}
		mongoMasterNodeInfo.AvailabilityZoneInfo = azInfo
	}
	*mongoNodeInfoList = append(*mongoNodeInfoList, mongoMasterNodeInfo)
	return
}

func (c *CtyunMongodbInstance) getReplicaNodeInfo(ctx context.Context, config *CtyunMongodbInstanceConfig, mongoNodeInfoList *[]mongodb.MongodbNodeInfoListRequest) (err error) {
	// 副本集需要一个master节点，和一个backup节点
	// master节点
	var mongoMasterNodeInfo mongodb.MongodbNodeInfoListRequest
	mongoMasterNodeInfo.NodeType = "ms"
	mongoMasterNodeInfo.InstSpec = "1"
	mongoMasterNodeInfo.StorageType = config.StorageType.ValueString()
	mongoMasterNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	mongoMasterNodeInfo.ProdPerformanceSpec = config.ProdPerformanceSpec.ValueString()

	// backup节点
	var mongoBackupNodeInfo mongodb.MongodbNodeInfoListRequest
	mongoBackupNodeInfo.NodeType = "backup"
	mongoBackupNodeInfo.InstSpec = "1"
	mongoBackupNodeInfo.StorageType = config.StorageType.ValueString()
	mongoBackupNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()

	// 处理azInfo，若azInfo不为空，用用户输入的azInfo
	if !config.AvailabilityZoneInfo.IsNull() && !config.AvailabilityZoneInfo.IsUnknown() {
		var azZoneInfoList []AvailabilityZoneModel
		var masterAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		var backupAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		diags := config.AvailabilityZoneInfo.ElementsAs(ctx, &azZoneInfoList, true)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		for _, azInfoItem := range azZoneInfoList {
			azZone := mongodb.AvailabilityZoneInfoRequest{
				AvailabilityZoneName:  azInfoItem.AvailabilityZoneName.ValueString(),
				AvailabilityZoneCount: azInfoItem.AvailabilityZoneCount.ValueInt32(),
				NodeType:              azInfoItem.NodeType.ValueString(),
			}
			if azInfoItem.NodeType.ValueString() == "master" {
				masterAzZoneInfo = append(masterAzZoneInfo, azZone)
			} else if azInfoItem.NodeType.ValueString() == "backup" {
				backupAzZoneInfo = append(backupAzZoneInfo, azZone)
			}
		}
		mongoMasterNodeInfo.AvailabilityZoneInfo = masterAzZoneInfo
		mongoBackupNodeInfo.AvailabilityZoneInfo = backupAzZoneInfo
	} else {
		// 若azInfo为空，则生成az信息
		// master节点
		azInfo, err2 := c.generateAzInfo(ctx, config, "replica", "master")
		if err2 != nil {
			err = err2
			return
		}
		mongoMasterNodeInfo.AvailabilityZoneInfo = azInfo

		// backup节点
		azInfo, err2 = c.generateAzInfo(ctx, config, "replica", "backup")
		if err2 != nil {
			err = err2
			return
		}
		mongoBackupNodeInfo.AvailabilityZoneInfo = azInfo
	}
	*mongoNodeInfoList = append(*mongoNodeInfoList, mongoMasterNodeInfo)
	*mongoNodeInfoList = append(*mongoNodeInfoList, mongoBackupNodeInfo)
	return
}

func (c *CtyunMongodbInstance) getClusterNodeInfo(ctx context.Context, config *CtyunMongodbInstanceConfig, mongoNodeInfoList *[]mongodb.MongodbNodeInfoListRequest) (err error) {
	// 副本集需要一个mongos, shard, config节点
	// mongos节点
	var mongoMongosNodeInfo mongodb.MongodbNodeInfoListRequest
	mongoMongosNodeInfo.NodeType = "mongos"
	mongoMongosNodeInfo.InstSpec = "1"
	mongoMongosNodeInfo.StorageType = config.StorageType.ValueString()
	mongoMongosNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	mongoMongosNodeInfo.ProdPerformanceSpec = config.ProdPerformanceSpec.ValueString()
	// shard节点
	var mongoShardNodeInfo mongodb.MongodbNodeInfoListRequest
	mongoShardNodeInfo.NodeType = "shard"
	mongoShardNodeInfo.InstSpec = "1"
	mongoShardNodeInfo.StorageType = config.StorageType.ValueString()
	mongoShardNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	mongoShardNodeInfo.ProdPerformanceSpec = config.ProdPerformanceSpec.ValueString()
	// config节点, config节点配置固定
	var mongoConfigNodeInfo mongodb.MongodbNodeInfoListRequest
	mongoConfigNodeInfo.NodeType = "config"
	mongoConfigNodeInfo.InstSpec = "1"
	mongoConfigNodeInfo.StorageType = "SSD"
	mongoConfigNodeInfo.StorageSpace = 100
	mongoConfigNodeInfo.ProdPerformanceSpec = "2C4G"
	// backup节点
	var mongoBackupNodeInfo mongodb.MongodbNodeInfoListRequest
	mongoBackupNodeInfo.NodeType = "backup"
	mongoBackupNodeInfo.InstSpec = "1"
	mongoBackupNodeInfo.StorageType = config.StorageType.ValueString()
	// backup节点磁盘空间 = shard数量*每个shard磁盘空间
	mongoBackupNodeInfo.StorageSpace = config.StorageSpace.ValueInt32() * config.ShardNum.ValueInt32()

	// 处理azInfo，若azInfo不为空，用用户输入的azInfo
	if !config.AvailabilityZoneInfo.IsNull() && !config.AvailabilityZoneInfo.IsUnknown() {
		var azZoneInfoList []AvailabilityZoneModel
		var mongosAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		var shardAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		var configAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		var backupAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		diags := config.AvailabilityZoneInfo.ElementsAs(ctx, &azZoneInfoList, true)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		for _, azInfoItem := range azZoneInfoList {
			azZone := mongodb.AvailabilityZoneInfoRequest{
				AvailabilityZoneName:  azInfoItem.AvailabilityZoneName.ValueString(),
				AvailabilityZoneCount: azInfoItem.AvailabilityZoneCount.ValueInt32(),
				NodeType:              azInfoItem.NodeType.ValueString(),
			}
			if azInfoItem.NodeType.ValueString() == "mongos" {
				mongosAzZoneInfo = append(mongosAzZoneInfo, azZone)
			} else if azInfoItem.NodeType.ValueString() == "backup" {
				backupAzZoneInfo = append(backupAzZoneInfo, azZone)
			} else if azInfoItem.NodeType.ValueString() == "config" {
				configAzZoneInfo = append(configAzZoneInfo, azZone)
			} else if azInfoItem.NodeType.ValueString() == "shard" {
				shardAzZoneInfo = append(shardAzZoneInfo, azZone)
			}
		}
		mongoMongosNodeInfo.AvailabilityZoneInfo = mongosAzZoneInfo
		mongoShardNodeInfo.AvailabilityZoneInfo = shardAzZoneInfo
		mongoConfigNodeInfo.AvailabilityZoneInfo = configAzZoneInfo
		mongoBackupNodeInfo.AvailabilityZoneInfo = backupAzZoneInfo
	} else {
		// 若azInfo为空，则生成az信息
		// mongos节点
		azInfo, err2 := c.generateAzInfo(ctx, config, "replica", "mongos")
		if err2 != nil {
			err = err2
			return
		}
		mongoMongosNodeInfo.AvailabilityZoneInfo = azInfo
		// shard
		azInfo, err2 = c.generateAzInfo(ctx, config, "replica", "shard")
		if err2 != nil {
			err = err2
			return
		}
		mongoShardNodeInfo.AvailabilityZoneInfo = azInfo

		// config
		azInfo, err2 = c.generateAzInfo(ctx, config, "replica", "config")
		if err2 != nil {
			err = err2
			return
		}
		mongoConfigNodeInfo.AvailabilityZoneInfo = azInfo

		// backup
		azInfo, err2 = c.generateAzInfo(ctx, config, "replica", "backup")
		if err2 != nil {
			err = err2
			return
		}
		mongoConfigNodeInfo.AvailabilityZoneInfo = azInfo
	}

	*mongoNodeInfoList = append(*mongoNodeInfoList, mongoMongosNodeInfo)
	*mongoNodeInfoList = append(*mongoNodeInfoList, mongoShardNodeInfo)
	*mongoNodeInfoList = append(*mongoNodeInfoList, mongoConfigNodeInfo)
	*mongoNodeInfoList = append(*mongoNodeInfoList, mongoBackupNodeInfo)
	return
}

func (c *CtyunMongodbInstance) upgradeStorage(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig) (err error) {
	// 若plan阶段存储空间与state阶段不一致，触发更新
	if !plan.StorageSpace.IsNull() && !plan.StorageSpace.Equal(state.StorageSpace) {
		// 确定实例处于running状态
		_, err = c.PreCheckUpdateLoop(ctx, state, 60)
		nodeType := "master"
		updateParams := &mongodb.MongodbUpgradeRequest{
			InstId:          state.ID.ValueString(),
			DiskVolume:      plan.StorageSpace.ValueInt32Pointer(),
			IsUpgradeBackup: plan.IsUpgradeBackUp.ValueBoolPointer(),
			NodeType:        &nodeType,
		}
		updateHeader := &mongodb.MongodbUpgradeRequestHeader{}
		if state.ProjectID.ValueString() != "" {
			updateHeader.ProjectID = state.ProjectID.ValueStringPointer()
		}
		resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbUpgradeApi.Do(ctx, c.meta.Credential, updateParams, updateHeader)
		if err2 != nil {
			err = err2
			return
		} else if resp == nil {
			err = errors.New("当进行磁盘升配操作中，执行结果返回为nil。请确认未升配成功后，再重试。")
			return
		} else if resp.StatusCode != 200 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}

		var planNodeInfo NodeInfoListModel
		planNodeInfo.NodeType = types.StringValue(nodeType)
		planNodeInfo.StorageSpace = types.Int32Value(plan.StorageSpace.ValueInt32())

		// 轮询确认是否已扩容完成
		err = c.UpgradeStorageLoop(ctx, state, plan, planNodeInfo, 60)
		if err != nil {
			return
		}
	}
	// 若plan阶段存储空间与state阶段不一致，触发更新
	flag, err := c.isBackupDiskUpgrade(ctx, state, plan)
	if !plan.backupStorageSpace.IsNull() && flag {
		// 确定实例处于running状态
		_, err = c.PreCheckUpdateLoop(ctx, state, 60)
		nodeType := "backup"
		updateParams := &mongodb.MongodbUpgradeRequest{
			InstId:          state.ID.ValueString(),
			DiskVolume:      plan.StorageSpace.ValueInt32Pointer(),
			IsUpgradeBackup: plan.IsUpgradeBackUp.ValueBoolPointer(),
			NodeType:        &nodeType,
		}
		updateHeader := &mongodb.MongodbUpgradeRequestHeader{}
		if state.ProjectID.ValueString() != "" {
			updateHeader.ProjectID = state.ProjectID.ValueStringPointer()
		}
		resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbUpgradeApi.Do(ctx, c.meta.Credential, updateParams, updateHeader)
		if err2 != nil {
			err = err2
			return
		} else if resp == nil {
			err = errors.New("当进行backup磁盘升配操作中，执行结果返回为nil。请确认未升配成功后，再重试。")
			return
		} else if resp.StatusCode != 200 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}

		var planNodeInfo NodeInfoListModel
		planNodeInfo.NodeType = types.StringValue(nodeType)
		planNodeInfo.StorageSpace = types.Int32Value(plan.backupStorageSpace.ValueInt32())

		// 轮询确认是否已扩容完成
		err = c.UpgradeStorageLoop(ctx, state, plan, planNodeInfo, 60)
		if err != nil {
			return
		}
	}
	return
}

func (c *CtyunMongodbInstance) isBackupDiskUpgrade(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig) (flag bool, err error) {

	detailResp, err := c.getMongoDetailInfo(ctx, state)
	if err != nil {
		return
	}
	if !plan.backupStorageSpace.IsNull() {
		backupStorageSpace := detailResp.Backup.Size[:len(detailResp.Backup.Size)-1]
		if fmt.Sprintf("%d", plan.backupStorageSpace.ValueInt32()) != backupStorageSpace {
			flag = true
			return
		}
	}
	return
}

func (c *CtyunMongodbInstance) upgradeMongoSpec(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig) (err error) {

	//upgradeParams := &mongodb.MongodbUpgradeRequest{
	//	InstId:              state.ID.ValueString(),
	//	NodeType:            ,
	//	ProdPerformanceSpec: plan.ProdPerformanceSpec.ValueStringPointer(),
	//	IsUpgradeBackup:     nil,
	//	AzList:              nil,
	//}
	//判断azInfo是否为空，
	if plan.AvailabilityZoneInfo.IsNull() || plan.AvailabilityZoneInfo.IsUnknown() {
		// 如果az信息为空，获取控制台信息
		var azInfo *[]mongodb.AvailabilityZoneInfo
		err = c.getMongodbNodeDistInfo(ctx, state, azInfo)
		if err != nil {
			return
		}
	} else {
		// 如果不为空，利用用户的输入作为输入

	}
	return
}

func (c *CtyunMongodbInstance) getMongoDetailInfo(ctx context.Context, config *CtyunMongodbInstanceConfig) (detail *mongodb.DetailRespReturnObj, err error) {
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
		return
	} else if resp == nil {
		err = errors.New("获取mongodb实例为nil，请稍后再试！")
		return
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	detail = resp.ReturnObj
	return
}

func (c *CtyunMongodbInstance) getRegionAzInfoList(ctx context.Context, state *CtyunMongodbInstanceConfig) (azList []mongodb.TeledbGetAvailabilityZoneResponseReturnObjData, err error) {
	params := &mongodb.TeledbGetAvailabilityZoneRequest{
		RegionId: state.RegionID.ValueString(),
	}
	header := &mongodb.TeledbGetAvailabilityZoneRequestHeader{}
	// 1. 获取az信息
	resp, err := c.meta.Apis.SdkMongodbApis.TeledbGetAvailabilityZone.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = errors.New("查询az信息时返回为nil，请稍后再试")
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj.Data == nil {
		err = common.InvalidReturnObjError
		return
	}
	azList = resp.ReturnObj.Data
	return
}

func (c *CtyunMongodbInstance) getMongodbNodeDistInfo(ctx context.Context, state *CtyunMongodbInstanceConfig, azInfo *[]mongodb.AvailabilityZoneInfo) (err error) {
	// 1. 获取当前实例分布，通过查询实例详情获取，并利用map存储（详情接口只有az display name， 当时请求接口需要 az id）
	azIdMap := make(map[string]string)

	// 1.1 获取az displayName : az id的映射表
	azList, err := c.getRegionAzInfoList(ctx, state)
	if err != nil {
		return
	}
	for _, azItem := range azList {
		azIdMap[azItem.DisplayName] = azItem.AvailabilityZoneId
	}
	// 1.2 获取实例节点详情
	detail, err := c.getMongoDetailInfo(ctx, state)
	if err != nil {
		return
	}
	nodeInfos := detail.NodeInfoVOS
	for _, nodeInfo := range nodeInfos {
		azId := azIdMap[nodeInfo.AzDisplayName]
		fmt.Println(azId)
	}

	// 2.
	return
}

type CtyunMongodbInstanceConfig struct {
	CycleType               types.String `tfsdk:"cycle_type"`                // 计费模式： 1是包周期，2是按需
	RegionID                types.String `tfsdk:"region_id"`                 // 资源池Id
	VpcID                   types.String `tfsdk:"vpc_id"`                    // 虚拟私有云Id
	HostType                types.String `tfsdk:"host_type"`                 // 主机类型 host type: S6 or S7
	SubnetID                types.String `tfsdk:"subnet_id"`                 // 子网Id
	SecurityGroupID         types.String `tfsdk:"security_group_id"`         // 安全组
	Name                    types.String `tfsdk:"name"`                      // 集群名称
	Password                types.String `tfsdk:"password"`                  // 管理员密码（RSA公钥加密）
	CycleCount              types.Int32  `tfsdk:"cycle_count"`               // 购买时长：单位月（范围：1-36）
	AutoRenew               types.Bool   `tfsdk:"auto_renew"`                // 自动续订状态（0-不自动续订，1-自动续订）
	ProdID                  types.String `tfsdk:"prod_id"`                   // 产品id
	NodeInfoList            types.List   `tfsdk:"node_info_list"`            //
	ProjectID               types.String `tfsdk:"project_id"`                // 项目ID
	MasterOrderID           types.String `tfsdk:"master_order_id"`           // 订单ID
	ID                      types.String `tfsdk:"id"`                        // 实例ID
	ReadPort                types.Int32  `tfsdk:"read_port"`                 // 读端口
	InnodbBufferPoolSize    types.String `tfsdk:"innodb_buffer_pool_size"`   // 缓存池大小
	InnodbThreadConcurrency types.Int64  `tfsdk:"innodb_thread_concurrency"` // 线程数
	ProdRunningStatus       types.Int32  `tfsdk:"prod_running_status"`       // 实例运行状态: 0->运行正常, 1->重启中, 2-备份操作中,3->恢复操作中,4->转换ssl,5->异常,6->修改参数组中,7->已冻结,8->已注销,9->施工中,10->施工失败,11->扩容中,12->主备切换中
	ProdRunningStatusDesc   types.String `tfsdk:"prod_running_status_desc"`  // prod_running_status的解释版本
	EipID                   types.String `tfsdk:"eip_id"`                    // eip id
	IsUpgradeBackUp         types.Bool   `tfsdk:"is_upgrade_back_up"`        // DDS模块磁盘扩容时候会使用 是否主磁盘与备磁盘一起扩容
	HostIp                  types.String `tfsdk:"host_ip"`                   // 主机ip
	ProdPerformanceSpec     types.String `tfsdk:"prod_performance_spec"`     // 主机配置
	StorageType             types.String `tfsdk:"storage_type"`              // 存储类型
	StorageSpace            types.Int32  `tfsdk:"storage_space"`             // 存储空间
	AvailabilityZoneInfo    types.List   `tfsdk:"availability_zone_info"`    // 节点可用区信息
	ShardNum                types.Int32  `tfsdk:"shard_num"`                 // 当实例为集群版，shard数量
	MongosNum               types.Int32  `tfsdk:"mongos_num"`                // 当实例为集群版，mongos节点数量
	InstanceSeries          types.String `tfsdk:"instance_series"`           // 实例规格（实例类型，1=通用型，2=计算增强型，3=内存优化型，4=直通（未用到））
	replicaNum              types.Int32  `tfsdk:"replica_num"`               // 副本集数量
	backupStorageSpace      types.Int32  `tfsdk:"backup_storage_space"`      // 备用节点磁盘空间，升配时使用
	UpgradeNodeType         types.String `tfsdk:"upgrade_node_type"`         // 集群版mongodb升配规格时，
}

type NodeInfoListModel struct {
	NodeType             types.String `tfsdk:"node_type"`              // 实例类型：master 或 readNode
	InstanceSeries       types.String `tfsdk:"instance_series"`        // 实例规格（实例类型，1=通用型，2=计算增强型，3=内存优化型，4=直通（未用到））
	StorageType          types.String `tfsdk:"storage_type"`           // 存储类型：SSD, SATA, SAS, SSD-genric, FAST-SSD
	StorageSpace         types.Int32  `tfsdk:"storage_space"`          // 存储空间（单位：GB，范围100到32768）
	ProdPerformanceSpec  types.String `tfsdk:"prod_performance_spec"`  // 规格（例：4C8G）
	AvailabilityZoneInfo types.List   `tfsdk:"availability_zone_info"` // 可用区信息
}

type AvailabilityZoneModel struct {
	AvailabilityZoneName  types.String `tfsdk:"availability_zone_name"`  // 资源池可用区名称
	AvailabilityZoneCount types.Int32  `tfsdk:"availability_zone_count"` // 资源池可用区总数
	NodeType              types.String `tfsdk:"node_type"`               // 表示分布AZ的节点类型，master/slave/readNode
}
