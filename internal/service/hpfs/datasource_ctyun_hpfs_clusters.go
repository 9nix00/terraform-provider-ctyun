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
	_ datasource.DataSource              = &ctyunHpfsClusters{}
	_ datasource.DataSourceWithConfigure = &ctyunHpfsClusters{}
)

type ctyunHpfsClusters struct {
	meta *common.CtyunMetadata
}

func (c *ctyunHpfsClusters) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunHpfsClusters) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_hpfs_clusters"
}

func (c *ctyunHpfsClusters) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *ctyunHpfsClusters) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunHpfsClustersConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("region id 为空")
		return
	}
	params := &hpfs.HpfsListClusterRequest{
		RegionID: regionId,
		PageNo:   config.pageNo.ValueInt32(),
		PageSize: config.pageSize.ValueInt32(),
	}
	if !config.SfsType.IsNull() {
		params.SfsType = config.SfsType.ValueString()
	}
	if !config.AzName.IsNull() {
		params.AzName = config.AzName.ValueString()
	}
	if !config.EbmDeviceType.IsNull() {
		params.EbmDeviceType = config.EbmDeviceType.ValueString()
	}

	resp, err := c.meta.Apis.SdkHpfsApis.HpfsListClusterApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = errors.New("查询hpfs 集群列表失败，返回为nil")
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var hpfsClusterList []CtyunHpfsClusterModel
	clusterList := resp.ReturnObj.ClusterList
	for _, clusterItem := range clusterList {
		var cluster CtyunHpfsClusterModel
		cluster.ClusterName = types.StringValue(clusterItem.ClusterName)
		cluster.RemainingStatus = types.BoolValue(*clusterItem.RemainingStatus)
		cluster.StorageType = types.StringValue(clusterItem.StorageType)
		cluster.AzName = types.StringValue(clusterItem.AzName)
		cluster.NetworkType = types.StringValue(clusterItem.NetworkType)
		protocolType, diags := types.SetValueFrom(ctx, types.StringType, clusterItem.ProtocolType)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		cluster.ProtocolType = protocolType
		baselines, diags := types.SetValueFrom(ctx, types.StringType, clusterItem.Baselines)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		cluster.Baselines = baselines

		ebmDeviceTypes, diags := types.SetValueFrom(ctx, types.StringType, clusterItem.EbmDeviceTypes)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		cluster.EbmDeviceTypes = ebmDeviceTypes

		hpfsClusterList = append(hpfsClusterList, cluster)
	}
}

type CtyunHpfsClusterModel struct {
	ClusterName     types.String `tfsdk:"cluster_name"`     // 集群名称
	RemainingStatus types.Bool   `tfsdk:"remaining_status"` // 是否可以售卖
	StorageType     types.String `tfsdk:"storage_type"`     // 集群的存储类型
	AzName          types.String `tfsdk:"az_name"`          // 多可用区下的可用区名字
	ProtocolType    types.Set    `tfsdk:"protocol_type"`    // 支持的协议列表
	Baselines       types.Set    `tfsdk:"baselines"`        // 性能基线列表
	NetworkType     types.String `tfsdk:"network_type"`     // 集群的网络类型
	EbmDeviceTypes  types.Set    `tfsdk:"ebm_device_types"` // 裸金属设备规格列表
}

type CtyunHpfsClustersConfig struct {
	RegionID      types.String            `tfsdk:"region_id"`       // 资源池 ID
	SfsType       types.String            `tfsdk:"sfs_type"`        // 文件系统类型
	AzName        types.String            `tfsdk:"az_name"`         // 可用区名称
	EbmDeviceType types.String            `tfsdk:"ebm_device_type"` // 裸金属设备规格
	PageNo        types.Int64             `tfsdk:"page_no"`         // 分页页码
	PageSize      types.Int64             `tfsdk:"page_size"`       // 每页元素数量
	pageSize      types.Int32             `tfsdk:"page_size"`       // 每页包含的元素个数范围(1-50)，默认值为10
	pageNo        types.Int32             `tfsdk:"page_no"`         // 列表的分页页码，默认值为1
	HpfsClusters  []CtyunHpfsClusterModel `tfsdk:"hpfs_clusters"`   // hpfs cluster列表
}
