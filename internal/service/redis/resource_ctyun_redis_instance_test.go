package redis_test

import (
	"fmt"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunVpce(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	and := utils.GenerateRandomString()

	resourceName := "ctyun_redis_instance." + rnd
	datasourceName := "data.ctyun_redis_instances." + dnd
	resourceFile := "resource_ctyun_redis_instance.tf"
	datasourceFile := "datasource_ctyun_redis_instances.tf"
	associationFile := "resource_ctyun_redis_association_eip.tf"

	initName := "tf-redis-" + utils.GenerateRandomString()

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
			{
				// 创建
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					dependence.redisVersion,
					dependence.redisEngineEdition,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", initName),
					resource.TestCheckResourceAttr(resourceName, "version", dependence.redisVersion),
					resource.TestCheckResourceAttr(resourceName, "edition", dependence.redisEngineEdition),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 绑定eip
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					dependence.redisVersion,
					dependence.redisEngineEdition,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				) + utils.LoadTestCase(
					associationFile, and,
					dependence.eipAddress,
					resourceName+".id",
				),
			},
			// 通过查询进行检查
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					dependence.redisVersion,
					dependence.redisEngineEdition,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				) + utils.LoadTestCase(
					associationFile, and,
					dependence.eipAddress,
					resourceName+".id",
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".instance_name",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.name", initName),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.eip_address", dependence.eipAddress),
				),
			},
			// 解绑并检查
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					dependence.redisVersion,
					dependence.redisEngineEdition,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".instance_name",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.name", initName),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.eip_address", dependence.eipAddress),
				),
			},

			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					dependence.redisVersion,
					dependence.redisEngineEdition,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".instance_name",
				),
				Destroy: true,
			},
		},
	})
}
