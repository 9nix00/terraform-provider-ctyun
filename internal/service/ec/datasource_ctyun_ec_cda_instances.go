package ec

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunEcCdaInstances{}
	_ datasource.DataSourceWithConfigure = &CtyunEcCdaInstances{}
)

func NewCtyunEcCdaInstances() datasource.DataSource {
	return &CtyunEcCdaInstances{}
}

type CtyunEcCdaInstances struct {
	meta *common.CtyunMetadata
}

type CtyunEcCdaInstancesConfig struct {
	ID         types.String             `tfsdk:"id"`
	EcID       types.String             `tfsdk:"ec_id"`
	CgwID      types.String             `tfsdk:"cgw_id"`
	CdaID      types.String             `tfsdk:"cda_id"`
	InstanceID types.String             `tfsdk:"instance_id"`
	Status     types.String             `tfsdk:"status"`
	Instances  []CtyunEcCdaInstanceInfo `tfsdk:"instances"`
}

type CtyunEcCdaInstanceInfo struct {
	EcID              types.String `tfsdk:"ec_id"`
	CgwID             types.String `tfsdk:"cgw_id"`
	CgwName           types.String `tfsdk:"cgw_name"`
	DcID              types.String `tfsdk:"dc_id"`
	CdaID             types.String `tfsdk:"cda_id"`
	CdaName           types.String `tfsdk:"cda_name"`
	InstanceID        types.String `tfsdk:"instance_id"`
	DefaultRtbID      types.String `tfsdk:"default_rtb_id"`
	DefaultRtbName    types.String `tfsdk:"default_rtb_name"`
	CIDR              types.List   `tfsdk:"cidr"`
	V6CIDR            types.List   `tfsdk:"v6_cidr"`
	CreateDate        types.String `tfsdk:"create_date"`
	Status            types.String `tfsdk:"status"`
	Weights           types.Int64  `tfsdk:"weights"`
	RedundantType     types.Int64  `tfsdk:"redundant_type"`
	RedundantInstUUID types.String `tfsdk:"redundant_inst_uuid"`
	RedundantInstName types.String `tfsdk:"redundant_inst_name"`
	RedundantInstType types.String `tfsdk:"redundant_inst_type"`
	RedundantInstID   types.String `tfsdk:"redundant_inst_id"`
	RouteLearn        types.Int64  `tfsdk:"route_learn"`
	RouteSync         types.Int64  `tfsdk:"route_sync"`
}

func (c *CtyunEcCdaInstances) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ec_cda_instances"
}

func (c *CtyunEcCdaInstances) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `**CDA网络实例列表数据源**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "数据源ID",
			},
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cgw_id": schema.StringAttribute{
				Optional:    true,
				Description: "云网关ID",
			},
			"cda_id": schema.StringAttribute{
				Optional:    true,
				Description: "云专线ID",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "网络实例ID",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "状态描述，取值包括：creating:加载中, running:已连接, removing:卸载中, flushing:路由待更新, error:失败",
			},
			"instances": schema.ListNestedAttribute{
				Computed:    true,
				Description: "返回的CDA网络实例列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ec_id": schema.StringAttribute{
							Computed:    true,
							Description: "云间高速ID",
						},
						"cgw_id": schema.StringAttribute{
							Computed:    true,
							Description: "云网关ID",
						},
						"cgw_name": schema.StringAttribute{
							Computed:    true,
							Description: "云网关名称",
						},
						"dc_id": schema.StringAttribute{
							Computed:    true,
							Description: "云网关资源池ID",
						},
						"cda_id": schema.StringAttribute{
							Computed:    true,
							Description: "云专线ID",
						},
						"cda_name": schema.StringAttribute{
							Computed:    true,
							Description: "云专线名称",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "网络实例的ID",
						},
						"default_rtb_id": schema.StringAttribute{
							Computed:    true,
							Description: "默认路由表ID",
						},
						"default_rtb_name": schema.StringAttribute{
							Computed:    true,
							Description: "默认路由表名称",
						},
						"cidr": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "v4 CIDR列表",
						},
						"v6_cidr": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "v6 CIDR列表",
						},
						"create_date": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态描述，取值包括：creating:加载中, running:已连接, removing:卸载中, flushing:路由待更新, error:失败",
						},
						"weights": schema.Int64Attribute{
							Computed:    true,
							Description: "权重，专线默认50",
						},
						"redundant_type": schema.Int64Attribute{
							Computed:    true,
							Description: "冗余类型，取值包括：1：主备, 2：负载, 0：无",
						},
						"redundant_inst_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "冗余负载实例的UUID",
						},
						"redundant_inst_name": schema.StringAttribute{
							Computed:    true,
							Description: "冗余负载实例的名称",
						},
						"redundant_inst_type": schema.StringAttribute{
							Computed:    true,
							Description: "冗余负载实例的类型，取值包括：2: 云专线, 3: sdwan, 4: vpn",
						},
						"redundant_inst_id": schema.StringAttribute{
							Computed:    true,
							Description: "冗余负载侧ID，即当前为cdaID, sdwanID，vpnID",
						},
						"route_learn": schema.Int64Attribute{
							Computed:    true,
							Description: "路由学习开关，开启后云网关自动学习网络实例路由，取值范围: 1:学习, 0:不学习, 默认学习",
						},
						"route_sync": schema.Int64Attribute{
							Computed:    true,
							Description: "路由同步开关，开启后云网关路由自动同步到网络实例，取值范围: 1:同步, 0:不同步, 默认同步",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunEcCdaInstances) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunEcCdaInstances) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunEcCdaInstancesConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 构造请求参数
	request := &ec.EcEcListCDANetworkRequest{
		EcID: plan.EcID.ValueString(),
	}

	if !plan.CgwID.IsNull() {
		request.CgwID = plan.CgwID.ValueString()
	}

	if !plan.CdaID.IsNull() {
		request.CdaID = plan.CdaID.ValueString()
	}

	if !plan.InstanceID.IsNull() {
		request.InstanceID = plan.InstanceID.ValueString()
	}

	if !plan.Status.IsNull() {
		request.Status = plan.Status.ValueString()
	}

	response, err := c.meta.Apis.SdkEcApis.EcEcListCDANetworkApi.Do(ctx, c.meta.SdkCredential, request)
	if err != nil {
		return
	} else if response == nil {
		err = fmt.Errorf("API return error. StatusCode is nil")
		return
	} else if *response.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *response.Message)
		return
	} else if response.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回的JSON字符串
	var results []map[string]interface{}
	if response.ReturnObj != nil && *response.ReturnObj != "" {
		err = json.Unmarshal([]byte(*response.ReturnObj), &results)
		if err != nil {
			err = fmt.Errorf("failed to unmarshal returnObj: %v", err)
			return
		}
	}

	// 处理返回结果
	var instances []CtyunEcCdaInstanceInfo
	for _, result := range results {
		instance := CtyunEcCdaInstanceInfo{}

		if val, ok := result["ecID"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.EcID = types.StringValue(str)
			}
		}

		if val, ok := result["cgwID"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.CgwID = types.StringValue(str)
			}
		}

		if val, ok := result["cgwName"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.CgwName = types.StringValue(str)
			}
		}

		if val, ok := result["dcID"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.DcID = types.StringValue(str)
			}
		}

		if val, ok := result["cdaID"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.CdaID = types.StringValue(str)
			}
		}

		if val, ok := result["cdaName"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.CdaName = types.StringValue(str)
			}
		}

		if val, ok := result["instanceID"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.InstanceID = types.StringValue(str)
			}
		}

		if val, ok := result["defaultRtbID"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.DefaultRtbID = types.StringValue(str)
			}
		}

		if val, ok := result["defaultRtbName"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.DefaultRtbName = types.StringValue(str)
			}
		}

		if val, ok := result["createDate"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.CreateDate = types.StringValue(str)
			}
		}

		if val, ok := result["status"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.Status = types.StringValue(str)
			}
		}

		if val, ok := result["weights"]; ok && val != nil {
			if f, ok := val.(float64); ok {
				instance.Weights = types.Int64Value(int64(f))
			}
		}

		if val, ok := result["redundantType"]; ok && val != nil {
			if f, ok := val.(float64); ok {
				instance.RedundantType = types.Int64Value(int64(f))
			}
		}

		if val, ok := result["redundantInstUUID"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.RedundantInstUUID = types.StringValue(str)
			}
		}

		if val, ok := result["redundantInstName"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.RedundantInstName = types.StringValue(str)
			}
		}

		if val, ok := result["redundantInstType"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.RedundantInstType = types.StringValue(str)
			}
		}

		if val, ok := result["redundantInstID"]; ok && val != nil {
			if str, ok := val.(string); ok {
				instance.RedundantInstID = types.StringValue(str)
			}
		}

		if val, ok := result["routeLearn"]; ok && val != nil {
			if f, ok := val.(float64); ok {
				instance.RouteLearn = types.Int64Value(int64(f))
			}
		}

		if val, ok := result["routeSync"]; ok && val != nil {
			if f, ok := val.(float64); ok {
				instance.RouteSync = types.Int64Value(int64(f))
			}
		}

		// 处理CIDR列表
		if val, ok := result["CIDR"]; ok && val != nil {
			if cidrList, ok := val.([]interface{}); ok {
				cidrStrings := make([]string, len(cidrList))
				for i, cidr := range cidrList {
					if cidrStr, ok := cidr.(string); ok {
						cidrStrings[i] = cidrStr
					}
				}
				instance.CIDR, _ = types.ListValueFrom(ctx, types.StringType, cidrStrings)
			}
		}

		// 处理V6 CIDR列表
		if val, ok := result["v6CIDR"]; ok && val != nil {
			if v6CidrList, ok := val.([]interface{}); ok {
				v6CidrStrings := make([]string, len(v6CidrList))
				for i, cidr := range v6CidrList {
					if cidrStr, ok := cidr.(string); ok {
						v6CidrStrings[i] = cidrStr
					}
				}
				instance.V6CIDR, _ = types.ListValueFrom(ctx, types.StringType, v6CidrStrings)
			}
		}

		instances = append(instances, instance)
	}

	plan.Instances = instances
	plan.ID = types.StringValue("cda_instances")

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}
