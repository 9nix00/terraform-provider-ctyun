package hpfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/hpfs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunHpfsInstances{}
	_ datasource.DataSourceWithConfigure = &ctyunHpfsInstances{}
)

type ctyunHpfsInstances struct {
	meta *common.CtyunMetadata
}

func (c *ctyunHpfsInstances) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_hpfs_instances"
}

func (c *ctyunHpfsInstances) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *ctyunHpfsInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunHpfsInstancesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("region id 为空")
		return
	}
	params := &hpfs.HpfsListSfsRequest{
		RegionID: regionId,
		PageSize: config.pageSize.ValueInt32(),
		PageNo:   config.pageNo.ValueInt32(),
	}

	if !config.SfsStatus.IsNull() {
		params.SfsStatus = config.SfsStatus.ValueString()
	}
	if !config.SfsProtocol.IsNull() {
		params.SfsProtocol = config.SfsProtocol.ValueString()
	}
	if !config.AzName.IsNull() {
		params.AzName = config.AzName.ValueString()
	}
	if !config.ProjectID.IsNull() {
		params.ProjectID = config.ProjectID.ValueString()
	}

	resp, err := c.meta.Apis.SdkHpfsApis.HpfsListSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = errors.New("查询hpfs 列表失败，返回为nil")
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var hpfsInstances []CtyunHpfsInstancesModel
	hpfsList := resp.ReturnObj.List
	for _, hpfsItem := range hpfsList {
		var hpfsInstance CtyunHpfsInstancesModel
		hpfsInstance.RegionID = types.StringValue(hpfsItem.RegionID)
		hpfsInstance.SfsName = types.StringValue(hpfsItem.SfsName)
		hpfsInstance.SfsType = types.StringValue(hpfsItem.SfsType)
		hpfsInstance.SfsID = types.StringValue(hpfsItem.SfsUID)
		hpfsInstance.SfsSize = types.Int32Value(hpfsItem.SfsSize)
		hpfsInstance.SfsStatus = types.StringValue(hpfsItem.SfsStatus)
		hpfsInstance.UsedSize = types.Int32Value(hpfsItem.UsedSize)
		hpfsInstance.CreateTime = types.Int64Value(hpfsItem.CreateTime)
		hpfsInstance.UpdateTime = types.Int64Value(hpfsItem.UpdateTime)
		hpfsInstance.ProjectID = types.StringValue(hpfsItem.ProjectID)
		hpfsInstance.OnDemand = types.BoolValue(*hpfsItem.OnDemand)
		hpfsInstance.AzName = types.StringValue(hpfsItem.AzName)
		hpfsInstance.ClusterName = types.StringValue(hpfsItem.ClusterName)
		hpfsInstance.Baseline = types.StringValue(hpfsItem.Baseline)
		hpfsInstance.HpfsSharePath = types.StringValue(hpfsItem.HpfsSharePath)
		hpfsInstance.SecretKey = types.StringValue(hpfsItem.SecretKey)
		hpfsInstance.DataflowCount = types.Int32Value(hpfsItem.DataflowCount)
		dataflowList, diagnostics := types.SetValueFrom(ctx, types.StringType, hpfsItem.DataflowList)
		if diagnostics.HasError() {
			err = errors.New(diagnostics[0].Detail())
			return
		}
		hpfsInstance.DataflowList = dataflowList
		hpfsInstances = append(hpfsInstances, hpfsInstance)
	}

}

func (c *ctyunHpfsInstances) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunHpfsInstancesModel struct {
	SfsName       types.String `tfsdk:"sfs_name"`        // 并行文件命名
	SfsID         types.String `tfsdk:"sfs_id"`          // 并行文件唯一ID
	SfsSize       types.Int32  `tfsdk:"sfs_size"`        // 大小(GB)
	SfsType       types.String `tfsdk:"sfs_type"`        // 文件系统类型
	SfsProtocol   types.String `tfsdk:"sfs_protocol"`    // 挂载协议
	SfsStatus     types.String `tfsdk:"sfs_status"`      // 文件系统状态
	UsedSize      types.Int32  `tfsdk:"used_size"`       // 已用大小(MB)
	CreateTime    types.Int64  `tfsdk:"create_time"`     // 创建时间戳(毫秒)
	UpdateTime    types.Int64  `tfsdk:"update_time"`     // 更新时间戳(毫秒)
	ProjectID     types.String `tfsdk:"project_id"`      // 企业项目ID
	OnDemand      types.Bool   `tfsdk:"on_demand"`       // 是否按需订购
	RegionID      types.String `tfsdk:"region_id"`       // 资源池ID
	AzName        types.String `tfsdk:"az_name"`         // 可用区名称
	ClusterName   types.String `tfsdk:"cluster_name"`    // 集群名称
	Baseline      types.String `tfsdk:"baseline"`        // 性能基线(MB/s/TB)
	HpfsSharePath types.String `tfsdk:"hpfs_share_path"` // HPFS共享路径
	SecretKey     types.String `tfsdk:"secret_key"`      // HPC挂载密钥
	DataflowList  types.Set    `tfsdk:"dataflow_list"`   // 数据流动策略ID列表
	DataflowCount types.Int32  `tfsdk:"dataflow_count"`  // 数据流动策略数量
}

type CtyunHpfsInstancesConfig struct {
	RegionID      types.String              `tfsdk:"region_id"`
	SfsStatus     types.String              `tfsdk:"sfs_status"`     // 并行文件状态。creating/available/unusable，不传为查询全部
	SfsProtocol   types.String              `tfsdk:"sfs_protocol"`   // 挂载协议。2 种，nfs/hpfs ，不传为查询全部
	AzName        types.String              `tfsdk:"az_name"`        // 多可用区下的可用区名字，不传为查询全部
	ProjectID     types.String              `tfsdk:"project_id"`     // 资源所属企业项目 ID，默认为"0"
	pageSize      types.Int32               `tfsdk:"page_size"`      // 每页包含的元素个数范围(1-50)，默认值为10
	pageNo        types.Int32               `tfsdk:"page_no"`        // 列表的分页页码，默认值为1
	HpfsInstances []CtyunHpfsInstancesModel `tfsdk:"hpfs_instances"` // hpfs列表
}
