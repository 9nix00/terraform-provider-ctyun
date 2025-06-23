package mongodb

import (
	"context"
	"errors"
	"net/http"
	ctyunsdk "terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

type MongodbUpdateInstanceNameApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbUpdateInstanceNameApi(client *ctyunsdk.CtyunClient) *MongodbUpdateInstanceNameApi {
	return &MongodbUpdateInstanceNameApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/openApi/modifyDBInstanceDescription",
		},
	}
}

func (this *MongodbUpdateInstanceNameApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbUpdateInstanceNameRequest, header *MongodbUpdateInstanceNameRequestHeader) (updatedNameResp *MongodbUpdateInstanceNameResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.RegionID == "" {
		err = errors.New("missing required field: RegionID")
		return
	}
	builder.AddHeader("regionId", header.RegionID)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameMongodb, builder)
	if err != nil {
		return
	}
	updatedNameResp = &MongodbUpdateInstanceNameResponse{}
	err = resp.Parse(updatedNameResp)
	if err != nil {
		return
	}
	return updatedNameResp, nil
}

type MongodbUpdateInstanceNameRequest struct {
	ProdInstId   string `json:"prodInstId"`   // 实例id
	ProdInstName string `json:"prodInstName"` // 修改后的实例名称
}

type MongodbUpdateInstanceNameRequestHeader struct {
	CustomInfo *string `json:"customInfo,omitempty"`
	AccountId  string  `json:"accountId"` // PaaS组件实例所归属天翼云账号标识
	RegionID   string  `json:"regionId"`
}

type MongodbUpdateInstanceNameResponse struct {
	StatusCode int32   `json:"statusCode"`
	Message    *string `json:"message"`
	Error      string  `json:"error"`
}
