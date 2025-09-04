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


data "ctyun_mongodb_instances" "test" {
  prod_inst_name = "db-aswd"
}

output "t" {
  value = data.ctyun_mongodb_instances.test
}
