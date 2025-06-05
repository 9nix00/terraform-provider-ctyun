package amqp

import (
	"context"
	"fmt"
	"net/http"
	ctyunsdk "terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

type AmqpInstanceQueryApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstanceQueryApi(client *ctyunsdk.CtyunClient) *AmqpInstanceQueryApi {
	return &AmqpInstanceQueryApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v3/instances/query",
		},
	}
}

func (this *AmqpInstanceQueryApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstanceQueryRequest) (res *AmqpInstanceQueryRequest, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	builder.AddParam("pageNum", fmt.Sprintf("%d", req.PageNum))
	builder.AddParam("pageSize", fmt.Sprintf("%d", req.PageSize))
	builder.AddHeader("regionId", req.RegionID)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointName, builder)
	if err != nil {
		return
	}
	res = &AmqpInstanceQueryRequest{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstanceQueryRequest struct {
	RegionID string `json:"regionId"`
	PageNum  int32  `json:"pageNum"`
	PageSize int32  `json:"pageSize"`
}

type AmqpInstanceQueryResponse struct {
	ReturnObj  *AmqpInstanceQueryResponseReturnObj `json:"returnObj"`
	Message    string                              `json:"message"`
	StatusCode string                              `json:"statusCode"`
}

type AmqpInstanceQueryResponseReturnObj struct {
	Total int32                                    `json:"total"`
	Data  []AmqpInstanceQueryResponseReturnObjData `json:"data"`
}

type AmqpInstanceQueryResponseReturnObjData struct {
	Cluster       string `json:"cluster"`
	Subnet        string `json:"subnet"`
	Prod          string `json:"prod"`
	EngineType    string `json:"engineType"`
	BillMode      string `json:"billMode"`
	SecurityGroup string `json:"securityGroup"`
	//ProdType      interface{} `json:"prodType"`
	Network     string `json:"network"`
	ExpireTime  string `json:"expireTime"`
	CreateTime  string `json:"createTime"`
	ClusterName string `json:"clusterName"`
	ProdInstId  string `json:"prodInstId"`
	Status      int32  `json:"status"`
}
