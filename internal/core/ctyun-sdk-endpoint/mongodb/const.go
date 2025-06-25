package mongodb

import ctyunsdk "terraform-provider-ctyun/internal/core/ctyun-sdk-core"

const (
	EndpointNameMongodb = "mongodb"
	UrlPordMongodb      = "mongodb-global.ctapi.ctyun.cn"
	UrlTestMongodb      = ""
)

var EndpointMongodbTest = ctyunsdk.Endpoint{
	EndpointName: EndpointNameMongodb,
	Url:          UrlTestMongodb,
}

var EndPointMongodbProd = ctyunsdk.Endpoint{
	EndpointName: EndpointNameMongodb,
	Url:          UrlPordMongodb,
}
