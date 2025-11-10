package dns

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunPrivateZoneRecords{}
	_ datasource.DataSourceWithConfigure = &CtyunPrivateZoneRecords{}
)

type CtyunPrivateZoneRecords struct {
	meta *common.CtyunMetadata
}

func NewCtyunPrivateZoneRecords() datasource.DataSource {
	return &CtyunPrivateZoneRecords{}
}
func (c *CtyunPrivateZoneRecords) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPrivateZoneRecords) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_zone_records"
}

func (c *CtyunPrivateZoneRecords) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunPrivateZoneRecords) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPrivateZoneRecordsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	config.RegionID = types.StringValue(regionId)

	params := &ctvpc.CtvpcListPrivateZoneRecordRequest{
		RegionID: config.RegionID.ValueString(),
		PageNo:   1,
		PageSize: 10,
	}
	if !config.Name.IsNull() {
		params.ZoneRecordName = config.Name.ValueStringPointer()
	}
	if !config.ZoneID.IsNull() {
		params.ZoneID = config.ZoneID.ValueStringPointer()
	}
	if !config.ID.IsNull() {
		params.ZoneRecordID = config.ID.ValueStringPointer()
	}
	if !config.PageNo.IsNull() {
		params.PageNo = config.PageNo.ValueInt32()
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListPrivateZoneRecordApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = errors.New("获取内网DNS记录列表失败！接口返回nil，请联系研发确认问题原因。")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var records []CtyunPrivateZoneRecordInfoModel
	for _, recordItem := range resp.ReturnObj.ZoneRecords {
		var record CtyunPrivateZoneRecordInfoModel
		record.ID = types.StringValue(recordItem.ZoneRecordID)
		record.Name = types.StringValue(recordItem.Name)
		record.ZoneID = types.StringValue(recordItem.ZoneID)
		record.ZoneName = types.StringValue(recordItem.ZoneName)
		record.Description = types.StringValue(recordItem.Description)
		record.Type = types.StringValue(recordItem.Type)
		record.TTL = types.Int32Value(int32(recordItem.TTL))
		record.CreateTime = types.StringValue(recordItem.CreatedAt)
		record.UpdateTime = types.StringValue(recordItem.UpdatedAt)
		record.Value = recordItem.Value
		records = append(records, record)
	}
	config.Records = records
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

}

type CtyunPrivateZoneRecordInfoModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	ZoneID      types.String `tfsdk:"zone_id"`
	ZoneName    types.String `tfsdk:"zone_name"`
	Description types.String `tfsdk:"description"`
	TTL         types.Int32  `tfsdk:"ttl"`
	Type        types.String `tfsdk:"type"`
	Value       []string     `tfsdk:"value"`
	CreateTime  types.String `tfsdk:"create_time"`
	UpdateTime  types.String `tfsdk:"update_time"`
}

type CtyunPrivateZoneRecordsConfig struct {
	RegionID types.String                      `tfsdk:"region_id"`
	Name     types.String                      `tfsdk:"name"`
	ZoneID   types.String                      `tfsdk:"zone_id"`
	ID       types.String                      `tfsdk:"id"`
	PageNo   types.Int32                       `tfsdk:"page_no"`
	PageSize types.Int32                       `tfsdk:"page_size"`
	Records  []CtyunPrivateZoneRecordInfoModel `tfsdk:"records"`
}
