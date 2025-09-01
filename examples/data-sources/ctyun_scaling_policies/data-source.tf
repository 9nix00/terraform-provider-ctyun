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
data "ctyun_scaling_policies" "scaling_policies_test" {
  group_id = 109737
}

