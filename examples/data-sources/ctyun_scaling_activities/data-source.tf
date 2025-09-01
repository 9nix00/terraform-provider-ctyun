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
data "ctyun_scaling_activities" "activities_test" {
  group_id = 109737
}
