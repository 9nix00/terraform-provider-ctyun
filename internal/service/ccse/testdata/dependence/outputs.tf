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

output "chart_name" {
  value = local.chart_name
}

output "chart_version1" {
  value = local.chart_version1
}

output "chart_version2" {
  value = local.chart_version2
}

output "chart_values_yaml" {
  value = jsonencode(data.ctyun_ccse_plugin_market.test1.values)
}

output "chart_values_json" {
  value = jsonencode(data.ctyun_ccse_plugin_market.test2.values)
}