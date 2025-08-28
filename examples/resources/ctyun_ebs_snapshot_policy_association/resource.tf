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

resource "ctyun_ebs_snapshot_policy_association" "test" {
  snapshot_policy_id = "6f017e65-5340-4348-a2da-07c9aae44e5f"
  disk_id_list = "ae432721-61bf-45b7-b207-7e3256c1c2d6"
}

