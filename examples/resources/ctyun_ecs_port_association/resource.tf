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


resource "ctyun_port" "port" {
  name = "port_for_association_test"
  network_id = "a0c0c0c0-c0c0-c0c0-c0c0-c0c0c0c0c0c0"
  security_group_id = "a0c0c0c0-c0c0-c0c0-c0c0-c0c0c0c0c0c0"
  admin_state_up = true
}
resource "ctyun_ecs_port_association" "ecs_port_for_association_test" {
  instance_id          =  "ae432721-61bf-45b7-b207-7e3256c1c2d6"
  network_interface_id = ctyun_port.port.id
}



