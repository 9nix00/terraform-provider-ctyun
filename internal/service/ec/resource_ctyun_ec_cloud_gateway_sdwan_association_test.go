package ec_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCtyunEcCloudGatewaySdwanAssociation_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_ec_cloud_gateway_sdwan_association." + rnd
	resourceFile := "resource_ctyun_ec_cloud_gateway_sdwan_association.tf"

	// 从环境变量获取测试依赖资源
	ecID := dependence.expressConnectID
	sdwanID := dependence.sdwanId
	cgwID := dependence.cloudGatewayId
	rtbID := dependence.rtbID
	cgw_list := fmt.Sprintf("{cgw_id = \"%s\",rtb_id= \"%s\"}",
		cgwID,
		rtbID,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建云网关与SDWAN关联测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, sdwanID, cgw_list,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ec_id", ecID),
					resource.TestCheckResourceAttr(resourceName, "sdwan_id", sdwanID),
					resource.TestCheckResourceAttr(resourceName, "cgw_list.0.cgw_id", cgwID),
					resource.TestCheckResourceAttr(resourceName, "cgw_list.0.rtb_id", rtbID),
				),
			},
			// 2. 更新关联测试（解绑）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, sdwanID, "",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ec_id", ecID),
					resource.TestCheckResourceAttr(resourceName, "sdwan_id", sdwanID),
					resource.TestCheckResourceAttr(resourceName, "cgw_list.#", "0"),
				),
			},
		},
	})
}
