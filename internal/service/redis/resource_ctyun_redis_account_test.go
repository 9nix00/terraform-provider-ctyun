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

func TestAccCtyunRedisAccounts(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_redis_account." + rnd
	datasourceName := "data.ctyun_redis_accounts." + dnd
	resourceFile := "resource_ctyun_redis_account.tf"
	datasourceFile := "datasource_ctyun_redis_accounts.tf"

	initName := "init_redis_account-" + rnd
	prodInstId := dependence.instanceId
	initPassword := "sad231Dwwww"
	updatePassword := "sad231Dwwasd"

	initPrivilege := "ro"
	updatePrivilege := "rw"

	initDescription := "Description1"
	updateDescription := "Description2"

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, initPassword, initPrivilege, initDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, updatePassword, updatePrivilege, updateDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "password", updatePassword),
					resource.TestCheckResourceAttr(resourceName, "privilege", updatePrivilege),
					resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, updatePassword, updatePrivilege, updateDescription) +
					utils.LoadTestCase(datasourceFile, dnd, prodInstId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "accounts.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "accounts.0.name", initName),
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
					password := ds.Attributes["password"]
					privilege := ds.Attributes["privilege"]
					return fmt.Sprintf("%s,%s,%s,%s,%s", prodInstId, regionId, name, password, privilege), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"id", "permission_info"},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, updatePassword, updatePrivilege, updateDescription) +
					utils.LoadTestCase(datasourceFile, dnd, prodInstId),
				Destroy: true,
			},
		},
	})
}
