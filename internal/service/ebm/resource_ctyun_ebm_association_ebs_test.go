package ebm_test

import (
	"terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunEbmAssociationEbs(t *testing.T) {
	resourceName := "ctyun_ebm_association_ebs.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "ctyun" {
  region_id            = "200000001852"
  az_name              = "cn-huabei2-tj-3a-public-ctcloud"
  env                  = "prod"
}

resource "ctyun_ebs" "ebs_test" {
  name       = "ebs-tf-test-0402"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
}

resource "ctyun_ebm_association_ebs" "test" {
  ebs_id = ctyun_ebs.ebs_test.id
  instance_id = "ss-uadmwtxinfp4tkbhvwp52vnzl2kn"
}
`,
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
