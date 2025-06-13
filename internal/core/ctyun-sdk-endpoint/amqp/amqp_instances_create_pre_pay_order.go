package amqp

import (
	"context"
	"net/http"
	ctyunsdk "terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

type AmqpInstancesCreatePrePayOrderApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstancesCreatePrePayOrderApi(client *ctyunsdk.CtyunClient) *AmqpInstancesCreatePrePayOrderApi {
	return &AmqpInstancesCreatePrePayOrderApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v3/instances/createPrePayOrder",
		},
	}
}

func (this *AmqpInstancesCreatePrePayOrderApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstancesCreatePrePayOrderRequest) (res *AmqpInstancesCreatePrePayOrderResponse, err error) {
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
	res = &AmqpInstancesCreatePrePayOrderResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstancesCreatePrePayOrderRequest struct {
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
	CycleCnt        int32  `json:"cycleCnt"`
}

type AmqpInstancesCreatePrePayOrderResponse struct {
	ReturnObj  *AmqpInstancesCreatePrePayOrderResponseReturnObj `json:"returnObj"`
	Message    string                                           `json:"message"`
	StatusCode string                                           `json:"statusCode"`
}

type AmqpInstancesCreatePrePayOrderResponseReturnObj struct {
	Data AmqpInstancesCreatePrePayOrderResponseReturnObjData `json:"data"`
}

type AmqpInstancesCreatePrePayOrderResponseReturnObjData struct {
	Submitted  bool   `json:"submitted"`
	NewOrderId string `json:"newOrderId"`
	NewOrderNo string `json:"newOrderNo"`
}
