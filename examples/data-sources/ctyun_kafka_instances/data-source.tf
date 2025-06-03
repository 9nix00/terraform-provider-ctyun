terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

data "ctyun_kafka_instances" "tbidgqvfbs" {
  instance_name = "123"
}
