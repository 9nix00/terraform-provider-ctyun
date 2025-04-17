package nat_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"os"
	"terraform-provider-ctyun/internal/service"
	"testing"
)

func TestAccNewCtyunNatDataSource(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
				 provider "ctyun" {
				   region_id = "200000002530"
				 }
				
				 data "ctyun_nats" "test"{
					 region_id = "bb9fdb42056f11eda1610242ac110002"
				 }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ctyun_nats.test", "nats.#", "1"),
					resource.TestCheckResourceAttr("data.ctyun_nats.test", "nats.0.name", "nat-codearts-dev"),
					resource.TestCheckResourceAttr("data.ctyun_nats.test", "nats.0.status", "2"),
				),
			},
		},
	})
}
