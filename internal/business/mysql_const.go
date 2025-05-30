package business

const (
	ProdIDSI57   = 10001003 // 单实例 single5.7版本
	ProdIDSI80   = 10001103 // 单实例 single8.0版本
	ProdIDSIRO57 = 10001005 // 单实例 single 只读5.7版本
	ProdIDSIRO80 = 10001105 // 单实例 single 只读8.0版本

	ProdIDMS57 = 10001001 // 一主一备 master-slave 5.7版本
	ProdIDMS80 = 10001101 // 一主一备 master-slave 8.0版本

	ProdIDM2S57 = 10001002 // 一主两备 master-2-slave 5.7版本
	ProdIDM2S80 = 10001102 // 一主两备 master-2-slave 8.0版本

	MysqlOrderStatusStarted      = 0  // 正常
	MysqlOrderStatusArrearage    = 1  //欠费暂停
	MysqlOrderStatusDestroyed    = 2  //已注销
	MysqlOrderStatusCreating     = 3  // 创建中
	MysqlOrderStatusFailed       = 4  // 施工失败
	MysqlOrderStatusExpire       = 5  //到期退订状态
	MysqlOrderStatusPause        = 6  // openapi暂停
	MysqlOrderStatusWaiting      = 7  // 创建完成等待变更单
	MysqlOrderStatusDestroy      = 8  // 待销毁
	MysqlOrderStatusManualPause  = 9  // 手动暂停
	MysqlOrderStatusManualRefund = 10 //手动退订

	MysqlRunningStatusStarted              = 0  // 正常
	MysqlRunningStatusRestarting           = 1  // 重启中
	MysqlRunningStatusBackup               = 2  // 备份中
	MysqlRunningStatusRecovering           = 3  // 恢复中
	MysqlRunningStatusModifying            = 4  // 修改参数中
	MysqlRunningStatusApplying             = 5  // 应用参数组中
	MysqlRunningStatusPreExpanding         = 6  // 扩容预处理中
	MysqlRunningStatusPreExpanded          = 7  // 扩容预处理完成
	MysqlRunningStatusUpdatePort           = 8  // 修改端口中
	MysqlRunningStatusMigrating            = 9  // 迁移中
	MysqlRunningStatusResetPassword        = 10 // 重置密码中
	MysqlRunningStatusUpdateCopyType       = 11 // 修改复制方式中
	MysqlRunningStatusPreShrinking         = 12 // 缩容预处理中
	MysqlRunningStatusPreShrinked          = 13 // 缩容预处理完成
	MysqlRunningStatusCoreUpgrade          = 15 // 内核小版本升级
	MysqlRunningStatusMigrateAz            = 17 // 迁移可用区中
	MysqlRunningStatusUpdateBackupConfig   = 18 // 修改备份配置中
	MysqlRunningStatusStopping             = 20 // 停止中
	MysqlRunningStatusStopped              = 21 // 已停止
	MysqlRunningStatusStarting             = 22 // 启动中
	MysqlRunningStatusWhiteListConfiguring = 26 // 白名单配置中

	MysqlBindEipStatusACTIVE              = "ACTIVE"                //已使用
	MysqlBindEipStatusDOWN                = "DOWN"                  //未使用
	MysqlBindEipStatusERROR               = "ERROR"                 //中间状态-异常
	MysqlBindEipStatusUPDATING            = "UPDATING"              //中间状态-更新中
	MysqlBindEipStatusBANDINGORUNBANGDING = "BANDING_OR_UNBANGDING" //中间状态-绑定或解绑中
	MysqlBindEipStatusDELETING            = "DELETING"              //中间状态-删除中
	MysqlBindEipStatusDELETED             = "DELETED"               //中间状态-已删除

	MysqlHostTypeS  = "1" // 通用型
	MysqlHostTypeC  = "2" // 计算增强型
	MysqlHostTypeM  = "3" // 内存增强型
	MysqlHostTypeHS = "4" //
	MysqlHostTypeHC = "5"
	MysqlHostTypeHM = "6"
	MysqlHostTypeKS = "7"
	MysqlHostTypeKM = "8"
	MysqlHostTypeKC = "9"

	MysqlBillModeCycle    = "1"
	MysqlBillModeOnDemand = "2"
	EipStatusUnbind       = 0 // 未绑定
	EipStatusBind         = 1 //已绑定
	EipStatusBinding      = 2 // 正在绑定

	ProdTypeUNKNOWN = "0" // UNKNOWN
	ProdTypeRDS     = "1" //RDS
	ProdTypeNoSql   = "2" // NoSql
	ProdTypeTOOL    = "3" // TOOL
	ProdTypeMemDB   = "4" // MemDB

	InstanceTypeNormal  = "1" // 通用型
	InstanceTypeCompute = "2" // 计算增强型
	InstanceTypeMemory  = "3" // 内存优化型
	InstanceTypeThrough = "4" //直通

	ProdCodeHBASE      = "HBASE"      // HBASE
	ProdCodeDDS        = "DDS"        // DDS
	ProdCodeMYSQL      = "MYSQL"      // MYSQL
	ProdCodePOSTGRESQL = "POSTGRESQL" //POSTGRESQL
	ProdCodeSQLSERVER  = "SQLSERVER"  // SQLSERVER
)

var ProdType = []string{
	ProdTypeUNKNOWN,
	ProdTypeRDS,
	ProdTypeNoSql,
	ProdTypeTOOL,
	ProdTypeMemDB,
}

var InstanceType = []string{
	InstanceTypeNormal,
	InstanceTypeCompute,
	InstanceTypeMemory,
	InstanceTypeThrough,
}

var ProdCode = []string{
	ProdCodeHBASE,
	ProdCodeDDS,
	ProdCodeMYSQL,
	ProdCodePOSTGRESQL,
	ProdCodeSQLSERVER,
}
var MysqlProdIDs = []int64{
	ProdIDSI57,
	ProdIDSI80,
	ProdIDSIRO57,
	ProdIDSIRO80,
	ProdIDMS57,
	ProdIDMS80,
	ProdIDM2S57,
	ProdIDM2S80,
}

var MysqlOrderStatus = []int32{
	MysqlOrderStatusStarted,
	MysqlOrderStatusArrearage,
	MysqlOrderStatusDestroyed,
	MysqlOrderStatusCreating,
	MysqlOrderStatusFailed,
	MysqlOrderStatusExpire,
	MysqlOrderStatusPause,
	MysqlOrderStatusWaiting,
	MysqlOrderStatusDestroy,
	MysqlOrderStatusManualPause,
	MysqlOrderStatusManualRefund,
}

var MysqlRunningStatus = []int32{
	MysqlRunningStatusStarted,
	MysqlRunningStatusRestarting,
	MysqlRunningStatusBackup,
	MysqlRunningStatusRecovering,
	MysqlRunningStatusModifying,
	MysqlRunningStatusApplying,
	MysqlRunningStatusPreExpanding,
	MysqlRunningStatusPreExpanded,
	MysqlRunningStatusUpdatePort,
	MysqlRunningStatusMigrating,
	MysqlRunningStatusResetPassword,
	MysqlRunningStatusUpdateCopyType,
	MysqlRunningStatusPreShrinking,
	MysqlRunningStatusPreShrinked,
	MysqlRunningStatusCoreUpgrade,
	MysqlRunningStatusMigrateAz,
	MysqlRunningStatusUpdateBackupConfig,
	MysqlRunningStatusStopping,
	MysqlRunningStatusStopped,
	MysqlRunningStatusStarting,
	MysqlRunningStatusWhiteListConfiguring,
}

var MysqlBindEipStatus = []string{
	MysqlBindEipStatusACTIVE,
	MysqlBindEipStatusDOWN,
	MysqlBindEipStatusERROR,
	MysqlBindEipStatusUPDATING,
	MysqlBindEipStatusBANDINGORUNBANGDING,
	MysqlBindEipStatusDELETING,
	MysqlBindEipStatusDELETED,
}

var MysqlHostType = []string{
	MysqlHostTypeS,
	MysqlHostTypeC,
	MysqlHostTypeM,
	MysqlHostTypeHS,
	MysqlHostTypeHC,
	MysqlHostTypeHM,
	MysqlHostTypeKS,
	MysqlHostTypeKM,
	MysqlHostTypeKC,
}

var MysqlBillMode = []string{
	MysqlBillModeCycle,
	MysqlBillModeOnDemand,
}
