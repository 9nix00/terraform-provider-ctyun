package ecs_test

import (
	"terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunAffinityGroup(t *testing.T) {
	resourceName := "ctyun_ecs_affinity_group.test"
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

resource "ctyun_ecs_affinity_group" "test" {
  affinity_group_name = "tf-test-group"
  affinity_group_policy = "anti-affinity"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "affinity_group_name", "tf-test-group"),
					resource.TestCheckResourceAttr(resourceName, "affinity_group_policy", "anti-affinity"),
					resource.TestCheckResourceAttrSet(resourceName, "affinity_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: `
provider "ctyun" {
  region_id  = "bb9fdb42056f11eda1610242ac110002"
  az_name    = "cn-huadong1-jsnj1A-public-ctcloud"
}

resource "ctyun_ecs_affinity_group" "test" {
  affinity_group_name = "tf-test-group1"
  affinity_group_policy = "anti-affinity"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "affinity_group_name", "tf-test-group1"),
					resource.TestCheckResourceAttr(resourceName, "affinity_group_policy", "anti-affinity"),
				),
			},
			{
				Config: `
provider "ctyun" {
  region_id  = "bb9fdb42056f11eda1610242ac110002"
  az_name    = "cn-huadong1-jsnj1A-public-ctcloud"
}

resource "ctyun_ecs_affinity_group" "test" {
  affinity_group_name = "tf-test-group2"
  affinity_group_policy = "anti-affinity"
}

data "ctyun_ecs_affinity_groups" "test" {
  affinity_group_id = ctyun_ecs_affinity_group.test.affinity_group_id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ctyun_ecs_affinity_groups.test", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.ctyun_ecs_affinity_groups.test", "groups.0.affinity_group_name", "tf-test-group2"),
					resource.TestCheckResourceAttr("data.ctyun_ecs_affinity_groups.test", "groups.0.affinity_group_policy", "anti-affinity"),
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
