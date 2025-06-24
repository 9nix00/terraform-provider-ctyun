resource "ctyun_elb_acl" "%[1]s" {
  name = "tf_acl"
  source_ips = ["127.0.0.1/32","192.168.0.0/16","192.168.10.0"]
}
