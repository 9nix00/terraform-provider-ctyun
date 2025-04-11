package business

const (
	// 订购周期常量
	MonthCycleType    = "month"
	YearCycleType     = "year"
	OnDemandCycleType = "on_demand"

	//资源状态常量
	NatStatusStarted   = "started"   //启用
	NatStatusRenewed   = "renewed"   //续订
	NatStatusRefunded  = "refunded"  //退订
	NatStatusDestroyed = "destroyed" //销毁
	NatStatusFailed    = "failed"    //失败
	NatStatusStarting  = "starting"  //正在启动
	NatStatusChanged   = "changed"   //变配
	NatStatusExpired   = "expired"   //过期
	NatStatusUnknown   = "unknown"   //未知

	//  Nat规格
	SpecSmall      = 1
	SpecMedium     = 2
	SpecLarge      = 3
	SpecExtraLarge = 4

	// protocol
	ProtocolTcp = "tcp"
	ProtocolUdp = "udp"

	//DNAT运行状态
	DNatStateActive   = "active"
	DNatStateFreezing = "freezing"
	DNatStateCreating = "creating"

	SubnetTypeVPC    = 1
	SubnetTypeCustom = 2

	// SNAT创建状态
	SNatCreateStatusING  = "in_progress"
	SnatCreateStatusDone = "done"
)

var NatOrderCycleTypes = []string{
	MonthCycleType,
	YearCycleType,
	OnDemandCycleType,
}

var NatStatus = []string{
	NatStatusStarted,
	NatStatusRenewed,
	NatStatusRefunded,
	NatStatusDestroyed,
	NatStatusFailed,
	NatStatusStarting,
	NatStatusChanged,
	NatStatusExpired,
	NatStatusUnknown,
}

var NatSpecs = []int64{
	SpecSmall,
	SpecMedium,
	SpecLarge,
	SpecExtraLarge,
}

var DNatProtocols = []string{
	ProtocolTcp,
	ProtocolUdp,
}
var DNatStatus = []string{
	DNatStateActive,
	DNatStateFreezing,
	DNatStateCreating,
}

var SNatSubnetTypes = []int32{
	SubnetTypeVPC,
	SubnetTypeCustom,
}
