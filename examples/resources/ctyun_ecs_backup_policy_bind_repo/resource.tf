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

resource "ctyun_ecs_backup_policy_bind_repo" "test" {
  id = "a4b793881bbd42edaa6a0002900e5819"
  repository_id = "0cd13a89-5ada-42a7-95e8-60fb9705eecc"
}
