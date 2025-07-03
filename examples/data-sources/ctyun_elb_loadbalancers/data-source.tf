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

data "ctyun_elb_loadbalancers" "test" {
}

 output "ctyun_elb_loadbalancers_test"{
     value = data.ctyun_elb_loadbalancers.test
 }