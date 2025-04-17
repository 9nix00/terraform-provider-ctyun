package nat

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctvpc"
	"terraform-provider-ctyun/internal/utils"
)

type ctyunSNats struct {
	meta *common.CtyunMetadata
}

func (c *ctyunSNats) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_nat_snats"
}

func (c *ctyunSNats) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://work.ctyun.cn/git/vnet/openapi-docs/src/branch/master/network/ctvpc/%E8%8E%B7%E5%8F%96SNAT%E5%88%97%E8%A1%A8.md`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    false,
				Description: "资源池id，如果不填这默认使用provider ctyun总region_id 或者环境变量",
			},
			"nat_gateway_id": schema.StringAttribute{
				Optional:    true,
				Description: "AT网关ID，选填",
			},
			"snat_id": schema.StringAttribute{
				Optional:    true,
				Description: "snat id，选填",
			},
			"subnet_id": schema.StringAttribute{
				Optional:    true,
				Description: "子网id，选填",
			},
			"page_number": schema.Int64Attribute{
				Optional:    true,
				Description: "列表的页码，默认值为1",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "分页查询时每页的行数，最大值为50，默认值为10。",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(50),
				},
			},
		},
	}
}

func (c *ctyunSNats) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunSNatsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	// region_id必不能为空
	if regionId == "" {
		msg := "regionID不能为空"
		response.Diagnostics.AddError(msg, msg)
		return
	}
	natGatewayId := config.NatGateWayID.ValueString()
	snatId := config.SubNetID.ValueString()
	subnetId := config.SubNetID.ValueString()
	// 分页信息先预留
	//pageNumber := c.ParseIntIfEmpty(config.PageNumber, types.Int64Value(1))
	//pageSize := c.ParseIntIfEmpty(config.PageSize, types.Int64Value(10))

	params := &ctvpc.CtvpcListSnatsRequest{
		RegionID:     regionId,
		NatGatewayID: &natGatewayId,
		SNatID:       &snatId,
		SubnetID:     &subnetId,
		//PageNumber:   pageNumber,
		//PageSize:     pageSize,
	}
	// 请求sdk,获取snat列表，list-snat接口返回值800为成功，900为失败
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListSnatsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	var snats []CtyunSNatsModel
	for _, page := range resp.ReturnObj {
		for _, snat := range page.Results {
			snatItem := CtyunSNatsModel{
				SNatID:       utils.SecStringValue(snat.SNatID),
				description:  utils.SecStringValue(snat.Description),
				SubNetCidr:   utils.SecStringValue(snat.SubnetID),
				SubNetType:   utils.SecStringValue(snat.SubnetID),
				CreationTime: utils.SecStringValue(snat.CreationTime),
				SubnetID:     utils.SecStringValue(snat.SubnetID),
				NatGatewayID: utils.SecStringValue(snat.NatGatewayID),
			}
			var eips []CtyunSNatsEipModel
			for _, eip := range snat.Eips {
				eipItem := CtyunSNatsEipModel{
					EipID:     utils.SecStringValue(eip.EipID),
					IpAddress: utils.SecStringValue(eip.IpAddress),
				}
				eips = append(eips, eipItem)
			}
			snatItem.Eips = eips
			snats = append(snats, snatItem)
		}
	}
	config.RegionID = types.StringValue(regionId)
	config.NatGateWayID = types.StringValue(natGatewayId)
	config.SNatID = types.StringValue(snatId)
	config.SubNetID = types.StringValue(subnetId)
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
}

// ParseIntIfEmpty 自定义方法，用于判断types.Int64类型字段是否为空，若为空则返回默认值。
func (c *ctyunSNats) ParseIntIfEmpty(value types.Int64, defaultValue types.Int64) types.Int64 {
	if value.IsNull() {
		return defaultValue
	}
	return value
}

type CtyunSNatsConfig struct {
	RegionID     types.String `tfsdk:"region_id"`      //区域id
	NatGateWayID types.String `tfsdk:"nat_gateway_id"` //要查询的NAT网关的ID。
	SNatID       types.String `tfsdk:"snat_id"`        // snat id
	SubNetID     types.String `tfsdk:"subnet_id"`      // 子网id
	PageNumber   types.Int64  `tfsdk:"page_number"`    //	列表的页码，默认值为1。
	PageSize     types.Int64  `tfsdk:"page_size"`      //分页查询时每页的行数，最大值为50，默认值为10。
}

type CtyunSNatsModel struct {
	SNatID       types.String         `tfsdk:"snat_id"`        //snat id
	description  types.String         `tfsdk:"description"`    //描述信息
	SubNetCidr   types.String         `tfsdk:"subnet_cidr"`    //要查询的NAT网关所属VPC子网的cidr
	SubNetType   types.String         `tfsdk:"subnet_type"`    //子网类型：1-有vpcID的子网，0-自定义
	CreationTime types.String         `tfsdk:"creation_time"`  //创建时间
	Eips         []CtyunSNatsEipModel `tfsdk:"eips"`           //绑定的 eip 信息
	SubnetID     types.String         `tfsdk:"subnet_id"`      //子网 ID
	NatGatewayID types.String         `tfsdk:"nat_gateway_id"` //ctvpc 网关 ID
}

type CtyunSNatsEipModel struct {
	EipID     types.String `tfsdk:"eip_id"`     //弹性 IP id
	IpAddress types.String `tfsdk:"ip_address"` //弹性 IP 地址
}
