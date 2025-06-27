terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  env = "prod"
}

data "ctyun_ebms" "test" {
  region_id            = "200000001852"
  az_name              = "cn-huabei2-tj-3a-public-ctcloud"
  # az_name =             "cn-huabei2-tj1A-public-ctcloud"
}

output "ctyun_ebms_test" {
  value = data.ctyun_ebms.test
}

