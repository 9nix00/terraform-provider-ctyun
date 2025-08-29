package scaling_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunScalingEcsProtection(t *testing.T) {
	//t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_scaling_ecs_protection." + rnd
	resourceFile := "resource_ctyun_scaling_ecs_protection.tf"

	// 从环境变量获取测试依赖资源
	scalingGroupID := 117717
	instanceIDList := fmt.Sprintf(`["%s","%s","%s"]`, "c8b8b345-fb9a-e9dd-c5b4-530263dcd74f", "dbe4e282-d38a-bcae-9a4d-38a6bc09171e", "bbe2c671-5fc1-610c-4560-30fd7e2bda14")

	protectStatus := true

	updatedProtectStatus := false

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 基础创建测试（开启保护）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, scalingGroupID, instanceIDList, protectStatus),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "true"),
				),
			},
			// 2. 资源更新测试（关闭保护并添加实例）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, scalingGroupID, instanceIDList, updatedProtectStatus),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "protect_status", "false"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, scalingGroupID, instanceIDList, updatedProtectStatus),
				Destroy: true,
			},
		},
	})
}
