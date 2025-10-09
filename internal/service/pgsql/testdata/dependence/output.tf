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

# output "eip_id" {
#   value = ctyun_eip.eip_test.id
# }
#
# output "pgsql_id" {
#   value = ctyun_postgresql_instance.test.id
# }

output "az_name" {
  value = data.ctyun_zones.az.zones[0]
}

# output "param_templates"{
#   value = data.ctyun_postgresql_param_templates.param_templates.parameter_templates[0].id
# }