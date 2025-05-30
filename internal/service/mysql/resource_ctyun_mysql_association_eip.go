package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
	"time"
)

var (
	_ resource.Resource                = &CtyunMysqlAssociationEip{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlAssociationEip{}
	_ resource.ResourceWithImportState = &CtyunMysqlAssociationEip{}
)

type CtyunMysqlAssociationEip struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMysqlAssociationEip) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_association_eip"
}
func NewCtyunMysqlAssociationEip() resource.Resource {
	return &CtyunMysqlAssociationEip{}
}

func (c *CtyunMysqlAssociationEip) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"eip_id": schema.StringAttribute{
				Required:    true,
				Description: "弹性id",
			},
			"eip": schema.StringAttribute{
				Required:    true,
				Description: "弹性ip",
			},
			"inst_id": schema.StringAttribute{
				Required:    true,
				Description: "实例id",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "项目id",
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
			"eip_status": schema.Int32Attribute{
				Computed:    true,
				Description: " 弹性ip状态 0->unbind，1->bind,2->binding",
				Validators: []validator.Int32{
					int32validator.Between(0, 2),
				},
			},
		},
	}
}

func (c *CtyunMysqlAssociationEip) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunAssociationEipConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 实例绑定弹性IP
	err = c.MysqlBindEip(ctx, &plan)
	if err != nil {
		return
	}
	// 轮询查看绑定状态
	err = c.BindLoop(ctx, &plan, business.EipStatusBind, business.EipStatusUnbind)
	if err != nil {
		return
	}
	// 查询实例详情，确认是否绑定成功
	err = c.getAndMergeBindEip(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlAssociationEip) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunAssociationEipConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeBindEip(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "is not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunMysqlAssociationEip) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	//暂无可更新内容
}

func (c *CtyunMysqlAssociationEip) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunAssociationEipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	unbindParams := &mysql.TeledbUnbindEipRequest{
		EipID:  state.EipID.ValueString(),
		Eip:    state.Eip.ValueString(),
		InstID: state.InstID.ValueString(),
	}
	unbindHeader := &mysql.TeledbUnbindEipRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		unbindHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbUnbindEipApi.Do(ctx, c.meta.Credential, unbindParams, unbindHeader)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 轮询确定解绑成功
	err = c.BindLoop(ctx, &state, business.EipStatusUnbind, business.EipStatusBind)
	if err != nil {
		return
	}
	return
}

func (c *CtyunMysqlAssociationEip) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	// todo
}

func (c *CtyunMysqlAssociationEip) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMysqlAssociationEip) MysqlBindEip(ctx context.Context, config *CtyunAssociationEipConfig) (err error) {
	params := &mysql.TeledbBindEipRequest{
		EipID:  config.EipID.ValueString(),
		Eip:    config.Eip.ValueString(),
		InstID: config.InstID.ValueString(),
	}
	header := &mysql.TeledbBindEipRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbBindEipApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s ", resp.Message)
		return
	}
	return
}

func (c *CtyunMysqlAssociationEip) getAndMergeBindEip(ctx context.Context, config *CtyunAssociationEipConfig) (err error) {
	detailParams := &mysql.TeledbQueryDetailRequest{
		OuterProdInstId: config.InstID.ValueString(),
	}
	header := &mysql.TeledbQueryDetailRequestHeaders{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, header)
	if err != nil {
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	returnObj := resp.ReturnObj
	config.Eip = types.StringValue(returnObj.EIP)
	config.EipStatus = types.Int32Value(returnObj.EIPStatus)
	config.ProjectID = types.StringValue(returnObj.ProjectId)
	return
}

func (c *CtyunMysqlAssociationEip) BindLoop(ctx context.Context, config *CtyunAssociationEipConfig, finalStatus int32, initStatus int32, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*10, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			detailParams := &mysql.TeledbQueryDetailRequest{
				OuterProdInstId: config.InstID.ValueString(),
			}
			header := &mysql.TeledbQueryDetailRequestHeaders{
				InstID:   config.InstID.ValueString(),
				RegionID: config.RegionID.ValueString(),
			}
			if config.ProjectID.ValueString() != "" {
				header.ProjectID = config.ProjectID.ValueStringPointer()
			}

			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, header)
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
			status := resp.ReturnObj.EIPStatus
			switch status {
			case business.EipStatusBinding:
				return true
			case initStatus:
				return true
			case finalStatus:
				return false
			default:
				err = errors.New("mysql绑定解绑eip时出现异常状态：" + fmt.Sprintf("%d", status))
				return false
			}
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，eip仍未绑定/解绑成功！")
	}
	return
}

type CtyunAssociationEipConfig struct {
	EipID     types.String `tfsdk:"eip_id"`     //弹性id
	Eip       types.String `tfsdk:"eip"`        //弹性ip
	InstID    types.String `tfsdk:"inst_id"`    //实例id
	ProjectID types.String `tfsdk:"project_id"` //项目id
	RegionID  types.String `tfsdk:"region_id"`  //区域Id
	EipStatus types.Int32  `tfsdk:"eip_status"` //弹性ip状态 0->unbind，1->bind,2->binding
}
