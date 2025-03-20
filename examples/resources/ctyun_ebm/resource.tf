terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  region_id            = "bb9fdb42056f11eda1610242ac110002"
  az_name              = "cn-huadong1-jsnj1A-public-ctcloud"
  env                  = "prod"
}

resource "ctyun_ebm" "ebm_test" {
  device_type = "physical.s5.2xlarge4"
  instance_name = "ebm-0312-tf"
  hostname = "ebm-0317-tf"
  image_uuid = "im-xevpi6apqilz1bixmogofyref9qm"
  password = "P@ss132345"
  security_group_id = "sg-vrp4x1lm7p"
  vpc_id = "vpc-5o8oe0oci6"
  ext_ip = "auto_assign"
  system_volume_raid_uuid = "r-wtzluqacgzzxgunnabdkpnpjew3d"
  data_volume_raid_uuid = "r-qytwf9r5h0yn9x4evjkyr0n1cwyb"
  instance_charge_type = "order_on_demand"
  # cycle_type = "month"
  # cycle_count = 2
  band_width = "100"
  network_card_list = [{
    master = true,
    subnet_id = "subnet-n7zbsy4b91"
  }]
}