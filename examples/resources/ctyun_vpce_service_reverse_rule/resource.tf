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

resource "ctyun_vpce_service_reverse_rule" "test" {
  endpoint_service_id = "xx"
  endpoint_id = "xxx"
  transit_ip = "192.168.1.3"
  transit_port = 1
  target_ip = "192.168.1.4"
  target_port = 2
  protocol = "TCP"
}