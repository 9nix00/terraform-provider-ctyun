package ebs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunBackup(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ebs_backup." + rnd
	datasourceName := "data.ctyun_ebs_backups." + dnd
	resourceFile := "resource_ctyun_ebs_backup.tf"
	datasourceFile := "datasource_ctyun_ebs_backups.tf"

	initName := "init-backup"
	updatedName := "updated-backup-" + rnd
	diskId := dependence.ebsID
	repositoryID := "671f67c4-6131-4154-8c1d-7c5b82edd1eb"

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
			// 创建
			{
				Config: utils.LoadTestCase(resourceFile, rnd, repositoryID, diskId, initName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, repositoryID, diskId, updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, repositoryID, diskId, updatedName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "backups.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "backups.0.name", updatedName),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"full_backup", "backup_status", "updated_time", "task_id",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, repositoryID, diskId, updatedName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	})
}
