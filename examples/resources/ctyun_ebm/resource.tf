terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}


provider "ctyun" {
  region_id            = "200000001852"
  # az_name              = "cn-huabei2-tj-3a-public-ctcloud"
  az_name =             "cn-huabei2-tj1A-public-ctcloud"
  env                  = "prod"
}

data "ctyun_ebm_device_types" "test" {
}

data "ctyun_ebm_device_raids" "system_raid" {
  device_type = data.ctyun_ebm_device_types.test.device_types[0].device_type
  volume_type = "system"
}


data "ctyun_ebm_device_raids" "data_raid" {
  device_type = data.ctyun_ebm_device_types.test.device_types[0].device_type
  volume_type = "data"
}

resource "ctyun_ebm" "ebm_test" {
  device_type = data.ctyun_ebm_device_types.test.device_types[0].device_type
  instance_name = "ebm-0321-tf"
  hostname = "ebm-03221-tf"
  image_uuid = "im-xevpi6apqilz1bixmogofyref9qm"
  password = "P@ss132345"
  security_group_id = "sg-t0ae11aig1"
  vpc_id = "vpc-6zxqwrg1r6"
  ext_ip = "not_use"
  system_volume_raid_uuid = length(data.ctyun_ebm_device_raids.system_raid.raids) > 0 ? data.ctyun_ebm_device_raids.system_raid.raids[0].uuid : ""
  data_volume_raid_uuid = length(data.ctyun_ebm_device_raids.data_raid.raids) > 0 ? data.ctyun_ebm_device_raids.data_raid.raids[0].uuid : ""
  instance_charge_type = "order_on_demand"
  status = "running"
  # cycle_type = "month"
  # cycle_count = 3
  # band_width = "100"
  disk_list = data.ctyun_ebm_device_types.test.device_types[0].cloud_boot ? [{
    disk_type = "system"
    size = "100"
    type = "sata"
  }] : []
  network_card_list = [{
    master = true,
    subnet_id = "subnet-43z7cqmjlp"
  }]
}