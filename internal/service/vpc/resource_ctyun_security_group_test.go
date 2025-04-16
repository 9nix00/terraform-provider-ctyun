package vpc_test

import (
	"fmt"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunSecurityGroup(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_security_group." + rnd
	datasourcName := "data.ctyun_security_groups." + dnd

	initName := "init"
	initDescription := "description"
	updatedName := "updated"
	updatedDescription := "updated-description"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: utils.LoadTestCase("ctyun_security_group.tf", rnd, dnd, initName, initDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase("ctyun_security_group.tf", rnd, dnd, updatedName, updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase("ctyun_security_group.tf", rnd, dnd, updatedName, updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourcName, "security_groups.#", "1"),
					resource.TestCheckResourceAttr(datasourcName, "security_groups.0.name", updatedName),
					resource.TestCheckResourceAttr(datasourcName, "security_groups.0.description", updatedDescription),
				),
			},
			{

				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
				},
			},
		},
	})
}
