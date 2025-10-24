package acl

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CtyunAclRule struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunAclRule) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_acl_rule"
}

func (c *CtyunAclRule) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunAclRule() resource.Resource {
	return &CtyunAclRule{}
}

func (c *CtyunAclRule) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunAclRuleConfig
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

func (c *CtyunAclRule) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunAclRule) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunAclRuleConfig
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

func (c *CtyunAclRule) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunAclRuleConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunAclRule) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunAclRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunAclRuleConfig
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

func (c *CtyunAclRule) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunAclRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunAclRule) create(ctx context.Context, config *CtyunAclRuleConfig) error {
	params := &ctvpc.CtvpcCreateAclRuleRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
		AclID:       config.AclID.ValueString(),
	}

	var rules []*ctvpc.CtvpcCreateAclRuleRulesRequest
	var rule *ctvpc.CtvpcCreateAclRuleRulesRequest
	rule.Direction = config.Direction.ValueString()
	rule.Priority = config.Priority.ValueInt32()
	rule.Protocol = config.Protocol.ValueString()
	rule.IpVersion = config.IpVersion.ValueString()
	if !config.DestinationPort.IsNull() {
		rule.DestinationPort = config.DestinationPort.ValueStringPointer()
	}
	if !config.SourcePort.IsNull() {
		rule.SourcePort = config.SourcePort.ValueStringPointer()
	}
	rule.SourceIpAddress = config.SourceIpAddress.ValueString()
	rule.DestinationIpAddress = config.DestinationIpAddress.ValueString()
	rule.Action = config.Action.ValueString()
	rule.Enabled = config.Enabled.ValueString()
	if !config.Description.IsNull() {
		rule.Description = config.Description.ValueString()
	}
	rules = append(rules, rule)
	params.Rules = rules
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateAclRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建acl规则失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	err = c.getRuleID(ctx, config)
	if err != nil {
		return err

	}
	return nil
}

func (c *CtyunAclRule) getAndMerge(ctx context.Context, config *CtyunAclRuleConfig) error {
	ingressDetail, egressDetail, err := c.getRuleDetail(ctx, config)
	if err != nil {
		return err
	}
	if ingressDetail == nil {
		config.Direction = types.StringValue(*egressDetail.Direction)
		config.Priority = types.Int32Value(egressDetail.Priority)
		config.Protocol = types.StringValue(*egressDetail.Protocol)
		config.IpVersion = types.StringValue(*egressDetail.IpVersion)
		config.DestinationPort = types.StringValue(*egressDetail.DestinationPort)
		config.SourcePort = types.StringValue(*egressDetail.SourcePort)
		config.SourceIpAddress = types.StringValue(*egressDetail.SourceIpAddress)
		config.DestinationIpAddress = types.StringValue(*egressDetail.DestinationIpAddress)
		config.Action = types.StringValue(*egressDetail.Action)
		config.Enabled = types.StringValue(*egressDetail.Enabled)
		config.Description = types.StringValue(*egressDetail.Description)
	} else {
		config.Direction = types.StringValue(*ingressDetail.Direction)
		config.Priority = types.Int32Value(ingressDetail.Priority)
		config.Protocol = types.StringValue(*ingressDetail.Protocol)
		config.IpVersion = types.StringValue(*ingressDetail.IpVersion)
		config.DestinationPort = types.StringValue(*ingressDetail.DestinationPort)
		config.SourcePort = types.StringValue(*ingressDetail.SourcePort)
		config.SourceIpAddress = types.StringValue(*ingressDetail.SourceIpAddress)
		config.DestinationIpAddress = types.StringValue(*ingressDetail.DestinationIpAddress)
		config.Action = types.StringValue(*ingressDetail.Action)
		config.Enabled = types.StringValue(*ingressDetail.Enabled)
		config.Description = types.StringValue(*ingressDetail.Description)
	}
	return nil
}

func (c *CtyunAclRule) getRuleID(ctx context.Context, config *CtyunAclRuleConfig) error {
	// 通过获取列表获取，rule id
	ruleList, err := c.getRuleList(ctx, config)
	if err != nil {
		return err
	}
	// 若刚刚创建的规则为入规则，遍历入规则列表
	if config.Direction.ValueString() == business.AclDirectionIngress {
		ingressRuleList := ruleList[0].InRules
		for _, ingressRule := range ingressRuleList {
			same := c.ingressCheckSame(ingressRule, config)
			if same {
				config.ID = types.StringValue(*ingressRule.AclRuleID)
			}
		}
	} else if config.Direction.ValueString() == business.AclDirectionEgress {
		egressRuleList := ruleList[0].OutRules
		for _, egressRule := range egressRuleList {
			same := c.egressCheckSame(egressRule, config)
			if same {
				config.ID = types.StringValue(*egressRule.AclRuleID)
			}
		}
	} else {
		err = fmt.Errorf("direction 取值有误！当前值为%s", config.Direction.ValueString())
		return err
	}
	return nil
}

func (c *CtyunAclRule) getRuleList(ctx context.Context, config *CtyunAclRuleConfig) ([]*ctvpc.CtvpcListAclRuleReturnObjResponse, error) {
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

func (c *CtyunAclRule) ingressCheckSame(rule *ctvpc.CtvpcListAclRuleReturnObjInRulesResponse, config *CtyunAclRuleConfig) bool {
	if *rule.Direction != config.Direction.ValueString() ||
		rule.Priority != config.Priority.ValueInt32() ||
		*rule.Protocol != config.Protocol.ValueString() ||
		*rule.IpVersion != config.IpVersion.ValueString() ||
		*rule.DestinationPort != config.DestinationPort.ValueString() ||
		*rule.SourcePort != config.SourcePort.ValueString() ||
		*rule.SourceIpAddress != config.SourceIpAddress.ValueString() ||
		*rule.DestinationIpAddress != config.DestinationIpAddress.ValueString() ||
		*rule.Action != config.Action.ValueString() ||
		*rule.Enabled != config.Enabled.ValueString() {
		return false
	}
	return true
}

func (c *CtyunAclRule) egressCheckSame(rule *ctvpc.CtvpcListAclRuleReturnObjOutRulesResponse, config *CtyunAclRuleConfig) bool {
	if *rule.Direction != config.Direction.ValueString() ||
		rule.Priority != config.Priority.ValueInt32() ||
		*rule.Protocol != config.Protocol.ValueString() ||
		*rule.IpVersion != config.IpVersion.ValueString() ||
		*rule.DestinationPort != config.DestinationPort.ValueString() ||
		*rule.SourcePort != config.SourcePort.ValueString() ||
		*rule.SourceIpAddress != config.SourceIpAddress.ValueString() ||
		*rule.DestinationIpAddress != config.DestinationIpAddress.ValueString() ||
		*rule.Action != config.Action.ValueString() ||
		*rule.Enabled != config.Enabled.ValueString() {
		return false
	}
	return true
}

func (c *CtyunAclRule) getRuleDetail(ctx context.Context, config *CtyunAclRuleConfig) (*ctvpc.CtvpcListAclRuleReturnObjInRulesResponse, *ctvpc.CtvpcListAclRuleReturnObjOutRulesResponse, error) {
	ruleList, err := c.getRuleList(ctx, config)
	if err != nil {
		return nil, nil, err
	}
	if config.Direction.ValueString() == business.AclDirectionIngress {
		ingressList := ruleList[0].InRules
		for _, ingressRule := range ingressList {
			if *ingressRule.AclRuleID == config.ID.ValueString() {
				return ingressRule, nil, nil
			}
		}
	} else if config.Direction.ValueString() == business.AclDirectionEgress {
		egressList := ruleList[0].OutRules
		for _, egressRule := range egressList {
			if *egressRule.AclRuleID == config.ID.ValueString() {
				return nil, egressRule, nil
			}
		}
	} else {
		err = fmt.Errorf("direction 取值有误！当前值为%s", config.Direction.ValueString())
		return nil, nil, err
	}
	return nil, nil, nil
}

func (c *CtyunAclRule) update(ctx context.Context, state *CtyunAclRuleConfig, plan *CtyunAclRuleConfig) error {
	uuidStr := uuid.NewString()
	params := &ctvpc.CtvpcUpdateAclRuleAttributeRequest{
		ClientToken: &uuidStr,
		RegionID:    state.RegionID.ValueString(),
		AclID:       state.AclID.ValueString(),
		Rules:       nil,
	}
	var rules []*ctvpc.CtvpcUpdateAclRuleAttributeRulesRequest
	var rule *ctvpc.CtvpcUpdateAclRuleAttributeRulesRequest
	rule.Direction = plan.Direction.ValueString()
	rule.Priority = plan.Priority.ValueInt32()
	rule.Protocol = plan.Protocol.ValueString()
	rule.IpVersion = plan.IpVersion.ValueString()
	if !plan.DestinationPort.IsNull() {
		rule.DestinationPort = plan.DestinationPort.ValueStringPointer()
	}
	if !plan.SourcePort.IsNull() {
		rule.SourcePort = plan.SourcePort.ValueStringPointer()
	}
	rule.SourceIpAddress = plan.SourceIpAddress.ValueString()
	rule.DestinationIpAddress = plan.DestinationIpAddress.ValueString()
	rule.Action = plan.Action.ValueString()
	rule.Enabled = plan.Enabled.ValueString()
	if !plan.Description.IsNull() {
		rule.Description = plan.Description.ValueStringPointer()
	}
	rules = append(rules, rule)
	params.Rules = rules
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdateAclRuleAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新acl规则失败（acl_id =%s,acl_rule_id =%s），接口返回nil，请联系研发确认问题原因！", state.AclID.ValueString(), plan.ID)
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil
}

func (c *CtyunAclRule) delete(ctx context.Context, config CtyunAclRuleConfig) error {
	aclRuleIDs := []string{config.ID.ValueString()}
	params := &ctvpc.CtvpcDeleteAclRuleRequest{
		ClientToken:   uuid.NewString(),
		RegionID:      config.RegionID.ValueString(),
		AclID:         config.ID.ValueString(),
		AclRuleIDList: aclRuleIDs,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteAclRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除acl规则失败（acl_id =%s,acl_rule_id =%s），接口返回nil，请联系研发确认问题原因！", config.AclID.ValueString(), config.ID)
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil
}

type CtyunAclRuleConfig struct {
	RegionID             types.String `tfsdk:"region_id"`
	ProjectID            types.String `tfsdk:"project_id"`
	AclID                types.String `tfsdk:"acl_id"`
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
	ID                   types.String `tfsdk:"id"`
}
