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

data "ctyun_nats" "test" {
  nat_gateway_id = "natgw-asdsmh8scy"
}

output "ctyun_nat_test"{
  value = data.ctyun_nats.test
}

