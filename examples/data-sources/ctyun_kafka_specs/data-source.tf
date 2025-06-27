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

data "ctyun_kafka_specs" "test" {
  
}

output "t" {
  value = data.ctyun_kafka_specs.test
}