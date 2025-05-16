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

data "ctyun_zos_buckets" "test" {
  page_no = 1
  page_size = 10
}

output "ctyun_test" {
  value = data.ctyun_zos_buckets.test
}

