output "vpc_id" {
  value = local.real_vpc_id
}

output "bandwidth_id" {
  value = ctyun_bandwidth.bandwidth_test.id
}

output "eip_id" {
  value = ctyun_eip.eip_test.id
}

output "ecs_id" {
  value = ctyun_ecs.ecs_test.id
}

output "security_group_id" {
  value = ctyun_security_group.security_group_test.id
}
# output "data_disk_id" {
#   value = ctyun_ebs.data_disk_test.id
# }
#

output "network_interface_id" {
  value =ctyun_port.port_test.id
}
output "subnet_id" {
  value = local.real_subnet_id
}

output "vip_id" {
  value =ctyun_vip.vip_test.id
}

