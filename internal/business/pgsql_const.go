package business

const (
	BillModeCycle          = "1"
	BillModelOnDemand      = "2"
	CaseSensitiveTrue      = "0"
	CaseSensitiveFalse     = "1"
	CaseSensitiveUnCertain = "2"
	OsTypePure             = "0"  // 裸机
	OsTypeWindows          = "1"  // Windows
	OsTypeCentos           = "2"  // Centos
	OsTypeUbuntu           = "3"  // Ubuntu
	OsTypeAndroid          = "4"  // Android
	OsTypeRedHat           = "5"  // RedHat
	OsTypeKylin            = "6"  // kylin
	OsTypeUos              = "7"  // Uos
	OsTypeSuse             = "8"  // Suse
	OsTypeAsianux          = "9"  // Asianus
	OsTypeOpenEuler        = "10" // OpenEuler
	OsTypeCtyunOS          = "11" // CtyunOS
	OsTypeEuler            = "12" // Euler

	PgsqlProdRunningStatusStarted             = 0
	pgsqlProdRunningStatusRestarting          = 1
	PgsqlProdRunningStatusBackup              = 2
	PgsqlProdRunningStatusRecovering          = 3
	PgsqlProdRunningStatusStopped             = 1001
	PgsqlProdRunningStatusRecoveryFailed      = 1006
	PgsqlProdRunningStatusVipUnavailable      = 1007
	PgsqlProdRunningStatusGatewayUnavailable  = 1008
	PgsqlProdRunningStatusMasterUnavailable   = 1009
	PgsqlProdRunningStatusSlaveUnavailable    = 1010
	PgsqlProdRunningStatusInstanceMaintenance = 1021
	PgsqlProdRunningStatusActivating          = 2000
	PgsqlProdRunningStatusUnsubscribed        = 2002
	PgsqlProdRunningStatusExpanding           = 2005
	PgsqlProdRunningStatusFreeze              = 2011

	PgsqlProdOrderStatusRunning    = 0
	PgsqlProdOrderStatusFreeze     = 1
	PgsqlProdOrderStatusDelete     = 2
	PgsqlProdOrderStatusProcessing = 3
	PgsqlProdOrderStatusFailure    = 4
	PgsqlProdOrderStatusExpanding  = 5

	PgsqlBindEipStatusACTIVE              = "ACTIVE"                //已使用
	PgsqlBindEipStatusDOWN                = "DOWN"                  //未使用
	PgsqlBindEipStatusERROR               = "ERROR"                 //中间状态-异常
	PgsqlBindEipStatusUPDATING            = "UPDATING"              //中间状态-更新中
	PgsqlBindEipStatusBANDINGORUNBANGDING = "BANDING_OR_UNBANGDING" //中间状态-绑定或解绑中
	PgsqlBindEipStatusDELETING            = "DELETING"              //中间状态-删除中
	PgsqlBindEipStatusDELETED             = "DELETED"

	PgsqlProdIDS1222  = 10003011
	PgsqlProdIDMS1222 = 10003012
	PgsqlProdIDS1417  = 10003013
)

var PgsqlBillModes = []string{
	BillModeCycle,
	BillModelOnDemand,
}

var PgsqlCaseSensitive = []string{
	CaseSensitiveTrue,
	CaseSensitiveFalse,
	CaseSensitiveUnCertain,
}

var PgsqlOsType = []string{
	OsTypePure,
	OsTypeWindows,
	OsTypeCentos,
	OsTypeUbuntu,
	OsTypeAndroid,
	OsTypeRedHat,
	OsTypeKylin,
	OsTypeUos,
	OsTypeSuse,
	OsTypeAsianux,
	OsTypeOpenEuler,
	OsTypeCtyunOS,
	OsTypeEuler,
}

var PgsqlProdOrderStatus = []int32{
	PgsqlProdOrderStatusRunning,
	PgsqlProdOrderStatusFreeze,
	PgsqlProdOrderStatusDelete,
	PgsqlProdOrderStatusProcessing,
	PgsqlProdOrderStatusFailure,
	PgsqlProdOrderStatusExpanding,
}

var PgsqlProdRunningStatus = []int32{
	PgsqlProdRunningStatusStarted,
	pgsqlProdRunningStatusRestarting,
	PgsqlProdRunningStatusBackup,
	PgsqlProdRunningStatusRecovering,
	PgsqlProdRunningStatusStopped,
	PgsqlProdRunningStatusRecoveryFailed,
	PgsqlProdRunningStatusVipUnavailable,
	PgsqlProdRunningStatusGatewayUnavailable,
	PgsqlProdRunningStatusMasterUnavailable,
	PgsqlProdRunningStatusSlaveUnavailable,
	PgsqlProdRunningStatusInstanceMaintenance,
	PgsqlProdRunningStatusActivating,
	PgsqlProdRunningStatusUnsubscribed,
	PgsqlProdRunningStatusExpanding,
	PgsqlProdRunningStatusFreeze,
}

var PgsqlBindEipStatus = []string{
	MysqlBindEipStatusACTIVE,
	MysqlBindEipStatusDOWN,
	MysqlBindEipStatusERROR,
	MysqlBindEipStatusUPDATING,
	MysqlBindEipStatusBANDINGORUNBANGDING,
	MysqlBindEipStatusDELETING,
	MysqlBindEipStatusDELETED,
}
