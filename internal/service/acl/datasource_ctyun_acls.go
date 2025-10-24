package acl

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunAcls{}
	_ datasource.DataSourceWithConfigure = &ctyunAcls{}
)

type ctyunAcls struct {
	meta *common.CtyunMetadata
}

func NewCtyunAcls() datasource.DataSource {
	return &ctyunAcls{}
}
func (c *ctyunAcls) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunAcls) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_acls"
}

func (c *ctyunAcls) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *ctyunAcls) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunAclsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}

	params := &ctvpc.CtvpcListAclRequest{
		RegionID: config.RegionID.ValueString(),
		PageNo:   1,
		PageSize: 10,
	}
	if !config.ID.IsNull() {
		params.AclID = config.ID.ValueStringPointer()
	}
	if !config.ProjectID.IsUnknown() && !config.ProjectID.IsNull() {
		params.ProjectID = config.ProjectID.ValueStringPointer()
	}
	if !config.Name.IsUnknown() && !config.Name.IsNull() {
		params.Name = config.Name.ValueStringPointer()
	}
	if !config.PageNo.IsNull() {
		params.PageNo = config.PageNo.ValueInt32()
	}
	if !config.pageSize.IsNull() {
		params.PageSize = config.pageSize.ValueInt32()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListAclApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("获取acl列表失败，接口返回nil，请联系研发确认问题原因！")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var acls []CtyunAclInfoModel
	for _, aclItem := range resp.ReturnObj.Acls {
		var acl CtyunAclInfoModel
		acl.ID = types.StringValue(*aclItem.AclID)
		acl.Name = types.StringValue(*aclItem.Name)
		acl.Enabled = types.StringValue(*aclItem.Enabled)
		acl.Description = types.StringValue(*aclItem.Description)
		acl.ApplyToPublicLb = types.BoolValue(*aclItem.ApplyToPublicLb)
		acl.CreateTime = types.StringValue(*aclItem.CreatedAt)
		acl.UpdateTime = types.StringValue(*aclItem.UpdatedAt)

		var diags diag.Diagnostics
		acl.OutPolicyID, diags = types.SetValueFrom(ctx, types.StringType, aclItem.OutPolicyID)
		if diags.HasError() {
			err = fmt.Errorf(diags[0].Detail())
			return
		}
		acl.InPolicyID, diags = types.SetValueFrom(ctx, types.StringType, aclItem.InPolicyID)
		if diags.HasError() {
			err = fmt.Errorf(diags[0].Detail())
			return
		}
		acl.SubnetIDs, diags = types.SetValueFrom(ctx, types.StringType, aclItem.SubnetIDs)
		if diags.HasError() {
			err = fmt.Errorf(diags[0].Detail())
			return
		}
		acls = append(acls, acl)
	}
	config.Acls = acls
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunAclInfoModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	ApplyToPublicLb types.Bool   `tfsdk:"applyToPublicLb"`
	VpcID           types.String `tfsdk:"vpcId"`
	Enabled         types.String `tfsdk:"enabled"`
	InPolicyID      types.Set    `tfsdk:"in_policy_id"`
	OutPolicyID     types.Set    `tfsdk:"out_policy_id"`
	CreateTime      types.String `tfsdk:"create_time"`
	UpdateTime      types.String `tfsdk:"update_time"`
	SubnetIDs       types.Set    `tfsdk:"subnet_ids"`
}

type CtyunAclsConfig struct {
	RegionID  types.String        `tfsdk:"region_id"`
	ID        types.String        `tfsdk:"id"`
	ProjectID types.String        `tfsdk:"project_id"`
	Name      types.String        `tfsdk:"name"`
	PageNo    types.Int32         `tfsdk:"page_no"`
	pageSize  types.Int32         `tfsdk:"page_size"`
	Acls      []CtyunAclInfoModel `tfsdk:"acls"`
}
