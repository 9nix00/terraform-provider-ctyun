package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseAttachClusterNodesApi
/* 调用该接口纳管节点。
 */type CcseAttachClusterNodesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseAttachClusterNodesApi(client *core.CtyunClient) *CcseAttachClusterNodesApi {
	return &CcseAttachClusterNodesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodes/attach",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseAttachClusterNodesApi) Do(ctx context.Context, credential core.Credential, req *CcseAttachClusterNodesRequest) (*CcseAttachClusterNodesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseAttachClusterNodesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseAttachClusterNodesRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	Instances []*CcseAttachClusterNodesInstancesRequest `json:"instances"`           /*  云主机ECS 信息  */
	VmType    string                                    `json:"vmType,omitempty"`    /*  默认ecs，值为：弹性云主机则是：ecs ；物理机则是：ebm  */
	Runtime   string                                    `json:"runtime,omitempty"`   /*  容器运行时，目前仅支持 CONTAINERD  */
	ImageUuid string                                    `json:"imageUuid,omitempty"` /*  镜像ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004475" target="_blank">节点规格和节点镜像</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&amp;api=4765&amp;data=89" target="_blank">创建私有镜像（云主机系统盘）</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&amp;api=5230&amp;data=89" target="_blank">创建私有镜像（云主机数据盘）</a>
	注：同一镜像名称在不同资源池的镜像ID是不同的，调用前需确认镜像ID是否归属当前资源池  */
	LoginType string `json:"loginType,omitempty"` /*  云主机密码登录类型：目前只支持 PASSWORD  */
	Password  string `json:"password,omitempty"`  /*  用户密码，满足以下规则：
	长度在8～30个字符
	必须包含大写字母、小写字母、数字以及特殊符号中的三项
	特殊符号可选：()`-!@#$%^&*_-+=｜{}[]:;'<>,.?/且不能以斜线号 / 开头
	不能包含3个及以上连续字符
	Linux镜像不能包含镜像用户名（root）、用户名的倒序（toor）、用户名大小写变化（如RoOt、rOot等）
	Windows镜像不能包含镜像用户名（Administrator）、用户名大小写变化（adminiSTrator等）  */
	Labels                   *CcseAttachClusterNodesLabelsRequest   `json:"labels"`                             /*  标签  */
	Taints                   []*CcseAttachClusterNodesTaintsRequest `json:"taints"`                             /*  节点污点，格式为 [{\"key\":\"{key}\",\"value\":\"{value}\",\"effect\":\"{effect}\"}]，上述的{key}、{value}、{effect}替换为所需字段。effect枚举包括NoSchedule、PreferNoSchedule、NoExecute  */
	VisibilityPostHostScript string                                 `json:"visibilityPostHostScript,omitempty"` /*  部署后执行自定义脚本 （输入的值需要经过Base64编码，方法如下：echo -n "待编码内容" \  */
	VisibilityHostScript     string                                 `json:"visibilityHostScript,omitempty"`     /*  部署前执行自定义脚本（输入的值需要经过Base64编码，方法如下：echo -n "待编码内容" \  */
}

type CcseAttachClusterNodesInstancesRequest struct {
	InstanceId string `json:"instanceId,omitempty"` /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs" target="_blank">弹性云主机</a>了解云主机的相关信息
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&amp;api=8309&amp;data=87">查询云主机列表</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&amp;api=8281&amp;data=87" target="_blank">创建一台按量付费或包年包月的云主机</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&amp;api=8282&amp;data=87">批量创建按量付费或包年包月云主机</a>

	物理机 instanceUUID，获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=6941&data=97&isNormal=1&vid=91">批量查询物理机</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=6942&data=97&isNormal=1&vid=91">物理机创建</a>  */
	AzName string `json:"azName,omitempty"` /*  可用区名称，纳管是物理机，此项必填，可用区信息可用区名称获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87&vid=81" target="_blank">资源池可用区查询</a>  */
}

type CcseAttachClusterNodesLabelsRequest struct{}

type CcseAttachClusterNodesTaintsRequest struct {
	Key    string `json:"key,omitempty"`    /*  键  */
	Value  string `json:"value,omitempty"`  /*  值  */
	Effect string `json:"effect,omitempty"` /*  策略  */
}

type CcseAttachClusterNodesResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  响应状态码  */
	RequestId  string `json:"requestId,omitempty"`  /*  请求ID  */
	Message    string `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *bool  `json:"returnObj"`            /*  响应对象  */
	Error      string `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}
