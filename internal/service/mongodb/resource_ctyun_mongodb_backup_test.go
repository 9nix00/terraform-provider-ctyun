package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMongodbBackup_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_mongodb_backup." + rnd
	resourceFile := "resource_ctyun_mongodb_backup.tf"

	instance_id := dependence.mongodbID
	backupName := utils.GenerateRandomString()
	description := "MongoDB备份测试"

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
			// 基本功能验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, backupName, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "backup_name", backupName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, backupName, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "backup_name", backupName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, backupName, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "backup_name", backupName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, instance_id, backupName, description),
				Destroy: true,
			},
		},
	})
}
