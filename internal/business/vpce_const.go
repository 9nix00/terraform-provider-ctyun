package business

const (
	VpcesIpVersionIpv4      = "ipv4"
	VpcesIpVersionIpv6      = "ipv6"
	VpcesIpVersionDualStack = "dual_stack"
)

var VPCEsIpVersionMap = map[string]int32{
	VpcesIpVersionIpv4:      0,
	VpcesIpVersionIpv6:      1,
	VpcesIpVersionDualStack: 2,
}

var VPCEsIpVersionMapReverse = map[int32]string{
	0: VpcesIpVersionIpv4,
	1: VpcesIpVersionIpv6,
	2: VpcesIpVersionDualStack,
}
