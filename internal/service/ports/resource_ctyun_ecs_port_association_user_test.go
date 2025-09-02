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

func TestAccCtyunEcsPortAssociation_all(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	name := "ctyun_ecs_port_association." + rnd
	configFile := "resource_ctyun_ecs_port_association.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				// 测试基本创建场景
				Config: utils.LoadTestCase(configFile, rnd, dependence.instanceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "region_id"),
					resource.TestCheckResourceAttrSet(name, "project_id"),
					resource.TestCheckResourceAttr(name, "instance_id", dependence.instanceID),
				),
			},
			{
				// 测试更新场景
				Config: utils.LoadTestCase(configFile, rnd, dependence.instanceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "region_id"),
					resource.TestCheckResourceAttrSet(name, "project_id"),
					resource.TestCheckResourceAttr(name, "instance_id", dependence.instanceID),
					resource.TestCheckResourceAttrSet(name, "network_interface_id"),
				),
			},
			{
				// 测试导入功能
				ResourceName:      name,
				ImportState:       true,
				ImportStateIdFunc: generateImportStateIdFunc(name),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"az_name",
					"project_id",
				},
			},
			{
				// 测试销毁解绑场景
				Config:  utils.LoadTestCase(configFile, rnd, dependence.instanceID),
				Destroy: true,
			},
		},
	})
}

// generateImportStateIdFunc 生成导入ID函数
func generateImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		regionId := rs.Primary.Attributes["region_id"]
		instanceId := rs.Primary.Attributes["instance_id"]
		networkInterfaceId := rs.Primary.Attributes["network_interface_id"]

		if regionId == "" || instanceId == "" || networkInterfaceId == "" {
			return "", fmt.Errorf("missing required attributes for import")
		}

		return fmt.Sprintf("%s,%s,%s", regionId, instanceId, networkInterfaceId), nil
	}
}
