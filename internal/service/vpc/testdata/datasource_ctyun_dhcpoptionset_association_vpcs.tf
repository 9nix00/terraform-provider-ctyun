data "ctyun_dhcpoptionset_association_vpcs" "test" {
  dhcp_option_sets_id  = "%[1]s"
  page_no              = 1
  page_size            = 10
}

output "vpcs" {
  value = data.ctyun_dhcpoptionset_association_vpcs.test.vpcs
}