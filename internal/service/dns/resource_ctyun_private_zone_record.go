package dns

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type CtyunPrivateZoneRecord struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunPrivateZoneRecord) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_zone_record"
}

func (c *CtyunPrivateZoneRecord) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunPrivateZoneRecord() resource.Resource {
	return &CtyunPrivateZoneRecord{}
}

func (c *CtyunPrivateZoneRecord) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunPrivateZoneRecordConfig
	var ID, regionId, projectId, vpcId, name string
	err = terraform_extend.Split(request.ID, &ID, &regionId, &projectId, &vpcId, &name)
	if err != nil {
		return
	}
	config.ID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionId)
	config.Name = types.StringValue(name)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunPrivateZoneRecord) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id,如果不填这默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "内网DNS id",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "内网DNS记录类型，支持: A / CNAME / MX / AAAA / TXT, 大小写不敏感",
				Validators: []validator.String{
					stringvalidator.OneOf("A", "CNAME", "MX", "AAAA", "TXT"),
				},
			},
			"value_list": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "最多同时支持 8 个",
			},
			"ttl": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "zone ttl, ",
				Default:     int32default.StaticInt32(300),
				Validators: []validator.Int32{
					int32validator.Between(300, 2147483647),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "DNS记录集的 name 长度",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "DNS记录集描述",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "内网DNS记录id",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间",
			},
		},
	}
}

func (c *CtyunPrivateZoneRecord) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunPrivateZoneRecordConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPrivateZoneRecord) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPrivateZoneRecordConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "不存在") {
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

func (c *CtyunPrivateZoneRecord) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunPrivateZoneRecordConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunPrivateZoneRecordConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPrivateZoneRecord) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunPrivateZoneRecordConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunPrivateZoneRecord) create(ctx context.Context, config *CtyunPrivateZoneRecordConfig) error {
	params := &ctvpc.CtvpcCreatePrivateZoneRecordRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
		ZoneID:      config.ZoneID.ValueString(),
		RawType:     config.Type.ValueString(),
		TTL:         config.TTL.ValueInt32(),
		Name:        config.Name.ValueStringPointer(),
	}
	var values []string
	diags := config.ValueList.ElementsAs(ctx, &values, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	params.ValueList = values
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreatePrivateZoneRecordApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建内网DNS记录失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	config.ID = types.StringValue(*resp.ReturnObj.ZoneRecordID)
	return nil
}

func (c *CtyunPrivateZoneRecord) getAndMerge(ctx context.Context, config *CtyunPrivateZoneRecordConfig) error {
	params := &ctvpc.CtvpcShowPrivateZoneRecordRequest{
		RegionID:     config.RegionID.ValueString(),
		ZoneRecordID: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowPrivateZoneRecordApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("获取内网DNS记录(id=%s)详情失败，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	config.Name = types.StringValue(*resp.ReturnObj.Name)
	config.Description = types.StringValue(*resp.ReturnObj.Description)
	config.ZoneID = types.StringValue(*resp.ReturnObj.ZoneRecordID)
	config.Type = types.StringValue(*resp.ReturnObj.RawType)
	config.TTL = types.Int32Value(resp.ReturnObj.TTL)
	config.CreatedTime = types.StringValue(*resp.ReturnObj.CreatedAt)
	config.UpdatedTime = types.StringValue(*resp.ReturnObj.UpdatedAt)
	// 处理value
	var values []string
	for _, value := range resp.ReturnObj.Value {
		values = append(values, *value)
	}
	valueTmp, diags := types.SetValueFrom(ctx, types.StringType, values)
	if diags.HasError() {
		err = fmt.Errorf(diags[0].Detail())
		return err
	}
	config.ValueList = valueTmp
	return nil
}

func (c *CtyunPrivateZoneRecord) update(ctx context.Context, state *CtyunPrivateZoneRecordConfig, plan *CtyunPrivateZoneRecordConfig) error {
	params := &ctvpc.CtvpcUpdatePrivateZoneRecordAttributeRequest{
		RegionID:     state.RegionID.ValueString(),
		ZoneRecordID: state.ID.ValueString(),
	}
	var values []string
	diags := plan.ValueList.ElementsAs(ctx, &values, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	params.ValueList = values
	if !plan.TTL.Equal(state.TTL) {
		params.TTL = plan.TTL.ValueInt32()
	}
	if !plan.Description.Equal(state.Description) {
		params.Description = plan.Description.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdatePrivateZoneRecordAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新内网DNS记录(id=%s)属性失败，接口返回nil，请联系研发确认问题原因！", state.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

func (c *CtyunPrivateZoneRecord) delete(ctx context.Context, config CtyunPrivateZoneRecordConfig) error {
	params := &ctvpc.CtvpcDeletePrivateZoneRecordRequest{
		RegionID:     config.RegionID.ValueString(),
		ZoneRecordID: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeletePrivateZoneRecordApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除内网DNS记录(id=%s)失败，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

type CtyunPrivateZoneRecordConfig struct {
	RegionID    types.String `tfsdk:"region_id"`
	ZoneID      types.String `tfsdk:"zone_id"`
	Type        types.String `tfsdk:"type"`
	ValueList   types.Set    `tfsdk:"value_list"`
	TTL         types.Int32  `tfsdk:"ttl"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ID          types.String `tfsdk:"id"`
	CreatedTime types.String `tfsdk:"create_time"`
	UpdatedTime types.String `tfsdk:"update_time"`
}
