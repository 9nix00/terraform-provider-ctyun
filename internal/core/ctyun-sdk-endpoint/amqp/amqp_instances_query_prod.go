package amqp

import (
	"context"
	"net/http"
	ctyunsdk "terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

type AmqpInstanceQueryProdApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstanceQueryProdApi(client *ctyunsdk.CtyunClient) *AmqpInstanceQueryProdApi {
	return &AmqpInstanceQueryProdApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v3/instances/queryProd",
		},
	}
}

func (this *AmqpInstanceQueryProdApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstanceQueryProdRequest) (res *AmqpInstanceQueryProdResponse, err error) {
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
	res = &AmqpInstanceQueryProdResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstanceQueryProdRequest struct {
	RegionID string `json:"regionId"`
}

type AmqpInstanceQueryProdResponse struct {
	ReturnObj  *AmqpInstanceQueryProdResponseReturnObj `json:"returnObj"`
	Message    string                                  `json:"message"`
	StatusCode string                                  `json:"statusCode"`
}

type AmqpInstanceQueryProdResponseReturnObj struct {
	Data []AmqpInstanceQueryProdResponseReturnObjData `json:"data"`
}

type AmqpInstanceQueryProdResponseReturnObjData struct {
	FlavorID      string      `json:"flavorID"`
	SpecName      string      `json:"specName"`
	FlavorType    string      `json:"flavorType"`
	FlavorName    string      `json:"flavorName"`
	CpuNum        int32       `json:"cpuNum"`
	MemSize       int32       `json:"memSize"`
	MultiQueue    int32       `json:"multiQueue"`
	Pps           int32       `json:"pps"`
	BandwidthBase float64     `json:"bandwidthBase"`
	BandwidthMax  int32       `json:"bandwidthMax"`
	CpuArch       interface{} `json:"cpuArch"`
	Series        string      `json:"series"`
	AzList        []string    `json:"azList"`
}
