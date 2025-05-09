package ecs_test

import (
	"terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunAffinityGroupAssociation(t *testing.T) {
	resourceName := "ctyun_ecs_affinity_group_association.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "ctyun" {
  region_id  = "bb9fdb42056f11eda1610242ac110002"
  az_name    = "cn-huadong1-jsnj1A-public-ctcloud"
}

resource "ctyun_ecs_affinity_group_association" "test" {
  instance_id = "25b2fc7e-2c09-4428-b5d3-81402ceaedfc"
  affinity_group_id = "e9d3239a-207a-4006-aa84-3945265bac27"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", "25b2fc7e-2c09-4428-b5d3-81402ceaedfc"),
					resource.TestCheckResourceAttr(resourceName, "affinity_group_id", "e9d3239a-207a-4006-aa84-3945265bac27"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}
