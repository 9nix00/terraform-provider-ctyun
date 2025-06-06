output "vpc_id" {
  value = local.real_vpc_id
}

output "subnet_id" {
  value = local.real_subnet_id
}

output "flavor_name" {
  value = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name
}

output "cluster_id" {
  value = ctyun_ccse_cluster.test.id
}