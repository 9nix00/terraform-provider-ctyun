package kafka_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunKafkaUsers(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_kafka_user." + rnd
	datasourceName := "data.ctyun_kafka_users." + dnd
	resourceFile := "resource_ctyun_kafka_user.tf"
	datasourceFile := "datasource_ctyun_kafka_users.tf"

	initName := "init-kafka-user-" + rnd
	prodInstId := dependence.instanceID
	initPassword := "sad231Dwwww"
	updatePassword := "sad231Dwwasd"
	topicName := dependence.topicName
	aclInfo := fmt.Sprintf(`permission_info = [{
operation = "READ"
topic = "%s"}]`, topicName)

	aclInfoUpdate := fmt.Sprintf(`permission_info = [{
permission = "DENY"
operation = "READ"
topic = "%s"},{
permission = "ALLOW"
operation = "WRITE"
topic = "%s"}]`, topicName, topicName)

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, initPassword, aclInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, updatePassword, aclInfoUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "password", updatePassword),
					resource.TestCheckResourceAttr(resourceName, "permission_info.#", "2"),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, updatePassword, aclInfoUpdate) +
					utils.LoadTestCase(datasourceFile, dnd, initName, prodInstId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "users.0.name", initName),
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
					return fmt.Sprintf("%s,%s,%s,%s", prodInstId, regionId, name, password), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"id", "permission_info"},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, updatePassword, aclInfoUpdate) +
					utils.LoadTestCase(datasourceFile, dnd, initName, prodInstId),
				Destroy: true,
			},
		},
	})
}
