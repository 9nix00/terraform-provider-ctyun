package iam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCtyunIamUsers() datasource.DataSource {
	return &CtyunIamUsers{}
}

type CtyunIamUsers struct {
	meta       *common.CtyunMetadata
	iamService *business.IamService
}

func (c *CtyunIamUsers) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_users"
}

func (c *CtyunIamUsers) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10345725/10355805`,
		Attributes: map[string]schema.Attribute{
			"page_size": schema.Int64Attribute{
				Required:    true,
				Description: "每页显示数量",
				Validators: []validator.Int64{
					int64validator.Between(1, 1000),
				},
			},
			"page_no": schema.Int64Attribute{
				Required:    true,
				Description: "当前页码",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"users": schema.SetNestedAttribute{
				Computed:    true,
				Description: "用户列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"user_id": schema.StringAttribute{
							Computed:    true,
							Description: "用户ID",
						},
						"email": schema.StringAttribute{
							Computed:    true,
							Description: "用户邮箱",
						},
						"account_id": schema.BoolAttribute{
							Computed:    true,
							Description: "账户ID",
						},
						"phone": schema.StringAttribute{
							Computed:    true,
							Description: "手机号",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "用户名",
						},
						"is_root": schema.BoolAttribute{
							Computed:    true,
							Description: "是否主用户",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "用户状态",
						},
						"groups": schema.SetNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:    true,
										Description: "用户组id",
									},
									"name": schema.StringAttribute{
										Computed:    true,
										Description: "用户组名称",
									},
									"description": schema.StringAttribute{
										Computed:    true,
										Description: "用户组信息",
									},
								},
							},
							Computed:    true,
							Description: "关联的用户组列表",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunIamUsers) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunIamUsersConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	aks, err := c.iamService.QueryList(ctx, config.UserID.ValueString())
	if err != nil {
		return
	}
	var sk string
	config.List = []CtyunIamUsersModel{}
	for _, a := range aks {
		if config.AK.ValueString() == "" || config.AK.ValueString() == utils.SecString(a.AccessKey) {
			item := CtyunIamUsersModel{
				AK:         utils.SecStringValue(a.AccessKey),
				Enabled:    types.BoolValue(map[string]bool{business.Enabled: true, business.Disabled: false}[utils.SecString(a.Status)]),
				CreateTime: types.StringValue(utils.FromUnixToUTC(a.CreatedTime)),
			}
			sk, err = c.iamService.DecryptSK(utils.SecString(a.SecretKey), c.meta.SdkCredential.GetAccessKey())
			if err != nil {
				return
			}
			item.SK = types.StringValue(sk)
			config.List = append(config.List, item)
		}
	}
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *CtyunIamUsers) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.iamService = business.NewIamService(meta)
}

type CtyunIamUsersConfig struct {
	UserID types.String         `tfsdk:"user_id"`
	AK     types.String         `tfsdk:"ak"`
	List   []CtyunIamUsersModel `tfsdk:"list"`
}

type CtyunIamUsersModel struct {
	AK         types.String `tfsdk:"ak"`
	SK         types.String `tfsdk:"sk"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	CreateTime types.String `tfsdk:"create_time"`
}
