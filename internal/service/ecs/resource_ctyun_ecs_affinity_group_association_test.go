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
			// Read testing
			{
				Config: `
provider "ctyun" {
  region_id  = "bb9fdb42056f11eda1610242ac110002"
  az_name    = "cn-huadong1-jsnj1A-public-ctcloud"
}

resource "ctyun_ecs_affinity_group_association" "test" {
  instance_id = "ae432721-61bf-45b7-b207-7e3256c1c2d6"
  affinity_group_id = "e9d3239a-207a-4006-aa84-3945265bac27"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", "ae432721-61bf-45b7-b207-7e3256c1c2d6"),
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
