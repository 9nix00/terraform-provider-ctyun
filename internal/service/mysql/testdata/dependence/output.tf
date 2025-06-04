output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}
output "security_group_id" {
  value = ctyun_security_group.test.id
}
output "eip_id" {
  value = ctyun_eip.eip_test.id
}

output "eip_address" {
  value = ctyun_eip.eip_test.address
}
output "mysql_id" {
  value = ctyun_mysql_instance.mysql_test.inst_id
}