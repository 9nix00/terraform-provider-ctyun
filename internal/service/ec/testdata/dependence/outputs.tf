
output "ctyun_express_connect_id" {
  value = ctyun_express_connect.express_connect_dependence.id
}


output "ctyun_ec_cloud_gateway_id"{
  value = ctyun_ec_cloud_gateway.cloud_gateway_dependence.id
}


output "vpc_id" {
  value = local.real_vpc_id
}

output "subnet_id" {
  value = local.real_subnet_id
}

output "subnet_id2" {
  value = local.real_subnet_id2
}
output "rtb_id"{
  value = ctyun_ec_cloud_gateway.cloud_gateway_dependence.rtb_id
}

output "vpc_instance_vpc_id"{
  value = ctyun_express_connect_vpc_instance.instance_test.vpc_id
}

output "cgw_id1" {
  value = ctyun_ec_cloud_gateway.cloud_gateway_xinan1.id
}
output "cgw_id2" {
  value = ctyun_ec_cloud_gateway.cloud_gateway_huabei2.id
}

output "packet_id" {
  value = ctyun_ec_packet.packet_test.id
}
# output "region_peer_id" {
#   value = ctyun_express_connect_region_peer.region_peer_test.id
# }