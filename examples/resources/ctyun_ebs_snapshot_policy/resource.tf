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


resource "ctyun_ebs_snapshot_policy" "test" {
    name           = "test"
    repeat_weekdays            = "0,1,2"
    repeat_times            = "0,1,2"
    retention_time        = 2
    is_enabled  = true
}

