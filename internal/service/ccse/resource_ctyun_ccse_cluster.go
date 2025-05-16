package ccse

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	ccse2 "terraform-provider-ctyun/internal/core/ccse"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
	"time"
)

var (
	_ resource.Resource                = &ctyunCcseCluster{}
	_ resource.ResourceWithConfigure   = &ctyunCcseCluster{}
	_ resource.ResourceWithImportState = &ctyunCcseCluster{}
)

type ctyunCcseCluster struct {
	meta *common.CtyunMetadata
}

func NewCtyunCcseCluster() resource.Resource {
	return &ctyunCcseCluster{}
}

func (c *ctyunCcseCluster) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ccse_cluster"
}

type CtyunCcseClusterConfig struct {
	ID         types.String             `tfsdk:"id"`
	RegionID   types.String             `tfsdk:"region_id"`
	BaseInfo   CtyunCcseClusterBaseInfo `tfsdk:"base_info"`
	SlaveHost  CtyunCcseClusterSlave    `tfsdk:"slave_host"`
	MasterHost *CtyunCcseClusterMaster  `tfsdk:"master_host"`
}

type CtyunCcseClusterBaseInfo struct {
	ProjectID                 types.String `tfsdk:"project_id"`
	VpcID                     types.String `tfsdk:"vpc_id"`
	SubnetID                  types.String `tfsdk:"subnet_id"`
	AutoGenerateSecurityGroup types.Bool   `tfsdk:"auto_generate_security_group"`
	SecurityGroupID           types.String `tfsdk:"security_group_id"`
	ClusterName               types.String `tfsdk:"cluster_name"`
	ClusterDomain             types.String `tfsdk:"cluster_domain"`
	NetworkPlugin             types.String `tfsdk:"network_plugin"`
	StartPort                 types.Int32  `tfsdk:"start_port"`
	EndPort                   types.Int32  `tfsdk:"end_port"`
	ElbProdCode               types.String `tfsdk:"elb_prod_code"`
	PodCidr                   types.String `tfsdk:"pod_cidr"`
	PodSubnetIdList           []string     `tfsdk:"pod_subnet_id_list"`
	CycleType                 types.String `tfsdk:"cycle_type"`
	CycleCount                types.Int64  `tfsdk:"cycle_count"`
	ContainerRuntime          types.String `tfsdk:"container_runtime"`
	TimeZone                  types.String `tfsdk:"timezone"`
	ClusterVersion            types.String `tfsdk:"cluster_version"`
	DeployType                types.String `tfsdk:"deploy_type"`
	KubeProxy                 types.String `tfsdk:"kube_proxy"`
	ClusterSeries             types.String `tfsdk:"cluster_series"`
}

type CtyunCcseClusterAzInfo struct {
	AzName types.String `tfsdk:"az_name"`
	Size   types.Int32  `tfsdk:"size"`
}
type CtyunCcseClusterMaster struct {
	ItemDefName types.String             `tfsdk:"item_def_name"`
	AzInfos     []CtyunCcseClusterAzInfo `tfsdk:"az_infos"`
	SysDisk     CtyunCcseClusterDisk     `tfsdk:"sys_disk"`
	DataDisks   []CtyunCcseClusterDisk   `tfsdk:"data_disks"`
}
type CtyunCcseClusterSlave struct {
	ItemDefName  types.String             `tfsdk:"item_def_name"`
	AzInfos      []CtyunCcseClusterAzInfo `tfsdk:"az_infos"`
	SysDisk      CtyunCcseClusterDisk     `tfsdk:"sys_disk"`
	DataDisks    []CtyunCcseClusterDisk   `tfsdk:"data_disks"`
	InstanceType types.String             `tfsdk:"instance_type"`
	MirrorID     types.String             `tfsdk:"mirror_id"`
	MirrorName   types.String             `tfsdk:"mirror_name"`
	MirrorType   types.Int32              `tfsdk:"mirror_type"`
}

type CtyunCcseClusterDisk struct {
	Type types.String `tfsdk:"type"`
	Size types.Int32  `tfsdk:"size"`
}

func (c *ctyunCcseCluster) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID",
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
			"base_info": schema.SingleNestedAttribute{
				Required:    true,
				Description: "支持拉丁字母、中文、数字，下划线，连字符，中文/英文字母开头，不能以http:/https:开头，长度2-32",
				Attributes: map[string]schema.Attribute{
					"project_id": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "企业项目ID",
						Default:     defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"vpc_id": schema.StringAttribute{
						Required:    true,
						Description: "vpc id",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"subnet_id": schema.StringAttribute{
						Required:    true,
						Description: "子网ID",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"auto_generate_security_group": schema.BoolAttribute{
						Optional:    true,
						Description: "是否自动生成安全组",
					},
					"security_group_id": schema.StringAttribute{
						Optional:    true,
						Description: "安全组ID",
						Validators: []validator.String{
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("base_info").AtName("auto_generate_security_group"),
								types.BoolValue(false),
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"cluster_name": schema.StringAttribute{
						Required:    true,
						Description: "集群名字",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"cluster_domain": schema.StringAttribute{
						Required:    true,
						Description: "集群本地域名",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},

					"network_plugin": schema.StringAttribute{
						Required:    true,
						Description: "网络插件",
						Validators: []validator.String{
							stringvalidator.OneOf("calico", "cubecni"),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"start_port": schema.Int32Attribute{
						Required:    true,
						Description: "节点服务开始端口，可选范围30000-65535",
						Validators: []validator.Int32{
							int32validator.Between(30000, 65535),
						},
						PlanModifiers: []planmodifier.Int32{
							int32planmodifier.RequiresReplace(),
						},
					},
					"end_port": schema.Int32Attribute{
						Required:    true,
						Description: "节点服务终止端口，可选范围30000-65535",
						Validators: []validator.Int32{
							int32validator.Between(30000, 65535),
						},
						PlanModifiers: []planmodifier.Int32{
							int32planmodifier.RequiresReplace(),
						},
					},
					"elb_prod_code": schema.StringAttribute{
						Required:    true,
						Description: "ApiServer的ELB类型，standardI（标准I型） ,standardII（标准II型）, enhancedI（增强I型）, enhancedII（增强II型） , higherI（高阶I型）",
						Validators: []validator.String{
							stringvalidator.OneOf("standardI", "standardII", "enhancedI", "enhancedII", "higherI"),
						},
					},
					"pod_cidr": schema.StringAttribute{
						Required:    true,
						Description: "pod网络cidr，使用cubecni作为网络插件时，podCidr传值为vpc cidr。使用calico作为网络插件时，podCidr与vpcCidr和serviceCidr不能重叠。",
						Validators: []validator.String{
							validator2.Cidr(),
						},
					},
					"pod_subnet_id_list": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "pod子网id列表，网络插件选择cubecni必传",
						Validators: []validator.Set{
							validator2.AlsoRequiresEqualSet(
								path.MatchRoot("base_info").AtName("network_plugin"),
								types.StringValue("cubecni"),
							),
						},
					},
					"cycle_type": schema.StringAttribute{
						Required:    true,
						Description: "订购周期类型，取值范围：month：按月，year：按年、on_demand：按需。当此值为month或者year时，cycle_count为必填",
						Validators: []validator.String{
							stringvalidator.OneOf(business.OrderCycleTypes...),
						},
					},
					"cycle_count": schema.Int64Attribute{
						Optional:    true,
						Description: "订购时长，该参数在cycle_type为month或year时才生效，当cycleType=month，支持续订1-11个月；当cycleType=year，支持续订1-5年",
						Validators: []validator.Int64{
							validator2.AlsoRequiresEqualInt64(
								path.MatchRoot("base_info").AtName("cycle_type"),
								types.StringValue(business.OrderCycleTypeMonth),
								types.StringValue(business.OrderCycleTypeYear),
							),
							validator2.ConflictsWithEqualInt64(
								path.MatchRoot("base_info").AtName("cycle_type"),
								types.StringValue(business.OrderCycleTypeOnDemand),
							),
							validator2.CycleCount(1, 11, 1, 3),
						},
					},
					"container_runtime": schema.StringAttribute{
						Required:    true,
						Description: "容器运行时,可选containerd、docker",
						Validators: []validator.String{
							stringvalidator.OneOf("containerd", "docker"),
						},
					},
					"timezone": schema.StringAttribute{
						Required:    true,
						Description: "时区，例如Asia/Shanghai (UTC+08:00)",
					},
					"cluster_version": schema.StringAttribute{
						Required:    true,
						Description: "集群版本，支持1.23.3 ，1.25.6 ，1.27.8，1.29.3",
						Validators: []validator.String{
							stringvalidator.OneOf("1.23.3", "1.25.6", "1.27.8", "1.29.3"),
						},
					},
					"deploy_type": schema.StringAttribute{
						Required:    true,
						Description: "部署模式，单可用区为single，多可用区为multi",
						Validators: []validator.String{
							stringvalidator.OneOf("single", "multi"),
						},
					},
					"kube_proxy": schema.StringAttribute{
						Required:    true,
						Description: "kubeProxy类型：iptables或ipvs。您可查看<a href=\"https://www.ctyun.cn/document/10083472/10915725\">iptables与IPVS如何选择</a>",
						Validators: []validator.String{
							stringvalidator.OneOf("iptables", "ipvs"),
						},
					},
					"cluster_series": schema.StringAttribute{
						Required:    true,
						Description: "集群系列，cce.standard，cce.managed，您可查看<a href=\"https://www.ctyun.cn/document/10083472/10892150\">产品定义</a>选择",
						Validators: []validator.String{
							stringvalidator.OneOf("cce.standard", "cce.managed"),
						},
					},
				},
			},
			"master_host": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "master节点基本信息，专有版必填，托管版时不传",
				Attributes: map[string]schema.Attribute{
					"item_def_name": schema.StringAttribute{
						Required:    true,
						Description: "规格名称",
					},
					"az_infos": schema.ListNestedAttribute{
						Required:    true,
						Description: "可用区信息",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"az_name": schema.StringAttribute{
									Required:    true,
									Description: "可用区编码",
								},
								"size": schema.Int32Attribute{
									Required:    true,
									Description: "节点数量",
								},
							},
						},
					},
					"sys_disk": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "系统盘",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:    true,
								Description: "盘规格",
								Validators: []validator.String{
									stringvalidator.OneOf("SATA", "SAS", "SSD"),
								},
							},
							"size": schema.Int32Attribute{
								Required:    true,
								Description: "盘大小，单位为G",
							},
						},
					},
					"data_disks": schema.ListNestedAttribute{
						Optional:    true,
						Description: "数据盘",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required:    true,
									Description: "盘规格",
									Validators: []validator.String{
										stringvalidator.OneOf("SATA", "SAS", "SSD"),
									},
								},
								"size": schema.Int32Attribute{
									Required:    true,
									Description: "盘大小，单位为G",
								},
							},
						},
					},
				},
				Validators: []validator.Object{
					validator2.AlsoRequiresEqualObject(
						path.MatchRoot("base_info").AtName("cluster_series"),
						types.StringValue("cce.standard-专有版"),
					),
					validator2.ConflictsWithEqualObject(
						path.MatchRoot("base_info").AtName("cluster_series"),
						types.StringValue("cce.managed-托管版"),
					),
				},
			},
			"slave_host": schema.SingleNestedAttribute{
				Required:    true,
				Description: "slave节点基本信息",
				Attributes: map[string]schema.Attribute{
					"instance_type": schema.StringAttribute{
						Required:    true,
						Description: "实例类型， ecs 或 ebm",
						Validators: []validator.String{
							stringvalidator.OneOf("ecs", "ebm"),
						},
					},
					"mirror_id": schema.StringAttribute{
						Optional:    true,
						Description: "云主机镜像id, 可查看<a href=\"https://www.ctyun.cn/document/10083472/11004475\">节点规格和节点镜像</a>",
						Validators: []validator.String{
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("slave_host").AtName("instance_type"),
								types.StringValue("ecs"),
							),
							validator2.ConflictsWithEqualString(
								path.MatchRoot("slave_host").AtName("instance_type"),
								types.StringValue("ebm"),
							),
						},
					},
					"mirror_name": schema.StringAttribute{
						Optional:    true,
						Description: "物理机镜像名称, 可查看<a href=\"https://www.ctyun.cn/document/10083472/11004475\">节点规格和节点镜像</a>",
						Validators: []validator.String{
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("slave_host").AtName("instance_type"),
								types.StringValue("ebm"),
							),
							validator2.ConflictsWithEqualString(
								path.MatchRoot("slave_host").AtName("instance_type"),
								types.StringValue("ecs"),
							),
						},
					},
					"mirror_type": schema.Int32Attribute{
						Required:    true,
						Description: "镜像类型，0-私有，1-公有",
						Validators: []validator.Int32{
							int32validator.Between(0, 1),
						},
					},
					"item_def_name": schema.StringAttribute{
						Required:    true,
						Description: "规格名称",
					},
					"az_infos": schema.ListNestedAttribute{
						Required:    true,
						Description: "可用区信息",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"az_name": schema.StringAttribute{
									Required:    true,
									Description: "可用区编码",
								},
								"size": schema.Int32Attribute{
									Required:    true,
									Description: "节点数量",
								},
							},
						},
					},
					"sys_disk": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "系统盘",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:    true,
								Description: "系统盘规格",
								Validators: []validator.String{
									stringvalidator.OneOf("SATA", "SAS", "SSD"),
								},
							},
							"size": schema.Int32Attribute{
								Required:    true,
								Description: "系统盘大小，单位为G",
							},
						},
					},
					"data_disks": schema.ListNestedAttribute{
						Optional:    true,
						Description: "数据盘",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required:    true,
									Description: "系统盘规格",
									Validators: []validator.String{
										stringvalidator.OneOf("SATA", "SAS", "SSD"),
									},
								},
								"size": schema.Int32Attribute{
									Required:    true,
									Description: "系统盘大小，单位为G",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunCcseCluster) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunCcseClusterConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}

	// 创建后检查
	id, err := c.checkAfterCreate(ctx, plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(id)
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunCcseCluster) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCcseClusterConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunCcseCluster) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunCcseClusterConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunCcseClusterConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunCcseCluster) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCcseClusterConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDelete(ctx, state)
	if err != nil {
		return
	}
	//response.State.RemoveResource(ctx)
}

func (c *ctyunCcseCluster) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [id],[regionID]
func (c *ctyunCcseCluster) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	err = fmt.Errorf("resource_ctyun_ccse_cluster 不支持 terraform import")
	//var cfg CtyunCcseClusterConfig
	//var id, regionID string
	//err = terraform_extend.Split(request.ID, &id, &regionID)
	//if err != nil {
	//	return
	//}
	//cfg.RegionID = types.StringValue(regionID)
	//// 查询远端
	//err = c.getAndMerge(ctx, &cfg)
	//if err != nil {
	//	return
	//}
	//response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeCreate 创建前检查
func (c *ctyunCcseCluster) checkBeforeCreate(ctx context.Context, plan CtyunCcseClusterConfig) (err error) {
	// 确保当前虚拟私有云存在，且子网与虚拟私有云存在对应关系
	vpc, regionID, projectID := plan.BaseInfo.VpcID.ValueString(), plan.RegionID.ValueString(), plan.BaseInfo.ProjectID.ValueString()
	err = business.NewVpcService(c.meta).MustExist(ctx, vpc, regionID, projectID)
	if err != nil {
		return err
	}
	return
}

// create 创建
func (c *ctyunCcseCluster) create(ctx context.Context, plan CtyunCcseClusterConfig) (err error) {
	params := &ccse2.CcseCreateClusterRequest{
		RegionId:  plan.RegionID.ValueString(),
		ResPoolId: plan.RegionID.ValueString(),
	}
	// 处理 clusterBaseInfo
	clusterBaseInfo := ccse2.CcseCreateClusterClusterBaseInfoRequest{
		ProjectId:                 plan.BaseInfo.ProjectID.ValueString(),
		VpcUuid:                   plan.BaseInfo.VpcID.ValueString(),
		SubnetUuid:                plan.BaseInfo.SubnetID.ValueString(),
		AutoGenerateSecurityGroup: plan.BaseInfo.AutoGenerateSecurityGroup.ValueBoolPointer(),
		SecurityGroupUuid:         plan.BaseInfo.SecurityGroupID.ValueString(),
		ClusterName:               plan.BaseInfo.ClusterName.ValueString(),
		ClusterDomain:             plan.BaseInfo.ClusterDomain.ValueString(),
		ClusterVersion:            plan.BaseInfo.ClusterVersion.ValueString(),
		ClusterSeries:             plan.BaseInfo.ClusterSeries.ValueString(),
		NetworkPlugin:             plan.BaseInfo.NetworkPlugin.ValueString(),
		StartPort:                 int64(plan.BaseInfo.StartPort.ValueInt32()),
		EndPort:                   int64(plan.BaseInfo.EndPort.ValueInt32()),
		ElbProdCode:               plan.BaseInfo.ElbProdCode.ValueString(),
		PodSubnetUuidList:         plan.BaseInfo.PodSubnetIdList,
		PodCidr:                   plan.BaseInfo.PodCidr.ValueString(),
		ContainerRuntime:          plan.BaseInfo.ContainerRuntime.ValueString(),
		Timezone:                  plan.BaseInfo.TimeZone.ValueString(),
		DeployType:                plan.BaseInfo.DeployType.ValueString(),
		KubeProxy:                 plan.BaseInfo.KubeProxy.ValueString(),
		AzInfos:                   []*ccse2.CcseCreateClusterClusterBaseInfoAzInfosRequest{},
	}
	switch plan.BaseInfo.CycleType.ValueString() {
	case business.OnDemandCycleType:
		clusterBaseInfo.BillMode = "2"
	case business.MonthCycleType:
		clusterBaseInfo.BillMode = "1"
		clusterBaseInfo.CycleType = "3"
		clusterBaseInfo.CycleCnt = int32(plan.BaseInfo.CycleCount.ValueInt64())
	case business.YearCycleType:
		clusterBaseInfo.BillMode = "1"
		clusterBaseInfo.CycleType = fmt.Sprintf("%d", plan.BaseInfo.CycleCount.ValueInt64()+4) // 1年传5，2年传6，3年传7
		clusterBaseInfo.CycleCnt = 1
	}

	// 处理masterHost
	if plan.MasterHost != nil {
		flavorName := plan.MasterHost.ItemDefName.ValueString()
		var flavor ctecs.EcsFlavorListFlavorListResponse
		flavor, err = business.NewEcsService(c.meta).GetFlavorByName(ctx, flavorName, plan.RegionID.ValueString())
		if err != nil {
			return
		}
		if flavor.FlavorCpu < 4 || flavor.FlavorRam < 8 {
			err = fmt.Errorf("master节点的规格至少需要4c8g")
		}

		masterHost := ccse2.CcseCreateClusterMasterHostRequest{
			Cpu:         int32(flavor.FlavorCpu),
			Mem:         int32(flavor.FlavorRam),
			ItemDefName: flavorName,
			ItemDefType: flavor.FlavorType,
			Size:        0,
			SysDisk: &ccse2.CcseCreateClusterMasterHostSysDiskRequest{
				ItemDefName: plan.MasterHost.SysDisk.Type.ValueString(),
				Size:        plan.MasterHost.SysDisk.Size.ValueInt32(),
			},
		}
		for _, az := range plan.MasterHost.AzInfos {
			clusterBaseInfo.AzInfos = append(clusterBaseInfo.AzInfos, &ccse2.CcseCreateClusterClusterBaseInfoAzInfosRequest{
				AzName: az.AzName.ValueString(),
				Size:   az.Size.ValueInt32(),
			})
			masterHost.Size += az.Size.ValueInt32()
		}
		for _, disk := range plan.MasterHost.DataDisks {
			masterHost.DataDisks = append(masterHost.DataDisks, &ccse2.CcseCreateClusterMasterHostDataDisksRequest{
				ItemDefName: disk.Type.ValueString(),
				Size:        disk.Size.ValueInt32(),
			})
		}
		params.MasterHost = &masterHost
	}

	// 处理slaveHost
	flavorName := plan.SlaveHost.ItemDefName.ValueString()
	flavor, err := business.NewEcsService(c.meta).GetFlavorByName(ctx, flavorName, plan.RegionID.ValueString())
	if err != nil {
		return
	}

	slaveHost := ccse2.CcseCreateClusterSlaveHostRequest{
		Cpu:         int32(flavor.FlavorCpu),
		Mem:         int32(flavor.FlavorRam),
		ItemDefName: flavorName,
		ItemDefType: flavor.FlavorType,
		Size:        0,
		SysDisk: &ccse2.CcseCreateClusterSlaveHostSysDiskRequest{
			ItemDefName: plan.SlaveHost.SysDisk.Type.ValueString(),
			Size:        plan.SlaveHost.SysDisk.Size.ValueInt32(),
		},
		MirrorType: plan.SlaveHost.MirrorType.ValueInt32(),
	}

	for _, az := range plan.SlaveHost.AzInfos {
		slaveHost.AzInfos = append(slaveHost.AzInfos, &ccse2.CcseCreateClusterSlaveHostAzInfosRequest{
			AzName: az.AzName.ValueString(),
			Size:   az.Size.ValueInt32(),
		})
		slaveHost.Size += az.Size.ValueInt32()
	}

	for _, disk := range plan.SlaveHost.DataDisks {
		slaveHost.DataDisks = append(slaveHost.DataDisks, &ccse2.CcseCreateClusterSlaveHostDataDisksRequest{
			ItemDefName: disk.Type.ValueString(),
			Size:        disk.Size.ValueInt32(),
		})
	}

	switch plan.SlaveHost.InstanceType.ValueString() {
	case "ecs":
		slaveHost.ForeignMirrorId = plan.SlaveHost.MirrorID.ValueString()
	case "ebm":
		slaveHost.MirrorName = plan.SlaveHost.MirrorName.ValueString()
	}

	params.ClusterBaseInfo = &clusterBaseInfo
	params.SlaveHost = &slaveHost

	resp, err := c.meta.Apis.SdkCcseApis.CcseCreateClusterApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// getAndMerge 从远端查询
func (c *ctyunCcseCluster) getAndMerge(ctx context.Context, plan *CtyunCcseClusterConfig) (err error) {
	params := &ccse2.CcseGetClusterRequest{
		RegionId:  plan.RegionID.ValueString(),
		ClusterId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseGetClusterApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	instance := resp.ReturnObj
	plan.BaseInfo.DeployType = types.StringValue(instance.DeployMode)
	plan.BaseInfo.VpcID = types.StringValue(instance.VpcId)
	plan.BaseInfo.SubnetID = types.StringValue(instance.SubnetUuid)
	plan.BaseInfo.NetworkPlugin = types.StringValue(instance.NetworkPlugin)
	plan.BaseInfo.PodCidr = types.StringValue(instance.PodCidr)
	plan.BaseInfo.ContainerRuntime = types.StringValue(instance.ContainerRuntime)
	plan.BaseInfo.TimeZone = types.StringValue(instance.Timezone)
	plan.BaseInfo.ClusterVersion = types.StringValue(instance.ClusterVersion)
	plan.BaseInfo.KubeProxy = types.StringValue(instance.KubeProxyPattern)
	switch instance.ClusterType {
	case 0:
		plan.BaseInfo.ClusterSeries = types.StringValue("cce.standard")
	case 2:
		plan.BaseInfo.ClusterSeries = types.StringValue("cce.managed")
	}
	plan.BaseInfo.StartPort = types.Int32Value(instance.StartPort)
	plan.BaseInfo.EndPort = types.Int32Value(instance.EndPort)
	return
}

// update 更新
func (c *ctyunCcseCluster) update(ctx context.Context, plan, state CtyunCcseClusterConfig) (err error) {
	//if plan.BaseInfo.ClusterVersion .Equal(state.BaseInfo.ClusterVersion) {
	//	return
	//}
	//params := ccse2.CcseUpgradeClusterRequest{
	//	ClusterId:   state.ID.ValueString(),
	//	RegionId:    state.RegionID.ValueString(),
	//	NextVersion: plan.BaseInfo.ClusterVersion.ValueString(),
	//	Version:     state.BaseInfo.ClusterVersion.ValueString(),
	//	Concurrency: 1,
	//}
	//resp, err := c.meta.Apis.SdkCcseApis.CcseDeleteClusterApi.Do(ctx, c.meta.SdkCredential, params)
	//if err != nil {
	//	return
	//}
	return
}

// delete 删除
func (c *ctyunCcseCluster) delete(ctx context.Context, plan CtyunCcseClusterConfig) (err error) {
	params := &ccse2.CcseDeleteClusterRequest{
		RegionId:   plan.RegionID.ValueString(),
		ResPoolId:  plan.RegionID.ValueString(),
		ProdInstId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseDeleteClusterApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// listByName 根据名称查询集群
func (c *ctyunCcseCluster) listByName(ctx context.Context, plan CtyunCcseClusterConfig) (clusters []*ccse2.CcseListClustersReturnObjRecordsResponse, err error) {
	params := &ccse2.CcseListClustersRequest{
		RegionId:    plan.RegionID.ValueString(),
		ClusterName: plan.BaseInfo.ClusterName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseListClustersApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	clusters = resp.ReturnObj.Records
	return
}

// checkAfterCreate 创建后检查
func (c *ctyunCcseCluster) checkAfterCreate(ctx context.Context, plan CtyunCcseClusterConfig) (id string, err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var clusters []*ccse2.CcseListClustersReturnObjRecordsResponse
			clusters, err = c.listByName(ctx, plan)
			if err != nil {
				return false
			}
			if len(clusters) == 0 || clusters[0].ClusterStatus == "creating" || clusters[0].Id == "" {
				return true
			}

			id = clusters[0].Id
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("创建时间过长")
	}
	return
}

// checkAfterDelete 删除后检查
func (c *ctyunCcseCluster) checkAfterDelete(ctx context.Context, plan CtyunCcseClusterConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var clusters []*ccse2.CcseListClustersReturnObjRecordsResponse
			clusters, err = c.listByName(ctx, plan)
			if err != nil {
				return false
			}
			if len(clusters) != 0 && clusters[0].ClusterStatus != "deleted" {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("删除时间过长")
	}
	return
}
