package nat_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"os"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccNewCtyunNatResource(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_nat." + rnd
	datasourceName := "data.ctyun_nats." + dnd
	initDescription := "terraform provider 开发测试"
	resourceFile := "resource_ctyun_nat.tf"
	datasourceFile := "datasource_ctyun_nat.tf"

	vpcId := "vpc-wf029jgx2d"
	spec := "1"
	cycle_type := "on_demand"
	az_name := "cn-huanan2-1A-public-ctcloud"
	initName := utils.GenerateRandomString()

	updatedName := utils.GenerateRandomString()
	updatedDescription := utils.GenerateRandomString()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1.resource create验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, spec, initName, initDescription, cycle_type, az_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"),
				),
			},
			// 2. resource update验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, spec, updatedName, updatedDescription, cycle_type, az_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"),
				),
			},
			// 3. datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, spec, updatedName, updatedDescription, cycle_type, az_name) +
					utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "nats.#", "1"),
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
