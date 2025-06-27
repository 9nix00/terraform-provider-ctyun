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

data "ctyun_zos_bucket_objects" "test" {

  bucket = "acc.te21fdsfdasfdsdwqedwed23e-asd.1"
}