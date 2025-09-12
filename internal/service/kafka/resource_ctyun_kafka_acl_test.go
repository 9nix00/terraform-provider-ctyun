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

func TestAccCtyunKafkaAcl(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_kafka_acl." + rnd
	datasourceName := "data.ctyun_kafka_acls." + dnd
	resourceFile := "resource_ctyun_kafka_acl.tf"
	datasourceFile := "datasource_ctyun_kafka_acls.tf"

	initName := "init-kafka-acl-" + rnd
	prodInstId := dependence.instanceID

	initUseNewTopic := "2"
	updateUseNewTopic := "1"
	userName := dependence.userName

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, initUseNewTopic, userName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "use_new_topic", initUseNewTopic),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, updateUseNewTopic, userName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "use_new_topic", updateUseNewTopic),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, updateUseNewTopic, userName) +
					utils.LoadTestCase(datasourceFile, dnd, initName, prodInstId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "acls.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "acls.0.name", initName),
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
					useNewTopic := ds.Attributes["use_new_topic"]
					return fmt.Sprintf("%s,%s,%s,%s", prodInstId, regionId, name, useNewTopic), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"id"},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, updateUseNewTopic, userName) + utils.LoadTestCase(datasourceFile, dnd, initName, prodInstId),
				Destroy: true,
			},
		},
	})
}
