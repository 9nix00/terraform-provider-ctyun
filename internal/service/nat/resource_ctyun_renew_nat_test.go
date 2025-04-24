package nat_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccNewCtyunRenewNatResource(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_nat." + rnd
	datasourceName := "data.ctyun_nats." + dnd
	initDescription := "terraform provider 开发测试"
	resourceFile := "resource_ctyun_renew_nat.tf"
	datasourceFile := "datasource_ctyun_nat.tf"

	vpcId := "vpc-wf029jgx2d"
	spec := "1"
	updatedSpec := "1"
	cycle_type := "month"
	cycle_count := "1"
	updated_cycle_count := "2"
	az_name := "cn-huanan2-1A-public-ctcloud"
	initName := utils.GenerateRandomString()

	updatedName := utils.GenerateRandomString()
	updatedDescription := utils.GenerateRandomString()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1.resource create验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, spec, initName, initDescription, cycle_type, cycle_count, az_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"),
				),
			},
			// 2. resource update验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, spec, updatedName, updatedDescription, cycle_type, cycle_count, az_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"),
				),
			},
			// 3. resource nat续费
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, updatedSpec, updatedName, updatedDescription, cycle_type, updated_cycle_count, az_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					//resource.TestCheckResourceAttr(resourceName, "expired_time", utils.getExpireTime()),
					resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"),
				),
			},
			// 4. datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, spec, updatedName, updatedDescription, cycle_type, cycle_count, az_name) +
					utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "nats.0.name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "nats.0.description", updatedDescription),
				),
			},
			// 4.import_state 验证
			/*
				{
					ResourceName: resourceName,
					ImportState:  true,
					ImportStateIdFunc: func(s *terraform.State) (string, error) {
						ds := s.RootModule().Resources[resourceName].Primary
						id := ds.ID
						regionId := ds.Attributes["region_id"]
						projectId := ds.Attributes["project_id"]
						if id == "" || regionId == "" {
							return "", fmt.Errorf("ID/projectID or regionID cannot be empty")
						}
						return fmt.Sprintf("%s,%s,%s", id, projectId, regionId), nil
					},
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{},
				},
			*/

		},
	})
}
