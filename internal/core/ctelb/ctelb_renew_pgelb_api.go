package apis

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/core"
)

// CtelbRenewPgelbApi
/* 保障型负载均衡续订
 */type CtelbRenewPgelbApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbRenewPgelbApi(client *core.CtyunClient) *CtelbRenewPgelbApi {
	return &CtelbRenewPgelbApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/renew-pgelb",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbRenewPgelbApi) Do(ctx context.Context, credential core.Credential, req *CtelbRenewPgelbRequest) (*CtelbRenewPgelbResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbRenewPgelbResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbRenewPgelbRequest struct {
	ClientToken     string `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID        string `json:"regionID,omitempty"`        /*  区域ID  */
	ElbID           string `json:"elbID,omitempty"`           /*  负载均衡 ID  */
	CycleType       string `json:"cycleType,omitempty"`       /*  订购类型：month（包月） / year（包年）  */
	CycleCount      int32  `json:"cycleCount,omitempty"`      /*  订购时长, 当 cycleType = month, 支持续订 1 - 11 个月; 当 cycleType = year, 支持续订 1 - 3 年  */
	PayVoucherPrice string `json:"payVoucherPrice,omitempty"` /*  代金券金额，支持到小数点后两位  */
}

type CtelbRenewPgelbResponse struct {
	StatusCode  int32                             `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbRenewPgelbReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                            `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbRenewPgelbReturnObjResponse struct{}
