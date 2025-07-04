package business

const (
	MongodbOrderStatusStarted = 0

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

	MongodbNodeTypeMongos = "mongos" // 集群节点
	MongodbNodeTypeShard  = "shard"  // 分片节点
	MongodbNodeTypeConfig = "config"
	MongodbNodeTypeMs     = "ms"     // 副本集
	MongodbNodeTypeS      = "s"      // 单机版
	MongodbNodeBackup     = "backup" // 备份节点

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

var MongodbProdIDDict = map[string]int64{
	"Single34":    MongodbProdIDS34,
	"Single40":    MongodbProdIDS40,
	"Replica3R34": MongodbProdID3R34, // 副本集，3副本，3.4版本
	"Replica3R40": MongodbProdID3R40,
	"Replica5R34": MongodbProdID5R34,
	"Replica5R40": MongodbProdID5R40,
	"Replica7R34": MongodbProdID7R34,
	"Replica7R40": MongodbProdID7R40,
	"Cluster34":   MongodbProdIDC34,
	"Cluster40":   MongodbProdIDC40,
	"Single42":    MongodbProdS42,
	"Replica3R42": MongodbProdID3R42,
	"Replica5R42": MongodbProdID5R42,
	"Replica7R42": MongodbProdID7R42,
	"Cluster42":   MongodbProdIDC42,
	"Single50":    MongodbProdIDS50,
	"Replica3R50": MongodbProdID3R50,
	"Replica5R50": MongodbProdID5R50,
	"Replica7R50": MongodbProdID7R50,
	"Cluster50":   MongodbProdIDC50,
	"Cluster60":   MongodbProdIDC60,
	"Replica3R60": MongodbProdID3R60,
	"Replica5R60": MongodbProdID5R60,
	"Replica7R60": MongodbProdID7R60,
	"Single60":    MongodbProdIDS60,
	"ReadOnly40":  MongodbProdIDRead40,
	"ReadOnly34":  MongodbProdIDRead34,
}

var MongodbProdIDRevDict = map[int64]string{
	MongodbProdIDS34:    "Single34",
	MongodbProdIDS40:    "Single40",
	MongodbProdID3R34:   "Replica3R34",
	MongodbProdID3R40:   "Replica3R40",
	MongodbProdID5R34:   "Replica5R34",
	MongodbProdID5R40:   "Replica5R40",
	MongodbProdID7R34:   "Replica7R34",
	MongodbProdID7R40:   "Replica7R40",
	MongodbProdIDC34:    "Cluster34",
	MongodbProdIDC40:    "Cluster40",
	MongodbProdS42:      "Single42",
	MongodbProdID3R42:   "Replica3R42",
	MongodbProdID5R42:   "Replica5R42",
	MongodbProdID7R42:   "Replica7R42",
	MongodbProdIDC42:    "Cluster42",
	MongodbProdIDS50:    "Single50",
	MongodbProdID3R50:   "Replica3R50",
	MongodbProdID5R50:   "Replica5R50",
	MongodbProdID7R50:   "Replica7R50",
	MongodbProdIDC50:    "Cluster50",
	MongodbProdIDC60:    "Cluster60",
	MongodbProdID3R60:   "Replica3R60",
	MongodbProdID5R60:   "Replica5R60",
	MongodbProdID7R60:   "Replica7R60",
	MongodbProdIDS60:    "Single60",
	MongodbProdIDRead40: "ReadOnly40",
	MongodbProdIDRead34: "ReadOnly34",
}

var MongodbProdIDs = []string{
	"Single34",
	"Single40",
	"Replica3R34",
	"Replica3R40",
	"Replica5R34",
	"Replica5R40",
	"Replica7R34",
	"Replica7R40",
	"Cluster34",
	"Cluster40",
	//"Single42",
	//"Replica3R42",
	//"Replica5R42",
	//"Replica7R42",
	//"Cluster42",
	//"Single50",
	//"Replica3R50",
	//"Replica5R50",
	//"Replica7R50",
	//"Cluster50",
	//"Cluster60",
	//"Replica3R60",
	//"Replica5R60",
	//"Replica7R60",
	//"Single60",
	//"ReadOnly40",
	//"ReadOnly34",
}

var MongodbNodeTypeDict = map[string]string{
	"Single34":    "master",
	"Single40":    "master",
	"Replica3R34": "master", // 副本集，3副本，3.4版本
	"Replica3R40": "master",
	"Replica5R34": "master",
	"Replica5R40": "master",
	"Replica7R34": "master",
	"Replica7R40": "master",
	"Cluster34":   "master",
	"Cluster40":   "master",
	"Single42":    "master",
	"Replica3R42": "master",
	"Replica5R42": "master",
	"Replica7R42": "master",
	"Cluster42":   "master",
	"Single50":    "master",
	"Replica3R50": "master",
	"Replica5R50": "master",
	"Replica7R50": "master",
	"Cluster50":   "master",
	"Cluster60":   "master",
	"Replica3R60": "master",
	"Replica5R60": "master",
	"Replica7R60": "master",
	"Single60":    "master",
	"ReadOnly40":  "ReadNode",
	"ReadOnly34":  "ReadNode",
}

var MongodbProdVersionDict = map[string]string{
	"Single34":    "3.4",
	"Single40":    "4.0",
	"Replica3R34": "3.4", // 副本集，3副本，3.4版本
	"Replica3R40": "3.4",
	"Replica5R34": "3.4",
	"Replica5R40": "4.0",
	"Replica7R34": "3.4",
	"Replica7R40": "4.0",
	"Cluster34":   "3.4",
	"Cluster40":   "4.0",
	"Single42":    "4.2",
	"Replica3R42": "4.2",
	"Replica5R42": "4.2",
	"Replica7R42": "4.2",
	"Cluster42":   "4.2",
	"Single50":    "5.0",
	"Replica3R50": "5.0",
	"Replica5R50": "5.0",
	"Replica7R50": "5.0",
	"Cluster50":   "5.0",
	"Cluster60":   "6.0",
	"Replica3R60": "6.0",
	"Replica5R60": "6.0",
	"Replica7R60": "6.0",
	"Single60":    "6.0",
	"ReadOnly40":  "4.0",
	"ReadOnly34":  "3.4",
}

var MongodbStorageType = []string{
	MongodbStorageTypeSSD,
	MongodbStorageTypeSAS,
	MongodbStorageTypeSATA,
	MongodbStorageTypeSSDGenric,
}

var MongodbInstanceSeriesDict = map[string]string{
	"S": "1",
	"C": "2",
	"M": "3",
}

var MongodbNodeType = []string{
	MongodbNodeTypeMongos,
	MongodbNodeTypeShard,
	MongodbNodeTypeConfig,
	MongodbNodeTypeMs,
	MongodbNodeTypeS,
	MongodbNodeBackup,
}
