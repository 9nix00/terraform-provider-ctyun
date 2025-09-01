output "vpc_id" {
  value = local.real_vpc_id
}

output "subnet_id" {
  value = local.real_subnet_id
}

output "security_group_id" {
  value = local.real_security_group_id
}
output "mongodb_id"{
  value = ctyun_mongodb_instance.mongodb_eip.id
}