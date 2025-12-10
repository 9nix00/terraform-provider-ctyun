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

resource "ctyun_acl" "example" {
  vpc_id        = "example-vpc-id"
  name          = "example-acl"
  description   = "Example ACL created for demonstration"
  enabled       = "enable"
  apply_to_public_lb = false
}