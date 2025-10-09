package mongodb_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccMongodbRestartDb_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_mongodb_restart_db." + rnd
	resourceFile := "resource_ctyun_mongodb_restart_db.tf"

	instance_id := "fafae8192ba049c9bc9de83f1cd8361b"
	instance_id2 := "1538d18eb4a94c11a99e584133594665"
	instance_id3 := "f68c4532bd4e4659986162f966fc34e0"
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
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},

			// 基本功能验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 基本功能验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id3),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, instance_id),
				Destroy: true,
			},
		},
	})
}
