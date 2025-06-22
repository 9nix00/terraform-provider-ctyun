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

	PgsqlProdIDS1222    = 10003011
	PgsqlProdIDMS1222   = 10003012
	PgsqlProdIDS1417    = 10003013
	PgsqlProdIDMS1417   = 10003014
	PgsqlProdIDS1320    = 10003015
	PgsqlProdIDMS1320   = 10003016
	PgsqlProdIDRead1222 = 10003017
	PgsqlProdIDRead1320 = 10003018
	PgsqlProdIDRead1417 = 10003019
	PgsqlProdIDS1512    = 10003021
	PgsqlProdIDMS1512   = 10003022
	PgsqlProdIDRead1512 = 10003023
	PgsqlProdIDM2S1222  = 10003024
	PgsqlProdIDM2S1417  = 10003025
	PgsqlProdIDM2S1320  = 10003026
	PgsqlProdIDM2S1512  = 10003027
	PgsqlProdIDS168     = 10003028
	PgsqlProdIDMS168    = 10003029
	PgsqlProdIDM2S168   = 10003031
	PgsqlProdIDRead168  = 10003030
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

var PgsqlProdID = []int64{
	PgsqlProdIDS1222,
	PgsqlProdIDMS1222,
	PgsqlProdIDS1417,
	PgsqlProdIDMS1417,
	PgsqlProdIDS1320,
	PgsqlProdIDMS1320,
	PgsqlProdIDRead1222,
	PgsqlProdIDRead1320,
	PgsqlProdIDRead1417,
	PgsqlProdIDS1512,
	PgsqlProdIDMS1512,
	PgsqlProdIDRead1512,
	PgsqlProdIDM2S1222,
	PgsqlProdIDM2S1417,
	PgsqlProdIDM2S1320,
	PgsqlProdIDM2S1512,
	PgsqlProdIDS168,
	PgsqlProdIDMS168,
	PgsqlProdIDM2S168,
	PgsqlProdIDRead168,
}
