terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  region_id            = "200000001852"
  az_name              = "cn-huabei2-tj-3a-public-ctcloud"
  env                  = "prod"
}

resource "ctyun_ebs" "ebs_test" {
  name       = "ebs-test"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
}

resource "ctyun_ebm_association_ebs" "test" {
  ebs_id = ctyun_ebs.ebs_test.id
  instance_id = "ss-3dcvp69qfw8enincfcwkhxeubjq0"
}
