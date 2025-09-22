package redis_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunRedisBackups(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_redis_backup." + rnd
	datasourceName := "data.ctyun_redis_backups." + dnd
	resourceFile := "resource_ctyun_redis_backup.tf"
	datasourceFile := "datasource_ctyun_redis_backups.tf"

	prodInstId := dependence.instanceId

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
				Config: utils.LoadTestCase(resourceFile, rnd, prodInstId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, prodInstId) +
					utils.LoadTestCase(datasourceFile, dnd, prodInstId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "rows.#", "1"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					regionId := ds.Attributes["region_id"]
					prodInstId := ds.Attributes["prod_inst_id"]
					name := ds.Attributes["name"]
					return fmt.Sprintf("%s,%s,%s", prodInstId, regionId, name), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"id", "download_urls"},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, prodInstId) +
					utils.LoadTestCase(datasourceFile, dnd, prodInstId),
				Destroy: true,
			},
		},
	})
}
