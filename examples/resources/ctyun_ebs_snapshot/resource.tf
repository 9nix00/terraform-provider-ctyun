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

resource "ctyun_ebs" "ebs_test" {
  name       = "ebs-test"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
}

resource "ctyun_ecs_snapshot" "test" {
  snapshot_name = "tf-test-group"
  disk_id = ctyun_ebs.ebs_test.id
  retention_policy = "forever"
}
