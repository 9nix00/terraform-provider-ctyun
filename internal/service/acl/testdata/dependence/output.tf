output "vpc_id" {
  value = local.real_vpc_id
}

output "subnet_id" {
  value = local.real_subnet_id
}
output "subnet_id2" {
  value = local.real_subnet_id2
}

output "acl_id" {
  value = ctyun_acl.acl_test.id
}

output "acl_id2" {
  value = ctyun_acl.acl_subnet_test.id
}