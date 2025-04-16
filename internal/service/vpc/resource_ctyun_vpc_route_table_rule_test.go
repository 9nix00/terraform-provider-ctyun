package vpc_test

import (
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunVpcRouteTableRule(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_vpc_route_table_rule." + rnd
	initDestination := "188.188.0.0/16"
	initDescription := "test"
	updatedDescription := "updated"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: utils.LoadTestCase("ctyun_vpc_route_table_rule.tf", rnd, initDestination, initDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckResourceAttr(resourceName, "destination", initDestination),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "rule_id"),
				),
			},
			{
				Config: utils.LoadTestCase("ctyun_vpc_route_table_rule.tf", rnd, initDestination, updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "destination", initDestination),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}
