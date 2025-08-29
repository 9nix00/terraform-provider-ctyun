package ecs_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunEcsFlavors_basic(t *testing.T) {
	rName := utils.GenerateRandomString()
	dataSourceName := "data.ctyun_ecs_flavors." + rName
	resourceFile := "datasource_ctyun_ecs_flavors.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCtyunEcsFlavorsConfig_basic(resourceFile, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.cpu"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.ram"),
				),
			},
		},
	})
}

func testAccCtyunEcsFlavorsConfig_basic(file, name string) string {
	return utils.LoadTestCase(file, name)
}

func TestAccCtyunEcsFlavors_byType(t *testing.T) {
	rName := utils.GenerateRandomString()
	dataSourceName := "data.ctyun_ecs_flavors." + rName
	resourceFile := "datasource_ctyun_ecs_flavors.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCtyunEcsFlavorsConfig_byType(resourceFile, rName, "CPU"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
					testAccCheckFlavorsByType(dataSourceName, "CPU"),
				),
			},
		},
	})
}

func testAccCheckFlavorsByType(n, flavorType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		// 验证所有返回的flavors都符合type过滤条件
		count := 0
		for key, value := range rs.Primary.Attributes {
			if fmt.Sprintf("flavors.%d.type", count) == key {
				if value != flavorType {
					return fmt.Errorf("expected type %s, got %s", flavorType, value)
				}
				count++
			}
		}

		return nil
	}
}

func testAccCtyunEcsFlavorsConfig_byType(file, name, flavorType string) string {
	return fmt.Sprintf(`
%s

data "ctyun_ecs_flavors" "%s_type" {
  type = "%s"
}`, utils.LoadTestCase(file, name), name, flavorType)
}
