resource "ctyun_ebm" "%[1]s" {
  instance_name = "%[2]s"
  hostname = "%[3]s"
  password = "%[4]s"
  status = "%[5]s"
  ext_ip = "not_use"
  cycle_type = "on_demand"
  device_type = "%[6]s"
  image_uuid = "%[7]s"
  security_group_ids = %[8]s
  vpc_id = "%[9]s"
  system_volume_raid_uuid = "%[10]s"
  data_volume_raid_uuid = "%[11]s"
  disk_list = %[12]v ? [{
    disk_type = "system"
    size = "100"
    type = "sata"
  }] : []
  network_card_list = [{
    master = true,
    subnet_id = "%[13]s"
  }]
}
