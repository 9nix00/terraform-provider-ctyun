package acl

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CtyunAclRules struct {
	meta *common.CtyunMetadata
}

func NewCtyunAclRules() datasource.DataSource {
	return &CtyunAclRules{}
}
func (c *CtyunAclRules) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunAclRules) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_acl_rules"
}

func (c *CtyunAclRules) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunAclRules) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunAclRulesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}

	rules, err := c.getRuleList(ctx, &config)
	if err != nil {
		return
	}
	rule := rules[0]
	config.Name = types.StringValue(*rule.Name)
	config.Description = types.StringValue(*rule.Description)
	config.VpcID = types.StringValue(*rule.VpcID)
	config.Enabled = types.StringValue(*rule.Enabled)
	config.InPolicyID = rule.InPolicyID
	config.OutPolicyID = rule.OutPolicyID
	config.SubnetIDs = rule.SubnetIDs
	var inRules []CtyunAclRuleModel
	for _, ruleItem := range rule.InRules {
		var inRule CtyunAclRuleModel
		inRule.AclRuleID = types.StringValue(*ruleItem.AclRuleID)
		inRule.Direction = types.StringValue(*ruleItem.Direction)
		inRule.Priority = types.Int32Value(ruleItem.Priority)
		inRule.IpVersion = types.StringValue(*ruleItem.IpVersion)
		inRule.DestinationPort = types.StringValue(*ruleItem.DestinationPort)
		inRule.SourcePort = types.StringValue(*ruleItem.SourcePort)
		inRule.SourceIpAddress = types.StringValue(*ruleItem.SourceIpAddress)
		inRule.DestinationIpAddress = types.StringValue(*ruleItem.DestinationIpAddress)
		inRule.Action = types.StringValue(*ruleItem.Action)
		inRule.Enabled = types.StringValue(*ruleItem.Enabled)
		inRule.Description = types.StringValue(*ruleItem.Description)
		inRules = append(inRules, inRule)
	}
	config.InRules = inRules

	var outRules []CtyunAclRuleModel
	for _, ruleItem := range rule.OutRules {
		var outRule CtyunAclRuleModel
		outRule.AclRuleID = types.StringValue(*ruleItem.AclRuleID)
		outRule.Direction = types.StringValue(*ruleItem.Direction)
		outRule.Priority = types.Int32Value(ruleItem.Priority)
		outRule.IpVersion = types.StringValue(*ruleItem.IpVersion)
		outRule.DestinationPort = types.StringValue(*ruleItem.DestinationPort)
		outRule.SourcePort = types.StringValue(*ruleItem.SourcePort)
		outRule.SourceIpAddress = types.StringValue(*ruleItem.SourceIpAddress)
		outRule.DestinationIpAddress = types.StringValue(*ruleItem.DestinationIpAddress)
		outRule.Action = types.StringValue(*ruleItem.Action)
		outRule.Enabled = types.StringValue(*ruleItem.Enabled)
		outRule.Description = types.StringValue(*ruleItem.Description)
		outRules = append(outRules, outRule)
	}
	config.OutRules = outRules
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *CtyunAclRules) getRuleList(ctx context.Context, config *CtyunAclRulesConfig) ([]*ctvpc.CtvpcListAclRuleReturnObjResponse, error) {
	params := &ctvpc.CtvpcListAclRuleRequest{
		RegionID: config.RegionID.ValueString(),
		AclID:    config.AclID.ValueString(),
	}
	if !config.ProjectID.IsUnknown() && !config.ProjectID.IsNull() {
		params.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListAclRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取acl规则列表失败（acl id=%s），接口返回nil，请联系研发确认问题原因！", config.AclID.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp.ReturnObj, nil
}

type CtyunAclRuleModel struct {
	AclRuleID            types.String `tfsdk:"acl_rule_id"`
	Direction            types.String `tfsdk:"direction"`
	Priority             types.Int32  `tfsdk:"priority"`
	Protocol             types.String `tfsdk:"protocol"`
	IpVersion            types.String `tfsdk:"ip_version"`
	DestinationPort      types.String `tfsdk:"destination_port"`
	SourcePort           types.String `tfsdk:"source_port"`
	SourceIpAddress      types.String `tfsdk:"source_ip_address"`
	DestinationIpAddress types.String `tfsdk:"destination_ip_address"`
	Action               types.String `tfsdk:"action"`
	Enabled              types.String `tfsdk:"enabled"`
	Description          types.String `tfsdk:"description"`
}

type CtyunAclRulesConfig struct {
	RegionID    types.String        `tfsdk:"region_id"`
	AclID       types.String        `tfsdk:"acl_id"`
	ProjectID   types.String        `tfsdk:"project_id"`
	Name        types.String        `tfsdk:"name"`
	Description types.String        `tfsdk:"description"`
	VpcID       types.String        `tfsdk:"vpc_id"`
	Enabled     types.String        `tfsdk:"enabled"`
	InPolicyID  []string            `tfsdk:"in_policy_id"`
	OutPolicyID []string            `tfsdk:"out_policy_id"`
	InRules     []CtyunAclRuleModel `tfsdk:"in_rules"`
	OutRules    []CtyunAclRuleModel `tfsdk:"out_rules"`
	SubnetIDs   []string            `tfsdk:"subnet_ids"`
}
