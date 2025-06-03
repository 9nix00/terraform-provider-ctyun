package kafka_test

import (
	"fmt"
	"os"
	"strconv"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunKafkaInstanceSingle(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_kafka_instance." + rnd
	datasourceName := "data.ctyun_kafka_instances." + dnd
	resourceFile := "resource_ctyun_kafka_instance.tf"
	datasourceFile := "datasource_ctyun_kafka_instances.tf"

	engineVersion := "3.6"
	nodeNum := 1
	zone := os.Getenv("CTYUN_AZ_NAME")
	extra := `cycle_type = "on_demand"`

	initName := "tf-kafka-init-" + utils.GenerateRandomString()
	initDiskSize := 100
	initRetentionHours := 80

	updatedName := "tf-kafka-updated-" + utils.GenerateRandomString()
	updatedDiskSize := 200
	updatedRetentionHours := 60

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
					engineVersion,
					dependence.kafkaSingleSpecName,
					nodeNum,
					zone,
					dependence.kafkaSingleDiskType,
					initDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					initRetentionHours,
					extra,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", initName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", engineVersion),
					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.kafkaSingleSpecName),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(nodeNum)),
					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
					resource.TestCheckResourceAttr(resourceName, "disk_type", dependence.kafkaSingleDiskType),
					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(initDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "retention_hours", strconv.Itoa(initRetentionHours)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			// 更新属性
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					engineVersion,
					dependence.kafkaSingleSpecName2,
					nodeNum,
					zone,
					dependence.kafkaSingleDiskType,
					updatedDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					updatedRetentionHours,
					extra,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", engineVersion),
					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.kafkaSingleSpecName2),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(nodeNum)),
					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
					resource.TestCheckResourceAttr(resourceName, "disk_type", dependence.kafkaSingleDiskType),
					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(updatedDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "retention_hours", strconv.Itoa(updatedRetentionHours)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},

			// 规格降级
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					engineVersion,
					dependence.kafkaSingleSpecName,
					nodeNum,
					zone,
					dependence.kafkaSingleDiskType,
					updatedDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					updatedRetentionHours,
					extra,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".id",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.instance_name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.engine_version", engineVersion),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.spec_name", dependence.kafkaSingleSpecName),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.node_num", strconv.Itoa(nodeNum)),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.disk_type", dependence.kafkaSingleDiskType),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.disk_size", strconv.Itoa(updatedDiskSize)),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.subnet_id", dependence.subnetID),
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
					"auto_renew_cycle_count",
					"auto_renew",
					"cycle_count",
					"cycle_type",
					"project_id",
					"security_group_id",
					"master_order_id",
				},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					engineVersion,
					dependence.kafkaSingleSpecName,
					nodeNum,
					zone,
					dependence.kafkaSingleDiskType,
					updatedDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					updatedRetentionHours,
					extra,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".id",
				),
				Destroy: true,
			},
		},
	})
}
