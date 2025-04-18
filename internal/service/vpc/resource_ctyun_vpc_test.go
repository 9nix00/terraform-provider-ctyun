package vpc_test

import (
	"fmt"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunVpc(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_vpc." + rnd
	datasourceName := "data.ctyun_vpcs." + dnd
	resourceFile := "resource_ctyun_vpc.tf"
	datasourceFile := "datasource_ctyun_vpcs.tf"

	initName := utils.GenerateRandomString()
	initCidr := "192.168.0.0/16"
	initDescription := utils.GenerateRandomString()

	updatedName := utils.GenerateRandomString()
	updatedDescription := utils.GenerateRandomString()

	var id string
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, initDescription, initCidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "cidr", initCidr),
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					func(state *terraform.State) error {
						rs, ok := state.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource not found")
						}
						id = rs.Primary.ID
						return nil
					},
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription, initCidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cidr", initCidr),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription, initCidr) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "vpcs.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "vpcs.0.name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "vpcs.0.cidr", initCidr),
					resource.TestCheckResourceAttr(datasourceName, "vpcs.0.description", updatedDescription),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					projectId := ds.Attributes["project_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,%s,%s", id, regionId, projectId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf(`"%s"`, id)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "vpcs.#", "0"),
				),
			},
		},
	})
}
