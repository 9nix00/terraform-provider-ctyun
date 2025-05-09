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
  bucket = "acc.te21fdsfdasfdsdwqedwed23e-asd.1"
}

resource "ctyun_zos_bucket_object" "test" {
  bucket = ctyun_zos_bucket.foo.bucket
  key = "aex.txt"
  content = "abcdw"
}