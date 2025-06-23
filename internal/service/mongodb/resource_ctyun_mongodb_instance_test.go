package mongodb_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccCtyunMongodbInstance(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance.tf"
	datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	cycleType := "on_demand"
	//cycleCount := 1
	//autoRenew := false
	vpcID := dependence.vpcID
	hostType := "S7"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb" + utils.GenerateRandomString()
	password := "BwLnOyOob+wbuDq0qI+EFMWmOtkFiJiIVIvlknaf1fJzfTfUAG7IwGM+368VIOZ9vuGX8QOUTVhzqT6ZKlaBzEHXUbHxQ8TeruEUZx5d4I8o4396WFZOu8+LuVklKDYx6V1tm8FDpXRXwAOsU+fYM8tQi3cu2oynET5TeQMVYxw="
	purchase_count := 1
	prodId := 10013001
	nodeInfoList := `{"node_type":"master","inst_spec":"1","storage_type":"SATA","storage_space":100,"prod_performance_spec":"2C4G","disks":1,"availability_zone_info":[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"master"}]}`
	updatedPort := "read_port=12345"
	updateName := "tf-mongodb-new" + utils.GenerateRandomString()
	updatedIsUpgradeBackUp := `is_upgrade_back_up=true`
	updateNodeInfoList := `{"node_type":"master","inst_spec":"1","storage_type":"SATA","storage_space":120,"prod_performance_spec":"2C4G","disks":1,"availability_zone_info":[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"master"}]}`
	updateBackUpDiskUpgradeNodeInfoList := `{"node_type":"backup","inst_spec":"1","storage_type":"SATA","storage_space":130,"prod_performance_spec":"2C4G","disks":1,"availability_zone_info":[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"master"}]}`
	updateSpecUpgradeNodeInfoList := `{"node_type":"master","inst_spec":"1","storage_type":"SATA","storage_space":130,"prod_performance_spec":"2C8G","disks":1,"availability_zone_info":[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"master"}]}`
	//updateProdIDUpgradeNodeInfoList := `{"node_type"="master","inst_spec"=1,"storage_type"="SATA","storage_space"=130,"prod_performance_spec":"2C8G","disks":1,"availability_zone_info":["availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"master"]}`
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
			// 创建mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, "", "", vpcID, hostType, subnetID, securityGroupID, name, password, purchase_count, prodId, nodeInfoList, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", prodId)),
				),
			},
			// 更新mongodb实例名称和端口号
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, "", "", vpcID, hostType, subnetID, securityGroupID, updateName, password, purchase_count, prodId, nodeInfoList, updatedPort, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", prodId)),
					resource.TestCheckResourceAttr(resourceName, "read_port", "12345"),
				),
			},
			// 升配mongodb-主+备份空间磁盘扩容,
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, "", "", vpcID, hostType, subnetID, securityGroupID, name, password, purchase_count, prodId, updateNodeInfoList, updatedPort, updatedIsUpgradeBackUp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			// 升配备份空间磁盘扩容
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, "", "", vpcID, hostType, subnetID, securityGroupID, name, password, purchase_count, prodId, updateBackUpDiskUpgradeNodeInfoList, updatedPort, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			// 升配规格升级
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, "", "", vpcID, hostType, subnetID, securityGroupID, name, password, purchase_count, prodId, updateSpecUpgradeNodeInfoList, updatedPort, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "prod_performance_spec", "2C8G"),
				),
			},
			//// 升级节点
			//{
			//	Config: utils.LoadTestCase(resourceFile, rnd, cycleType, cycleCount, autoRenew, vpcID, hostType, subnetID, securityGroupID, name, password, purchase_count, prodId, updateSpecUpgradeNodeInfoList, updatedPort, ""),
			//},
			// datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, "", "", vpcID, hostType, subnetID, securityGroupID, name, password, purchase_count, prodId, updateSpecUpgradeNodeInfoList, updatedPort, "") +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("prod_inst_name=%s", updateName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "mongodb_instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "mongodb_instances.0.name", name),
					resource.TestCheckResourceAttr(datasourceName, "mongodb_instances.0.prod_performance_spec", "2C8G"),
				),
			},
			// destroy
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleType, "", "", vpcID, hostType, subnetID, securityGroupID, name, password, purchase_count, prodId, updateSpecUpgradeNodeInfoList, updatedPort, ""),
				Destroy: true,
			},
		},
	})
}
