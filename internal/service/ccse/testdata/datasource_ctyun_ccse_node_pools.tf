data "ctyun_ccse_node_pools" "%[1]s" {
  cluster_id               = "%[2]s"
  node_pool_name           = %[3]s
}
