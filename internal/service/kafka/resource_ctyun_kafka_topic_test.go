package kafka_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunKafkaTopics(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_kafka_topic." + rnd
	datasourceName := "data.ctyun_kafka_topics." + dnd
	resourceFile := "resource_ctyun_kafka_topic.tf"
	datasourceFile := "datasource_ctyun_kafka_topics.tf"

	initName := "init-kafka-topic-" + rnd
	prodInstId := dependence.instanceID
	initPartitionNum := 1
	updatePartitionNum := 2

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, initPartitionNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, updatePartitionNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "partition_num", strconv.Itoa(updatePartitionNum)),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, updatePartitionNum) +
					utils.LoadTestCase(datasourceFile, dnd, initName, prodInstId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "topics.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "topics.0.name", initName),
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
				ImportStateVerifyIgnore: []string{"id"},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, initName, prodInstId, updatePartitionNum) + utils.LoadTestCase(datasourceFile, dnd, initName, prodInstId),
				Destroy: true,
			},
		},
	})
}
