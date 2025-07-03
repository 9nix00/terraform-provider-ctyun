package ccse

import (
	"context"
	"net/http"
	"strconv"
	"terraform-provider-ctyun/internal/core/core"
)

// CcseListPluginInstancesApi
/* 调用该接口可查询插件实例列表。
 */type CcseListPluginInstancesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseListPluginInstancesApi(client *core.CtyunClient) *CcseListPluginInstancesApi {
	return &CcseListPluginInstancesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/plugininstance/page",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseListPluginInstancesApi) Do(ctx context.Context, credential core.Credential, req *CcseListPluginInstancesRequest) (*CcseListPluginInstancesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.PageNow != 0 {
		ctReq.AddParam("pageNow", strconv.FormatInt(int64(req.PageNow), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Namespace != "" {
		ctReq.AddParam("namespace", req.Namespace)
	}
	if req.ChartName != "" {
		ctReq.AddParam("chartName", req.ChartName)
	}
	if req.PluginName != "" {
		ctReq.AddParam("instanceName", req.PluginName)
	}
	if req.ListAll != nil {
		ctReq.AddParam("listAll", strconv.FormatBool(*req.ListAll))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseListPluginInstancesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseListPluginInstancesRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	PageNow    int32  /*  当前页码  */
	PageSize   int32  /*  每页条数  */
	Namespace  string /*  命名空间名称  */
	ChartName  string /*  chart名称  */
	PluginName string /*  插件名称  */
	ListAll    *bool  /*  是否列举全部  */
}

type CcseListPluginInstancesResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  状态码  */
	RequestId  string                                    `json:"requestId,omitempty"`  /*  请求id  */
	Message    string                                    `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseListPluginInstancesReturnObjResponse `json:"returnObj"`            /*  分页查询结果对象  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码  */
}

type CcseListPluginInstancesReturnObjResponse struct {
	Records []*CcseListPluginInstancesReturnObjRecordsResponse `json:"records"`           /*  记录列表  */
	Total   int32                                              `json:"total,omitempty"`   /*  总条数  */
	Size    int32                                              `json:"size,omitempty"`    /*  每页条数  */
	Current int32                                              `json:"current,omitempty"` /*  当前页码  */
	Pages   int32                                              `json:"pages,omitempty"`   /*  总页数  */
}

type CcseListPluginInstancesReturnObjRecordsResponse struct {
	Name           string      `json:"name,omitempty"`           /*  实例名称。  */
	Revision       string      `json:"revision,omitempty"`       /*  版本  */
	Namespace      string      `json:"namespace,omitempty"`      /*  命名空间  */
	Updated        string      `json:"updated,omitempty"`        /*  更新时间  */
	Status         string      `json:"status,omitempty"`         /*  状态  */
	Chart          string      `json:"chart,omitempty"`          /*  Chart名称和版本  */
	AppVersion     string      `json:"appVersion,omitempty"`     /*  版本  */
	ClusterId      string      `json:"clusterId,omitempty"`      /*  集群ID  */
	RepositoryId   interface{} `json:"repositoryId,omitempty"`   /*  仓库ID  */
	ChartName      string      `json:"chartName,omitempty"`      /*  Chart名称  */
	ChartVersion   string      `json:"chartVersion,omitempty"`   /*  Chart版本  */
	ChartUrl       string      `json:"chartUrl,omitempty"`       /*  Chart地址  */
	Icon           string      `json:"icon,omitempty"`           /*  icon地址  */
	Description    string      `json:"description,omitempty"`    /*  描述  */
	NextVersion    string      `json:"nextVersion,omitempty"`    /*  下一版本  */
	TemplateName   string      `json:"templateName,omitempty"`   /*  模板类型  */
	KubeConfigPath string      `json:"kubeConfigPath,omitempty"` /*  kubeConfig路径  */
}
