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


data "ctyun_mysql_specs" "test" {
  instance_series = "S"
}


