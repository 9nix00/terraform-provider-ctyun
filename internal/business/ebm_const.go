package business

import "terraform-provider-ctyun/internal/utils"

const (
	EbmExtIpNotUse     = "not_use"
	EbmExtIpAutoAssign = "auto_assign"
	EbmExtIpUseExist   = "use_exist"

	EbmOrderOnCycle  = "ORDER_ON_CYCLE"
	EbmOrderOnDemand = "ORDER_ON_DEMAND"

	EbmStatusCreating          = "creating"
	EbmStatusStarting          = "starting"
	EbmStatusRunning           = "running"
	EbmStatusStopping          = "stopping"
	EbmStatusStopped           = "stopped"
	EbmStatusRestarting        = "restarting"
	EbmStatusError             = "error"
	EbmStatusResettingPassword = "resetting_password"
	EbmStatusResettingHostname = "resetting_hostname"

	EbmSystemDiskType = "system"
	EbmDataDiskType   = "data"
)

const (
	EbmExtIpMapScene1 = iota
)

var EbmExtIp = []string{
	EbmExtIpNotUse,
	EbmExtIpAutoAssign,
	EbmExtIpUseExist,
}

var EbmExtIpMap = utils.Must(
	[]any{
		EbmExtIpNotUse,
		EbmExtIpAutoAssign,
		EbmExtIpUseExist,
	},
	map[utils.Scene][]any{
		EbmExtIpMapScene1: {
			"0",
			"1",
			"2",
		},
	},
)

var EbmDiskTypes = []string{
	EbsDiskTypeSata,
	EbsDiskTypeSas,
	EbsDiskTypeSsd,
	EbsDiskTypeSsdGenric,
}
