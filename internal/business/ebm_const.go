package business

import "terraform-provider-ctyun/internal/utils"

const (
	EbmExtIpNotUse     = "not_use"
	EbmExtIpAutoAssign = "auto_assign"
	EbmExtIpUseExist   = "use_exist"

	EbmOrderOnCycle  = "order_on_cycle"
	EbmOrderOnDemand = "order_on_demand"

	EbmCycleTypeMonth = "month"
	EbmCycleTypeYear  = "year"

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
	EbmDataDiskType   = "system"
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
