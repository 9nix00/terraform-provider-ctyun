package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcPacketListPacketApi
/* 查询带宽包 */
type EcEcPacketListPacketApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcPacketListPacketApi(client *core.CtyunClient) *EcEcPacketListPacketApi {
	return &EcEcPacketListPacketApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/bandwidth-packet/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcPacketListPacketApi) Do(ctx context.Context, credential core.Credential, req *EcEcPacketListPacketRequest) (*EcEcPacketListPacketResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("ecID", req.EcID)
	if req.PacketID != nil && *req.PacketID != "" {
		ctReq.AddParam("packetID", *req.PacketID)
	}
	if req.ResourceID != nil && *req.ResourceID != "" {
		ctReq.AddParam("resourceID", *req.ResourceID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcPacketListPacketResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcPacketListPacketRequest struct {
	EcID       string  `json:"ecID"`                 /*  云间高速实例ID  */
	PacketID   *string `json:"packetID,omitempty"`   /*  带宽包ID  */
	ResourceID *string `json:"resourceID,omitempty"` /*  带宽包 resource ID  */
}

type EcEcPacketListPacketResponse struct {
	StatusCode  *int32  `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *string `json:"returnObj"`   /*  返回参数  */
}
