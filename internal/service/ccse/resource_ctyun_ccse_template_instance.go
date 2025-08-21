package ccse

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
//	ccse2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/crs"
//	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
//	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
//	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
//	"github.com/hashicorp/terraform-plugin-framework/path"
//	"github.com/hashicorp/terraform-plugin-framework/resource"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
//	"github.com/hashicorp/terraform-plugin-framework/types"
//	"time"
//)
//
//var (
//	_ resource.Resource                = &ctyunCcseTemplateInstance{}
//	_ resource.ResourceWithConfigure   = &ctyunCcseTemplateInstance{}
//	_ resource.ResourceWithImportState = &ctyunCcseTemplateInstance{}
//)
//
//type ctyunCcseTemplateInstance struct {
//	meta *common.CtyunMetadata
//}
//
//func NewCtyunCcseTemplateInstance() resource.Resource {
//	return &ctyunCcseTemplateInstance{}
//}
//
//func (c *ctyunCcseTemplateInstance) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
//	response.TypeName = request.ProviderTypeName + "_ccse_template_instance"
//}
//
//type CtyunCcseTemplateInstanceConfig struct {
//	ID         types.String `tfsdk:"id"`
//	ClusterID  types.String `tfsdk:"cluster_id"`
//	RegionID   types.String `tfsdk:"region_id"`
//	Namespace  types.String `tfsdk:"namespace"`
//	Name       types.String `tfsdk:"name"`
//	TplName    types.String `tfsdk:"tpl_name"`
//	TplVersion types.String `tfsdk:"tpl_version"`
//	ValuesYaml types.String `tfsdk:"values_yaml"`
//	ValuesJson types.String `tfsdk:"values_json"`
//
//	namespaceID  int64
//	repositoryID int64
//}
//
//func (c *ctyunCcseTemplateInstance) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
//	response.Schema = schema.Schema{
//		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10083472/10102631**`,
//		Attributes: map[string]schema.Attribute{
//			"id": schema.StringAttribute{
//				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
//				Computed:      true,
//				Description:   "ID",
//			},
//			"region_id": schema.StringAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
//				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"cluster_id": schema.StringAttribute{
//				Required:    true,
//				Description: "集群ID",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"namespace": schema.StringAttribute{
//				Required:    true,
//				Description: "集群命名空间",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"name": schema.StringAttribute{
//				Required:    true,
//				Description: "实例名称",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"tpl_name": schema.StringAttribute{
//				Required:    true,
//				Description: "模板名称",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"tpl_version": schema.StringAttribute{
//				Required:    true,
//				Description: "模板版本号，可通过ctyun_ccse_template_market查询",
//			},
//			"values_yaml": schema.StringAttribute{
//				Optional:    true,
//				Description: "模板配置参数(YAML格式)，与values_json二选一。",
//				Validators: []validator.String{
//					stringvalidator.ConflictsWith(path.MatchRoot("values_json")),
//					validator2.Yaml(),
//				},
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"values_json": schema.StringAttribute{
//				Optional:    true,
//				Description: "插件配置参数(JSON格式)，与values_yaml二选一。",
//				Validators: []validator.String{
//					stringvalidator.ConflictsWith(path.MatchRoot("values_yaml")),
//					validator2.Json(),
//				},
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//		},
//	}
//}
//
//func (c *ctyunCcseTemplateInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var plan CtyunCcseTemplateInstanceConfig
//	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//
//	err = c.checkBeforeCreate(ctx, plan)
//	if err != nil {
//		return
//	}
//	// 创建
//	err = c.create(ctx, plan)
//	if err != nil {
//		return
//	}
//	// 创建后检查
//	err = c.checkAfterCreate(ctx, plan)
//	if err != nil {
//		return
//	}
//	// 反查信息
//	err = c.getAndMerge(ctx, &plan)
//	if err != nil {
//		return
//	}
//
//	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
//}
//
//func (c *ctyunCcseTemplateInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var state CtyunCcseTemplateInstanceConfig
//	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// 查询远端
//	err = c.getAndMerge(ctx, &state)
//	if err != nil {
//		return
//	}
//
//	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
//}
//
//func (c *ctyunCcseTemplateInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	// tf文件中的
//	var plan CtyunCcseTemplateInstanceConfig
//	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// state中的
//	var state CtyunCcseTemplateInstanceConfig
//	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// 更新
//	err = c.update(ctx, &plan, &state)
//	if err != nil {
//		return
//	}
//
//	// 查询远端信息
//	err = c.getAndMerge(ctx, &state)
//	if err != nil {
//		return
//	}
//
//	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
//}
//
//func (c *ctyunCcseTemplateInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var state CtyunCcseTemplateInstanceConfig
//	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// 删除
//	err = c.delete(ctx, state)
//	if err != nil {
//		return
//	}
//	//response.State.RemoveResource(ctx)
//}
//
//func (c *ctyunCcseTemplateInstance) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
//	if request.ProviderData == nil {
//		return
//	}
//	meta := request.ProviderData.(*common.CtyunMetadata)
//	c.meta = meta
//}
//
//// 导入命令：terraform import [配置标识].[导入配置名称] [pluginName],[clusterID],[regionID]
//func (c *ctyunCcseTemplateInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var cfg CtyunCcseTemplateInstanceConfig
//	var tplName, clusterID, regionID string
//	err = terraform_extend.Split(request.ID, &tplName, &clusterID, &regionID)
//	if err != nil {
//		return
//	}
//	cfg.TplName = types.StringValue(tplName)
//	cfg.RegionID = types.StringValue(regionID)
//	cfg.ClusterID = types.StringValue(clusterID)
//	// 查询远端
//	err = c.getAndMerge(ctx, &cfg)
//	if err != nil {
//		return
//	}
//	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
//}
//
//// checkBeforeCreate 创建前检查
//func (c *ctyunCcseTemplateInstance) checkBeforeCreate(ctx context.Context, plan CtyunCcseTemplateInstanceConfig) (err error) {
//	exist, err := c.checkExist(ctx, plan)
//	if err != nil {
//		return err
//	}
//	if exist {
//		return fmt.Errorf("模板实例 %s 已经存在", plan.Name)
//	}
//	tpl, err := c.getTpl(ctx, plan)
//	if err != nil {
//		return err
//	}
//
//	return
//}
//
//func (c *ctyunCcseTemplateInstance) getTpl(ctx context.Context, plan CtyunCcseTemplateInstanceConfig) (tpl *crs.CrsListTemplateReturnObjRecordsResponse, err error) {
//	params := &crs.CrsListTemplateRequest{
//		RegionIdHeader: plan.RegionID.ValueString(),
//		RegionIdParam:  plan.RegionID.ValueString(),
//		RepositoryName: plan.TplName.ValueString(),
//	}
//	resp, err := c.meta.Apis.SdkCrsApis.CrsListTemplateApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.ReturnObj == nil {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if len(resp.ReturnObj.Records) == 0 {
//		err = fmt.Errorf("模板 %s 版本 %s 不存在", plan.TplName.ValueString(), plan.TplVersion.ValueString())
//		return
//	}
//	tpl = resp.ReturnObj.Records[0]
//	return
//}
//
//// checkExist 检查模板实例是否存在
//func (c *ctyunCcseTemplateInstance) checkExist(ctx context.Context, plan CtyunCcseTemplateInstanceConfig) (exist bool, err error) {
//	params := &ccse2.CcseHasTemplateInstanceExistedRequest{
//		ClusterId:            plan.ClusterID.ValueString(),
//		NamespaceName:        plan.Namespace.ValueString(),
//		TemplateInstanceName: plan.Name.ValueString(),
//		RegionId:             plan.RegionID.ValueString(),
//	}
//	resp, err := c.meta.Apis.SdkCcseApis.CcseHasTemplateInstanceExistedApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCode {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	}
//	exist = utils.SecBool(resp.ReturnObj)
//	return
//}
//
//// checkAfterCreate 创建后检查
//func (c *ctyunCcseTemplateInstance) checkAfterCreate(ctx context.Context, plan CtyunCcseTemplateInstanceConfig) (err error) {
//	var executeSuccessFlag bool
//	var failedCnt int
//	retryer, _ := business.NewRetryer(time.Second*10, 30)
//	retryer.Start(
//		func(currentTime int) bool {
//			var plugin *ccse2.CcseListTemplateInstanceInstancesReturnObjRecordsResponse
//			plugin, err = c.getByTplName(ctx, plan)
//			if err != nil {
//				if errors.Is(err, common.InvalidReturnObjResultsError) {
//					return true
//				}
//				return false
//			}
//			if plugin.Status == "failed" {
//				failedCnt++
//			}
//			if failedCnt > 1 {
//				err = fmt.Errorf("安装失败，请登录控制台查看原因")
//				return false
//			}
//			if plugin.Status != "deployed" {
//				return true
//			}
//			executeSuccessFlag = true
//			return false
//		})
//	if err != nil {
//		return
//	}
//	if !executeSuccessFlag {
//		err = fmt.Errorf("插件安装超时，请登录控制台查看原因")
//	}
//	return
//}
//
//// create 创建
//func (c *ctyunCcseTemplateInstance) create(ctx context.Context, plan CtyunCcseTemplateInstanceConfig) (err error) {
//	params := &ccse2.CcseDeployTemplateInstanceRequest{
//		ClusterId:     plan.ClusterID.ValueString(),
//		NamespaceName: plan.Namespace.ValueString(),
//		RegionId:      plan.RegionID.ValueString(),
//		ChartName:     plan.TplName.ValueString(),
//		ChartVersion:  plan.TplVersion.ValueString(),
//		CrNamespaceId: plan.namespaceID,
//		InstanceName:  plan.Name.ValueString(),
//		RepositoryId:  plan.repositoryID,
//	}
//	if plan.ValuesJson.ValueString() != "" {
//		params.InstanceValue = plan.ValuesJson.ValueString()
//	} else {
//		params.InstanceValue = plan.ValuesYaml.ValueString()
//	}
//
//	resp, err := c.meta.Apis.SdkCcseApis.CcseDeployTemplateInstanceApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCode {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//
//	return
//}
//
//// getAndMerge 从远端查询
//func (c *ctyunCcseTemplateInstance) getAndMerge(ctx context.Context, plan *CtyunCcseTemplateInstanceConfig) (err error) {
//	plugin, err := c.getByTplName(ctx, *plan)
//	if err != nil {
//		return
//	}
//	plan.Namespace = types.StringValue(plugin.Namespace)
//	plan.TplName = types.StringValue(plugin.TplName)
//	plan.TplVersion = types.StringValue(plugin.TplVersion)
//	plan.ClusterID = types.StringValue(plugin.ClusterId)
//	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", plugin.TplName, plugin.ClusterId, plan.RegionID.ValueString()))
//
//	return
//}
//
//// update 更新
//func (c *ctyunCcseTemplateInstance) update(ctx context.Context, plan, state *CtyunCcseTemplateInstanceConfig) (err error) {
//	err = c.updateTplVersion(ctx, *plan, *state)
//	if err != nil {
//		return
//	}
//	return
//}
//
//// updateTplVersion 更新tpl_version
//func (c *ctyunCcseTemplateInstance) updateTplVersion(ctx context.Context, plan, state CtyunCcseTemplateInstanceConfig) (err error) {
//	if plan.TplVersion.Equal(state.TplVersion) {
//		return
//	}
//	params := &ccse2.CcseUpgradeTemplateInstanceInstanceRequest{
//		ClusterId:  plan.ClusterID.ValueString(),
//		RegionId:   plan.RegionID.ValueString(),
//		TplName:    plan.TplName.ValueString(),
//		TplVersion: plan.TplVersion.ValueString(),
//		Values:     state.ValuesYaml.ValueString(),
//		ValuesJson: state.ValuesJson.ValueString(),
//	}
//
//	resp, err := c.meta.Apis.SdkCcseApis.CcseUpgradeTemplateInstanceInstanceApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCode {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//
//	return c.checkAfterTplVersion(ctx, plan)
//}
//
//// checkAfterTplVersion 升级插件版本后检查
//func (c *ctyunCcseTemplateInstance) checkAfterTplVersion(ctx context.Context, plan CtyunCcseTemplateInstanceConfig) (err error) {
//	var executeSuccessFlag bool
//	var failedCnt int
//	retryer, _ := business.NewRetryer(time.Second*10, 30)
//	retryer.Start(
//		func(currentTime int) bool {
//			var plugin *ccse2.CcseListTemplateInstanceInstancesReturnObjRecordsResponse
//			plugin, err = c.getByTplName(ctx, plan)
//			if err != nil {
//				return false
//			}
//
//			if plugin.Status == "failed" {
//				failedCnt++
//			}
//			if failedCnt > 1 {
//				err = fmt.Errorf("升级插件版本失败")
//				return false
//			}
//
//			if plugin.Status != "deployed" || plugin.TplVersion != plan.TplVersion.ValueString() {
//				return true
//			}
//			executeSuccessFlag = true
//			return false
//		})
//	if err != nil {
//		return
//	}
//	if !executeSuccessFlag {
//		err = fmt.Errorf("插件升级超时")
//	}
//	return
//}
//
//// delete 删除
//func (c *ctyunCcseTemplateInstance) delete(ctx context.Context, plan CtyunCcseTemplateInstanceConfig) (err error) {
//	params := &ccse2.CcseDeleteTemplateInstanceInstanceRequest{
//		ClusterId:    plan.ClusterID.ValueString(),
//		RegionId:     plan.RegionID.ValueString(),
//		InstanceName: plan.TplName.ValueString(),
//		Namespace:    plan.Namespace.ValueString(),
//	}
//	resp, err := c.meta.Apis.SdkCcseApis.CcseDeleteTemplateInstanceInstanceApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCode {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	}
//	return
//}
//
//// checkAfterDelete 删除后检查
//func (c *ctyunCcseTemplateInstance) checkAfterDelete(ctx context.Context, plan CtyunCcseTemplateInstanceConfig) (err error) {
//	var executeSuccessFlag bool
//	retryer, _ := business.NewRetryer(time.Second*10, 30)
//	retryer.Start(
//		func(currentTime int) bool {
//			var plugin *ccse2.CcseListTemplateInstanceInstancesReturnObjRecordsResponse
//			plugin, err = c.getByTplName(ctx, plan)
//			if err != nil {
//				if errors.Is(err, common.InvalidReturnObjResultsError) {
//					err = nil
//					executeSuccessFlag = true
//				}
//				return false
//			}
//			if plugin.Status != "uninstalled" {
//				return true
//			}
//			executeSuccessFlag = true
//			return false
//		})
//	if err != nil {
//		return
//	}
//	if !executeSuccessFlag {
//		err = fmt.Errorf("插件卸载超时")
//	}
//	return
//}
//
//// getByTplName通过插件名称查询
//func (c *ctyunCcseTemplateInstance) getByTplName(ctx context.Context, plan CtyunCcseTemplateInstanceConfig) (plugin *ccse2.CcseListTemplateInstanceInstancesReturnObjRecordsResponse, err error) {
//	params := &ccse2.CcseListTemplateInstanceInstancesRequest{
//		ClusterId: plan.ClusterID.ValueString(),
//		RegionId:  plan.RegionID.ValueString(),
//		TplName:   plan.TplName.ValueString(),
//	}
//	resp, err := c.meta.Apis.SdkCcseApis.CcseListTemplateInstanceInstancesApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCode {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if len(resp.ReturnObj.Records) == 0 {
//		err = common.InvalidReturnObjResultsError
//		return
//	}
//	plugin = resp.ReturnObj.Records[0]
//	return
//}
