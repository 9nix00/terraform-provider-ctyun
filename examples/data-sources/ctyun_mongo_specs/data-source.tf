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

data "ctyun_mongodb_specs" "test" {
  instance_type = "1"
}

output "t" {
  value = data.ctyun_mongodb_specs.test
}