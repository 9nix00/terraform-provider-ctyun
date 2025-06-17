terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

data "ctyun_rabbitmq_instances" "tbidgqvfbs" {
  instance_id ="8d839e64a4314edb8121d0d1f69b8b19"
}

output "list" {
  value = data.ctyun_rabbitmq_instances.tbidgqvfbs
}