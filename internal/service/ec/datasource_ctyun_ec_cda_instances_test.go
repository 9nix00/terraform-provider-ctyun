package ec_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCtyunEcCdaInstances_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dataSourceName := "data.ctyun_ec_cda_instances." + rnd
	dataSourceFile := "datasource_ctyun_ec_cda_instances.tf"

	// 从环境变量获取测试依赖资源
	ecID := dependence.expressConnectID

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(dataSourceFile, rnd, ecID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.ec_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.cgw_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.cda_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.default_rtb_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.create_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.status"),
				),
			},
		},
	})
}

func TestAccCtyunEcCdaInstances_withCdaID(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dataSourceName := "data.ctyun_ec_cda_instances." + rnd
	dataSourceFile := "datasource_ctyun_ec_cda_instances.tf"

	// 从环境变量获取测试依赖资源
	ecID := dependence.expressConnectID
	cdaID := "fake-cda-id-for-test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(dataSourceFile, rnd, ecID) + fmt.Sprintf(`
  cda_id = "%s"
`, cdaID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.#"),
				),
			},
		},
	})
}
