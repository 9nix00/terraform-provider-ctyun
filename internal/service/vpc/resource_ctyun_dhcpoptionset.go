package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

func NewCtyunDhcpOptionSet() resource.Resource {
	return &ctyunDhcpOptionSet{}
}

type ctyunDhcpOptionSet struct {
	meta *common.CtyunMetadata
}

func (c *ctyunDhcpOptionSet) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_dhcpoptionset"
}

func (c *ctyunDhcpOptionSet) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `DHCP选项集资源，用于管理DHCP选项集配置`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "DHCP选项集ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "集合名，支持拉丁字母、中文、数字，下划线，连字符，必须以中文/英文字母开头，不能以数字、_和-、http:/https:开头，长度2-32",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "描述信息，支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&**()_-+= <>?:\"{},./;'[**\r\n\r\n**]·~！@#￥%……&**（） —— -+={}《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128",
			},
			"domain_name": schema.StringAttribute{
				Required:    true,
				Description: "整个域名的总长度不能超过 255 个字符，每个子域名（包括顶级域名）的长度不能超过 63 个字符，域名中的字符集包括大写字母、小写字母、数字和连字符（减号），连字符不能位于域名的开头",
			},
			"dns_list": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "服务ip地址列表，最多只能4个IP地址",
			},
		},
	}
}

func (c *ctyunDhcpOptionSet) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunDhcpOptionSet) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunDhcpOptionSetConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.create(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunDhcpOptionSet) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunDhcpOptionSetConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (c *ctyunDhcpOptionSet) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan CtyunDhcpOptionSetConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.update(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunDhcpOptionSet) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunDhcpOptionSetConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.delete(ctx, &state)
	if err != nil {
		return
	}
}

// create 创建dhcpoptionset
func (c *ctyunDhcpOptionSet) create(ctx context.Context, plan *CtyunDhcpOptionSetConfig) (err error) {
	// 准备请求参数
	req := &ctvpc.CtvpcDhcpoptionsetscreateRequest{
		RegionID:   plan.RegionId.ValueString(),
		Name:       plan.Name.ValueString(),
		DomainName: plan.DomainName.ValueString(),
	}

	// 设置可选参数
	if !plan.Description.IsNull() {
		description := plan.Description.ValueString()
		req.Description = &description
	}

	// 设置DNS列表
	var dnsList []string
	for _, dns := range plan.DnsList {
		dnsList = append(dnsList, dns.ValueString())
	}
	req.DnsList = dnsList

	// 调用API创建DHCP选项集
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcpoptionsetscreateApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 设置资源ID
	plan.Id = types.StringValue(*resp.ReturnObj.DhcpOptionSetsID)

	return nil
}

// getAndMerge 查询dhcpoptionset并合并状态
func (c *ctyunDhcpOptionSet) getAndMerge(ctx context.Context, state *CtyunDhcpOptionSetConfig) (err error) {
	// 调用API获取DHCP选项集详情
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcpoptionsetsShowApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDhcpoptionsetsShowRequest{
		RegionID:         state.RegionId.ValueString(),
		DhcpOptionSetsID: state.Id.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 更新状态
	state.Id = utils.SecStringValue(resp.ReturnObj.DhcpOptionSetsID)
	state.Name = utils.SecStringValue(resp.ReturnObj.Name)
	state.Description = utils.SecStringValue(resp.ReturnObj.Description)

	if len(resp.ReturnObj.DomainName) > 0 && resp.ReturnObj.DomainName[0] != nil {
		state.DomainName = types.StringValue(*resp.ReturnObj.DomainName[0])
	}

	// 更新DNS列表
	state.DnsList = []types.String{}
	for _, dns := range resp.ReturnObj.DnsList {
		if dns != nil {
			state.DnsList = append(state.DnsList, types.StringValue(*dns))
		}
	}

	return
}

// update 更新dhcpoptionset
func (c *ctyunDhcpOptionSet) update(ctx context.Context, plan *CtyunDhcpOptionSetConfig) (err error) {
	// 准备请求参数
	req := &ctvpc.CtvpcUpdatedhcpoptionsetsRequest{
		RegionID:         plan.RegionId.ValueString(),
		DhcpOptionSetsID: plan.Id.ValueString(),
	}

	// 设置可选参数
	if !plan.Name.IsNull() {
		name := plan.Name.ValueString()
		req.Name = &name
	}

	if !plan.Description.IsNull() {
		description := plan.Description.ValueString()
		req.Description = &description
	}

	if !plan.DomainName.IsNull() {
		domainName := plan.DomainName.ValueString()
		req.DomainName = &domainName
	}

	// 设置DNS列表
	var dnsList []*string
	for _, dns := range plan.DnsList {
		dnsValue := dns.ValueString()
		dnsList = append(dnsList, &dnsValue)
	}
	req.DnsList = dnsList

	// 调用API更新DHCP选项集
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdatedhcpoptionsetsApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return nil
}

// delete 删除dhcpoptionset
func (c *ctyunDhcpOptionSet) delete(ctx context.Context, state *CtyunDhcpOptionSetConfig) (err error) {
	// 调用API删除DHCP选项集
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcpoptionsetsdeleteApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDhcpoptionsetsdeleteRequest{
		RegionID:         state.RegionId.ValueString(),
		DhcpOptionSetsID: state.Id.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return nil
}

type CtyunDhcpOptionSetConfig struct {
	Id          types.String   `tfsdk:"id"`
	RegionId    types.String   `tfsdk:"region_id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	DomainName  types.String   `tfsdk:"domain_name"`
	DnsList     []types.String `tfsdk:"dns_list"`
}
