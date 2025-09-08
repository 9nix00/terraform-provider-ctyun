package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpQueueCreateV3Api
/* 创建队列v3
 */type AmqpQueueCreateV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpQueueCreateV3Api(client *core.CtyunClient) *AmqpQueueCreateV3Api {
	return &AmqpQueueCreateV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/queue/create",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpQueueCreateV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpQueueCreateV3Request) (*AmqpQueueCreateV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpQueueCreateV3Request
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpQueueCreateV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpQueueCreateV3Request struct {
	RegionId    string `json:"regionId,omitempty"`   /*  资源池id  */
	ProdInstId  string `json:"prodInstId,omitempty"` /*  实例ID  */
	Vhost       string `json:"vhost,omitempty"`      /*  vhost名称  */
	Name        string `json:"name,omitempty"`       /*  队列名称  */
	Durable     *bool  `json:"durable"`              /*  是否持久化  */
	Auto_delete *bool  `json:"auto_delete"`          /*  是否自动删除  */
}

type AmqpQueueCreateV3Response struct {
	ReturnObj  *AmqpQueueCreateV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                              `json:"message"`    /*  描述状态  */
	StatusCode string                              `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                              `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpQueueCreateV3ReturnObjResponse struct {
	Data string `json:"data"` /*  返回数据  */
}
