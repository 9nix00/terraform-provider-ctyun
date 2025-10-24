package pgsql_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCtyunPostgresqlReadOnlyInstance(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_readonly_instance." + rnd
	resourceFile := "resource_ctyun_postgresql_readonly_instance_on_demand.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"                 // 默认项目ID
	instanceID := dependence.PgsqlID // 主实例ID
	cycleType := "on_demand"

	// 测试数据
	instanceName := "test-pg-ro-" + rnd
	flavorName := "s7.large.2" // PostgreSQL 规格

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建按需计费只读实例测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceID, cycleType, flavorName,
					projectID, instanceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cycle_type", "on_demand"),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			//// 2. 创建包月计费只读实例测试
			//{
			//	Config: utils.LoadTestCase(
			//		resourceFile, rnd,
			//		instanceID, "month",
			//		1, true, flavorName, // 添加 cycle_count 和 auto_renew
			//		projectID, storageType, storageSpace, instanceName,
			//	),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr(resourceName, "cycle_type", "month"),
			//		resource.TestCheckResourceAttr(resourceName, "cycle_count", "1"),
			//		resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
			//	),
			//},
			//// 3. 更新存储空间测试
			//{
			//	Config: utils.LoadTestCase(
			//		resourceFile, rnd,
			//		instanceID, "month",
			//		1, true, flavorName,
			//		projectID, storageType, storageSpace+50, // 增加存储空间
			//		instanceName,
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
					projectID, instanceName,
				),
				Destroy: true,
			},
		},
	})
}
