package vpc_test

import (
	"fmt"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunVpcRouteTable(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_vpc_route_table." + rnd
	datasourceName := "data.ctyun_vpc_route_tables." + dnd
	resourceFile := "resource_ctyun_vpc_route_table.tf"
	datasourceFile := "datasource_ctyun_vpc_route_tables.tf"

	initName := "terraform-unit"
	updatedName := "terraform-route-table"
	var id string
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
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
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "route_tables.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "route_tables.0.name", updatedName),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
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
					resource.TestCheckResourceAttr(datasourceName, "route_tables.#", "0"),
				),
			},
		},
	})
}
