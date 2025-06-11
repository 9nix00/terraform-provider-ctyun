package amqp

import (
	"context"
	"net/http"
	ctyunsdk "terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

type AmqpInstanceCreatePrePayOrderApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstanceCreatePrePayOrderApi(client *ctyunsdk.CtyunClient) *AmqpInstanceCreatePrePayOrderApi {
	return &AmqpInstanceCreatePrePayOrderApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v3/instances/createPrePayOrder",
		},
	}
}

func (this *AmqpInstanceCreatePrePayOrderApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstanceCreatePrePayOrderRequest) (res *AmqpInstanceCreatePrePayOrderRequest, err error) {
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
	res = &AmqpInstanceCreatePrePayOrderRequest{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstanceCreatePrePayOrderRequest struct {
	RegionID        string `json:"regionId"`
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
	CycleCnt        int32  `json:"cycleCnt"`
}

type AmqpInstanceCreatePrePayOrderResponse struct {
	ReturnObj  *AmqpInstanceCreatePrePayOrderResponseReturnObj `json:"returnObj"`
	Message    string                                          `json:"message"`
	StatusCode string                                          `json:"statusCode"`
}

type AmqpInstanceCreatePrePayOrderResponseReturnObj struct {
	Data AmqpInstanceCreatePrePayOrderResponseReturnObjData `json:"data"`
}

type AmqpInstanceCreatePrePayOrderResponseReturnObjData struct {
	Submitted  bool   `json:"submitted"`
	NewOrderId string `json:"newOrderId"`
	NewOrderNo string `json:"newOrderNo"`
}
