package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCtyunMysqlReadOnlyInstance(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_readonly_instance." + rnd
	resourceFile := "resource_ctyun_mysql_read_node.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	instanceID := dependence.mysqlID
	cycleType := "on_demand"
	//cycleCount := 1
	// 测试数据
	instanceName := "test-ro-" + rnd
	flavorName := "s7.large.2" // 示例规格，根据实际情况调整
	storageType := "SATA"      // 存储类型
	storageSpace := 100        // 存储空间(GB)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建按需计费只读实例测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceID, cycleType, flavorName,
					projectID, storageType, storageSpace, instanceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cycle_type", "on_demand"),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			//// 2. 创建包月计费只读实例测试
			//{
			//	Config: testAccCtyunMysqlReadOnlyInstanceConfig(
			//		rnd, resourceFile,
			//		masterInstanceID, "month",
			//		1, true, flavorName,
			//		regionID, projectID, storageType,
			//		storageSpace, instanceName, availabilityZone,
			//	),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr(resourceName, "cycle_type", "month"),
			//		resource.TestCheckResourceAttr(resourceName, "cycle_count", "1"),
			//		resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
			//		resource.TestCheckResourceAttrSet(resourceName, "id"),
			//	),
			//},
			//// 3. 更新存储空间测试
			//{
			//	Config: testAccCtyunMysqlReadOnlyInstanceConfig(
			//		rnd, resourceFile,
			//		masterInstanceID, "month",
			//		1, true, flavorName,
			//		regionID, projectID, storageType,
			//		storageSpace+50, // 增加存储空间
			//		instanceName, availabilityZone,
			//	),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", storageSpace+50)),
			//	),
			//},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceID, cycleType, flavorName,
					projectID, storageType, storageSpace, instanceName,
				),
				Destroy: true,
			},
		},
	})
}
