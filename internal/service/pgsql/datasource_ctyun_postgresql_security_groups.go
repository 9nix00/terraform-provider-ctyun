package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunPgsqlSecurityGroups{}
	_ datasource.DataSourceWithConfigure = &CtyunPgsqlSecurityGroups{}
)

type CtyunPgsqlSecurityGroups struct {
	meta *common.CtyunMetadata
}

func NewCtyunPgsqlSecurityGroups() *CtyunPgsqlSecurityGroups {
	return &CtyunPgsqlSecurityGroups{}
}
func (c *CtyunPgsqlSecurityGroups) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPgsqlSecurityGroups) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_security_groups"
}
func (c *CtyunPgsqlSecurityGroups) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id",
			},
			"inst_id": schema.StringAttribute{
				Required:    true,
				Description: "实例id",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "项目id",
			},
			"security_groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "安全组列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{},
				},
			},
		},
	}

}

func (c *CtyunPgsqlSecurityGroups) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPgsqlSecurityGroupsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	params := &pgsql.PgsqlSecurityGroupListRequest{
		RegionID: regionId,
		InstID:   config.InstID.ValueString(),
	}
	header := &pgsql.PgsqlSecurityGroupListRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlSecurityGroupListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s ", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	returnObj := resp.ReturnObj.Data
	var securityGroups []CtyunPgsqlSecurityGroupInfoModel
	for _, securityGroupItem := range returnObj {
		var securityGroup CtyunPgsqlSecurityGroupInfoModel
		securityGroup.SecurityGroupName = types.StringValue(securityGroupItem.SecurityGroupName)
		securityGroup.SecurityGroupID = types.StringValue(securityGroupItem.ID)
		securityGroup.VmNum = types.StringValue(securityGroupItem.VmNum)
		securityGroup.Origin = types.StringValue(securityGroupItem.Origin)
		securityGroup.VpcName = types.StringValue(securityGroupItem.VpcName)
		securityGroup.CreationTime = types.StringValue(securityGroupItem.CreationTime)
		securityGroup.RegionID = types.StringValue(securityGroupItem.RegionID)
		securityGroup.Description = types.StringValue(securityGroupItem.Description)
		securityGroup.VpcID = types.StringValue(securityGroupItem.VpcID)
		// SecurityGroupRules
		var securityGroupRules []CtyunPgsqlSecurityGroupRuleModel
		for _, ruleItem := range securityGroupItem.SecurityGroupRules {
			var securityGroupRule CtyunPgsqlSecurityGroupRuleModel
			securityGroupRule.Direction = types.StringValue(ruleItem.Direction)
			securityGroupRule.Priority = types.Int32Value(ruleItem.Priority)
			securityGroupRule.EtherType = types.StringValue(ruleItem.EtherType)
			securityGroupRule.Protocol = types.StringValue(ruleItem.Protocol)
			securityGroupRule.DestCidrIp = types.StringValue(ruleItem.DestCidrIp)
			securityGroupRule.Description = types.StringValue(ruleItem.Description)
			securityGroupRule.SecurityGroupId = types.StringValue(ruleItem.SecurityGroupID)
			securityGroupRule.Origin = types.StringValue(ruleItem.Origin)
			securityGroupRule.CreateTime = types.StringValue(ruleItem.CreateTime)
			securityGroupRules = append(securityGroupRules, securityGroupRule)
		}
		securityGroup.SecurityGroupRules = securityGroupRules
		securityGroups = append(securityGroups, securityGroup)
	}
	config.SecurityGroups = securityGroups
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunPgsqlSecurityGroupRuleModel struct {
	Direction       types.String `tfsdk:"direction"`
	Priority        types.Int32  `tfsdk:"priority"`
	EtherType       types.String `tfsdk:"ether_type"`
	Protocol        types.String `tfsdk:"protocol"`
	DestCidrIp      types.String `tfsdk:"dest_cidr_ip"`
	Description     types.String `tfsdk:"description"`
	ID              types.String `tfsdk:"id"`
	SecurityGroupId types.String `tfsdk:"security_group_id"`
	Origin          types.String `tfsdk:"origin"`
	CreateTime      types.String `tfsdk:"create_time"`
}

type CtyunPgsqlSecurityGroupInfoModel struct {
	SecurityGroupName  types.String                       `tfsdk:"security_group_name"`
	SecurityGroupID    types.String                       `tfsdk:"security_group_id"`
	VmNum              types.String                       `tfsdk:"vm_num"`
	Origin             types.String                       `tfsdk:"origin"`
	VpcName            types.String                       `tfsdk:"vpc_name"`
	CreationTime       types.String                       `tfsdk:"creation_time"`
	RegionID           types.String                       `tfsdk:"region_id"`
	Description        types.String                       `tfsdk:"description"`
	VpcID              types.String                       `tfsdk:"vpc_id"`
	SecurityGroupRules []CtyunPgsqlSecurityGroupRuleModel `tfsdk:"security_group_rules"`
}

type CtyunPgsqlSecurityGroupsConfig struct {
	RegionID       types.String                       `tfsdk:"region_id"`
	InstID         types.String                       `tfsdk:"inst_id"`
	ProjectID      types.String                       `tfsdk:"project_id"`
	SecurityGroups []CtyunPgsqlSecurityGroupInfoModel `tfsdk:"security_groups"`
}
