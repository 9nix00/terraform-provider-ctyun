package acl

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CtyunPrefix struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunPrefix) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_acl"
}

func (c *CtyunPrefix) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunPrefix() resource.Resource {
	return &CtyunPrefix{}
}

func (c *CtyunPrefix) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunPrefixConfig
	var ID, regionId, projectId, vpcId, name string
	err = terraform_extend.Split(request.ID, &ID, &regionId, &projectId, &vpcId, &name)
	if err != nil {
		return
	}

	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *CtyunPrefix) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunPrefix) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunPrefixConfig
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

func (c *CtyunPrefix) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunPrefix) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunPrefix) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunPrefix) getAndMerge(ctx context.Context, config *CtyunPrefixConfig) error {
	params := &ctvpc.CtvpcPrefixlistShowRequest{
		RegionID:     config.RegionID.ValueString(),
		PrefixListID: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcPrefixlistShowApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("获取prefix失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}

}

func (c *CtyunPrefix) create(ctx context.Context, config *CtyunPrefixConfig) error {
	params := &ctvpc.CtvpcPrefixlistCreateRequest{
		RegionID:        config.RegionID.ValueString(),
		Name:            config.Name.ValueString(),
		Limit:           config.Limit.ValueInt32(),
		AddressType:     business.PrefixAddressTypeMap[config.AddressType.ValueString()],
		PrefixListRules: nil,
	}
	var prefixListRules []*ctvpc.CtvpcPrefixlistCreatePrefixListRulesRequest
	var prefixes []CtyunPrefixModel
	diags := config.PrefixListRules.ElementsAs(ctx, &prefixes, true)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	for _, rule := range prefixes {
		var prefix *ctvpc.CtvpcPrefixlistCreatePrefixListRulesRequest
		prefix.Cidr = rule.Cidr.ValueString()
		prefix.Description = rule.Description.ValueStringPointer()
		prefixListRules = append(prefixListRules, prefix)
	}
	params.PrefixListRules = prefixListRules
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcPrefixlistCreateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建prefix失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	config.ID = types.StringValue(*resp.ReturnObj.PrefixListID)
	return nil
}

type CtyunPrefixModel struct {
	Cidr        types.String `ctyun:"cidr"`
	Description types.String `ctyun:"description"`
}
type CtyunPrefixConfig struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	RegionID        types.String `tfsdk:"region_id"`
	Limit           types.Int32  `tfsdk:"limit"`
	AddressType     types.String `tfsdk:"address_type"`
	PrefixListRules types.List   `tfsdk:"prefix_list_rules"`
}
