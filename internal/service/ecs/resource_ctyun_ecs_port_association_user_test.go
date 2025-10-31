package ecs_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunEcsPortAssociation_all(t *testing.T) {
	rnd := utils.GenerateRandomString()
	name := "ctyun_ecs_port_association." + rnd
	configFile := "resource_ctyun_ecs_port_association.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				// 测试基本创建场景
				Config: utils.LoadTestCase(configFile, rnd, dependence.instanceID, dependence.ecsPortForAssociationId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "region_id"),
					resource.TestCheckResourceAttr(name, "instance_id", dependence.instanceID),
				),
			},
			{
				// 测试更新场景
				Config: utils.LoadTestCase(configFile, rnd, dependence.instanceID, dependence.ecsPortForAssociationId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "region_id"),
					resource.TestCheckResourceAttr(name, "instance_id", dependence.instanceID),
					resource.TestCheckResourceAttrSet(name, "port_id"),
				),
			},
			{
				// 测试导入功能
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"az_name",
					"project_id",
				},
			},
			{
				// 测试销毁解绑场景
				Config:  utils.LoadTestCase(configFile, rnd, dependence.instanceID, dependence.ecsPortForAssociationId),
				Destroy: true,
			},
		},
	})
}
