output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "security_group_id" {
  value = data.ctyun_ebm_device_types.test.device_types[0].smart_nic_exist ? format("[\"%s\"]",ctyun_security_group.security_group_test.id)  : "[]"
}

output "device_type" {
  value = data.ctyun_ebm_device_types.test.device_types[0].device_type
}

output "smart_nic_exist" {
  value = data.ctyun_ebm_device_types.test.device_types[0].smart_nic_exist ? "true" : "false"
}

output "support_cloud" {
  value = data.ctyun_ebm_device_types.test.device_types[0].support_cloud ? "true" : "false"
}

output "cloud_boot" {
  value = data.ctyun_ebm_device_types.test.device_types[0].cloud_boot ? "true" : "false"
}

output "system_raid" {
  value  = length(data.ctyun_ebm_device_raids.system_raid.raids) > 0 ? data.ctyun_ebm_device_raids.system_raid.raids[0].uuid : ""
}

output "data_raid" {
  value  = length(data.ctyun_ebm_device_raids.data_raid.raids) > 0 ? data.ctyun_ebm_device_raids.data_raid.raids[0].uuid : ""
}

output "image_uuid" {
  value = data.ctyun_ebm_device_images.test.images[0].image_uuid
}

output "ebs_id" {
  value = ctyun_ebs.ebs_test.id
}
