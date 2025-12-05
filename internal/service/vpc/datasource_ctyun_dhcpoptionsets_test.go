package vpc_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunDhcpOptionSetsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpOptionSetsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ctyun_dhcpoptionsets.test", "dhcpoptionsets.#"),
					resource.TestCheckResourceAttrSet("data.ctyun_dhcpoptionsets.test", "total_count"),
				),
			},
		},
	})
}

func TestAccCtyunDhcpOptionSetsDataSource_withQuery(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpOptionSetsDataSourceConfigWithQuery(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ctyun_dhcpoptionsets.test", "dhcpoptionsets.#"),
				),
			},
		},
	})
}

func testAccDhcpOptionSetsDataSourceConfig() string {
	return `
data "ctyun_dhcpoptionsets" "test" {
  page_no   = 1
  page_size = 10
}
`
}

func testAccDhcpOptionSetsDataSourceConfigWithQuery() string {
	return `
data "ctyun_dhcpoptionsets" "test" {
  query_content = "test"
  page_no       = 1
  page_size     = 10
}
`
}
