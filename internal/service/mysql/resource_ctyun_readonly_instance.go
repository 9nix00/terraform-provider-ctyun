package mysql

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CtyunMysqlReadOnlyInstance{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlReadOnlyInstance{}
	_ resource.ResourceWithImportState = &CtyunMysqlReadOnlyInstance{}
)

type CtyunMysqlReadOnlyInstance struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMysqlReadOnlyInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunMysqlReadOnlyInstanceConfig
	var IDStr, regionId, projectId, description, engine, name string
	err = terraform_extend.Split(request.ID, &IDStr, &regionId, &projectId, &description, &engine, &name)
	if err != nil {
		return
	}
	//ID, err := strconv.ParseInt(IDStr, 10, 64)
	//if err != nil {
	//	fmt.Println("id转换失败，输入有误:", err)
	//	return
	//}
	//cfg.ID = types.Int64Value(ID)
	//cfg.RegionID = types.StringValue(regionId)
	//cfg.ProjectID = types.StringValue(projectId)
	//cfg.Description = types.StringValue(description)
	//cfg.Engine = types.StringValue(engine)
	//cfg.Name = types.StringValue(name)
	//err = c.getAndMergeMysqlParameterTemplate(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *CtyunMysqlReadOnlyInstance) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_readonly_instance"
}
func NewCtyunMysqlReadOnlyInstance() resource.Resource {
	return &CtyunMysqlReadOnlyInstance{}
}

func (c *CtyunMysqlReadOnlyInstance) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMysqlReadOnlyInstance) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunMysqlReadOnlyInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunMysqlReadOnlyInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunMysqlReadOnlyInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunMysqlReadOnlyInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

type CtyunMysqlReadOnlyInstanceConfig struct {
	ParentProdInstId types.String `tfsdk:"parent_prod_inst_id"`
	CycleType        types.String `tfsdk:"cycle_type"`        // 计费模式： 支持on_demand和month
	CycleCount       types.Int32  `tfsdk:"cycle_count"`       // 购买时长：单位月（范围：1-12，24，36）
	FlavorName       types.String `tfsdk:"flavor_name"`       // 规格名称
	RegionID         types.String `tfsdk:"region_id"`         // 资源池id
	StorageType      types.String `tfsdk:"storage_type"`      // 存储类型
	StorageSpace     types.Int32  `tfsdk:"storage_space"`     // 存储空间, 磁盘大小100G-2T 步长10G
	Name             types.String `tfsdk:"name"`              // 只读实例名称
	Password         types.String `tfsdk:"password"`          // 只读实例密码
	SubnetID         types.String `tfsdk:"subnet_id"`         // 子网ID
	VpcID            types.String `tfsdk:"vpc_id"`            // vpc id
	SecurityGroupID  types.String `tfsdk:"security_group_id"` // 安全组

}
