output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "security_group_id" {
  value = ctyun_security_group.security_group_test.id
}

output "eip_address" {
  value = ctyun_eip.eip_test.address
}

output "redis_version" {
  value = local.spec.version
}

output "redis_engine_edition" {
  value = local.spec.series_code
}