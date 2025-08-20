package ports_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunNetworkInterface_basic(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_port." + rnd
	resourceFile := "resource_ctyun_network_interface.tf"

	name := "tf-port-test-" + utils.GenerateRandomString()
	updatedName := "tf-port-test-updated-" + utils.GenerateRandomString()

	description := "测试网络接口描述"
	updatedDescription := description + "-updated"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				// 测试创建
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					name,
					description,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface_id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
				),
			},
			{
				// 测试更新
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					updatedDescription,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
			{
				// 测试导入
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}

					regionId := rs.Primary.Attributes["region_id"]
					if regionId == "" {
						return "", fmt.Errorf("region_id is not set")
					}

					return fmt.Sprintf("%s,%s", rs.Primary.ID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"subnet_id",
					"security_group_ids",
					"secondary_private_ip_count",
				},
			},
			{
				// 测试销毁
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					updatedDescription,
				),
				Destroy: true,
			},
		},
	})
}
