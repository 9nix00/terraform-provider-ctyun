package amqp

import (
	"context"
	"net/http"
	ctyunsdk "terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

type AmqpInstanceCreatePostPayOrderApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstanceCreatePostPayOrderApi(client *ctyunsdk.CtyunClient) *AmqpInstanceCreatePostPayOrderApi {
	return &AmqpInstanceCreatePostPayOrderApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v3/instances/createPostPayOrder",
		},
	}
}

func (this *AmqpInstanceCreatePostPayOrderApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstanceCreatePostPayOrderRequest) (res *AmqpInstanceCreatePostPayOrderRequest, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	builder.AddHeader("regionId", req.RegionID)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointName, builder)
	if err != nil {
		return
	}
	res = &AmqpInstanceCreatePostPayOrderRequest{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstanceCreatePostPayOrderRequest struct {
	RegionID        string `json:"regionId"`
	HostType        string `json:"hostType"`
	DiskType        string `json:"diskType"`
	DiskSize        string `json:"diskSize"`
	CpuNum          int    `json:"cpuNum"`
	MemSize         int    `json:"memSize"`
	NodeNum         int    `json:"nodeNum"`
	EngineType      string `json:"engineType"`
	CluterName      string `json:"cluterName"`
	VpcId           string `json:"vpcId"`
	AzInfo          string `json:"azInfo"`
	SubnetId        string `json:"subnetId"`
	SecurityGroupId string `json:"securityGroupId"`
	AutoPay         bool   `json:"autoPay"`
}

type AmqpInstanceCreatePostPayOrderResponse struct {
	ReturnObj  *AmqpInstanceCreatePostPayOrderResponseReturnObj `json:"returnObj"`
	Message    string                                           `json:"message"`
	StatusCode string                                           `json:"statusCode"`
}

type AmqpInstanceCreatePostPayOrderResponseReturnObj struct {
	Data AmqpInstanceCreatePostPayOrderResponseReturnObjData `json:"data"`
}

type AmqpInstanceCreatePostPayOrderResponseReturnObjData struct {
	Submitted  bool   `json:"submitted"`
	NewOrderId string `json:"newOrderId"`
	NewOrderNo string `json:"newOrderNo"`
}
