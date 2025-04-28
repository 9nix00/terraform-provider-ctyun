output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "ecs_id" {
  value = ctyun_ecs.ecs_test.id
}

output "vpce_server_id" {
  value = ctyun_vpce_server.vpce_server_test.id
}