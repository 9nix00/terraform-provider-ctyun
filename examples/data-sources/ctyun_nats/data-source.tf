
terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

 provider "ctyun" {
   region_id            = "200000002530"
 }

 data "ctyun_nats" "test"{
     region_id = "200000002530"
 }

 output "ctyun_nat_test"{
     value = data.ctyun_nats.test
 }