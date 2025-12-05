data "ctyun_dhcpoptionsets" "test" {
  page_no   = 1
  page_size = 10
}

output "dhcpoptionsets" {
  value = data.ctyun_dhcpoptionsets.test.dhcpoptionsets
}