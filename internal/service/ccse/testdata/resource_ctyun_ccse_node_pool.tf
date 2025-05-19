resource "ctyun_ccse_node_pool" "%[1]s" {
  node_pool_name           = "%[2]s"
  cycle_type              = "%[10]s"
  %[11]s
  auto_renew_status        = %[3]d
  visibility_post_host_script = "%[4]s"
  visibility_host_script = "%[5]s"
  instance_type            = "ecs"
  mirror_name             = "CTyunOS-23.01-CCND_CCSE_40_08-x86_64"
  mirror_id                = "3f80d8c0-8eb5-4afa-a506-13ba68b61872"
  mirror_type              = 1
  password                 = "P@ss2wsx"
  max_pod_num              = 110
  item_def_name            = "%[12]s"
  cluster_id               = "%[13]s"
  sys_disk = {
    type = "%[6]s"
    size = %[7]d
  }

  data_disks = [
    {
      type = "%[8]s"
      size = %[9]d
    }
  ]
}
