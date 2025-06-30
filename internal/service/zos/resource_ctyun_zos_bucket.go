package zos

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctzos"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	"terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
	"terraform-provider-ctyun/internal/utils"
	"time"
)

var (
	_ resource.Resource                = &ctyunZosBucket{}
	_ resource.ResourceWithConfigure   = &ctyunZosBucket{}
	_ resource.ResourceWithImportState = &ctyunZosBucket{}
)

type ctyunZosBucket struct {
	meta *common.CtyunMetadata
}

func NewCtyunZosBucket() resource.Resource {
	return &ctyunZosBucket{}
}

func (c *ctyunZosBucket) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_zos_bucket"
}

type CtyunZosBucketConfig struct {
	ID          types.String `tfsdk:"id"`
	RegionID    types.String `tfsdk:"region_id"`
	ACL         types.String `tfsdk:"acl"`
	Bucket      types.String `tfsdk:"bucket"`
	ProjectID   types.String `tfsdk:"project_id"`
	StorageType types.String `tfsdk:"storage_type"`
	IsEncrypted types.Bool   `tfsdk:"is_encrypted"`
	CmkUUID     types.String `tfsdk:"cmk_uuid"`
	AzPolicy    types.String `tfsdk:"az_policy"`
}

func (c *ctyunZosBucket) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026735/10181237**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bucket": schema.StringAttribute{
				Required:    true,
				Description: "桶名称，不可为空。长度3-63个字符内（含）字符只能有大小写字母、数字以及英文句号（.）和中划线（-）。禁止两个英文句号（.）或英文句号（.）中划线（-）相邻。禁止英文句号（.）和中划线（-）作为开头或结尾。",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 63),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z0-9](?:[a-zA-Z0-9]|[.-][a-zA-Z0-9])+$"), "桶名称不符合规则"),
				},
			},
			"acl": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "桶权限，可选值为'private'、'public-read'、'public-read-write'，分别表示私有、公共读、公共读写，默认为'private'",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.ZosAclPrivate, business.ZosAclPublicRead, business.ZosAclPublicReadWrite),
				},
				Default: stringdefault.StaticString(business.ZosAclPrivate),
			},
			"cmk_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "密钥管理服务中创建的密钥ID，使用此参数时，is_encrypted必须为true。当is_encrypted为true但未指定此参数时，会自动创建密钥",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.ConflictsWithEqualString(
						path.MatchRoot("is_encrypted"),
						types.BoolValue(false),
					),
				},
			},
			"is_encrypted": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "加密状态，默认false",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"storage_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "存储类型，可选的值STANDARD、STANDARD_IA、GLACIER，分别表示标准、低频、归档，默认STANDARD",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.ZosStorageTypeStandard, business.ZosStorageTypeStandardIA, business.ZosStorageTypeGlacier),
				},
				Default: stringdefault.StaticString(business.ZosStorageTypeStandard),
			},
			"az_policy": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "az策略，可选值为single-az、multi-az，分别表示单az、多az，默认为single-az",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.ZosAzPolicySingle, business.ZosAzPolicyMulti),
				},
				Default: stringdefault.StaticString(business.ZosAzPolicySingle),
			},
		},
	}
}

func (c *ctyunZosBucket) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunZosBucketConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}
	c.getAndMerge(ctx, &plan)
	time.Sleep(30 * time.Second)
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunZosBucket) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunZosBucketConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found bucket") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunZosBucket) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (c *ctyunZosBucket) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunZosBucketConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
	//response.State.RemoveResource(ctx)
}

func (c *ctyunZosBucket) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [bucket],[regionID]
func (c *ctyunZosBucket) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunZosBucketConfig
	var bucket, regionID string
	err = terraform_extend.Split(request.ID, &bucket, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.Bucket = types.StringValue(bucket)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeCreate 创建前检查
func (c *ctyunZosBucket) checkBeforeCreate(ctx context.Context, plan CtyunZosBucketConfig) (err error) {
	params := &ctzos.ZosGetOssServiceStatusRequest{
		RegionID: plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosGetOssServiceStatusApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	if resp.ReturnObj.State != "true" {
		err = fmt.Errorf("您尚未在该资源池开通对象存储服务，请前往控制台开通后使用")
	}
	return
}

// create 创建
func (c *ctyunZosBucket) create(ctx context.Context, plan CtyunZosBucketConfig) (err error) {
	params := &ctzos.ZosCreateBucketRequest{
		RegionID:    plan.RegionID.ValueString(),
		ACL:         plan.ACL.ValueString(),
		Bucket:      plan.Bucket.ValueString(),
		ProjectID:   plan.ProjectID.ValueString(),
		CmkUUID:     plan.CmkUUID.ValueString(),
		IsEncrypted: plan.IsEncrypted.ValueBoolPointer(),
		StorageType: plan.StorageType.ValueString(),
		AZPolicy:    plan.AzPolicy.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtZosApis.ZosCreateBucketApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

// getAndMerge 从远端查询
func (c *ctyunZosBucket) getAndMerge(ctx context.Context, plan *CtyunZosBucketConfig) (err error) {
	b, err := business.NewZosService(c.meta).GetZosBucketInfo(ctx, plan.Bucket.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		return
	}
	plan.AzPolicy = types.StringValue(b.AZPolicy)
	plan.StorageType = types.StringValue(b.StorageType)
	plan.CmkUUID = utils.SecStringValue(b.CmkUUID)
	if b.CmkUUID != nil {
		plan.IsEncrypted = types.BoolValue(true)
	} else {
		plan.IsEncrypted = types.BoolValue(false)
	}
	plan.ID = plan.Bucket
	// 以下字段无接口查询
	// plan.Acl
	return
}

// delete 删除
func (c *ctyunZosBucket) delete(ctx context.Context, plan CtyunZosBucketConfig) (err error) {
	params := &ctzos.ZosDeleteBucketRequest{
		Bucket:   plan.Bucket.ValueString(),
		RegionID: plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosDeleteBucketApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}
