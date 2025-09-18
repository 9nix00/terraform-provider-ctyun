package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
	"time"
)

// 创建备份集 + 取消备份任务
func TestAccCtyunMysqlBackup(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_backup." + rnd
	resourceFile := "resource_ctyun_mysql_backup.tf"

	cancelResourceName := "ctyun_mysql_backup_cancel." + rnd
	cancelResourceFile := "resource_ctyun_mysql_backup_cancel.tf"

	dnd := utils.GenerateRandomString()
	backupsDatasourceName := "data.ctyun_mysql_backups." + dnd
	backupsDatasourceFile := "datasource_ctyun_mysql_backups.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	mysqlInstanceID := "e5ad1c553e394bc891c5bf8fc58be191"
	// 备份描述信息
	description := "Test backup created by Terraform"

	wait10Seconds := func(s *terraform.State) error {
		t.Logf("等待20秒...")
		time.Sleep(20 * time.Second)
		return nil
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
			// 1. 创建备份测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					mysqlInstanceID, projectID,
					description, "full",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "inst_id", mysqlInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "task_type", "full"),
					resource.TestCheckResourceAttrSet(resourceName, "backup_name"),
					wait10Seconds,
				),
			},
			// 2. 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}

					// 构造导入ID: ID|region_id|project_id|inst_id
					return fmt.Sprintf("%s,%s,%s,%s",
						rs.Primary.ID,
						rs.Primary.Attributes["region_id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["inst_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description", "task_type"}, // 不需要忽略任何字段
			},

			// 验证 backups datasource
			{
				Config: utils.LoadTestCase(resourceFile, rnd, mysqlInstanceID, projectID, description, "full") +
					utils.LoadTestCase(backupsDatasourceFile, dnd, mysqlInstanceID, fmt.Sprintf("%s.backup_name", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(backupsDatasourceName, "backup_list.#"),
					resource.TestCheckResourceAttr(backupsDatasourceName, "backup_list.0.inst_id", mysqlInstanceID),
				),
			},
			// datasource获取backup_record_id
			{

				Config: utils.LoadTestCase(resourceFile, rnd, mysqlInstanceID, projectID, description, "full") +
					utils.LoadTestCase(backupsDatasourceFile, dnd, mysqlInstanceID, fmt.Sprintf("%s.backup_name", resourceName)) +
					utils.LoadTestCase(cancelResourceFile, rnd, mysqlInstanceID, projectID, fmt.Sprintf("%s.backup_list.0.records.0.backup_record_id", backupsDatasourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(cancelResourceName, "inst_id", mysqlInstanceID),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					mysqlInstanceID, projectID,
					description, "full",
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunMysqlBackupRecovery(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_backup_recovery." + rnd
	resourceFile := "resource_ctyun_mysql_backup_recovery.tf"

	timePointDatasourceName := "data.ctyun_mysql_recoverable_time_points." + dnd
	timePointDatasourceFile := "datasource_ctyun_mysql_recoverable_time_points.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	srcInstanceID := "e5ad1c553e394bc891c5bf8fc58be191"
	dstInstanceID := "e5ad1c553e394bc891c5bf8fc58be191"
	//toTimepoint := os.Getenv("CTYUN_MYSQL_BACKUP_RECOVERY_TIME")

	// 等待函数
	//wait10Seconds := func() {
	//	t.Logf("等待10秒...")
	//	time.Sleep(10 * time.Second)
	//}

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
			// 1. 验证datasource
			{
				Config: utils.LoadTestCase(timePointDatasourceFile, dnd, srcInstanceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(timePointDatasourceName, "backup_time_points.#"),
				),
			},
			// 1. 创建备份恢复任务
			{
				Config: utils.LoadTestCase(timePointDatasourceFile, dnd, srcInstanceID) +
					utils.LoadTestCase(
						resourceFile, rnd, srcInstanceID, projectID, srcInstanceID, dstInstanceID,
						fmt.Sprintf("%s.backup_time_points.0.end_timestamp", timePointDatasourceName),
					),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "inst_id", srcInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "src_inst_id", srcInstanceID),
					resource.TestCheckResourceAttr(resourceName, "dst_inst_id", dstInstanceID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 销毁
			{
				Config: utils.LoadTestCase(timePointDatasourceFile, dnd, srcInstanceID) +
					utils.LoadTestCase(
						resourceFile, rnd, srcInstanceID, projectID, srcInstanceID, dstInstanceID,
						fmt.Sprintf("%s.backup_time_points.0.end_timestamp", timePointDatasourceName),
					),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunMysqlBackupRecoveryByTaskID(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_backup_recovery." + rnd
	resourceFile := "resource_ctyun_mysql_backup_recovery_task_id.tf"

	backupsDatasourceName := "data.ctyun_mysql_backups." + dnd
	backupsDatasourceFile := "datasource_ctyun_mysql_backups_none_backup_name.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	srcInstanceID := "e5ad1c553e394bc891c5bf8fc58be191"
	dstInstanceID := "e5ad1c553e394bc891c5bf8fc58be191"

	taskID := dependence.taskID
	pageSize := 10
	pageNo := 1
	//toTimepoint := os.Getenv("CTYUN_MYSQL_BACKUP_RECOVERY_TIME")

	// 等待函数
	//wait10Seconds := func() {
	//	t.Logf("等待10秒...")
	//	time.Sleep(10 * time.Second)
	//}

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
			// 1. 验证datasource
			{
				Config: utils.LoadTestCase(backupsDatasourceFile, dnd, srcInstanceID, pageNo, pageSize),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(backupsDatasourceName, "backup_list.#"),
				),
			},
			// 1. 创建备份恢复任务
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, srcInstanceID, projectID, srcInstanceID, dstInstanceID, taskID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "inst_id", srcInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "src_inst_id", srcInstanceID),
					resource.TestCheckResourceAttr(resourceName, "dst_inst_id", dstInstanceID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 销毁
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, srcInstanceID, projectID, srcInstanceID, dstInstanceID, taskID),
				Destroy: true,
			},
		},
	})
}
