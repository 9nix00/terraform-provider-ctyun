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

resource "ctyun_sfs_permission_group" "sfs_permission_group_test" {
  name = "permission-group_example"
  description = "创建sfs规则组"
}

data "ctyun_sfs_permission_rules" "%[1]s" {
  permission_group_fuid = ctyun_sfs_permission_group.sfs_permission_group_test.id
}
