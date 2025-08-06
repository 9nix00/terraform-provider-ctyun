package ecs_test

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

	resourceName := "ctyun_ecs_backup." + rnd
	datasourceName := "data.ctyun_ecs_backups." + dnd
	resourceFile := "resource_ctyun_ecs_backup.tf"
	datasourceFile := "datasource_ctyun_ecs_backups.tf"

	initName := "init-backup"
	updatedName := "updated-backup"
	instanceId := dependence.ecsID
	repositoryID := "a94ee44e-8a86-42f6-9986-4e9df5cb3c18"

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
				Config: utils.LoadTestCase(resourceFile, rnd, repositoryID, instanceId, initName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_backup_name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, repositoryID, instanceId, updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_backup_name", updatedName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "backups.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "backups.0.instance_backup_name", updatedName),
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
					"project_id",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	})
}
