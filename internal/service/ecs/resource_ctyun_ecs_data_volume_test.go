package ecs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunEcsDataVolume(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_ecs_data_volume." + rnd
	resourceFile := "resource_ctyun_ecs_data_volume.tf"
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
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ecsID,
					dependence.ebsID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.ecsID),
					resource.TestCheckResourceAttr(resourceName, "ebs_ids.0", dependence.ebsID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ecsID,
					dependence.ebsID2,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.ecsID),
					resource.TestCheckResourceAttr(resourceName, "ebs_ids.0", dependence.ebsID2),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ecsID,
					dependence.ebsID2,
				),
				Destroy: true,
			},
		},
	},
	)
}
