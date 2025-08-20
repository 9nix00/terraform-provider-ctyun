package ports_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunEcsPortAssociation_basic(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	name := "ctyun_ecs_port_association." + rnd
	configFile := "resource_ctyun_ecs_port_association.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(configFile, rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "region_id"),
					resource.TestCheckResourceAttrSet(name, "project_id"),
					resource.TestCheckResourceAttr(name, "instance_id", "7fe3c19f-d364-8450-72c0-32ddd818fd3c"),
					resource.TestCheckResourceAttr(name, "network_interface_id", "port-gyvm3rewlx"),
				),
			},
		},
	})
}
