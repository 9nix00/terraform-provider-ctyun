package amqp

import (
	ctyunsdk "terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

type Apis struct {
	AmqpInstanceQueryProdApi *AmqpInstanceQueryProdApi
	AmqpInstanceQueryApi     *AmqpInstanceQueryApi
}

func NewApis(client *ctyunsdk.CtyunClient) *Apis {
	builder := ctyunsdk.NewApiHookBuilder()
	for _, hook := range client.Config.ApiHooks {
		builder.AddHooks(hook)
	}

	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentProd, EndPointProd)
	return &Apis{
		AmqpInstanceQueryProdApi: NewAmqpInstanceQueryProdApi(client),
		AmqpInstanceQueryApi:     NewAmqpInstanceQueryApi(client),
	}
}
