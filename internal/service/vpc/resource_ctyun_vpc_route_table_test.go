package vpc_test

import (
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunVpcRouteTable(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_vpc_route_table." + rnd
	datasourcName := "data.ctyun_vpc_route_tables." + dnd
	initName := "terraform-unit"
	updatedName := "terraform-route-table"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: utils.LoadTestCase("ctyun_vpc_route_table.tf", rnd, dnd, initName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
				),
			},
			{
				Config: utils.LoadTestCase("ctyun_vpc_route_table.tf", rnd, dnd, updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
			{
				Config: utils.LoadTestCase("ctyun_vpc_route_table.tf", rnd, dnd, updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourcName, "route_tables.#", "1"),
					resource.TestCheckResourceAttr(datasourcName, "route_tables.0.name", updatedName),
				),
			},
			{

				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
				},
			},
		},
	})
}
