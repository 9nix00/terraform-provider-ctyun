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

resource "ctyun_zos_bucket" "foo" {
  bucket = "acc.te21"
  acl = "public-read"
  storage_type = "STANDARD_IA"
}