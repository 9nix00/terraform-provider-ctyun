resource "ctyun_vpc_peer_connection_route" "%[1]s" {
  ip_version   = "%[2]s"
  next_hop_id  = "%[3]s"
  vpc_id       = "%[4]s"
  destination  = "%[5]s"
}

resource "ctyun_vpc_route_table_rule" "%[1]s"{
  ip_version   = "%[2]s"
  next_hop_id = "%[3]s"
  destination = "%[2]s"
  description = "%[3]s"
  next_hop_type = "igw"
  route_table_id = ctyun_vpc_route_table.route.id
}

