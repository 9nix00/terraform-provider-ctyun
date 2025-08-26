package sfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunSfsPermissionGroupAssociation(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_sfs_permission_group_association." + rnd
	resourceFile := "resource_ctyun_sfs_permission_group_association.tf"

	// 从环境变量获取测试依赖资源
	sfsUID := dependence.SfsUID
	vpcID := dependence.vpcID
	vpcID1 := dependence.vpcID1
	permissionGroupFuid1 := dependence.sfsPermissionGroupID
	permissionGroupFuid2 := dependence.sfsPermissionGroupID1

	// 验证环境变量是否设置
	if sfsUID == "" || vpcID == "" || permissionGroupFuid1 == "" || permissionGroupFuid2 == "" {
		t.Skip("Skipping test: required environment variables not set")
	}

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
			// 1. 基础创建测试（绑定第一个权限组）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid1, sfsUID, vpcID1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "permission_group_fuid", permissionGroupFuid1),
					resource.TestCheckResourceAttr(resourceName, "sfs_uid", sfsUID),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID1),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_cidr"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_description"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_is_default"),
				),
			},
			// 2. 资源更新测试（更换为第二个权限组）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid2, sfsUID, vpcID1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttr(resourceName, "permission_group_fuid", permissionGroupFuid2),
					resource.TestCheckResourceAttr(resourceName, "sfs_uid", sfsUID),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID1),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_cidr"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_description"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_is_default"),
				),
			},
			// 3. 资源导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}

					// 构造导入ID: "id,region_id"
					return fmt.Sprintf("%s,%s,%s,%s",
						rs.Primary.Attributes["region_id"],
						rs.Primary.Attributes["vpc_id"],
						rs.Primary.Attributes["permission_group_fuid"],
						rs.Primary.Attributes["sfs_uid"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{}, // 不需要忽略任何字段
			},
			// 4. 清理资源
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid2, sfsUID, vpcID),
				Destroy: true,
			},
		},
	})
}
