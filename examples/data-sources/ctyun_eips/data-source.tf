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

resource "ctyun_eip" "eip_test1" {
  name        = "tf-eip-test1"
  bandwidth   = 2
  cycle_type = "on_demand"
  demand_billing_type = "upflowc"
}

data "ctyun_eips" "test" {
  ids = ctyun_eip.eip_test1.id
}

output "ctyun_test" {
  value = data.ctyun_eips.test
}

