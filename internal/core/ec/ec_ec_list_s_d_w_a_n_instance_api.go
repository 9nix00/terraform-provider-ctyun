package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EcEcListSDWANInstanceApi
/* 查询sdwan网络实例 */
type EcEcListSDWANInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcListSDWANInstanceApi(client *core.CtyunClient) *EcEcListSDWANInstanceApi {
	return &EcEcListSDWANInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/sdwan-instance/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcListSDWANInstanceApi) Do(ctx context.Context, credential core.Credential, req *EcEcListSDWANInstanceRequest) (*EcEcListSDWANInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("ecID", req.EcID)
	if req.SdwanID != nil && *req.SdwanID != "" {
		ctReq.AddParam("sdwanID", *req.SdwanID)
	}
	if req.CgwID != nil && *req.CgwID != "" {
		ctReq.AddParam("cgwID", *req.CgwID)
	}
	if req.InstanceID != nil && *req.InstanceID != "" {
		ctReq.AddParam("instanceID", *req.InstanceID)
	}
	if req.QueryContent != nil && *req.QueryContent != "" {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	if *req.IsAuth != 0 {
		ctReq.AddParam("isAuth", strconv.FormatInt(int64(*req.IsAuth), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcListSDWANInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcListSDWANInstanceRequest struct {
	EcID         string  `json:"ecID"`                   /*  云间高速ID  */
	SdwanID      *string `json:"sdwanID,omitempty"`      /*  sdwan ID  */
	CgwID        *string `json:"cgwID,omitempty"`        /*  云网关ID  */
	InstanceID   *string `json:"instanceID,omitempty"`   /*  网络实例ID  */
	QueryContent *string `json:"queryContent,omitempty"` /*  模糊查询，支持sdwanID，网络实例ID，云网关ID三个属性  */
	IsAuth       *int32  `json:"isAuth,omitempty"`       /*  是否是跨账号实例<br/>取值包括：<br/>0:本账号<br/>1:跨账号  */
}

type EcEcListSDWANInstanceResponse struct {
	StatusCode  *int32  `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *string `json:"returnObj"`   /*  返回参数  */
}
