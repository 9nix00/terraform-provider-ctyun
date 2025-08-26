output "vpc_id" {
  value = local.real_vpc_id
}

output "vpc_id1" {
  value = local.read_iaas_vpc_id
}

output "subnet_id" {
  value = local.real_subnet_id
}

output "sfs_uid" {
  value = ctyun_sfs.sfs_test.id
}

output "sfs_permission_group_id" {
  value = ctyun_sfs_permission_group.group_test.id
}

output "sfs_permission_group_id1" {
  value = ctyun_sfs_permission_group.group_test1.id
}