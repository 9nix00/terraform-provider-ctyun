package business

const (
	LbResourceTypeInternal = "internal" //内网负载均衡
	LbResourceTypeExternal = "external" //公网负载均衡

	AdminStatusDown   = "down"   //管理状态: DOWN
	AdminStatusActive = "active" //管理状态: ACTIVE
)

var LbResourceType = []string{LbResourceTypeInternal, LbResourceTypeExternal}
var AdminStatusName = []string{AdminStatusActive, AdminStatusDown}
