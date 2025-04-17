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

resource "ctyun_nat" "nat_test"{
    region_id = "200000002530"
    vpc_id = "vpc-wf029jgx2d"
    spec = 1
    name = "nat-terraform-test"
    description = "terraform测试"
    cycle_type = "on_demand"
 #   cycle_count = ,
    az_name = "az1"
 #   pay_voucher_price = "",
 #   project_id = "",
}