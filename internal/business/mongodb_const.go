package business

const (
	MongodbRunningStatusStarted     = 0  // 已启动
	MongodbRunningStatusRestarting  = 1  // 重启中
	MongodbRunningStatusBackup      = 2  // 备份操作中
	MongodbRunningStatusRecovery    = 3  // 操作恢复中
	MongodbRunningStatusTransferSSl = 4  // 转换ssl
	MongodbRunningStatusException   = 5  // 异常
	MongodbRunningStatusModify      = 6  // 修改参数组中
	MongodbRunningStatusFrozen      = 7  // 已冻结
	MongodbRunningStatusLogout      = 8  // 已注销
	MongodbRunningStatusProcessing  = 9  // 施工中
	MongodbRunningStatusFailed      = 10 // 施工失败
	MongodbRunningStatusUpgrading   = 11 // 扩容中
	MongodbRunningStatusSwitch      = 12 //主备切换中

	MongodbProdIDS34    = 10013001 // 3.4单机版
	MongodbProdIDS40    = 10013002 // 4.0单机版
	MongodbProdID3R34   = 10013003 // 3.4副本集三副本
	MongodbProdID3R40   = 10013004 // 4.0副本集三副本
	MongodbProdID5R34   = 10013005 // 3.4副本集五副本
	MongodbProdID5R40   = 10013006 // 4.0副本集五副本
	MongodbProdID7R34   = 10013007 // 3.4副本集七副本
	MongodbProdID7R40   = 10013008 // 4.0副本集七副本
	MongodbProdIDC34    = 10013009 // 3.4集群版
	MongodbProdIDC40    = 10013010 // 4.0集群版
	MongodbProdS42      = 10013011 // 4.2单机版
	MongodbProdID3R42   = 10013012 // 4.2副本集三副本
	MongodbProdID5R42   = 10013013 // 4.2副本集五副本
	MongodbProdID7R42   = 10013014 // 4.2副本集七副本
	MongodbProdIDC42    = 10013015 // 4.2集群版
	MongodbProdIDS50    = 10013016 // 5.0单机版
	MongodbProdID3R50   = 10013017 // 5.0副本集三副本
	MongodbProdID5R50   = 10013018 // 5.0副本集五副本
	MongodbProdID7R50   = 10013019 // 5.0副本集七副本
	MongodbProdIDC50    = 10013020 // 5.0集群版
	MongodbProdIDC60    = 10013021 // 6.0集群版
	MongodbProdID3R60   = 10013022 // 6.0副本集三副本
	MongodbProdID5R60   = 10013023 // 6.0副本集五副本
	MongodbProdID7R60   = 10013024 // 6.0副本集七副本
	MongodbProdIDS60    = 10013025 // 6.0单机版
	MongodbProdIDRead40 = 10013110 // 4.0只读版
	MongodbProdIDRead34 = 10013119 // 3.4只读版

	MongodbStorageTypeSSD       = "SSD"        // 超高IO
	MongodbStorageTypeSAS       = "SAS"        // 高IO
	MongodbStorageTypeSATA      = "SATA"       // 普通IO
	MongodbStorageTypeSSDGenric = "SSD-Genric" // 通用型SSD
)

var MongodbRunningStatus = []int32{
	MongodbRunningStatusStarted,
	MongodbRunningStatusRestarting,
	MongodbRunningStatusBackup,
	MongodbRunningStatusRecovery,
	MongodbRunningStatusTransferSSl,
	MongodbRunningStatusException,
	MongodbRunningStatusModify,
	MongodbRunningStatusFrozen,
	MongodbRunningStatusLogout,
	MongodbRunningStatusProcessing,
	MongodbRunningStatusFailed,
	MongodbRunningStatusUpgrading,
	MongodbRunningStatusSwitch,
}

var MongodbProdID = []int64{
	MongodbProdIDS34,
	MongodbProdIDS40,
	MongodbProdID3R34,
	MongodbProdID3R40,
	MongodbProdID5R34,
	MongodbProdID5R40,
	MongodbProdID7R34,
	MongodbProdID7R40,
	MongodbProdIDC34,
	MongodbProdIDC40,
	MongodbProdS42,
	MongodbProdID3R42,
	MongodbProdID5R42,
	MongodbProdID7R42,
	MongodbProdIDC42,
	MongodbProdIDS50,
	MongodbProdID3R50,
	MongodbProdID5R50,
	MongodbProdID7R50,
	MongodbProdIDC50,
	MongodbProdIDC60,
	MongodbProdID3R60,
	MongodbProdID5R60,
	MongodbProdID7R60,
	MongodbProdIDS60,
	MongodbProdIDRead40,
	MongodbProdIDRead34,
}

var MongodbStorageType = []string{
	MongodbStorageTypeSSD,
	MongodbStorageTypeSAS,
	MongodbStorageTypeSATA,
	MongodbStorageTypeSSDGenric,
}
