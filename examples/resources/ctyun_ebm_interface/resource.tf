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

resource "ctyun_ebm_interface" "test" {
  security_group_ids = ["sg-t0ae11aig1"]
  instance_id = "ss-uadmwtxinfp4tkbhvwp52vnzl2kn"
  ipv4 = "192.168.0.13"
  subnet_id = "subnet-43z7cqmjlp"
}
