package nat_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccCtyun(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_snat." + rnd
	datasourceName := "data.ctyun_snat." + dnd
	resourceFile := "resource_ctyun_snat.tf"
	datasourceFile := "datasource_ctyun_snat.tf"

	initDescription1 := "terraform provider 开发测试" + utils.GenerateRandomString()
	initDescription2 := "terraform provider 开发测试" + utils.GenerateRandomString()
	initSourceCidr := ""
	updatedSourceCidr := ""
	sourceSubnetId := ""

	updateDescription := utils.GenerateRandomString()
	//var id string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// 1.resource create验证1:
				// subnetType = 1(有vpcId 的子网情况),sourceSubnetId必传
				Config: utils.LoadTestCase(resourceFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr(resourceName, "source_cidr", initSourceCidr),
					resource.TestCheckResourceAttr(resourceName, "description", initDescription1),
					resource.TestCheckResourceAttr(resourceName, "source_subnet_id", sourceSubnetId),
					resource.TestCheckResourceAttrSet(resourceName, "snatIps"),
				),
			},
			{
				// 2.resource create验证2:
				// subnetType = 0(自定义情况),sourceCIDR必传
				Config: utils.LoadTestCase(resourceFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "description", initDescription2),
					resource.TestCheckResourceAttr(datasourceName, "source_cidr", initSourceCidr),
					resource.TestCheckResourceAttrSet(resourceName, "snat_ips"),
					resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"),
				),
			},
			{
				// 3. resource update验证
				Config: utils.LoadTestCase(resourceFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(resourceName, "source_cidr", updatedSourceCidr),
				),
			},
			{
				// 4. datasource验证
				Config: utils.LoadTestCase(resourceFile) + utils.LoadTestCase(datasourceFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "nats.#", "2"),
					resource.TestCheckResourceAttr(datasourceName, "nats.0.description", initDescription1),
					resource.TestCheckResourceAttr(datasourceName, "nats.0.source_subnet_id", sourceSubnetId),
					resource.TestCheckResourceAttr(datasourceName, "nats.1.description", updateDescription),
					resource.TestCheckResourceAttr(datasourceName, "nats.1.source_cidr", updatedSourceCidr),
				),
			},
		},
	})
}
