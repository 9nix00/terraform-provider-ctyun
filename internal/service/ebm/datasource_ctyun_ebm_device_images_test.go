package ebm_test

import (
	"terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunEbmDeviceImages(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "ctyun" {
  env = "prod"
}

data "ctyun_ebm_device_images" "test" {
  region_id = "200000001852"
  az_name = "cn-huabei2-tj1A-public-ctcloud"
  device_type = "physical.s5.2xlarge4"
  os_type = "linux"
  image_type = "standard"
  image_uuid = "im-idxitiryuxevcr87wknzxadj0nvk"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_images.test", "images.#", "1"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_images.test", "images.0.image_uuid", "im-idxitiryuxevcr87wknzxadj0nvk"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_images.test", "images.0.name_en", "CTyunOS"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_images.test", "images.0.image_type", "standard"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_images.test", "images.0.os.uuid", "o-ryj4xogjs2mcqbgfngwwffsx0vid"),
				),
			},
		},
	})
}
