terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

data "ctyun_rabbitmq_specs" "test"{

}
