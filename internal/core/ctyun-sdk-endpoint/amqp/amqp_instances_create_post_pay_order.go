package amqp

import (
	"context"
	"net/http"
	ctyunsdk "terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

type AmqpInstancesCreatePostPayOrderApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstancesCreatePostPayOrderApi(client *ctyunsdk.CtyunClient) *AmqpInstancesCreatePostPayOrderApi {
	return &AmqpInstancesCreatePostPayOrderApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v3/instances/createPostPayOrder",
		},
	}
}

func (this *AmqpInstancesCreatePostPayOrderApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstancesCreatePostPayOrderRequest) (res *AmqpInstancesCreatePostPayOrderResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	builder.AddHeader("regionId", req.RegionId)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointName, builder)
	if err != nil {
		return
	}
	res = &AmqpInstancesCreatePostPayOrderResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstancesCreatePostPayOrderRequest struct {
	RegionId        string `json:"regionId"`
	HostType        string `json:"hostType"`
	DiskType        string `json:"diskType"`
	DiskSize        string `json:"diskSize"`
	CpuNum          int32  `json:"cpuNum"`
	MemSize         int32  `json:"memSize"`
	NodeNum         int32  `json:"nodeNum"`
	EngineType      string `json:"engineType"`
	ClusterName     string `json:"clusterName"`
	VpcId           string `json:"vpcId"`
	AzInfo          string `json:"azInfo"`
	SubnetId        string `json:"subnetId"`
	SecurityGroupId string `json:"securityGroupId"`
	AutoPay         bool   `json:"autoPay"`
}

type AmqpInstancesCreatePostPayOrderResponse struct {
	ReturnObj  *AmqpInstancesCreatePostPayOrderResponseReturnObj `json:"returnObj"`
	Message    string                                            `json:"message"`
	StatusCode string                                            `json:"statusCode"`
}

type AmqpInstancesCreatePostPayOrderResponseReturnObj struct {
	Data AmqpInstancesCreatePostPayOrderResponseReturnObjData `json:"data"`
}

type AmqpInstancesCreatePostPayOrderResponseReturnObjData struct {
	Submitted  bool   `json:"submitted"`
	NewOrderId string `json:"newOrderId"`
	NewOrderNo string `json:"newOrderNo"`
}
