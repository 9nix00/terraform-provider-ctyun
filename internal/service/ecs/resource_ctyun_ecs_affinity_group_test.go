package ecs_test

import (
	"fmt"
	"terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunAffinityGroup(t *testing.T) {
	resourceName := "ctyun_ecs_affinity_group.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
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

				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}
