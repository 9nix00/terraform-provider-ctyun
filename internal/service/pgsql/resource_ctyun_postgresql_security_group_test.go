package pgsql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunPostgresqlSecurityGroup(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_security_group." + rnd
	resourceFile := "resource_ctyun_postgresql_security_group.tf"

	//instanceId := dependence.PgsqlID
	instanceId := "24c876ba30c04b59a5417a0a39500797"
	securityGroupIds := "\"sg-0fabnzup10\""
	securityGroupIdsUpdate := "\"sg-0fabnzup10\", \"sg-mabpconbi8\""

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
				Config: utils.LoadTestCase(resourceFile, rnd, instanceId, securityGroupIds),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile, rnd, instanceId, securityGroupIdsUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// destroy
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, instanceId, securityGroupIdsUpdate),
				Destroy: true,
			},
		},
	})
}
