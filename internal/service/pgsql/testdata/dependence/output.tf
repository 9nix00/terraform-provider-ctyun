output "vpc_id" {
  value = local.real_vpc_id
}

output "subnet_id" {
  value = local.real_subnet_id
}

output "security_group_id1" {
  value = local.real_security_group_id1
}
output "security_group_id2" {
  value = local.real_security_group_id2
}

output "eip_id" {
  value = ctyun_eip.eip_test.id
}

output "eip_address" {
  value = ctyun_eip.eip_test.address
}
