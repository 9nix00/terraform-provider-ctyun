package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMongodbBackupParamTemplate_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_mongodb_param_template." + rnd
	resourceFile := "resource_ctyun_mongodb_param_template.tf"

	parameter_group_name := utils.GenerateRandomString()
	engine := "4.0"
	description := "MongoDB备份测试"
	node_type := "Mongod"

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
				Config: utils.LoadTestCase(resourceFile, rnd, parameter_group_name, engine, description, node_type),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "engine", engine),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, parameter_group_name, engine, description, node_type),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "engine", engine),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile, rnd, parameter_group_name, engine, description, node_type),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "engine", engine),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, parameter_group_name, engine, description, node_type),
				Destroy: true,
			},
		},
	})
}
